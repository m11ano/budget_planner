<script setup lang="ts">
    import type { NavigationMenuItem } from '@nuxt/ui';
    import { getMenu as budgetGetMenu } from '~/modules/budget/sysHooks/getMenu';
    import type { IAppData } from '~/plugins/app/model/types/types';
    import { menuItemsToNav } from './model/lib/menuItemsToNav';

    const appData = useState<IAppData>('app');

    const items = computed<NavigationMenuItem[][]>(() => {
        const menu: NavigationMenuItem[][] = [
            [
                {
                    label: 'Меню',
                    type: 'label',
                },
            ],
        ];

        menu.push(menuItemsToNav(budgetGetMenu(), appData.value.menuSel, appData.value.subMenuSel));

        return menu;
    });
</script>

<template>
    <div>
        <UNavigationMenu
            orientation="vertical"
            :items="items"
        />
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    .wrapper {
        display: flex;
        flex-direction: column;
        gap: 30px;
    }

    .parent {
        > a {
            width: 100%;
        }
    }

    .kids {
        margin-top: 15px;
        display: flex;
        flex-direction: column;
        gap: 10px;

        > div {
            > a {
                text-decoration: underline;
            }

            &.active {
                > a {
                    font-weight: bold;
                }
            }
        }
    }
</style>
