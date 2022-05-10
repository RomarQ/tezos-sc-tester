package michelson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMichelineOfJSON(t *testing.T) {

	runTest := func(t *testing.T, file string) {
		jsonBytes, err := getTestData(fmt.Sprintf("%s.json", file))
		assert.NoError(t, err)
		michelineBytes, err := getTestData(fmt.Sprintf("%s.tz", file))
		assert.NoError(t, err)

		micheline, err := MichelineOfJSON(jsonBytes)
		assert.Nil(t, err, "Must not fail")

		assert.Equal(t, micheline, strings.Trim(string(michelineBytes), "\n"), "Decode JSON to Micheline")
	}

	t.Run("Convert JSON to Micheline", func(t *testing.T) {
		runTest(t, "positive_int")
		runTest(t, "negative_int")
		runTest(t, "string")
		runTest(t, "bytes")
		runTest(t, "sequence")

		// Contracts

		runTest(t, "simple_contract")
		runTest(t, "contract_with_multiply_entrypoints")
		runTest(t, "fa2_contract")
	})
}

func TestJSONOfMicheline(t *testing.T) {

	runTest := func(t *testing.T, file string) {
		jsonBytes, err := getTestData(fmt.Sprintf("%s.json", file))
		assert.NoError(t, err)
		michelineBytes, err := getTestData(fmt.Sprintf("%s.tz", file))
		assert.NoError(t, err)

		j, err := JSONOfMicheline(string(michelineBytes))
		assert.Nil(t, err, "Must not fail")

		assert.Equal(t, PrettifyJSON(j), PrettifyJSON(json.RawMessage(jsonBytes)), "Decode JSON to Micheline")
	}

	t.Run("Convert Micheline to JSON", func(t *testing.T) {

		runTest(t, "positive_int")
		runTest(t, "negative_int")
		runTest(t, "string")
		runTest(t, "bytes")
		runTest(t, "sequence")

		// Contracts

		runTest(t, "simple_contract")
		runTest(t, "contract_with_multiply_entrypoints")
		runTest(t, "fa2_contract")
	})
}

func getTestData(fileName string) ([]byte, error) {
	wd, _ := os.Getwd()
	contract_file_path := path.Join(wd, "__test_data__", fileName)
	contract_file, err := os.Open(contract_file_path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(contract_file)
}

// Utility for saving test snapshots
func saveSnapshot(fileName string, bytes []byte) {
	wd, _ := os.Getwd()
	filePath := path.Join(wd, "__test_data__", fileName)
	os.WriteFile(filePath, bytes, 0644)
}

// PrettifyJSON
func PrettifyJSON(o interface{}) string {
	prettyJSON, _ := json.MarshalIndent(o, "", "  ")
	return string(prettyJSON)
}
