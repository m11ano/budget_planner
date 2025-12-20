import type { IAuthData } from '../types/auth.types';

export const clearAuthUserData = async () => {
    const authData = useState<IAuthData>('authData');
    authData.value.isAuth = false;
    authData.value.userData = null;
};
