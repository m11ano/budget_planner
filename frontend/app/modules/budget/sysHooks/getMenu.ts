import type { IMenu } from '~/core/domain/model/types/menu';
import { module } from '../const';

export function getMenu(): IMenu[] {
    return [
        {
            name: module.title,
            icon: module.icon,
            menuSel: module.urlName,
            defaultOpen: true,
            kids: [
                {
                    name: 'Бюджеты',
                    to: `/${module.urlName}/budgets`,
                    subMenuSel: 'budgets',
                },
                {
                    name: 'Транзакции',
                    to: `/${module.urlName}/transactions`,
                    subMenuSel: 'transactions',
                },
                {
                    name: 'Отчеты',
                    to: `/${module.urlName}/reports`,
                    subMenuSel: 'reports',
                },
                {
                    name: 'Импорт',
                    to: `/${module.urlName}/import`,
                    subMenuSel: 'import',
                },
            ],
        },
    ];
}
