import type { IBudgetListItem } from '~/modules/budget/domain/model/types/budget';
import { tryToCatchApiErrors } from '~/shared/errors/errors';

const DEFAULT_LIMIT = 20;

export const fetchBudgetsList = async (page?: number, limit?: number) => {
    if (!page) page = 1;
    if (!limit) limit = DEFAULT_LIMIT;

    try {
        return await useNuxtApp().$apiFetch<{ items: IBudgetListItem[]; total: number }>('v1/ledger/budgets', {
            params: {
                offset: (page - 1) * limit,
                limit,
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
