enum AppTheme {
    LIGHT = 'light',
    DARK = 'dark',
}

export default defineNuxtPlugin({
    name: 'theme',
    enforce: 'pre',
    async setup(nuxtApp) {
        const appTheme = useState<AppTheme>('appTheme', () => ref(AppTheme.LIGHT));

        if (import.meta.client) {
            //Подменяем тему если надо
        }

        return {
            provide: {
                appTheme: appTheme,
                toogleAppTheme: () => {
                    appTheme.value = appTheme.value === AppTheme.LIGHT ? AppTheme.DARK : AppTheme.LIGHT;
                },
            },
        };
    },
    env: {
        islands: true,
    },
});
