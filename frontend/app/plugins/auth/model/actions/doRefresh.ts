import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IAuthTokens } from '../types/auth.types';
import { setTokens } from './setAuthData';

export const doRefresh = async (): Promise<void> => {
    try {
        const result = await $fetch<IAuthTokens>('/bff-api/auth/refresh', {
            method: 'POST',
        });

        setTokens(result);
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
