import type { ITransactionListItem } from '~/modules/budget/domain/model/types/transaction';
import { tryToCatchApiErrors } from '~/shared/errors/errors';

const DEFAULT_LIMIT = 20;

export const fetchTransactionsList = async (page?: number, limit?: number) => {
    if (!page) page = 1;
    if (!limit) limit = DEFAULT_LIMIT;

    try {
        return await useNuxtApp().$apiFetch<{ items: ITransactionListItem[]; total: number }>('v1/ledger/transactions', {
            params: {
                offset: (page - 1) * limit,
                limit,
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
