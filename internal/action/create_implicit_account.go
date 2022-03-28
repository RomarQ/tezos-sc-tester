package action

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type CreateImplicitAccountAction struct {
	Name string `json:"name"`
}

func (action *CreateImplicitAccountAction) validate() error {
	if action.Name == "" {
		return fmt.Errorf("Actions of kind (%s) must contain a name field.", CreateImplicitAccount)
	}
	return nil
}

// Unmarshal action
func (action *CreateImplicitAccountAction) Unmarshal(bytes json.RawMessage) error {
	if err := json.Unmarshal([]byte(gjson.GetBytes(bytes, `action`).Raw), &action); err != nil {
		return err
	}

	// Validate action
	return action.validate()
}
