export interface ITransactionItem {
    id: string;
    accountID: string;
    isIncome: boolean;
    amount: string;
    occurredOn: string;
    categoryID: number;
    description: string;
    createdAt: string;
    updatedAt: string;
}

export type ITransactionListItem = ITransactionItem;

export interface ITransactionItemState {
    isIncome: boolean;
    amount: string;
    occurredOn: string;
    categoryID: number;
    description: string;
}
