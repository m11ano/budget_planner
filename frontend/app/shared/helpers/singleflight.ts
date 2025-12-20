type AnyFn<T> = () => Promise<T> | T;

export interface DoResult<T> {
    value: T;
    shared: boolean;
}

export class SingleFlight {
    private inFlight = new Map<
        string,
        {
            promise: Promise<unknown>;
            waiters: number;
        }
    >();

    do<T>(key: string, fn: AnyFn<T>): Promise<DoResult<T>> {
        const existing = this.inFlight.get(key);
        if (existing) {
            existing.waiters++;
            return existing.promise.then(
                (value) => ({ value: value as T, shared: true }),
                (err) => Promise.reject(err),
            );
        }

        let resolveOuter!: (v: T) => void;
        let rejectOuter!: (e: unknown) => void;

        const p = new Promise<T>((resolve, reject) => {
            resolveOuter = resolve;
            rejectOuter = reject;
        });

        this.inFlight.set(key, { promise: p, waiters: 0 });

        (async () => {
            try {
                const v = await fn();
                resolveOuter(v);
            } catch (e) {
                rejectOuter(e);
            } finally {
                this.inFlight.delete(key);
            }
        })();

        return p.then(
            (value) => ({ value, shared: false }),
            (err) => Promise.reject(err),
        );
    }

    forget(key: string) {
        this.inFlight.delete(key);
    }

    waiters(key: string): number {
        return this.inFlight.get(key)?.waiters ?? 0;
    }

    has(key: string): boolean {
        return this.inFlight.has(key);
    }
}
