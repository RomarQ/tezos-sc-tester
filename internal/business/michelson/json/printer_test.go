package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	t.Run("Test indentation", func(t *testing.T) {
		ast := ast.Sequence{
			Elements: []ast.Node{
				ast.Prim{
					Prim: "storage",
					Arguments: []ast.Node{
						ast.Prim{Prim: "unit"},
					},
				},
				ast.Prim{
					Prim: "parameter",
					Arguments: []ast.Node{
						ast.Prim{
							Prim: "unit",
							Annotations: []ast.Annotation{
								{
									Value: "%do_something",
								},
							},
						},
					},
				},
				ast.Prim{
					Prim: "code",
					Arguments: []ast.Node{
						ast.Sequence{
							Elements: []ast.Node{
								ast.Prim{Prim: "DROP"},
								ast.Prim{Prim: "UNIT"},
								ast.Prim{
									Prim: "NIL",
									Arguments: []ast.Node{
										ast.Prim{Prim: "operation"},
									},
								},
								ast.Prim{Prim: "PAIR"},
							},
						},
					},
				},
			},
		}
		bytes, err := getTestData("print_snapshot.json")
		assert.NoError(t, err, "Must not fail")
		j, err := Print(ast, "", "    ")
		assert.NoError(t, err, "Must not fail")
		assert.Equal(t, string(j), strings.Trim(string(json.RawMessage(bytes)), "\n"), "Validate snapshot")
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
