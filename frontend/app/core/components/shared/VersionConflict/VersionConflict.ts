import type { ApiError } from '~/shared/errors/errors';
import Confirm from './modals/Confirm.vue';

interface Output<T> {
    isConflict: boolean;
    isCancel: boolean;
    result?: T;
}

export async function OnVersionConflict<T>(err: ApiError, onContinue: () => Promise<T>): Promise<Output<T>> {
    if (err.code === 409 && err.textCode === 'VERSION_CONFLICT') {
        const modal = useOverlay().create(Confirm, {
            props: {
                text: `Кто-то изменил данные, пока вы редактировали форму. Если вы продолжите, предыдущие изменения могут быть потеряны. Рекомендуется обновить страницу и заново внести изменения. Вы действительно хотите продолжить?`,
            },
            destroyOnClose: true,
        });

        const instance = modal.open();

        const shouldDo = await instance.result;
        if (shouldDo) {
            const contRes = await onContinue();

            return {
                isConflict: true,
                isCancel: false,
                result: contRes,
            };
        }

        return {
            isConflict: true,
            isCancel: true,
        };
    }

    return {
        isConflict: false,
        isCancel: false,
    };
}
