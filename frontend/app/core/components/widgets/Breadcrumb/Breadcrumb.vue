<script setup lang="ts">
    import type { BreadcrumbItem } from '@nuxt/ui';
    import type { IAppData } from '~/plugins/app/model/types/types';

    const appData = useState<IAppData>('app');

    const items = computed(() => {
        const items: BreadcrumbItem[] = [
            {
                label: 'Главная',
                icon: 'i-lucide-house',
                to: '/',
            },
            ...appData.value.breadcrumbs.map((i) => ({
                label: i.name,
                icon: i.icon,
                to: i.to,
            })),
        ];

        return items;
    });
</script>

<template>
    <UBreadcrumb :items="items" />
</template>

<style lang="less" module>
    @import '@styles/includes';
</style>
