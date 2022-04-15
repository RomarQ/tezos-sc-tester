package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/romarq/visualtez-testing/internal/utils"
)

type CreateImplicitAccountAction struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

const (
	revealFee = 0.5
)

// Unmarshal action
func (action *CreateImplicitAccountAction) Unmarshal(bytes json.RawMessage) error {
	if err := json.Unmarshal(bytes, &action); err != nil {
		return err
	}

	// Validate action
	return action.validate()
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
	balance := action.Balance + revealFee // Increments revealFee which will be debited when revealing the wallet
	if err = mockup.Transfer(business.CallContractArgument{
		Recipient: address,
		Source:    "bootstrap1",
		Amount:    balance,
	}); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return action.buildFailureResult("Could not fund wallet.")
	}

	// Reveal wallet
	if err = mockup.RevealWallet(action.Name, revealFee); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return action.buildFailureResult("Could not reveal wallet.")
	}

	// Confirm that the wallet was funded with the expected amount
	walletBalance, err := mockup.GetBalance(action.Name)
	if err != nil {
		return action.buildFailureResult("Failed to confirm balance.")
	}
	if walletBalance != action.Balance {
		err := fmt.Sprintf("Account balance mismatch %f <> %f.", action.Balance, walletBalance)
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
	if action.Name == "" {
		missingFields = append(missingFields, "name")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.Name); err != nil {
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
		Action: action,
	}
}

func (action CreateImplicitAccountAction) buildFailureResult(details string) ActionResult {
	return ActionResult{
		Status: Failure,
		Kind:   CreateImplicitAccount,
		Action: action,
		Result: map[string]interface{}{
			"details": details,
		},
	}
}
