export interface IMenu {
    name: string;
    icon: string;
    to?: string;
    menuSel: string;
    kids?: ISubMenu[];
    defaultOpen?: boolean;
}

export interface ISubMenu {
    name: string;
    to: string;
    subMenuSel: string;
}
