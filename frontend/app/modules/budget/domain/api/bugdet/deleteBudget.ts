import { tryToCatchApiErrors } from '~/shared/errors/errors';

export async function deleteBudget(id: string) {
    try {
        return await useNuxtApp().$apiFetch<Response>(`v1/ledger/budgets/${id}`, {
            method: 'DELETE',
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
