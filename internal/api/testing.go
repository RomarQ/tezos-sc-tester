package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	Action "github.com/romarq/visualtez-testing/internal/action"
	Logger "github.com/romarq/visualtez-testing/internal/logger"
	"github.com/tidwall/gjson"
)

type TestRequest struct {
	Kind   Action.ActionKind `json:"kind"`
	Action interface{}       `json:"action"`
}

type TestingAPI struct{}

func InitTestingAPI() TestingAPI {
	api := TestingAPI{}
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
	actions, err := getActions(ctx.Request().Body)
	if err != nil {
		return HTTPError(ctx, http.StatusBadRequest, err.Error())
	}

	Logger.Debug("%v", actions...)

	return nil
}

// Unmarshal actions
func getActions(body io.ReadCloser) ([]interface{}, error) {
	rawActions := make([]json.RawMessage, 0)

	err := json.NewDecoder(body).Decode(&rawActions)
	if err != nil {
		return nil, err
	}

	actions := make([]interface{}, 0)
	for _, rawAction := range rawActions {
		kind := gjson.GetBytes(rawAction, `kind`)
		switch kind.String() {
		default:
			return nil, fmt.Errorf("Unexpected action kind (%s).", kind)
		case string(Action.CreateImplicitAccount):
			action := Action.CreateImplicitAccountAction{}
			err = action.Unmarshal(rawAction)
			if err != nil {
				return nil, err
			}
			actions = append(actions, action)
		}
	}

	return actions, err
}
