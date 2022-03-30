package action

import (
	"encoding/json"
	"fmt"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/logger"
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

// Perform the action
func (action *CreateImplicitAccountAction) Run(mockup business.Mockup) error {
	keyPair, err := business.GenerateKey()
	if err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return fmt.Errorf("Could not generate wallet.")
	}

	// Import private key
	privateKey := keyPair.String()
	if err = mockup.ImportSecret(privateKey, action.Name); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return fmt.Errorf("Could not import wallet.")
	}

	// Fund wallet
	balance := action.Balance + revealFee // Increments revealFee which will be debited when revealing the wallet
	if err = mockup.Transfer(balance, "bootstrap1", keyPair.Address().String()); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return fmt.Errorf("Could not fund wallet.")
	}

	// Reveal wallet
	if err = mockup.RevealWallet(action.Name, revealFee); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return fmt.Errorf("Could not reveal wallet.")
	}

	// Confirm that the wallet was funded with the expected amount
	walletBalance, err := mockup.GetBalance(action.Name)
	if err != nil {
		return fmt.Errorf("Failed to confirm balance.")
	}
	if walletBalance != action.Balance {
		return fmt.Errorf("Account balance mismatch %f <> %f.", action.Balance, walletBalance)
	}

	return nil
}

func (action *CreateImplicitAccountAction) validate() error {
	if action.Name == "" {
		return fmt.Errorf("Actions of kind (%s) must contain a name field.", CreateImplicitAccount)
	}
	return nil
}
