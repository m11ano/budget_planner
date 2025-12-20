import type { ITransactionItemState } from '../model/types/transaction';

export function checkTransactionState(data: ITransactionItemState): string[] {
    const result: string[] = [];

    if (!data.occurredOn) {
        result.push('Дата не указана');
    }

    if (data.amount === '') {
        result.push('Сумма не указана');
    }

    return result;
}
