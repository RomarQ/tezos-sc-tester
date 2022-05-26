package micheline

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/romarq/tezos-sc-tester/internal/business/michelson/ast"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	codeAST := ast.Sequence{
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

	t.Run("With indentation", func(t *testing.T) {
		bytes, err := getTestData("print_with_indent.tz")
		assert.NoError(t, err, "Must not fail")
		micheline := Print(codeAST, "    ")
		assert.Equal(t, micheline, strings.Trim(string(bytes), "\n"), "Validate snapshot")
	})

	t.Run("Without indentation", func(t *testing.T) {
		bytes, err := getTestData("print_without_indent.tz")
		assert.NoError(t, err, "Must not fail")
		micheline := Print(codeAST, "")
		assert.Equal(t, micheline, strings.Trim(string(bytes), "\n"), "Validate snapshot")
	})

	t.Run("Print sequences", func(t *testing.T) {
		seq := ast.Sequence{
			Elements: []ast.Node{
				ast.Prim{
					Prim: "Pair",
					Arguments: []ast.Node{
						ast.Int{
							Value: "1",
						},
						ast.Int{
							Value: "2",
						},
					},
				},
				ast.Prim{
					Prim: "Pair",
					Arguments: []ast.Node{
						ast.Int{
							Value: "3",
						},
						ast.Int{
							Value: "4",
						},
					},
				},
			},
		}
		micheline := Print(seq, "")
		assert.Equal(t, micheline, "{ Pair 1 2; Pair 3 4 }", "Validate snapshot")
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
