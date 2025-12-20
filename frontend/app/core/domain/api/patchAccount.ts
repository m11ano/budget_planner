import { OnVersionConflict } from '~/core/components/shared/VersionConflict/VersionConflict';
import { ApiError, tryToCatchApiErrors } from '~/shared/errors/errors';

interface Input {
    _version: number;
    _skipVersionCheck?: boolean;
    isBlocked?: boolean;
    password?: string;
    roleID?: number;
    profileName?: string;
    profileSurname?: string;
    profilePhotos?: {
        photoOriginalFileID?: string | null;
        photo100x100FileID?: string | null;
    };
}

export async function patchAccount(id: string, input: Input) {
    try {
        await useNuxtApp().$apiFetch(`v1/users/${id}`, {
            method: 'PATCH',
            body: input,
        });
    } catch (e: unknown) {
        const err = tryToCatchApiErrors(e);
        if (err instanceof ApiError) {
            const res = await OnVersionConflict(err, () => patchAccount(id, { ...input, _skipVersionCheck: true }));
            if (res.isConflict && !res.isCancel) {
                return;
            }
        }

        throw err;
    }
}
