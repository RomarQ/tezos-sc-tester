// Action kinds
export enum ActionKind {
    CreateImplicitAccount = 'create_implicit_account',
    OriginateContract = 'originate_contract',
    CallContract = 'call_contract',
    AssertAccountBalance = 'assert_account_balance',
    AssertContractStorage = 'assert_contract_storage',
    ModifyChainID = 'modify_chain_id',
    ModifyBlockLevel = 'modify_block_level',
    ModifyBlockTimestamp = 'modify_block_timestamp',
    PackData = 'pack_data',
}

// Action result status
export enum ActionResultStatus {
    Success = 'success',
    Failure = 'failure',
}

export type IAction =
    | ICreateImplicitAccountAction
    | IOriginateContractAction
    | ICallContractAction
    | IAssertAccountBalanceAction
    | IAssertContractStorageAction
    | IModifyChainIDAction
    | IModifyBlockLevelAction
    | IModifyBlockTimestampAction
    | IPackDataAction;

export interface IActionResult {
    status: ActionResultStatus;
    action: IAction;
    result: Record<string, unknown>;
}

// create_implicit_account

export interface ICreateImplicitAccountPayload {
    name: string;
    balance: string;
}
export interface ICreateImplicitAccountAction {
    kind: ActionKind.CreateImplicitAccount;
    payload: ICreateImplicitAccountPayload;
}

// originate_contract

export interface IOriginateContractPayload {
    name: string;
    balance: string;
    code: Record<string, unknown> | Record<string, unknown>[];
    storage: Record<string, unknown> | Record<string, unknown>[];
}
export interface IOriginateContractAction {
    kind: ActionKind.OriginateContract;
    payload: IOriginateContractPayload;
}

// call_contract

export interface ICallContractPayload {
    recipient: string;
    sender: string;
    amount: string;
    level?: number;
    timestamp?: string;
    entrypoint: string;
    parameter: Record<string, unknown> | Record<string, unknown>[];
    expect_failwith?: Record<string, unknown> | Record<string, unknown>[];
}
export interface ICallContractAction {
    kind: ActionKind.CallContract;
    payload: ICallContractPayload;
}

// assert_account_balance

export interface IAssertAccountBalancePayload {
    account_name: string;
    balance: string;
}
export interface IAssertAccountBalanceAction {
    kind: ActionKind.AssertAccountBalance;
    payload: IAssertAccountBalancePayload;
}

// assert_contract_storage

export interface IAssertContractStoragePayload {
    contract_name: string;
    storage: Record<string, unknown> | Record<string, unknown>[];
}
export interface IAssertContractStorageAction {
    kind: ActionKind.AssertContractStorage;
    payload: IAssertContractStoragePayload;
}

// modify_chain_id

export interface IModifyChainIDPayload {
    chain_id: string;
}
export interface IModifyChainIDAction {
    kind: ActionKind.ModifyChainID;
    payload: IModifyChainIDPayload;
}

// modify_block_level
export interface IModifyBlockLevelPayload {
    level: number;
}
export interface IModifyBlockLevelAction {
    kind: ActionKind.ModifyBlockLevel;
    payload: IModifyBlockLevelPayload;
}

// modify_block_timestamp

export interface IModifyBlockTimestampPayload {
    timestamp: string;
}
export interface IModifyBlockTimestampAction {
    kind: ActionKind.ModifyBlockTimestamp;
    payload: IModifyBlockTimestampPayload;
}

// pack_data

export interface IPackDataPayload {
    data: Record<string, unknown> | Record<string, unknown>[];
    type: Record<string, unknown> | Record<string, unknown>[];
}
export interface IPackDataAction {
    kind: ActionKind.PackData;
    payload: IPackDataPayload;
}
