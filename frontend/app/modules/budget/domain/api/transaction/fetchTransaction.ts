import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { ITransactionItem } from '../../model/types/transaction';

export const fetchTransaction = async (id: string, controller?: AbortController) => {
    try {
        const abortController = controller ?? new AbortController();

        return await useNuxtApp().$apiFetch<{ item: ITransactionItem }>(`v1/ledger/transactions/${id}`, {
            signal: abortController.signal,
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
