package business

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	Config "github.com/romarq/visualtez-testing/internal/config"
	Logger "github.com/romarq/visualtez-testing/internal/logger"
)

const cmd_tezos_client = "tezos-client"

type TezosClientArgumentKind int8

const (
	COMMAND TezosClientArgumentKind = iota
	Mode
	BaseDirectory
	Protocol
	ProtocolConstants
	BootstrapAccounts
	BurnCap
)

type TezosClientArgument struct {
	Kind       TezosClientArgumentKind
	Parameters []string
}

type Mockup struct {
	TaskID string
	Config Config.Config
}

func InitMockup(taskID string, config Config.Config) Mockup {
	return Mockup{
		TaskID: taskID,
		Config: config,
	}
}

// Bootstrap a mockup environment for the task
func (m *Mockup) Bootstrap() error {
	temporaryDirectory := m.getTaskDirectory()
	Logger.Debug("[Task #%s] - Creating task directory (%s).", m.TaskID, temporaryDirectory)

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
func (m *Mockup) Teardown() error {
	temporaryDirectory := m.getTaskDirectory()
	Logger.Debug("[Task #%s] - Deleting task directory (%s).", m.TaskID, temporaryDirectory)

	return os.RemoveAll(temporaryDirectory)
}

// Generate Wallet
func (m *Mockup) GenerateWallet(walletName string) error {
	Logger.Debug("[Task #%s] - Generating wallet (%s).", m.TaskID, walletName)

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

func (m *Mockup) ImportSecret(privateKey string, walletName string) error {
	Logger.Debug("[Task #%s] - Importing secret key (%s).", m.TaskID, walletName)

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

func (m *Mockup) Transfer(amount int64, source string, recipient string) error {
	Logger.Debug("[Task #%s] - Transfering %dꜩ from %s to %s.", m.TaskID, amount, source, recipient)

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
			Parameters: []string{"transfer", fmt.Sprint(amount), "from", source, "to", recipient},
		},
		TezosClientArgument{
			Kind:       BurnCap,
			Parameters: []string{"0.1"},
		},
	)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	return err
}

func (m *Mockup) RevealWallet(walletName string) error {
	Logger.Debug("[Task #%s] - Revealing wallet (%s).", m.TaskID, walletName)

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
	)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	return err
}

// Execute a "tezos-client" command
func (m *Mockup) runTezosClient(command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)

	var errBuffer bytes.Buffer
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		if errBuffer.Len() > 0 {
			msg := string(errBuffer.Bytes()[:])
			Logger.Error("Got the following error:\n\n%s\nwhen executing command: %s.", msg, cmd.Args)
		}
		return nil, err
	}

	output := outBuffer.Bytes()
	if len(output) > 0 {
		Logger.Debug("Got the following output:\n\n%s\nwhen executing command: %s.", string(output[:]), cmd.Args)
	}

	return output, nil
}

func (m *Mockup) getTaskDirectory() string {
	return fmt.Sprintf("%s/_tmp/%s", m.Config.Tezos.BaseDirectory, m.TaskID)
}

func (m *Mockup) getTezosClientPath() string {
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
		}
		arguments = append(arguments, argument.Parameters...)
	}
	return arguments
}
