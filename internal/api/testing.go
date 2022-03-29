package api

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	Mockup "github.com/romarq/visualtez-testing/internal/business"
	Action "github.com/romarq/visualtez-testing/internal/business/action"
	Config "github.com/romarq/visualtez-testing/internal/config"
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
// @ID get-run-a-test
// @Description Run a test
// @Produce json
// @Consumes json
// @Success 200
// @Failure default {object} Error
// @Router /testing [post]
func (api *TestingAPI) RunTest(ctx echo.Context) error {
	actions, err := Action.GetActions(ctx.Request().Body)
	if err != nil {
		return HTTPError(http.StatusBadRequest, err.Error())
	}

	prime, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		return HTTPError(http.StatusInternalServerError, "Something went wrong.")
	}

	taskID := fmt.Sprintf("task_%d", prime)
	mockup := Mockup.InitMockup(taskID, api.Config)
	// Teardown on exit
	defer mockup.Teardown()

	// Boostrap mockup
	err = mockup.Bootstrap()
	if err != nil {
		return HTTPError(http.StatusInternalServerError, "Could not bootstrap test environment.")
	}

	Logger.Debug("%s %v", fmt.Sprintf("%d", prime), actions)

	return ctx.JSON(http.StatusOK, Action.ApplyActions(mockup, taskID, actions))
}
