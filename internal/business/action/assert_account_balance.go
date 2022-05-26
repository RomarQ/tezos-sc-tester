package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/tezos-sc-tester/internal/business"
	"github.com/romarq/tezos-sc-tester/internal/utils"
)

type AssertAccountBalanceAction struct {
	json struct {
		Kind    ActionKind `json:"kind"`
		Payload struct {
			AccountName string `json:"account_name"`
			Balance     string `json:"balance"`
		} `json:"payload"`
	}
	AccountName string
	Balance     business.Mutez
}

// Unmarshal action
func (action *AssertAccountBalanceAction) Unmarshal(ac Action) error {
	action.json.Kind = ac.Kind
	err := json.Unmarshal(ac.Payload, &action.json.Payload)
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

// Marshal returns the JSON of the action (cached)
func (action AssertAccountBalanceAction) Action() interface{} {
	return action.json
}

// Perform the action
func (action AssertAccountBalanceAction) Run(mockup business.Mockup) (interface{}, bool) {
	balance := mockup.GetBalance(action.AccountName)

	if balance.String() != action.Balance.String() {
		return map[string]string{
			"expected": action.Balance.String(),
			"actual":   balance.String(),
		}, false
	}

	return map[string]string{
		"balance": balance.String(),
	}, true
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
