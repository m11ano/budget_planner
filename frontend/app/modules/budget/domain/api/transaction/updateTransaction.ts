import { OnVersionConflict } from '~/core/components/shared/VersionConflict/VersionConflict';
import { ApiError, tryToCatchApiErrors } from '~/shared/errors/errors';
import type { ITransactionItemState } from '../../model/types/transaction';

type Input = {
    _version?: number;
    _skipVersionCheck?: boolean;
} & ITransactionItemState;

type Payload = Input;

const mapDataToRequest = (data: Input): Payload => {
    return data;
};

export async function updateTransaction(id: string, input: Input) {
    try {
        await useNuxtApp().$apiFetch(`v1/ledger/transactions/${id}`, {
            method: 'PATCH',
            body: mapDataToRequest(input),
        });
    } catch (e: unknown) {
        const err = tryToCatchApiErrors(e);
        if (err instanceof ApiError) {
            const res = await OnVersionConflict(err, () => updateTransaction(id, { ...input, _skipVersionCheck: true }));
            if (res.isConflict && !res.isCancel) {
                return;
            }
        }

        throw err;
    }
}
