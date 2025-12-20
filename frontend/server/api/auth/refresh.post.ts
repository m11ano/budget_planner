import { FetchError } from 'ofetch';
import { AUTH_REFRESH_TOKEN_KEY } from '~/plugins/auth/model/const/const';
import { GetDefaultHeaders } from '~/shared/api/headers';

type BackendResponseDTO = {
    accessJWT: string;
    refreshJWT: string;
    [key: string]: any;
};

export default defineEventHandler(async (event) => {
    const config = useRuntimeConfig();

    const refreshFromCookie = getCookie(event, AUTH_REFRESH_TOKEN_KEY);

    if (!refreshFromCookie) {
        throw createError({
            statusCode: 401,
        });
    }

    try {
        const baseURL = `${config.public.apiServerPrefix}${config.public.apiBase}`;

        const response = await $fetch<BackendResponseDTO>('/v1/auth/refresh', {
            baseURL: baseURL,
            method: 'POST',
            headers: GetDefaultHeaders(event),
            body: {
                refreshJWT: refreshFromCookie,
            },
        });

        const { accessJWT, refreshJWT, ...otherResponse } = response;

        const host = getRequestHeader(event, 'host') || undefined;
        const cookieDomain = host?.split(':')[0];

        setCookie(event, AUTH_REFRESH_TOKEN_KEY, refreshJWT, {
            httpOnly: true,
            secure: config.public.isProd,
            sameSite: 'strict',
            path: '/bff-api/auth',
            maxAge: 60 * 60 * 24 * 30,
            domain: cookieDomain,
        });

        return { ...otherResponse, accessJWT };
    } catch (e: unknown) {
        if (e instanceof FetchError && e.statusCode) {
            setResponseStatus(event, e.statusCode);
            return e.data;
        }

        throw e;
    }
});
