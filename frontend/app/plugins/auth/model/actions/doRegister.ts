import { tryToCatchApiErrors } from '~/shared/errors/errors';

export const doRegister = async (email: string, password: string, profileName: string, profileSurname: string): Promise<void> => {
    try {
        await useNuxtApp().$apiFetch('v1/auth/register', {
            method: 'POST',
            body: {
                email,
                password,
                profileName,
                profileSurname,
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
