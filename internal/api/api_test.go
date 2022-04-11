package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/romarq/visualtez-testing/internal/business/action"
	"github.com/romarq/visualtez-testing/internal/config"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestRunTest(t *testing.T) {
	const TESTING_URL = "/testing"
	t.Run("Create Implicit Account",
		func(t *testing.T) {
			api := InitTestingAPI(config.Config{
				Log: config.LogConfig{
					Location: "../../tmp_test/api.log",
				},
				Tezos: config.TezosConfig{
					DefaultProtocol: "ProtoALphaALphaALphaALphaALphaALphaALphaALphaDdp3zK",
					BaseDirectory:   "../../tezos-bin",
				},
			})
			logger.SetupLogger(api.Config.Log.Location, api.Config.Log.Level)

			CreateImplicitAccountAction := map[string]interface{}{
				"kind": action.CreateImplicitAccount,
				"payload": map[string]interface{}{
					"name":    "bob",
					"balance": float64(10),
				},
			}
			actions, _ := json.Marshal([]map[string]interface{}{CreateImplicitAccountAction})
			e := echo.New()
			req := httptest.NewRequest(echo.POST, TESTING_URL, bytes.NewReader(actions))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ctx := e.NewContext(req, rec)

			err := api.RunTest(ctx)
			assert.Nil(t, err, "Must not fail")
			assert.Equal(t, rec.Code, 200)

			var result []action.ActionResult
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result), "Unmarshal should not fail")
			assert.Len(t, result, 1, "Must only contain a single action result")

			actionResult := result[0]
			assert.Equal(t, actionResult.Status, action.Success, "Action result must be (success)")
			assert.Equal(t, actionResult.Kind, action.CreateImplicitAccount, "Validate action kind")
			assert.Equal(t, actionResult.Action, CreateImplicitAccountAction["payload"], "Validate action payload")
			assert.Equal(t, fmt.Sprintf("%v", actionResult.Result["address"])[0:3], "tz1", "Validate result payload")
		})
}
