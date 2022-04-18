package business

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/romarq/visualtez-testing/internal/config"
	"github.com/romarq/visualtez-testing/internal/logger"
)

const cmd_tezos_client = "tezos-client"

type (
	TezosClientArgumentKind int8
	MichelsonFormat         string
	TezosClientArgument     struct {
		Kind       TezosClientArgumentKind
		Parameters []string
	}
	CallContractArgument struct {
		Recipient  string
		Source     string
		Entrypoint string
		Parameter  string
		Amount     *TMutez
	}
	Mockup struct {
		TaskID    string
		Config    config.Config
		Addresses map[string]string
	}
)

const (
	COMMAND TezosClientArgumentKind = iota
	Mode
	BaseDirectory
	Protocol
	ProtocolConstants
	BootstrapAccounts
	BurnCap
	Fee
	Init
	Arg
	Entrypoint
	//
	Michelson MichelsonFormat = "michelson"
	JSON      MichelsonFormat = "json"
)

func InitMockup(taskID string, cfg config.Config) Mockup {
	return Mockup{
		TaskID: taskID,
		Config: cfg,
		Addresses: map[string]string{
			"bootstrap1": "tz1",
			"bootstrap2": "tz1",
			"bootstrap3": "tz1",
			"bootstrap4": "tz1",
			"bootstrap5": "tz1",
		},
	}
}

// Bootstrap a mockup environment for the task
func (m Mockup) Bootstrap() error {
	temporaryDirectory := m.getTaskDirectory()
	logger.Debug("[Task #%s] - Creating task directory (%s).", m.TaskID, temporaryDirectory)

	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{temporaryDirectory},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"create", "mockup"},
		},
		TezosClientArgument{
			Kind:       ProtocolConstants,
			Parameters: []string{fmt.Sprintf("%s/protocol-constants.json", m.Config.Tezos.BaseDirectory)},
		},
		TezosClientArgument{
			Kind:       BootstrapAccounts,
			Parameters: []string{fmt.Sprintf("%s/bootstrap-accounts.json", m.Config.Tezos.BaseDirectory)},
		},
	)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	return err
}

// Clear task artifacts
func (m Mockup) Teardown() error {
	temporaryDirectory := m.getTaskDirectory()
	logger.Debug("[Task #%s] - Deleting task directory (%s).", m.TaskID, temporaryDirectory)

	return os.RemoveAll(temporaryDirectory)
}

// Generate Wallet
func (m Mockup) GenerateWallet(walletName string) error {
	logger.Debug("[Task #%s] - Generating wallet (%s).", m.TaskID, walletName)

	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"gen", "keys", walletName},
		},
	)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	return err
}

func (m Mockup) ImportSecret(privateKey string, walletName string) error {
	logger.Debug("[Task #%s] - Importing secret key (%s).", m.TaskID, walletName)

	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"import", "secret", "key", walletName, fmt.Sprintf("unencrypted:%s", privateKey)},
		},
	)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	return err
}

func (m Mockup) Transfer(arg CallContractArgument) error {
	logger.Debug("[Task #%s] - Calling contract %s. %v", m.TaskID, arg.Recipient, arg)

	args := make([]TezosClientArgument, 0)
	args = append(
		args,
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"transfer", TezOfMutez(arg.Amount).Text('f', 6), "from", arg.Source, "to", arg.Recipient},
		},
	)
	if arg.Entrypoint != "" {
		args = append(
			args,
			TezosClientArgument{
				Kind:       Entrypoint,
				Parameters: []string{arg.Entrypoint},
			},
		)
	}
	if arg.Parameter != "" {
		args = append(
			args,
			TezosClientArgument{
				Kind:       Arg,
				Parameters: []string{arg.Parameter},
			},
		)
	}

	args = append(
		args,
		TezosClientArgument{
			Kind:       BurnCap,
			Parameters: []string{"1"},
		},
	)
	arguments := composeArguments(args...)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	if err != nil {
		return err
	}

	return nil
}

// Reveal wallet
func (m Mockup) RevealWallet(walletName string, revealFee *TMutez) error {
	logger.Debug("[Task #%s] - Revealing wallet (%s).", m.TaskID, walletName)

	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"reveal", "key", "for", walletName},
		},
		TezosClientArgument{
			Kind:       Fee,
			Parameters: []string{TezOfMutez(revealFee).Text('f', 6)},
		},
	)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	return err
}

func (m Mockup) GetBalance(name string) (*TMutez, error) {
	logger.Debug("[Task #%s] - Get balance of (%s).", m.TaskID, name)

	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"get", "balance", "for", name},
		},
	)

	// Execute command
	output, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	if err != nil {
		return nil, err
	}

	// Extract balance in ꜩ
	pattern := regexp.MustCompile(`(\d*.?\d*)\sꜩ`)
	match := pattern.FindStringSubmatch(output)
	if len(match) < 2 {
		return nil, fmt.Errorf("Could not get the balance for account %s.", name)
	}

	balance, ok := new(TTez).SetString(match[1])
	if ok {
		return MutezOfTez(balance), nil
	}

	return nil, fmt.Errorf("Could not get contract balance.")
}

func (m *Mockup) Originate(sender string, contractName string, balance *TMutez, code string, storage string) (string, error) {
	logger.Debug("[Task #%s] - Originating contract (%s).", m.TaskID, contractName)

	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind: COMMAND,
			Parameters: []string{
				"originate", "contract", contractName,
				"transferring", TezOfMutez(balance).Text('f', 6), "from", sender,
				"running", code,
			},
		},
		TezosClientArgument{
			Kind:       Init,
			Parameters: []string{storage},
		},
		TezosClientArgument{
			Kind:       BurnCap,
			Parameters: []string{"1"},
		},
	)

	output, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	if err != nil {
		return "", err
	}

	// Extract balance in ꜩ
	pattern := regexp.MustCompile(`New\scontract\s(\w+)\soriginated`)
	match := pattern.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", fmt.Errorf("Could not extract the contract address from origination output.")
	}

	return match[1], nil
}

func (m Mockup) GetContractStorage(contractName string) (string, error) {
	logger.Debug("[Task #%s] - Get storage from contract (%s).", m.TaskID, contractName)

	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind: COMMAND,
			Parameters: []string{
				"get", "contract", "storage", "for", contractName,
			},
		},
	)

	output, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Convert script format between "michelson" and "json"
func (m Mockup) ConvertScript(script string, from MichelsonFormat, to MichelsonFormat) (string, error) {
	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind: COMMAND,
			Parameters: []string{
				"convert", "script", script, "from", string(from), "to", string(to),
			},
		},
	)

	return m.runTezosClient(m.getTezosClientPath(), arguments)
}

// Checks if address exists
func (m Mockup) ContainsAddress(name string) bool {
	return m.Addresses[name] != ""
}

// Set address
func (m Mockup) SetAddress(name string, address string) {
	m.Addresses[name] = address
}

// Convert data format between "michelson" and "json"
func (m Mockup) ConvertData(data string, from MichelsonFormat, to MichelsonFormat) (string, error) {
	arguments := composeArguments(
		TezosClientArgument{
			Kind:       Mode,
			Parameters: []string{"mockup"},
		},
		TezosClientArgument{
			Kind:       BaseDirectory,
			Parameters: []string{m.getTaskDirectory()},
		},
		TezosClientArgument{
			Kind:       Protocol,
			Parameters: []string{m.Config.Tezos.DefaultProtocol},
		},
		TezosClientArgument{
			Kind: COMMAND,
			Parameters: []string{
				"convert", "data", data, "from", string(from), "to", string(to),
			},
		},
	)

	return m.runTezosClient(m.getTezosClientPath(), arguments)
}

// Execute a "tezos-client" command
func (m Mockup) runTezosClient(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)

	var errBuffer bytes.Buffer
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		if errBuffer.Len() > 0 {
			msg := errBuffer.String()
			logger.Error("Got the following error:\n\n%s\nwhen executing command: %s.", msg, cmd.Args)
			err = fmt.Errorf(msg)
		}
		return "", err
	}

	output := outBuffer.String()
	if len(output) > 0 {
		logger.Debug("Got the following output:\n\n%s\nwhen executing command: %s.", string(output), cmd.Args)
	}

	return output, nil
}

func (m Mockup) getTaskDirectory() string {
	return fmt.Sprintf("%s/_tmp/%s", m.Config.Tezos.BaseDirectory, m.TaskID)
}

func (m Mockup) getTezosClientPath() string {
	return fmt.Sprintf("%s/%s", m.Config.Tezos.BaseDirectory, cmd_tezos_client)
}

func composeArguments(args ...TezosClientArgument) []string {
	arguments := make([]string, 0)
	for _, argument := range args {
		switch argument.Kind {
		case Mode:
			arguments = append(arguments, "-M")
		case Protocol:
			arguments = append(arguments, "-p")
		case BaseDirectory:
			arguments = append(arguments, "-d")
		case ProtocolConstants:
			arguments = append(arguments, "--protocol-constants")
		case BootstrapAccounts:
			arguments = append(arguments, "--bootstrap-accounts")
		case BurnCap:
			arguments = append(arguments, "--burn-cap")
		case Fee:
			arguments = append(arguments, "--fee")
		case Init:
			arguments = append(arguments, "--init")
		case Arg:
			arguments = append(arguments, "--arg")
		case Entrypoint:
			arguments = append(arguments, "--entrypoint")
		}
		arguments = append(arguments, argument.Parameters...)
	}
	return arguments
}
