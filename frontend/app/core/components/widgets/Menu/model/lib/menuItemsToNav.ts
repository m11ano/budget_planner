import type { NavigationMenuItem } from '@nuxt/ui';
import type { IMenu } from '~/core/domain/model/types/menu';

export function menuItemsToNav(items: IMenu[], menuSel: string, subMenuSel: string): NavigationMenuItem[] {
    const result: NavigationMenuItem[] = [];

    items.forEach((item) => {
        result.push({
            label: item.name,
            icon: item.icon,
            defaultOpen: item.defaultOpen,
            active: menuSel === item.menuSel,
            to: item.to,
            children: item.kids?.map((subItem) => {
                return {
                    label: subItem.name,
                    active: subMenuSel === subItem.subMenuSel,
                    to: subItem.to,
                };
            }),
        });
    });

    return result;
}
