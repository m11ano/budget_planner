import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { ICategory } from '../../model/types/category';

export const fetchCategoriesList = async () => {
    try {
        return await useNuxtApp().$apiFetch<{ items: ICategory[] }>('v1/ledger/categories');
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};

export const loadCategoriesList = () => {
    return useApiFetch<{ items: ICategory[] }>('v1/ledger/categories');
};
