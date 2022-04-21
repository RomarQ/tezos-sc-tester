package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	Mockup "github.com/romarq/visualtez-testing/internal/business"
	Action "github.com/romarq/visualtez-testing/internal/business/action"
	Config "github.com/romarq/visualtez-testing/internal/config"
	Error "github.com/romarq/visualtez-testing/internal/error"
	Logger "github.com/romarq/visualtez-testing/internal/logger"
)

type TestingAPI struct {
	Config Config.Config
}

type testSuiteRequest struct {
	Protocol string            `json:"protocol"`
	Actions  []json.RawMessage `json:"actions"`
}

func InitTestingAPI(config Config.Config) TestingAPI {
	api := TestingAPI{
		Config: config,
	}
	return api
}

// RunTest - Run a test (`/testing`)
// @ID post-testing
// @Description Run test actions
// @Consumes json
// @Produce json
// @Success 200
// @Failure default {object} Error
// @Router /testing [post]
func (api *TestingAPI) RunTest(ctx echo.Context) error {
	var request testSuiteRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&request); err != nil {
		return Error.HttpError(http.StatusBadRequest, "request body is invalid.")
	}

	actions, err := Action.GetActions(request.Actions)
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
		return Error.HttpError(http.StatusInternalServerError, "Something went wrong.")
	}

	taskID := fmt.Sprintf("task_%d", prime)
	mockup := Mockup.InitMockup(taskID, request.Protocol, api.Config)
	defer func() {
		err := recover()
		if err != nil {
			Logger.Debug("Panic detected:", err)
		}
		// Teardown on exit
		mockup.Teardown()
	}()

	// Bootstrap mockup
	err = mockup.Bootstrap()
	if err != nil {
		Logger.Debug("Something went wrong: %s", err)
		return Error.HttpError(http.StatusInternalServerError, "Could not bootstrap test environment.")
	}

	return ctx.JSON(http.StatusOK, Action.ApplyActions(mockup, actions))
}
