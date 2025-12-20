export interface IReportItemItem {
    categoryID: number;
    sum?: string;
    spentBudget?: string;
    itemBudget?: string;
}

export interface IReportListItem {
    periodStart: string;
    periodEnd: string;
    items: IReportItemItem[];
}
