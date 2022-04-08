package action

type ActionKind string

const (
	CreateImplicitAccount ActionKind = "create_implicit_account"
	OriginateContract                = "originate_contract"
	CallContract                     = "call_contract"
)
