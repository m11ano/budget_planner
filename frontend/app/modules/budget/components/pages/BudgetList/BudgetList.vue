<script setup lang="ts">
    import type { DropdownMenuItem, TableColumn } from '@nuxt/ui';
    import Confirm from '~/core/components/shared/Confirm/modals/Confirm.vue';
    import { showErrors, showSuccess } from '~/core/components/shared/inform/toast';
    import { module } from '~/modules/budget/const';
    import { setModuleBreadcrums } from '~/modules/budget/domain/actions/setModuleBreadcrums';
    import { deleteBudget } from '~/modules/budget/domain/api/bugdet/deleteBudget';
    import { fetchBudgetsList } from '~/modules/budget/domain/api/bugdet/fetchBudgetList';
    import { loadCategoriesList } from '~/modules/budget/domain/api/categories/fetchCategoriesList';
    import { MONTHES } from '~/modules/budget/domain/model/const/const';
    import type { IBudgetListItem } from '~/modules/budget/domain/model/types/budget';
    import { setMenu } from '~/plugins/app/model/actions/setMenu';
    import { ApiError } from '~/shared/errors/errors';
    import { coolNumber } from '~/shared/helpers/functions';

    useSeoMeta({
        title: 'Список бюджетов',
    });

    setMenu(module.urlName, 'budgets');

    setModuleBreadcrums([]);

    const route = useRoute();

    let routePage = Number(route.query.page);
    if (!routePage || isNaN(routePage)) {
        routePage = 1;
    }

    const page = ref(routePage);

    const { data: categories, status: categoriesStatus } = loadCategoriesList();

    const getCategory = (id: number) => {
        return categories.value?.items.find((c) => c.id === id);
    };

    const isDataLoading = ref(true);

    const isAnyLoading = computed(() => isDataLoading.value || categoriesStatus.value === 'pending');

    const list = ref<IBudgetListItem[]>([]);

    const defaultLimit = 20;

    let routeLimit = Number(route.query.limit);
    if (!routeLimit || isNaN(routeLimit)) {
        routeLimit = defaultLimit;
    }

    const limit = ref(routeLimit);

    const total = ref(0);

    const fetchData = async () => {
        isDataLoading.value = true;

        try {
            const data = await fetchBudgetsList(page.value, limit.value < 1 ? 1 : limit.value);
            if (data.items) {
                list.value = data.items;
                total.value = data.total;

                await navigateTo({
                    query: {
                        ...route.query,
                        page: page.value > 1 ? page.value : undefined,
                        limit: limit.value !== defaultLimit ? limit.value : undefined,
                    },
                });
            }
        } catch (e: unknown) {
            if (e instanceof ApiError) {
                showErrors(e.formHints());
            }
        } finally {
            isDataLoading.value = false;
        }
    };

    const onPageUpdate = (p: number) => {
        if (isDataLoading.value) return;
        page.value = p;
        fetchData();
    };

    watch(limit, () => {
        page.value = 1;
        setTimeout(() => {
            fetchData();
        }, 100);
    });

    const removeBudget = async (id: string): Promise<boolean> => {
        const modal = useOverlay().create(Confirm, {
            props: {
                text: 'Вы действительно хотите удалить объект?',
            },
            destroyOnClose: true,
        });

        const instance = modal.open();

        const shouldDelete = await instance.result;
        if (shouldDelete) {
            try {
                await deleteBudget(id);

                showSuccess('Объект удален');

                return true;
            } catch (e) {
                if (e instanceof ApiError) {
                    showErrors(e.formHints());
                }
            }
        }

        return false;
    };

    const columns: TableColumn<IBudgetListItem>[] = [
        {
            id: 'id',
            header: 'ID',
        },
        {
            id: 'info',
            header: 'Информация',
        },
        {
            id: 'action',
        },
    ];

    function getDropdownActions(item: IBudgetListItem): DropdownMenuItem[][] {
        return [
            [
                {
                    label: 'Редактировать',
                    icon: 'i-lucide-edit',
                    to: `/${module.urlName}/budgets/${item.id}`,
                },
                {
                    label: 'Удалить',
                    icon: 'i-lucide-trash',
                    color: 'error',
                    onSelect: async () => {
                        const result = await removeBudget(item.id);
                        if (result) {
                            list.value = list.value.filter((p) => p.id !== item.id);
                            total.value -= 1;
                        }
                    },
                },
            ],
        ];
    }

    const columnPinning = ref({ left: [], right: ['action'] });

    onMounted(() => {
        fetchData();
    });
</script>

<template>
    <div>
        <div class="flex justify-end">
            <UButton :to="`/${module.urlName}/budgets/new`">Создать бюджет</UButton>
        </div>
        <UTable
            v-model:column-pinning="columnPinning"
            :data="list"
            :columns="columns"
            :loading="isAnyLoading"
            :ui="{ td: 'whitespace-normal' }"
        >
            <template #id-cell="{ row }">
                <div style="font-size: 10px">{{ row.original.id }}</div>
            </template>
            <template #info-cell="{ row }">
                <div class="flex gap-2 flex-col">
                    <div>
                        <b>{{ MONTHES[row.original.period.month.toString()!] }} {{ row.original.period.year }}</b>
                    </div>
                    <div>Категория: {{ getCategory(row.original.categoryID)?.title }}</div>
                    <div>
                        Лимит: <b>{{ coolNumber(Number(row.original.amount)) }} руб.</b>
                    </div>
                </div>
            </template>
            <template #action-cell="{ row }">
                <div class="flex justify-end">
                    <UDropdownMenu :items="getDropdownActions(row.original)">
                        <UButton
                            icon="i-lucide-ellipsis-vertical"
                            color="neutral"
                            variant="ghost"
                            aria-label="Действия"
                        />
                    </UDropdownMenu>
                </div>
            </template>
        </UTable>
        <SharedPaginator
            v-model="limit"
            :disabled="isDataLoading"
        >
            <UPagination
                :page="page"
                :items-per-page="limit"
                :total="total"
                @update:page="onPageUpdate"
            />
        </SharedPaginator>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';
</style>
