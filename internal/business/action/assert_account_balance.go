package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/utils"
)

type AssertAccountBalanceAction struct {
	raw  json.RawMessage
	json struct {
		Kind    string `json:"kind"`
		Payload struct {
			AccountName string `json:"account_name"`
			Balance     string `json:"balance"`
		} `json:"payload"`
	}
	AccountName string
	Balance     business.Mutez
}

// Unmarshal action
func (action *AssertAccountBalanceAction) Unmarshal() error {
	err := json.Unmarshal(action.raw, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "account_name" field
	action.AccountName = action.json.Payload.AccountName

	// "balance" field
	action.Balance, err = business.MutezOfString(action.json.Payload.Balance)
	if err != nil {
		return err
	}

	return nil
}

func (action AssertAccountBalanceAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.AccountName == "" {
		missingFields = append(missingFields, "account_name")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Payload.AccountName); err != nil {
		return err
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", AssertAccountBalance, strings.Join(missingFields, ", "))
	}

	return nil
}
