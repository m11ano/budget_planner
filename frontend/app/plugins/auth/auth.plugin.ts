import { doAuth } from './model';
import type { IAuthData } from './model/types/auth.types';

export default defineNuxtPlugin({
    name: 'auth',
    enforce: 'pre',
    async setup() {
        const authData = useState<IAuthData>('authData', () =>
            reactive({
                isLoading: false,
                isAuth: false,
                userData: null,
            }),
        );

        if (import.meta.client) {
            authData.value.isLoading = true;

            let loaderDiv;
            if (window) {
                loaderDiv = window.document.createElement('DIV');
                loaderDiv.id = 'pageLoader';
                loaderDiv.innerHTML = '<div>Загрузка...</div>';
                window.document.body.append(loaderDiv);
            }

            try {
                await doAuth();
            } catch (e: unknown) {
                //
            }

            authData.value.isLoading = false;

            if (loaderDiv) {
                loaderDiv.remove();
            }
        }

        return {
            provide: {
                authData: authData.value,
            },
        };
    },
    env: {
        islands: true,
    },
});
