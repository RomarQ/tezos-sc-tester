package action

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/romarq/visualtez-testing/internal/utils"
)

type CreateImplicitAccountAction struct {
	raw  json.RawMessage
	json struct {
		Kind    string `json:"kind"`
		Payload struct {
			Name    string `json:"name"`
			Balance string `json:"balance"`
		} `json:"payload"`
	}
	Name    string
	Balance business.Mutez
}

// Unmarshal action
func (action *CreateImplicitAccountAction) Unmarshal() error {
	err := json.Unmarshal(action.raw, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "name" field
	action.Name = action.json.Payload.Name

	// "balance" field
	action.Balance, err = business.MutezOfString(action.json.Payload.Balance)
	if err != nil {
		return err
	}

	return nil
}

// Marshal returns the JSON of the action (cached)
func (a CreateImplicitAccountAction) Marshal() json.RawMessage {
	return a.raw
}

// Perform action (Creates an implicit account)
func (action CreateImplicitAccountAction) Run(mockup business.Mockup) (interface{}, bool) {
	if mockup.ContainsAddress(action.Name) {
		return fmt.Sprintf("Name (%s) is already in use.", action.Name), false
	}

	keyPair, err := utils.GenerateKey()
	if err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return "Could not generate wallet.", false
	}

	// Import private key
	privateKey := keyPair.String()
	if err = mockup.ImportSecret(privateKey, action.Name); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return "Could not import wallet.", false
	}

	// Fund wallet
	address := keyPair.Address().String()
	revealCost := business.MutezOfFloat(big.NewFloat(mockup.Config.Tezos.RevealFee))
	if err = mockup.Transfer(business.CallContractArgument{
		Recipient: address,
		Source:    mockup.Config.Tezos.Originator,
		Amount:    business.AddMutez(action.Balance, revealCost), // Increments revealFee which will be debited when revealing the wallet
	}); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return "Could not fund wallet.", false
	}

	// Reveal wallet
	if err = mockup.RevealWallet(action.Name, revealCost); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return "Could not reveal wallet.", false
	}

	// Confirm that the wallet was funded with the expected amount
	walletBalance := mockup.GetBalance(action.Name)
	// Verify the wallet balance
	if walletBalance.String() != action.Balance.String() {
		err := fmt.Sprintf("Account balance mismatch %s <> %s.", action.Balance, walletBalance.String())
		return err, false
	}

	// Cache contract address
	mockup.CacheAccountAddress(action.Name, address)

	return map[string]interface{}{
		"address": address,
	}, true
}

func (action CreateImplicitAccountAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.Name == "" {
		missingFields = append(missingFields, "name")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Payload.Name); err != nil {
		return err
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", CreateImplicitAccount, strings.Join(missingFields, ", "))
	}
	return nil
}
