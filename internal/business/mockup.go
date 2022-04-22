package business

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"regexp"

	"github.com/romarq/visualtez-testing/internal/business/michelson"
	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	"github.com/romarq/visualtez-testing/internal/config"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/tidwall/sjson"
)

const cmd_tezos_client = "tezos-client"

type (
	TezosClientArgumentKind int8
	ParsingMode             string
	MichelsonFormat         string
	TezosClientArgument     struct {
		Kind       TezosClientArgumentKind
		Parameters []string
	}
	CallContractArgument struct {
		Recipient  string
		Source     string
		Entrypoint string
		Amount     Mutez
		Parameter  string
	}
	ContractCache struct {
		StorageType ast.Node
	}
	Mockup struct {
		TaskID    string
		Protocol  string
		Config    config.Config
		Addresses map[string]string
		contracts map[string]ContractCache
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
	UnparsingMode
	// Parsing modes
	Readable  ParsingMode = "Readable"
	Optimized ParsingMode = "Optimized"
	// Michelson Formats
	Michelson MichelsonFormat = "michelson"
	JSON      MichelsonFormat = "json"
)

func InitMockup(taskID string, protocol string, cfg config.Config) Mockup {
	return Mockup{
		TaskID:    taskID,
		Protocol:  protocol,
		Config:    cfg,
		contracts: map[string]ContractCache{},
	}
}

// Bootstrap bootstraps a mockup environment for the task
func (m *Mockup) Bootstrap() error {
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
			Parameters: []string{m.getProtocol()},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"create", "mockup"},
		},
		TezosClientArgument{
			Kind:       BootstrapAccounts,
			Parameters: []string{fmt.Sprintf("%s/bootstrap-accounts.json", m.Config.Tezos.BaseDirectory)},
		},
		TezosClientArgument{
			Kind:       ProtocolConstants,
			Parameters: []string{fmt.Sprintf("%s/protocol-constants.json", m.Config.Tezos.BaseDirectory)},
		},
	)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	if err != nil {
		return fmt.Errorf("could not bootstrap mockup. %s", err)
	}

	// Populate the address cache map
	m.Addresses = m.fetchKnownAddresses()

	return nil
}

// Teardown clears task artifacts
func (m Mockup) Teardown() error {
	temporaryDirectory := m.getTaskDirectory()
	logger.Debug("[Task #%s] - Deleting task directory (%s).", m.TaskID, temporaryDirectory)

	return os.RemoveAll(temporaryDirectory)
}

// UpdateChainID updates the chain identifier in the mockup context
func (m Mockup) UpdateChainID(chainID string) error {
	logger.Debug("[Task #%s] - Updating chain_id to (%s).", m.TaskID, chainID)
	contextPath := fmt.Sprintf("%s/mockup/context.json", m.getTaskDirectory())

	errorMsg := fmt.Errorf("could not modify chain_id.")

	bytes, err := os.ReadFile(contextPath)
	if err != nil {
		logger.Debug("could not open %s: %s", contextPath, err)
		return errorMsg
	}
	bytes, err = sjson.SetBytes(bytes, "chain_id", chainID)
	if err != nil {
		logger.Debug(`could not modify "chain_id" field. %s`, err)
		return errorMsg
	}

	err = os.WriteFile(contextPath, bytes, 644)
	if err != nil {
		logger.Debug("could not write to %s: %s", contextPath, err)
		return errorMsg
	}

	return nil
}

// UpdateHeadBlockLevel updates the level of the head block in the mockup context
func (m Mockup) UpdateHeadBlockLevel(level int32) error {
	logger.Debug("[Task #%s] - Updating block level to (%s).", m.TaskID, level)
	contextPath := fmt.Sprintf("%s/mockup/context.json", m.getTaskDirectory())

	errorMsg := fmt.Errorf("could not modify block level.")

	bytes, err := os.ReadFile(contextPath)
	if err != nil {
		logger.Debug("could not open %s: %s", contextPath, err)
		return errorMsg
	}
	bytes, err = sjson.SetBytes(bytes, "context.shell_header.level", level)
	if err != nil {
		logger.Debug(`could not modify "context.shell_header.level" field. %s`, err)
		return errorMsg
	}

	err = os.WriteFile(contextPath, bytes, 644)
	if err != nil {
		logger.Debug("could not write to %s: %s", contextPath, err)
		return errorMsg
	}

	return nil
}

// GenerateWallet generates a new wallet (uses ed25519 curve)
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
			Parameters: []string{m.getProtocol()},
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
			Parameters: []string{m.getProtocol()},
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
			Parameters: []string{m.getProtocol()},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"transfer", arg.Amount.ToTez().String(), "from", arg.Source, "to", arg.Recipient},
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
func (m Mockup) RevealWallet(walletName string, revealFee Mutez) error {
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
			Parameters: []string{m.getProtocol()},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"reveal", "key", "for", walletName},
		},
		TezosClientArgument{
			Kind:       Fee,
			Parameters: []string{revealFee.ToTez().String()},
		},
	)

	_, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	return err
}

func (m Mockup) GetBalance(name string) Mutez {
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
			Parameters: []string{m.getProtocol()},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"get", "balance", "for", name},
		},
	)

	// Execute command
	output, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	if err != nil {
		return MutezOfFloat(big.NewFloat(0))
	}

	// Extract balance in ꜩ
	pattern := regexp.MustCompile(`(\d*.?\d*)\sꜩ`)
	match := pattern.FindStringSubmatch(output)
	if len(match) < 2 {
		return MutezOfFloat(big.NewFloat(0))
	}

	balance, err := TezOfString(match[1])
	if err != nil {
		return MutezOfFloat(big.NewFloat(0))
	}

	return balance.ToMutez()
}

func (m *Mockup) Originate(sender string, contractName string, amount Mutez, code string, storage string) (string, error) {
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
			Parameters: []string{m.getProtocol()},
		},
		TezosClientArgument{
			Kind: COMMAND,
			Parameters: []string{
				"originate", "contract", contractName,
				"transferring", amount.ToTez().String(), "from", sender,
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
		logger.Debug("could originate contract. %s", err)
		return "", err
	}

	// Extract contract address
	pattern := regexp.MustCompile(`New\scontract\s(\w+)\soriginated`)
	match := pattern.FindStringSubmatch(output)
	if len(match) < 2 || len(match[1]) < 36 || match[1][0:3] != "KT1" {
		return "", fmt.Errorf("Could not extract the contract address from origination output.")
	}

	return match[1], nil
}

func (m Mockup) GetContractStorage(contractName string) (ast.Node, error) {
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
			Parameters: []string{m.getProtocol()},
		},
		TezosClientArgument{
			Kind: COMMAND,
			Parameters: []string{
				"get", "contract", "storage", "for", contractName,
			},
		},
		TezosClientArgument{
			Kind:       UnparsingMode,
			Parameters: []string{string(Readable)},
		},
	)

	output, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	if err != nil {
		return nil, fmt.Errorf("could not fetch storage from contract (%s). %s", contractName, err)
	}

	ast, err := michelson.ParseMicheline(output)
	if err != nil {
		return nil, fmt.Errorf("could parse contract (%s) storage from 'micheline' format. %s", contractName, err)
	}

	return ast, nil
}

// Checks if address exists
func (m Mockup) ContainsAddress(name string) bool {
	return m.Addresses[name] != ""
}

// CacheAddress caches the address of a contract by name
func (m Mockup) CacheAccountAddress(name string, address string) {
	m.Addresses[name] = address
}

// CacheContract caches contract information
func (m Mockup) CacheContract(name string, code ast.Node) error {
	switch seq := code.(type) {
	case ast.Sequence:
		for _, node := range seq.Elements {
			switch prim := node.(type) {
			case ast.Prim:
				if prim.Prim == "storage" {
					m.contracts[name] = ContractCache{
						StorageType: prim.Arguments[0],
					}
				}
				return nil
			}
		}
	}
	return fmt.Errorf("could not cache contract. michelson is invalid.")
}

// GetCachedContract
func (m Mockup) GetCachedContract(name string) ContractCache {
	return m.contracts[name]
}

// ConvertScript converts script format between "michelson" and "json"
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
			Parameters: []string{m.getProtocol()},
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

// ConvertData converts data format between "michelson" and "json"
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
			Parameters: []string{m.getProtocol()},
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

// NormalizeData normalize a data expression against a gicen type
func (m Mockup) NormalizeData(data string, dataType string, mode ParsingMode) (ast.Node, error) {
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
			Parameters: []string{m.getProtocol()},
		},
		TezosClientArgument{
			Kind: COMMAND,
			Parameters: []string{
				"normalize", "data", data, "of", "type", dataType,
			},
		},
		TezosClientArgument{
			Kind:       UnparsingMode,
			Parameters: []string{string(mode)},
		},
	)

	output, err := m.runTezosClient(m.getTezosClientPath(), arguments)
	if err != nil {
		return nil, fmt.Errorf("could not normalize data %s against type %s. %s", data, dataType, err)
	}

	ast, err := michelson.ParseMicheline(output)
	if err != nil {
		return nil, fmt.Errorf("could parse normalized data %s. %s", output, err)
	}

	return ast, nil
}

// runTezosClient executes a "tezos-client" command
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
			return "", err
		}
		return "", err
	}

	output := outBuffer.String()
	if len(output) > 0 {
		logger.Debug("Got the following output:\n\n%s\nwhen executing command: %s.", string(output), cmd.Args)
	}

	return output, nil
}

// fetchKnownAddresses gets all accounts known by 'tezos-client"
func (m Mockup) fetchKnownAddresses() map[string]string {
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
			Parameters: []string{m.getProtocol()},
		},
		TezosClientArgument{
			Kind:       COMMAND,
			Parameters: []string{"list", "known", "addresses"},
		},
	)

	output, _ := m.runTezosClient(m.getTezosClientPath(), arguments)

	// Extract addresses
	pattern := regexp.MustCompile(`(\w+):\s(\w+)\s`)
	match := pattern.FindAllStringSubmatch(output, -1)

	addresses := map[string]string{}
	for _, m := range match {
		// m[1] = name
		// m[2] = address
		addresses[m[1]] = m[2]
	}

	return addresses
}

// getTaskDirectory gives the path to the temporary folder that
// will be used to store all artifacts produced by the mockup
func (m Mockup) getTaskDirectory() string {
	return fmt.Sprintf("%s/_tmp/%s", m.Config.Tezos.BaseDirectory, m.TaskID)
}

// getTezosClientPath gives the path to the 'tezos-client' binary
func (m Mockup) getTezosClientPath() string {
	return fmt.Sprintf("%s/%s", m.Config.Tezos.BaseDirectory, cmd_tezos_client)
}

// getProtocol gives the protocol being used in the mockup
func (m Mockup) getProtocol() string {
	if m.Protocol == "" {
		return m.Config.Tezos.DefaultProtocol
	}
	return m.Protocol
}

// composeArguments prepares the arguments for using with 'tezos-client'
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
		case UnparsingMode:
			arguments = append(arguments, "--unparsing-mode")
		}
		arguments = append(arguments, argument.Parameters...)
	}
	return arguments
}
