import type { IUserAccount } from '~/core/domain/model/types/users';

export interface IAuthUserData {
    account: IUserAccount;
}

export interface IAuthTokens {
    accessJWT: string;
}

export interface IAuthData {
    isLoading: boolean;
    isAuth: boolean;
    userData: IAuthUserData | null;
}
