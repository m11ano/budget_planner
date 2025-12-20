import { AUTH_ACCESS_TOKEN_KEY } from '../const/const';
import type { IAuthData, IAuthTokens, IAuthUserData } from '../types/auth.types';

export const setAuthUserData = async (userData: IAuthUserData) => {
    const authData = useState<IAuthData>('authData');
    authData.value.isAuth = true;
    authData.value.userData = userData;
};

export const setTokens = async (tokens: IAuthTokens) => {
    localStorage.setItem(AUTH_ACCESS_TOKEN_KEY, tokens.accessJWT);
};

export const clearTokens = async () => {
    localStorage.removeItem(AUTH_ACCESS_TOKEN_KEY);
};
