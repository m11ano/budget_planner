import type { EventHandlerRequest, H3Event } from 'h3';

export function GetDefaultHeaders(event: H3Event<EventHandlerRequest> | undefined) {
    const headers = new Headers({
        'X-Check': 'true',
    });

    try {
        if (event) {
            const req = event.node.req;

            const xRealIP = getRequestHeader(event, 'x-real-ip')?.trim() || '';
            const xffHeader = getRequestHeader(event, 'x-forwarded-for') || '';
            const socketIP = req.socket.remoteAddress || '';

            const clientIP = xRealIP || xffHeader || socketIP || '';

            headers.set('X-Real-IP', clientIP || '');
            headers.set('X-Forwarded-For', clientIP || '');

            headers.set('X-Forwarded-Proto', getRequestProtocol(event) || 'http');
            headers.set('X-Forwarded-Host', getRequestHeader(event, 'host') || '');
        }

        let currentURL = '';

        if (import.meta.client) {
            currentURL = window.location.href;
        } else if (event) {
            const proto = getRequestProtocol(event) || 'http';
            const host = getRequestHeader(event, 'host') || '';
            const path = event.node.req.url || '/';
            currentURL = `${proto}://${host}${path}`;
        }

        if (currentURL) {
            const safeURL = currentURL.slice(0, 2048);
            headers.set('X-Referer', safeURL);
        }
    } catch (e) {
        console.error(e);
    }

    return headers;
}
