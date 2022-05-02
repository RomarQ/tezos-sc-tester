package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/romarq/visualtez-testing/internal/business/action"
	"github.com/romarq/visualtez-testing/internal/config"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/romarq/visualtez-testing/internal/utils"
	"github.com/stretchr/testify/assert"
)

// Utility for saving test snapshots
func saveSnapshot(fileName string, bytes []byte) {
	wd, _ := os.Getwd()
	filePath := path.Join(wd, "__test_data__", fileName)
	os.WriteFile(filePath, bytes, 0644)
}

func TestRunTest(t *testing.T) {
	const TESTING_URL = "/testing"

	api := InitTestingAPI(config.Config{
		Log: config.LogConfig{
			Location: "../../.tmp_test/api.log",
		},
		Tezos: config.TezosConfig{
			TezosClient:     "../../tezos-bin/amd64/tezos-client",
			DefaultProtocol: "ProtoALphaALphaALphaALphaALphaALphaALphaALphaDdp3zK",
			BaseDirectory:   "../../tezos-bin",
			RevealFee:       1000,
			Originator:      "bootstrap2",
		},
	})
	logger.SetupLogger(api.Config.Log.Location, api.Config.Log.Level)

	t.Run("Perform a valid request and validate the response", func(t *testing.T) {
		request, err := getTestData("valid_request.json")
		assert.Nil(t, err, "Must not fail")

		e := echo.New()
		req := httptest.NewRequest(echo.POST, TESTING_URL, bytes.NewReader(request))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)

		err = api.RunTest(ctx)
		assert.Nil(t, err, "Must not fail")
		assert.Equal(t, rec.Code, 200)

		var actionResponses []action.ActionResult
		err = json.Unmarshal(rec.Body.Bytes(), &actionResponses)
		assert.Nil(t, err, "Must not fail")

		assert.Equal(t, len(actionResponses), 7, "Expects 6 action results")

		for _, response := range actionResponses {
			assert.Equal(t, response.Status, action.Success, response.Result)
		}
	})

	t.Run("Create Implicit Account",
		func(t *testing.T) {
			CreateImplicitAccountAction := map[string]interface{}{
				"kind": action.CreateImplicitAccount,
				"payload": map[string]interface{}{
					"name":    "bob",
					"balance": "10",
				},
			}
			actions, _ := json.Marshal(
				map[string]interface{}{
					"actions": []map[string]interface{}{CreateImplicitAccountAction},
				},
			)
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
			assert.Equal(t, actionResult.Status, action.Success, "Action status must be (success)")
			assert.Equal(t, utils.PrettifyJSON(actionResult.Action), utils.PrettifyJSON(CreateImplicitAccountAction), "Validate action request")
			assert.Contains(t, fmt.Sprintf("%v", actionResult.Result), "tz1", actionResult.Result)
		})

	t.Run("Originate Contract",
		func(t *testing.T) {
			OriginateContractAction := map[string]interface{}{
				"kind": action.OriginateContract,
				"payload": map[string]interface{}{
					"name":    "contract_1",
					"balance": "10",
					"code": json.RawMessage(`
						[
							{
								"args": [
									{
										"prim": "unit"
									}
								],
								"prim": "storage"
							},
							{
								"args": [
									{
										"prim": "unit"
									}
								],
								"prim": "parameter"
							},
							{
								"args": [
									[
										{
											"prim": "DROP"
										},
										{
											"prim": "UNIT"
										},
										{
											"args": [
												{
													"prim": "operation"
												}
											],
											"prim": "NIL"
										},
										{
											"prim": "PAIR"
										}
									]
								],
								"prim": "code"
							}
						]
					`),
					"storage": map[string]string{
						"prim": "Unit",
					},
				},
			}
			actions, _ := json.Marshal(
				map[string]interface{}{
					"actions": []map[string]interface{}{OriginateContractAction},
				},
			)
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
			assert.Equal(t, utils.PrettifyJSON(actionResult.Action), utils.PrettifyJSON(OriginateContractAction), "Validate action request")
			assert.Contains(t, fmt.Sprintf("%v", actionResult.Result), "KT1", actionResult.Result)
		})
}

func getTestData(fileName string) ([]byte, error) {
	wd, _ := os.Getwd()
	contract_file_path := path.Join(wd, "__test_data__", fileName)
	return os.ReadFile(contract_file_path)
}
