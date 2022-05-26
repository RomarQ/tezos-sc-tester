package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/romarq/tezos-sc-tester/internal/business/action"
	"github.com/romarq/tezos-sc-tester/internal/config"
	"github.com/romarq/tezos-sc-tester/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
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

	t.Run("Run FA2 test actions", func(t *testing.T) {
		request, err := getTestData("fa2_actions.json")
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

		for _, response := range actionResponses {
			assert.Equal(t, response.Status, action.Success, response.Result)
		}

		snapshotBytes, err := json.MarshalIndent(actionResponses, "", "  ")
		assert.Nil(t, err, "Must not fail")

		assert.NoError(t, saveSnapshot("fa2_actions_response.json", snapshotBytes))
	})
}

func getTestData(fileName string) ([]byte, error) {
	wd, _ := os.Getwd()
	contract_file_path := path.Join(wd, "__test_data__", fileName)
	return os.ReadFile(contract_file_path)
}

// Utility for saving test snapshots
func saveSnapshot(fileName string, bytes []byte) error {
	wd, _ := os.Getwd()
	filePath := path.Join(wd, "__test_data__/snapshots", fileName)
	return os.WriteFile(filePath, bytes, 0644)
}
