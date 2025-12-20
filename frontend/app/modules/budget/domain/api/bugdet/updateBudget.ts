import { OnVersionConflict } from '~/core/components/shared/VersionConflict/VersionConflict';
import { ApiError, tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IBudgetItemState } from '../../model/types/budget';

type Input = {
    _version?: number;
    _skipVersionCheck?: boolean;
} & IBudgetItemState;

type Payload = Input;

const mapDataToRequest = (data: Input): Payload => {
    return data;
};

export async function updateBudget(id: string, input: Input) {
    try {
        await useNuxtApp().$apiFetch(`v1/ledger/budgets/${id}`, {
            method: 'PATCH',
            body: mapDataToRequest(input),
        });
    } catch (e: unknown) {
        const err = tryToCatchApiErrors(e);
        if (err instanceof ApiError) {
            const res = await OnVersionConflict(err, () => updateBudget(id, { ...input, _skipVersionCheck: true }));
            if (res.isConflict && !res.isCancel) {
                return;
            }
        }

        throw err;
    }
}
