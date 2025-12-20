import type { IAppData } from '../types/types';

export function setMenu(value: string, subMenu?: string) {
    const appData = useState<IAppData>('app');
    appData.value.menuSel = value || '';
    appData.value.subMenuSel = subMenu || '';
}
