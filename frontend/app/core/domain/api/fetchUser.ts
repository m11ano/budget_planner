import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IUserAccount } from '../model/types/users';

export const fetchUser = async (profileID: string) => {
    try {
        return await useNuxtApp().$apiFetch<IUserAccount>(`v1/users/${profileID}`);
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
