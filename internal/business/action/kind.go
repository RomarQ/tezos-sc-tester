package action

type ActionKind string

const (
	CreateImplicitAccount ActionKind = "create_implicit_account"
	OriginateContract     ActionKind = "originate_contract"
	CallContract          ActionKind = "call_contract"
	AssertAccountBalance  ActionKind = "assert_account_balance"
	AssertContractStorage ActionKind = "assert_contract_storage"
)
