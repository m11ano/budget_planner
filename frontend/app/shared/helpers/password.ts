export interface GeneratePasswordOptions {
    lowercase?: boolean;
    uppercase?: boolean;
    digits?: boolean;
    symbols?: boolean;
    avoidAmbiguous?: boolean;
    customChars?: string;
    blacklist?: string;
    enforceEach?: boolean;
}

export function generatePassword(length: number, opts: GeneratePasswordOptions = {}): string {
    if (!Number.isInteger(length) || length < 1) {
        throw new Error('length должен быть целым числом ≥ 1');
    }

    const options: Required<GeneratePasswordOptions> = {
        lowercase: true,
        uppercase: true,
        digits: true,
        symbols: true,
        avoidAmbiguous: false,
        customChars: '',
        blacklist: '',
        enforceEach: true,
        ...opts,
    };

    let pools: string[][] = [];
    if (options.lowercase) pools.push([...'abcdefghijklmnopqrstuvwxyz']);
    if (options.uppercase) pools.push([...'ABCDEFGHIJKLMNOPQRSTUVWXYZ']);
    if (options.digits) pools.push([...'0123456789']);
    if (options.symbols) pools.push([...'!@#$%^&*()-_=+[]{};:,<.>/?\\|~`']);
    if (options.customChars) pools.push([...options.customChars]);

    if (options.avoidAmbiguous) {
        const amb = new Set('0Oo1Il|`\'";:.,');
        pools = pools.map((arr) => arr.filter((ch) => !amb.has(ch)));
    }

    if (options.blacklist) {
        const bl = new Set(options.blacklist);
        pools = pools.map((arr) => arr.filter((ch) => !bl.has(ch)));
    }

    const all: string[] = Array.from(new Set(pools.flat()));
    if (all.length === 0) throw new Error('Пустой алфавит — проверьте опции.');

    if (options.enforceEach) {
        const active = pools.filter((a) => a.length > 0).length;
        if (length < active) {
            throw new Error(`length (${length}) меньше числа активных пулов (${active}).`);
        }
    }

    if (typeof globalThis.crypto?.getRandomValues !== 'function') {
        throw new Error('Web Crypto недоступен (crypto.getRandomValues).');
    }

    function randInt(max: number): number {
        if (max <= 0) throw new Error('max должен быть > 0');
        const range = 0x100000000;
        const limit = Math.floor(range / max) * max;
        const buf = new Uint32Array(1);
        let x = 0;
        do {
            globalThis.crypto.getRandomValues(buf);
            x = buf[0]!;
        } while (x >= limit);
        return x % max;
    }

    function sample<T>(arr: T[]): T {
        if (arr.length === 0) throw new Error('Нельзя выбирать из пустого массива');
        return arr[randInt(arr.length)]!;
    }

    function shuffle<T>(arr: T[]): T[] {
        for (let i = arr.length - 1; i > 0; i--) {
            const j = randInt(i + 1);
            const tmp = arr[i];
            arr[i] = arr[j]!;
            arr[j] = tmp!;
        }
        return arr;
    }

    const out: string[] = [];

    if (options.enforceEach) {
        for (const s of pools) {
            if (s.length === 0) continue;
            out.push(sample(s));
        }
    }

    while (out.length < length) {
        out.push(sample(all));
    }

    return shuffle(out).slice(0, length).join('');
}

export default generatePassword;
