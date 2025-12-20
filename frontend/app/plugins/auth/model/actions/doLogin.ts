import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IAuthTokens, IAuthUserData } from '../types/auth.types';
import { setAuthUserData, setTokens } from './setAuthData';

export const doLogin = async (email: string, password: string): Promise<void> => {
    try {
        const result = await $fetch<IAuthUserData & { tokens: IAuthTokens }>('/bff-api/auth/login', {
            method: 'POST',
            body: {
                email,
                password,
            },
        });

        setAuthUserData({
            account: result.account,
        });

        setTokens({
            accessJWT: result.tokens.accessJWT,
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
