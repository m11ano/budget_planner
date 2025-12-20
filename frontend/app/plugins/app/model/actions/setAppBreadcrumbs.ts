import type { IAppData, IAppDataBreadcrumb } from '../types/types';

export function setAppBreadcrumbs(breadcrumbs: IAppDataBreadcrumb[]) {
    const appData = useState<IAppData>('app');
    appData.value.breadcrumbs = breadcrumbs;
}
