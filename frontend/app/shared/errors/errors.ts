import { FetchError } from 'ofetch';
import { declOfNum } from '../helpers/functions';

export class ApiError extends Error {
    hints: string[] = [];
    details: Record<string, any> = {};
    code: number = 400;
    retryAfter: number = 0;
    textCode: string = 'BAD_REQUEST';

    constructor(text: string = 'bad request to api') {
        super(text);
    }

    setCode(code: number): ApiError {
        this.code = code;
        return this;
    }

    setTextCode(textCode: string): ApiError {
        this.textCode = textCode;
        return this;
    }

    setHints(hints: string[]): ApiError {
        this.hints = hints;
        return this;
    }

    setDetails(details: Record<string, any>): ApiError {
        this.details = details;
        return this;
    }

    setRetryAfter(retryAfter: number): ApiError {
        if (retryAfter < 1) {
            retryAfter = 1;
        }

        this.retryAfter = retryAfter;
        return this;
    }

    errorTextByData() {
        let hint = '';

        if (this.code === 429) {
            if (this.textCode === 'BACKOFF') {
                hint = `Ограничение запросов. Вы можете повторить попытку через ${formatRetryAfter(this.retryAfter)}`;
            } else {
                hint = `От вас поступает слишком много запросов. Повторите попытку позже`;
            }
        } else if (this.code === 409) {
            hint = `Конфликт данных`;
        } else if (this.code >= 500) {
            hint = `Внутренняя ошибка на сервере, код: ${this.code}`;
        } else if (this.code === 403) {
            hint = `Нет доступа к ресурсу`;
        } else if (this.code === 401) {
            hint = `Вы не авторизованы`;
        } else if (this.code >= 400 && this.code < 500) {
            if (this.textCode === 'INVALID_IMAGE_SIZE') {
                hint = `Некорректное разрешение изображения`;
            } else {
                hint = `Некорректные данные запроса`;
            }
        }

        return hint;
    }

    formHints(): string[] {
        const extraHints: string[] = [];

        const dataHint = this.errorTextByData();
        if (dataHint) {
            extraHints.push(dataHint);
        }

        return [...extraHints, ...this.hints];
    }
}

export function tryToCatchApiErrors(e: unknown): any {
    if (e instanceof FetchError && e.statusCode) {
        const err = new ApiError().setCode(e.statusCode);

        if (e.data) {
            if (typeof e.data.textCode == 'string') {
                err.setTextCode(e.data.textCode);
            }

            if (typeof e.data.details == 'object' && e.data.details !== null) {
                err.setDetails(e.data.details);
            }

            if (typeof e.data.hints == 'object' && e.data.hints !== null) {
                err.setHints(e.data.hints);
            }
        }

        if (err.details['try_after_sec'] !== undefined) {
            err.setRetryAfter(Number(err.details['try_after_sec']));
        }

        return err;
    }

    return e;
}

function formatRetryAfter(seconds: number): string {
    if (seconds < 60) {
        return `${seconds} ${declOfNum(seconds, ['секунда', 'секунды', 'секунд'])}`;
    }
    const minutes = Math.floor(seconds / 60);
    const secRem = seconds % 60;
    const minuteStr = `${minutes} ${declOfNum(minutes, ['минута', 'минуты', 'минут'])}`;

    if (secRem === 0) {
        return minuteStr;
    }

    const secondStr = `${secRem} ${declOfNum(secRem, ['секунда', 'секунды', 'секунд'])}`;

    return `${minuteStr} ${secondStr}`;
}
