import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IBudgetItem, IBudgetItemState } from '../../model/types/budget';

type Input = IBudgetItemState;

type Payload = Input;

const mapDataToRequest = (data: Input): Payload => {
    return data;
};

export async function createBudget(input: Input) {
    try {
        return await useNuxtApp().$apiFetch<{ item: IBudgetItem }>(`v1/ledger/budgets`, {
            method: 'POST',
            body: mapDataToRequest(input),
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
