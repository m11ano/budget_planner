import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IUserAccount } from '../model/types/users';

interface Input {
    email: string;
    roleID: number;
    password: string;
    profileName: '';
    profileSurname: '';
    profilePhotoOriginalFileID: null | string;
    profilePhoto100x100FileID: null | string;
}

export async function createUser(input: Input) {
    try {
        return await useNuxtApp().$apiFetch<IUserAccount>(`v1/users`, {
            method: 'POST',
            body: input,
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
