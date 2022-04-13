package michelson

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMichelineOfJSON(t *testing.T) {
	t.Run("Convert JSON to Micheline (Positive Int)", func(t *testing.T) {
		micheline, err := MichelineOfJSON(json.RawMessage(`{ "int": "91" }`))
		assert.Nil(t, err, "Must not fail")
		assert.Equal(t, micheline, "91", "Verify micheline")
	})
	t.Run("Convert JSON to Micheline (String)", func(t *testing.T) {
		micheline, err := MichelineOfJSON(json.RawMessage(`{ "string": "Hello World" }`))
		assert.Nil(t, err, "Must not fail")
		assert.Equal(t, micheline, "\"Hello World\"", "Verify micheline")
	})
	t.Run("Convert JSON to Micheline (Bool)", func(t *testing.T) {
		micheline, err := MichelineOfJSON(json.RawMessage(`{ "bool": "True" }`))
		assert.Nil(t, err, "Must not fail")
		assert.Equal(t, micheline, "True", "Verify micheline")
	})
	t.Run("Convert JSON to Micheline (Bytes)", func(t *testing.T) {
		micheline, err := MichelineOfJSON(json.RawMessage(`{ "bytes": "0x01" }`))
		assert.Nil(t, err, "Must not fail")
		assert.Equal(t, micheline, "0x01", "Verify micheline")
	})
	t.Run("Convert JSON to Micheline (Sequence)", func(t *testing.T) {
		micheline, err := MichelineOfJSON(json.RawMessage(`[{ "int": "1" }, { "int": "2" }]`))
		assert.Nil(t, err, "Must not fail")
		assert.Equal(t, micheline, "{ 1 ; 2 }", "Verify micheline")
	})
	t.Run("Convert JSON to Micheline (Contract)", func(t *testing.T) {
		bytes, err := getTestData("simple_contract.json")
		assert.NoError(t, err)
		micheline, err := MichelineOfJSON(json.RawMessage(bytes))
		assert.Nil(t, err, "Must not fail")
		assert.Equal(t, micheline, "{ storage unit ; parameter (unit %do_something) ; code { DROP ; UNIT ; NIL operation ; PAIR } }", "Verify micheline")
	})
	t.Run("Convert JSON to Micheline (Contract with multiple entrypoints)",
		func(t *testing.T) {
			jsonBytes, err := getTestData("contract_with_multiply_entrypoints.json")
			assert.NoError(t, err)
			michelineBytes, err := getTestData("contract_with_multiply_entrypoints.tz")
			assert.NoError(t, err)
			micheline, err := MichelineOfJSON(json.RawMessage(jsonBytes))
			assert.Nil(t, err, "Must not fail")
			assert.Equal(t, []byte(micheline), michelineBytes, "Decode JSON to Micheline")
		})
	t.Run("Convert JSON to Micheline (FA2 Contract)",
		func(t *testing.T) {
			jsonBytes, err := getTestData("fa2_contract.json")
			assert.NoError(t, err)
			michelineBytes, err := getTestData("fa2_contract.tz")
			assert.NoError(t, err)
			micheline, err := MichelineOfJSON(json.RawMessage(jsonBytes))
			assert.Nil(t, err, "Must not fail")
			assert.Equal(t, []byte(micheline), michelineBytes, "Decode JSON to Micheline")
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
