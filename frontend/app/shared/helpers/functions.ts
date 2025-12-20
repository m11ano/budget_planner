export function nextAnimFrame(f: () => void) {
    if (!window) {
        f();
        return;
    }
    const raf = window.requestAnimationFrame || (window as any).webkitRequestAnimationFrame || (window as any).msRequestAnimationFrame;
    if (raf) {
        raf(() => {
            raf(() => {
                f();
            });
        });
    } else {
        f();
    }
}

export function waitNextAnimFrame(): Promise<void> {
    return new Promise((resolve) => {
        nextAnimFrame(resolve);
    });
}

interface elementAnimateOptions {
    transition: string;
    properties: Record<string, string>;
}

export function elementAnimate(elem: HTMLElement, options: elementAnimateOptions) {
    return new Promise<() => void>((resolve) => {
        const clearProperties = () => {
            Object.keys(options.properties).forEach((k) => {
                elem.style.removeProperty(k);
            });
        };

        const endHandler = (e: TransitionEvent) => {
            if (e.target === elem) {
                elem.removeEventListener('transitionend', endHandler);
                elem.style.removeProperty('transition');
                resolve(clearProperties);
            }
        };
        if (elem) {
            elem.addEventListener('transitionend', endHandler);
            elem.style.setProperty('transition', options.transition, 'important');
            nextAnimFrame(() => {
                Object.entries(options.properties).forEach(([k, value]) => {
                    elem.style.setProperty(k, value);
                });
            });
        }
    });
}

export function windowScrollToBlock(v: string, offset: number | undefined = undefined, is_smooth = true) {
    const elem = document.getElementById(v);
    if (!elem) return;
    if (offset === undefined) offset = window.innerWidth < 576 ? 70 : 100;
    windowScrollToElem(elem, offset, is_smooth);
}

export function windowScrollToElem(elem: HTMLElement, offset: number | undefined = undefined, is_smooth = true) {
    if (offset === undefined) offset = window.innerWidth < 576 ? 70 : 100;
    const s = elem.getBoundingClientRect().top + scrollY - offset;
    if (is_smooth) {
        windowScrollTo(s);
    } else {
        window.scrollTo({ top: s });
    }
}

export function windowScrollTo(v: number) {
    window.scrollTo({ top: v, behavior: 'smooth' });
}

export function numberFormat(input: number | string, decimals: number = 0, decPoint: string = '.', thousandsSep: string = ','): string {
    const cleaned = `${input}`.replace(/[^0-9+\-Ee.]/g, '');
    const n = Number(cleaned);
    if (!Number.isFinite(n)) {
        return '0';
    }

    const prec = Math.abs(Number.isFinite(decimals) ? decimals : 0);

    const toFixedFix = (value: number, precision: number): string => {
        const k = 10 ** precision;
        return (Math.round(value * k) / k).toFixed(precision);
    };

    const formatted = prec > 0 ? toFixedFix(n, prec) : `${Math.round(n)}`;

    const [intPart = '', fracPart = ''] = formatted.split('.');

    const s: [string, string] = [intPart, fracPart];

    if (s[0].length > 3) {
        s[0] = s[0].replace(/\B(?=(?:\d{3})+(?!\d))/g, thousandsSep);
    }

    if (s[1].length < prec) {
        s[1] = s[1].padEnd(prec, '0');
    }

    return s.join(decPoint);
}

export function coolNumber(v: number): string {
    const isInt = Number.isInteger(v);

    return new Intl.NumberFormat('ru-RU', {
        minimumFractionDigits: isInt ? 0 : 2,
        maximumFractionDigits: isInt ? 0 : 2,
    }).format(v);
}

export function declOfNum(n: number, titles: [string, string, string]): string {
    const cases = [2, 0, 1, 1, 1, 2] as const;

    const idx = n % 100 > 4 && n % 100 < 20 ? 2 : cases[Math.min(n % 10, 5)]!;

    return titles[idx];
}

export function parseDate(dateStr: string): Date {
    if (typeof dateStr !== 'string') {
        throw new Error(`parseDate: expected string, got ${typeof dateStr}`);
    }

    const trimmed = dateStr.trim();
    if (trimmed.length === 0) {
        throw new Error('parseDate: empty string is not a valid date');
    }

    const direct = new Date(trimmed);
    if (!isNaN(direct.getTime())) {
        return direct;
    }

    const ymdMatch = trimmed.match(/^(\d{4})-(\d{2})-(\d{2})$/);
    if (ymdMatch) {
        const [_, y, m, d] = ymdMatch;
        return new Date(Date.UTC(Number(y), Number(m) - 1, Number(d), 0, 0, 0, 0));
    }

    throw new Error(`parseDate: unsupported date format "${dateStr}"`);
}

export function calcAge(birthDate: Date, compareDate: Date = new Date()): number {
    const birthYear = birthDate.getUTCFullYear();
    const birthMonth = birthDate.getUTCMonth();
    const birthDay = birthDate.getUTCDate();

    const currentYear = compareDate.getUTCFullYear();
    const currentMonth = compareDate.getUTCMonth();
    const currentDay = compareDate.getUTCDate();

    let age = currentYear - birthYear;

    if (currentMonth < birthMonth || (currentMonth === birthMonth && currentDay < birthDay)) {
        age--;
    }

    return age;
}
