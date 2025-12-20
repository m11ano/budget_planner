import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IReportListItem } from '../../model/types/report';

export const fetchReportsList = async (dateFrom?: string, dateTo?: string) => {
    try {
        return await useNuxtApp().$apiFetch<{ reports: IReportListItem[] }>('v1/ledger/reports', {
            params: {
                date_from: dateFrom,
                date_to: dateTo,
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
