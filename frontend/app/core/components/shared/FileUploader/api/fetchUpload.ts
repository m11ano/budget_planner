import { ApiError, tryToCatchApiErrors } from '~/shared/errors/errors';

export interface FetchUploadedFile {
    id: string;
    url: string;
    filename: string;
}

export interface UploadFileResponse {
    files: Record<string, FetchUploadedFile>;
}

export async function fetchUpload(url: string, file: File) {
    const formData = new FormData();
    formData.append('file', file);

    try {
        return await useNuxtApp().$apiFetch<UploadFileResponse>(url, {
            method: 'POST',
            body: formData,
        });
    } catch (e: unknown) {
        const err = tryToCatchApiErrors(e);
        if (err instanceof ApiError) {
            if (err.code === 429) {
                await new Promise<void>((resolve) => setTimeout(resolve, 1000));
                return fetchUpload(url, file);
            }
        }

        throw err;
    }
}
