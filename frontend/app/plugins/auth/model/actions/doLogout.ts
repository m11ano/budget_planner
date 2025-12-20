import { AUTH_ACCESS_TOKEN_KEY } from '../const/const';
import { clearAuthUserData } from './clearAuthData';
import { clearTokens } from './setAuthData';

export const doLogout = async (): Promise<void> => {
    setTimeout(() => {
        const accessToken = localStorage.getItem(AUTH_ACCESS_TOKEN_KEY);

        clearAuthUserData();
        clearTokens();

        $fetch('/bff-api/auth/logout', {
            method: 'POST',
            headers: {
                Authorization: accessToken || '',
            },
        });
    }, 100);

    navigateTo('/');
};
