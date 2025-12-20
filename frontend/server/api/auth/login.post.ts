import { FetchError } from 'ofetch';
import { AUTH_REFRESH_TOKEN_KEY } from '~/plugins/auth/model/const/const';
import { GetDefaultHeaders } from '~/shared/api/headers';

interface RequestDTO {
    email: string;
    password: string;
}

type BackendResponseDTO = {
    tokens: {
        accessJWT: string;
        refreshJWT: string;
    };
    [key: string]: any;
};

export default defineEventHandler(async (event) => {
    const config = useRuntimeConfig();

    const request = await readBody<RequestDTO>(event);
    if (typeof request?.email !== 'string' || typeof request?.password !== 'string') {
        throw createError({
            statusCode: 400,
        });
    }

    try {
        const baseURL = `${config.public.apiServerPrefix}${config.public.apiBase}`;

        const response = await $fetch<BackendResponseDTO>('/v1/auth/login', {
            baseURL: baseURL,
            method: 'POST',
            headers: GetDefaultHeaders(event),
            body: {
                email: request.email,
                password: request.password,
            },
        });

        const { tokens, ...otherResponse } = response;

        const host = getRequestHeader(event, 'host') || undefined;
        const cookieDomain = host?.split(':')[0];

        setCookie(event, AUTH_REFRESH_TOKEN_KEY, tokens.refreshJWT, {
            httpOnly: true,
            secure: config.public.isProd,
            sameSite: 'strict',
            path: '/bff-api/auth',
            maxAge: 60 * 60 * 24 * 30,
            domain: cookieDomain,
        });

        return { ...otherResponse, tokens: { accessJWT: tokens.accessJWT } };
    } catch (e: unknown) {
        if (e instanceof FetchError && e.statusCode) {
            setResponseStatus(event, e.statusCode);
            return e.data;
        }

        throw e;
    }
});
