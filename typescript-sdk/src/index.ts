import {
    ActionKind,
    IAction,
    ICallContractPayload,
    ICreateImplicitAccountPayload,
    IOriginateContractPayload,
} from './action';

/**
 * Builds an action responsible for creating an implicit account
 *
 * @param payload action payload
 * @returns action
 */
export const buildCreateImplicitAccountAction = (payload: ICreateImplicitAccountPayload): IAction => ({
    kind: ActionKind.CreateImplicitAccount,
    payload,
});

/**
 * Builds an action responsible for originating a contract
 *
 * @param payload action payload
 * @returns action
 */
export const buildOriginateContractAction = (payload: IOriginateContractPayload): IAction => ({
    kind: ActionKind.OriginateContract,
    payload,
});

/**
 * Builds an action responsible for calling a contract
 *
 * @param payload action payload
 * @returns action
 */
export const buildCallContractAction = (payload: ICallContractPayload): IAction => ({
    kind: ActionKind.CallContract,
    payload,
});
