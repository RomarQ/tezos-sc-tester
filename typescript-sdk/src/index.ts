import { ActionKind, IAction } from './action';

export interface TestSuite {
    protocol?: string;
    actions: IAction[];
}

/**
 * Builds an action
 *
 * @param payload action payload
 * @returns IAction
 */
export const buildAction = <T extends ActionKind, A extends IAction = Extract<IAction, { kind: T }>>(
    kind: T,
    payload: A['payload'],
): IAction =>
    ({
        kind,
        payload,
    } as IAction);
