package action

import (
	"encoding/json"
	"fmt"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/logger"
)

type CreateImplicitAccountAction struct {
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

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
	if err = mockup.Transfer(action.Balance, "bootstrap1", keyPair.Address().String()); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return fmt.Errorf("Could not fund wallet.")
	}

	// Reveal wallet
	if err = mockup.RevealWallet(action.Name); err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return fmt.Errorf("Could not reveal wallet.")
	}

	return nil
}

func (action *CreateImplicitAccountAction) validate() error {
	if action.Name == "" {
		return fmt.Errorf("Actions of kind (%s) must contain a name field.", CreateImplicitAccount)
	}
	return nil
}
