// https://nuxt.com/docs/api/configuration/nuxt-config

import modules from './modules.config';

export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: { enabled: true },

    nitro: {
        routeRules: {
            '/bff-api/**': {
                proxy: '/api/**',
            },
        },
    },

    $production: {
        nitro: {
            compressPublicAssets: true,
        },
    },

    devServer: {
        host: '0.0.0.0',
        port: Number(process.env.NUXT_PORT) || 3000,
    },

    ssr: false,

    alias: {
        '@styles': '/assets/styles',
    },

    dir: {
        pages: 'routes',
    },

    modules: ['@nuxt/eslint', '@pinia/nuxt', '@vueuse/nuxt', '@nuxt/ui'],

    pinia: {
        storesDirs: ['shared/stores/**'],
    },

    imports: {
        dirs: ['shared/hooks/**/*.ts'],
    },

    plugins: ['~/plugins/api-fetch.plugin.ts', '~/plugins/auth/auth.plugin.ts', '~/plugins/theme/theme.plugin.ts', '~/plugins/app/app.plugin.ts'],

    components: [
        {
            path: '~/core/components/shared',
            extensions: ['.vue'],
            prefix: 'Shared',
        },
        {
            path: '~/core/components/features',
            extensions: ['.vue'],
            prefix: 'Feature',
        },
        {
            path: '~/core/components/entities',
            extensions: ['.vue'],
            prefix: 'Entity',
        },
        {
            path: '~/core/components/widgets',
            extensions: ['.vue'],
            prefix: 'Widget',
        },
        {
            path: '~/core/components/pages',
            extensions: ['.vue'],
            prefix: 'Page',
        },
        ...modules.flatMap((m) => [
            {
                path: `~/modules/${m.name}/components/shared`,
                extensions: ['.vue'],
                prefix: `${m.prefixName}Shared`,
            },
            {
                path: `~/modules/${m.name}/components/features`,
                extensions: ['.vue'],
                prefix: `${m.prefixName}Feature`,
            },
            {
                path: `~/modules/${m.name}/components/entities`,
                extensions: ['.vue'],
                prefix: `${m.prefixName}Entity`,
            },
            {
                path: `~/modules/${m.name}/components/widgets`,
                extensions: ['.vue'],
                prefix: `${m.prefixName}Widget`,
            },
            {
                path: `~/modules/${m.name}/components/pages`,
                extensions: ['.vue'],
                prefix: `${m.prefixName}Page`,
            },
        ]),
    ],

    css: ['@/assets/styles/index.less', '@/assets/styles/nuxt-ui.css'],

    app: {
        baseURL: '/',
        head: {
            viewport: 'width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no',
        },
    },

    runtimeConfig: {
        public: {
            isProd: process.env.NODE_ENV === 'production',
            apiBase: 'http://localhost:8080/api',
            apiServerPrefix: '',
        },
    },

    ui: {
        colorMode: false,
        theme: {
            colors: ['primary', 'secondary', 'graylight', 'info', 'success', 'warning', 'error'],
        },
    },

    experimental: {
        entryImportMap: false,
    },
});
