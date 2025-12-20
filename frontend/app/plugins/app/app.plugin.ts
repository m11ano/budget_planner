import type { IAppData } from './model/types/types';

export default defineNuxtPlugin({
    name: 'appProvider',
    enforce: 'pre',
    async setup() {
        const appProvider = useState<IAppData>('app', () =>
            reactive({
                title: 'Budget Planner: панель управления',
                breadcrumbs: [],
                menuSel: '',
                subMenuSel: '',
            }),
        );

        return {
            provide: {
                appProvider: appProvider.value,
            },
        };
    },
    env: {
        islands: true,
    },
});
