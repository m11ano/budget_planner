import { FetchError } from 'ofetch';
import { AUTH_REFRESH_TOKEN_KEY } from '~/plugins/auth/model/const/const';
import { GetDefaultHeaders } from '~/shared/api/headers';

type BackendResponseDTO = {
    tokens: {
        accessJWT: string;
        refreshJWT: string;
    };
    [key: string]: any;
};

export default defineEventHandler(async (event) => {
    const config = useRuntimeConfig();

    const authHeader = getRequestHeader(event, 'authorization');

    if (!authHeader) {
        throw createError({
            statusCode: 401,
        });
    }

    try {
        const baseURL = `${config.public.apiServerPrefix}${config.public.apiBase}`;

        const headers = GetDefaultHeaders(event);
        headers.set('Authorization', authHeader);

        const response = await $fetch<BackendResponseDTO>('/v1/auth/logout', {
            baseURL: baseURL,
            method: 'POST',
            headers,
        });

        const host = getRequestHeader(event, 'host') || undefined;
        const cookieDomain = host?.split(':')[0];

        deleteCookie(event, AUTH_REFRESH_TOKEN_KEY, {
            httpOnly: true,
            secure: config.public.isProd,
            sameSite: 'strict',
            path: '/bff-api/auth',
            domain: cookieDomain,
        });

        return response;
    } catch (e: unknown) {
        if (e instanceof FetchError && e.statusCode) {
            setResponseStatus(event, e.statusCode);
            return e.data;
        }

        throw e;
    }
});
