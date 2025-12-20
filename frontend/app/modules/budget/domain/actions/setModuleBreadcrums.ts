import { setAppBreadcrumbs } from '~/plugins/app/model/actions/setAppBreadcrumbs';
import type { IAppDataBreadcrumb } from '~/plugins/app/model/types/types';
import { module } from '../../const';

export function setModuleBreadcrums(items: IAppDataBreadcrumb[]) {
    setAppBreadcrumbs([
        {
            name: module.title,
            to: `/${module.urlName}`,
            icon: module.icon,
        },
        ...items.map((item) => ({
            name: item.name,
            to: item.to ? `/${module.urlName}${item.to}` : undefined,
            icon: item.icon,
        })),
    ]);
}
