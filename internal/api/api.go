package api

import (
	"crypto/rand"
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
	actions, err := Action.GetActions(ctx.Request().Body)
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
	mockup := Mockup.InitMockup(taskID, api.Config)
	defer func() {
		err := recover()
		if err != nil {
			Logger.Debug("Panic detected: %v", err)
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
