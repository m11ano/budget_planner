import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IUserAccount } from '../model/types/users';

const DEFAULT_LIMIT = 20;

export const fetchAccountsList = async (page?: number, limit?: number) => {
    if (!page) page = 1;
    if (!limit) limit = DEFAULT_LIMIT;

    try {
        return await useNuxtApp().$apiFetch<{ items: IUserAccount[]; total: number }>('v1/users', {
            params: {
                offset: (page - 1) * limit,
                limit,
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
