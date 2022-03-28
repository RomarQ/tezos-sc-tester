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
)

type TezosClientArgument struct {
	Kind       TezosClientArgumentKind
	Parameters []string
}

type Mockup struct {
	Config Config.Config
}

func InitMockup(config Config.Config) Mockup {
	return Mockup{Config: config}
}

func (m *Mockup) Bootstrap(taskID string) error {
	temporaryDirectory := fmt.Sprintf("%s/_tmp/%s", m.Config.Tezos.BaseDirectory, taskID)
	Logger.Debug("Creating mockup directory: %s.", temporaryDirectory)

	command := fmt.Sprintf("%s/%s", m.Config.Tezos.BaseDirectory, cmd_tezos_client)
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

	_, err := m.runTezosClient(command, arguments)
	return err
}

func (m *Mockup) Teardown(taskID string) error {
	temporaryDirectory := fmt.Sprintf("%s/_tmp/%s", m.Config.Tezos.BaseDirectory, taskID)
	Logger.Debug("Deleting mockup directory: %s.", temporaryDirectory)
	return os.RemoveAll(temporaryDirectory)
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
		Logger.Debug("Got the following output:\n\n%s\nwhen executing command: %s.", output, cmd.Args)
	}

	return output, nil
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
		}
		arguments = append(arguments, argument.Parameters...)
	}
	return arguments
}
