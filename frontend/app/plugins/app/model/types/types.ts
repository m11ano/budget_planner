export interface IAppDataBreadcrumb {
    name: string;
    to?: string;
    icon?: string;
}

export interface IAppData {
    title: string;
    breadcrumbs: IAppDataBreadcrumb[];
    menuSel: string;
    subMenuSel: string;
}
