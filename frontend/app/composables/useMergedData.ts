import type { AsyncDataOptions, NuxtError } from '#app';
import type { UnwrapNestedRefs } from 'vue';

type PartialNull<T> = { [P in keyof T]: T[P] | null };

type AsyncData<T> = {
    error: Ref<NuxtError | null>;
    data: UnwrapNestedRefs<PartialNull<T>>;
    isFetching: Ref<boolean>;
    isAllNotNull: Ref<boolean>;
};

export function useMergedData<T extends Record<string, any>>(
    keyPrefix: string,
    fetchers: Record<string, () => Promise<any>>,
    makeRefetch: boolean = true,
    fetchOptions: AsyncDataOptions<any> = {},
): AsyncData<T> {
    fetchOptions = { lazy: true, ...fetchOptions };
    if (!makeRefetch) {
        fetchOptions.getCachedData = (key) => {
            return useNuxtApp().payload.data[key] || useNuxtApp().static.data[key];
        };
    }

    const asyncResults = reactive<Record<string, null | ReturnType<typeof useAsyncData>>>(
        Object.keys(fetchers).reduce(
            (acc, key) => {
                acc[key] = useAsyncData(`${keyPrefix}:${key}`, fetchers[key], fetchOptions);
                return acc;
            },
            {} as Record<string, null | ReturnType<typeof useAsyncData>>,
        ),
    );

    const isAllNotNull = computed(() => {
        for (const key in asyncResults) {
            if (asyncResults[key] === null || asyncResults[key].data.value === null) {
                return false;
            }
        }

        return true;
    });

    const error = computed(() => {
        for (const key in asyncResults) {
            if (asyncResults[key] !== null && asyncResults[key].error !== null) {
                return asyncResults[key].error.value;
            }
        }
        return null;
    });

    const isFetching = computed(() => {
        for (const key in asyncResults) {
            if (asyncResults[key] === null || asyncResults[key].status.value === 'pending') {
                return true;
            }
        }

        return false;
    });

    const data = reactiveComputed<PartialNull<T>>(() => {
        return Object.keys(asyncResults).reduce((acc, key) => {
            //@ts-ignore
            acc[key] = asyncResults[key]?.data;
            return acc;
        }, {} as PartialNull<T>);
    });

    return {
        data,
        isAllNotNull,
        error,
        isFetching,
    };
}
