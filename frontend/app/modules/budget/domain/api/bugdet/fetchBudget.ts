import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IBudgetItem } from '../../model/types/budget';

export const fetchBudget = async (id: string, controller?: AbortController) => {
    try {
        const abortController = controller ?? new AbortController();

        return await useNuxtApp().$apiFetch<{ item: IBudgetItem }>(`v1/ledger/budgets/${id}`, {
            signal: abortController.signal,
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
