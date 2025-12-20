import { tryToCatchApiErrors } from '~/shared/errors/errors';

export async function deleteTransaction(id: string) {
    try {
        return await useNuxtApp().$apiFetch<Response>(`v1/ledger/transactions/${id}`, {
            method: 'DELETE',
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
