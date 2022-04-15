package json

type (
	MichelsonJSON struct {
		Prim   string        `json:"prim,omitempty"`
		Int    string        `json:"int,omitempty"`
		String string        `json:"string,omitempty"`
		Bytes  string        `json:"bytes,omitempty"`
		Args   []interface{} `json:"args,omitempty"`
		Annots []string      `json:"annots,omitempty"`
	}
)

func (json MichelsonJSON) isInt() bool {
	return json.Int != ""
}
func (json MichelsonJSON) isString() bool {
	return json.String != ""
}
func (json MichelsonJSON) isBytes() bool {
	return json.Bytes != ""
}
func (json MichelsonJSON) isPrim() bool {
	return json.Prim != ""
}
