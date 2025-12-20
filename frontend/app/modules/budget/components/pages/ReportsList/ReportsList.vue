<script setup lang="ts">
    import type { TableColumn } from '@nuxt/ui';
    import { showErrors } from '~/core/components/shared/inform/toast';
    import { module } from '~/modules/budget/const';
    import { setModuleBreadcrums } from '~/modules/budget/domain/actions/setModuleBreadcrums';
    import { loadCategoriesList } from '~/modules/budget/domain/api/categories/fetchCategoriesList';
    import { fetchReportsList } from '~/modules/budget/domain/api/report/fetchReportList';
    import type { IReportListItem } from '~/modules/budget/domain/model/types/report';
    import { setMenu } from '~/plugins/app/model/actions/setMenu';
    import { ApiError } from '~/shared/errors/errors';
    import { coolNumber } from '~/shared/helpers/functions';

    useSeoMeta({
        title: 'Список отчетов',
    });

    setMenu(module.urlName, 'reports');

    setModuleBreadcrums([
        {
            name: 'Список отчетов',
        },
    ]);

    const { data: categories, status: categoriesStatus } = loadCategoriesList();

    const getCategory = (id: number) => {
        return categories.value?.items.find((c) => c.id === id);
    };

    const dateFromFilter = ref('');
    const dateToFilter = ref('');

    const isDataLoading = ref(true);

    const isAnyLoading = computed(() => isDataLoading.value || categoriesStatus.value === 'pending');

    const list = ref<IReportListItem[]>([]);

    const fetchData = async () => {
        isDataLoading.value = true;

        try {
            const data = await fetchReportsList(dateFromFilter.value, dateToFilter.value);
            if (data.reports) {
                list.value = data.reports;
            }
        } catch (e: unknown) {
            if (e instanceof ApiError) {
                showErrors(e.formHints());
            }
        } finally {
            isDataLoading.value = false;
        }
    };

    const columns: TableColumn<IReportListItem>[] = [
        {
            id: 'period',
            header: 'Период',
        },
        {
            id: 'info',
            header: 'Информация',
        },
        {
            id: 'action',
        },
    ];

    watch([dateFromFilter, dateToFilter], () => {
        fetchData();
    });

    onMounted(() => {
        fetchData();
    });

    const download = async () => {
        const nuxtApp = useNuxtApp();

        const blob = await nuxtApp.$apiFetch<Blob>('v1/ledger/transactions/export', {
            method: 'GET',
            params: {
                date_from: dateFromFilter.value,
                date_to: dateToFilter.value,
            },
            responseType: 'blob',
        });

        const from = dateFromFilter.value || 'all';
        const to = dateToFilter.value || 'all';
        const filename = `transactions_${from}_${to}.csv`;

        const url = URL.createObjectURL(blob.type ? blob : new Blob([blob], { type: 'text/csv;charset=utf-8' }));

        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        a.style.display = 'none';
        document.body.appendChild(a);
        a.click();
        a.remove();

        URL.revokeObjectURL(url);
    };
</script>

<template>
    <div>
        <div class="flex justify-end">
            <UButton @click="download">Скачать CSV</UButton>
        </div>
        <div class="mt-4 flex gap-4 items-center flex-wrap">
            <div>
                Дата, от:
                <UInput
                    v-model="dateFromFilter"
                    type="date"
                    size="xl"
                />
            </div>
            <div>
                Дата, до:
                <UInput
                    v-model="dateToFilter"
                    type="date"
                    size="xl"
                />
            </div>
        </div>
        <UTable
            :data="list"
            :columns="columns"
            :loading="isAnyLoading"
            :ui="{ td: 'whitespace-normal' }"
        >
            <template #period-cell="{ row }"> {{ row.original.periodStart }} — {{ row.original.periodEnd }} </template>
            <template #info-cell="{ row }">
                <div class="flex gap-2 flex-col">
                    <template
                        v-for="(item, i) in row.original.items"
                        :key="i"
                    >
                        <div>
                            <div>Категория: {{ getCategory(item.categoryID)?.title }}</div>
                            <div>
                                Баланс:
                                <template v-if="item.sum"
                                    ><b>{{ coolNumber(Number(item.sum)) }} руб.</b></template
                                ><template v-else><i>Нет данных</i></template>
                            </div>
                            <div v-if="item.itemBudget">
                                Бюджет: <b>{{ coolNumber(Number(item.itemBudget)) }} руб.</b>
                            </div>
                            <div v-if="item.spentBudget">
                                Освоено бюджета: <b>{{ coolNumber(Number(item.spentBudget)) }} %</b>
                            </div>
                        </div>
                    </template>
                </div>
            </template>
        </UTable>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';
</style>
