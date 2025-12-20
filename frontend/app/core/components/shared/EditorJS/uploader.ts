import { ApiError } from '~/shared/errors/errors';
import { fetchUpload } from '../FileUploader/api/fetchUpload';
import { showErrors } from '../inform/toast';

export async function uploadToApi(url: string, file: File) {
    return new Promise((ready, reject) => {
        const limit = useNuxtApp().$config.public.maxFilesize;
        if (file.size > limit) {
            showErrors([`Файл ${file.name} слишком большой, лимит ${limit / 1024 / 1024} МБ`]);
            reject();
            return;
        }

        fetchUpload(url, file)
            .then((uploaded) => {
                const data = Object.entries(uploaded.files)[0]!;
                setTimeout(() => {
                    ready(data);
                }, 1000);
            })
            .catch((e) => {
                if (e instanceof ApiError) {
                    showErrors(e.formHints());
                }
                reject();
            });
    })
        .then((data: any) => {
            return {
                success: 1,
                file: {
                    url: data[1].url,
                    type: 'file',
                    fileID: data[1].id,
                    filename: data[1].filename,
                    fileTarget: data[0],
                },
            };
        })
        .catch(() => {
            return {
                success: 0,
            };
        });
}
