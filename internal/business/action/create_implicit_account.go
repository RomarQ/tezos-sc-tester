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
	json struct {
		Name    string `json:"name"`
		Balance string `json:"balance"`
	}
	Name    string
	Balance business.Mutez
}

// Unmarshal action
func (action *CreateImplicitAccountAction) Unmarshal(bytes json.RawMessage) error {
	err := json.Unmarshal(bytes, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "name" field
	action.Name = action.json.Name

	// "balance" field
	action.Balance, err = business.MutezOfString(action.json.Balance)
	if err != nil {
		return err
	}

	return nil
}

// Perform action (Creates an implicit account)
func (action CreateImplicitAccountAction) Run(mockup business.Mockup) ActionResult {
	if mockup.ContainsAddress(action.Name) {
		return action.buildFailureResult(fmt.Sprintf("Name (%s) is already in use.", action.Name))
	}

	keyPair, err := utils.GenerateKey()
	if err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return action.buildFailureResult("Could not generate wallet.")
	}

	// Import private key
	privateKey := keyPair.String()
	if err = mockup.ImportSecret(privateKey, action.Name); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return action.buildFailureResult("Could not import wallet.")
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
		return action.buildFailureResult("Could not fund wallet.")
	}

	// Reveal wallet
	if err = mockup.RevealWallet(action.Name, revealCost); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return action.buildFailureResult("Could not reveal wallet.")
	}

	// Confirm that the wallet was funded with the expected amount
	walletBalance, err := mockup.GetBalance(action.Name)
	if err != nil {
		return action.buildFailureResult("Failed to confirm balance.")
	}
	// Verify the wallet balance
	if walletBalance.String() != action.Balance.String() {
		err := fmt.Sprintf("Account balance mismatch %s <> %s.", action.Balance, walletBalance.String())
		return action.buildFailureResult(err)
	}

	// Save new address
	mockup.SetAddress(action.Name, address)

	return action.buildSuccessResult(map[string]interface{}{
		"address": address,
	})
}

func (action CreateImplicitAccountAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Name == "" {
		missingFields = append(missingFields, "name")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Name); err != nil {
		return err
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", CreateImplicitAccount, strings.Join(missingFields, ", "))
	}
	return nil
}

func (action CreateImplicitAccountAction) buildSuccessResult(result map[string]interface{}) ActionResult {
	return ActionResult{
		Status: Success,
		Kind:   CreateImplicitAccount,
		Result: result,
		Action: action.json,
	}
}

func (action CreateImplicitAccountAction) buildFailureResult(details string) ActionResult {
	return ActionResult{
		Status: Failure,
		Kind:   CreateImplicitAccount,
		Action: action.json,
		Result: map[string]interface{}{
			"details": details,
		},
	}
}
