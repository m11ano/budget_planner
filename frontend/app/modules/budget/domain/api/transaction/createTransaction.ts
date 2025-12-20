import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { ITransactionItem, ITransactionItemState } from '../../model/types/transaction';

type Input = ITransactionItemState;

type Payload = Input;

const mapDataToRequest = (data: Input): Payload => {
    return data;
};

export async function createTransaction(input: Input) {
    try {
        return await useNuxtApp().$apiFetch<{ item: ITransactionItem }>(`v1/ledger/transactions`, {
            method: 'POST',
            body: mapDataToRequest(input),
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
