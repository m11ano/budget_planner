import type { IBudgetItemState } from '../model/types/budget';

export function checkBudgetState(data: IBudgetItemState): string[] {
    const result: string[] = [];

    if (!data.period.year) {
        result.push('Год не указан');
    }

    if (data.amount === '') {
        result.push('Лимит не указан');
    }

    return result;
}
