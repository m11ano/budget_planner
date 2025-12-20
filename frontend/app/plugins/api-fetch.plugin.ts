import { GetDefaultHeaders } from '~/shared/api/headers';
import { ApiError } from '~/shared/errors/errors';
import { SingleFlight } from '~/shared/helpers/singleflight';
import { doLogout } from './auth/model';
import { doRefresh } from './auth/model/actions/doRefresh';
import { AUTH_ACCESS_TOKEN_KEY } from './auth/model/const/const';

export default defineNuxtPlugin({
    name: 'api-fetch',
    enforce: 'pre',
    async setup(nuxtApp) {
        const singleflight = new SingleFlight();

        const apiFetch = $fetch.create({
            baseURL: nuxtApp.$config.public.apiBase,
            retry: 1,
            retryStatusCodes: [401],
            retryDelay: 0,
            headers: GetDefaultHeaders(undefined),

            async onRequest({ options }) {
                if ((options as any)._meta === undefined) {
                    (options as any)._meta = {};
                }
                const accessToken = localStorage.getItem(AUTH_ACCESS_TOKEN_KEY);
                if (accessToken) {
                    options.headers.set('Authorization', `Bearer ${accessToken}`);
                }
            },

            async onResponseError({ response, error, options }) {
                if (response.status === 401 && (options as any)._meta.refresh_tried === undefined) {
                    const accessToken = localStorage.getItem(AUTH_ACCESS_TOKEN_KEY);
                    if (!accessToken) {
                        throw error;
                    }

                    try {
                        await singleflight.do('refresh', async () => {
                            await doRefresh();
                        });
                    } catch (e: unknown) {
                        if (e instanceof ApiError) {
                            if (e.code === 401) {
                                doLogout();
                            }
                        }

                        throw error;
                    } finally {
                        (options as any)._meta.refresh_tried = true;
                    }
                }
            },
        });

        return {
            provide: {
                apiFetch,
            },
        };
    },
});
