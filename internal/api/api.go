package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	Mockup "github.com/romarq/tezos-sc-tester/internal/business"
	"github.com/romarq/tezos-sc-tester/internal/business/action"
	"github.com/romarq/tezos-sc-tester/internal/config"
	Error "github.com/romarq/tezos-sc-tester/internal/error"
	"github.com/romarq/tezos-sc-tester/internal/logger"
)

type testingAPI struct {
	Config config.Config
}

type testSuiteRequest struct {
	Protocol string          `json:"protocol"`
	Actions  []action.Action `json:"actions"`
}

// InitTestingAPI initializes the testing API
func InitTestingAPI(config config.Config) testingAPI {
	api := testingAPI{
		Config: config,
	}
	return api
}

// RunTest - Run a test (`/testing`) godoc
// @Summary  Run a test
// @ID       post-testing
// @Accept   json
// @Produce  json
// @Param    request  body      testSuiteRequest     true  "Test Request"
// @Success  200      {array}   action.ActionResult  "Success"
// @Failure  409      {object}  Error.Error          "Fail"
// @Router   /testing [post]
func (api *testingAPI) RunTest(ctx echo.Context) error {
	var mockup Mockup.Mockup
	defer func() {
		err := recover()
		if err != nil {
			logger.Debug("got an unexpected panic: %s", err)
		}
		// Teardown on exit
		mockup.Teardown()
	}()

	var request testSuiteRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&request); err != nil {
		return Error.HttpError(http.StatusBadRequest, "request body is invalid.")
	}

	// Parse test actions
	actions, err := action.GetActions(request.Actions)
	if err != nil {
		switch err.(type) {
		default:
			return Error.HttpError(http.StatusBadRequest, err.Error())
		case *echo.HTTPError:
			return err
		}
	}

	prime, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		logger.Debug("could not generate random prime. %s", err.Error())
		return Error.HttpError(http.StatusInternalServerError, "Something went wrong.")
	}

	taskID := fmt.Sprintf("task_%d", prime)
	mockup = Mockup.InitMockup(taskID, request.Protocol, api.Config)

	// Bootstrap mockup
	err = mockup.Bootstrap()
	if err != nil {
		logger.Debug("something went wrong. %s", err.Error())
		return Error.HttpError(http.StatusInternalServerError, "Could not bootstrap test environment.")
	}

	return ctx.JSON(http.StatusOK, action.ApplyActions(mockup, actions))
}
