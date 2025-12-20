export interface IBudgetItem {
    id: string;
    accountID: string;
    amount: string;
    period: {
        month: number;
        year: number;
    };
    categoryID: number;
    createdAt: string;
    updatedAt: string;
}

export type IBudgetListItem = IBudgetItem;

export interface IBudgetItemState {
    amount: string;
    period: {
        month: number;
        year: number;
    };
    categoryID: number;
}
