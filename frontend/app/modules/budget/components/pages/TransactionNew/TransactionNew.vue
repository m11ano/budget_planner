<script setup lang="ts">
    import type { BudgetWidgetTransactionForm } from '#components';
    import { showSuccess } from '~/core/components/shared/inform/toast';
    import { module } from '~/modules/budget/const';
    import { setModuleBreadcrums } from '~/modules/budget/domain/actions/setModuleBreadcrums';
    import { loadCategoriesList } from '~/modules/budget/domain/api/categories/fetchCategoriesList';
    import { createTransaction } from '~/modules/budget/domain/api/transaction/createTransaction';
    import { checkTransactionState } from '~/modules/budget/domain/hooks/checkTransactionState';
    import type { ITransactionItemState } from '~/modules/budget/domain/model/types/transaction';
    import { setMenu } from '~/plugins/app/model/actions/setMenu';
    import { ApiError } from '~/shared/errors/errors';

    useSeoMeta({
        title: 'Создание объекта',
    });

    setMenu(module.urlName, 'transactions');

    setModuleBreadcrums([
        {
            name: 'Список транзакций',
            to: `/transactions`,
        },
        {
            name: 'Создание объекта',
        },
    ]);

    const form = ref<InstanceType<typeof BudgetWidgetTransactionForm> | null>(null);

    const initState: ITransactionItemState = {
        amount: '',
        categoryID: 1,
        isIncome: false,
        occurredOn: String(new Date().toISOString().split('T')[0]),
        description: '',
    };

    const { data: categories, status: categoriesStatus } = loadCategoriesList();

    const itemObject = ref<null>(null);

    const itemState = ref<ITransactionItemState>(initState);

    const isLoading = ref(false);

    const errors = ref<string[]>([]);

    const isLoadingAnything = computed(() => isLoading.value || categoriesStatus.value === 'pending');

    const save = async () => {
        if (isLoadingAnything.value || !itemState.value) return;

        if (form.value) {
            await form.value.syncAllData();
        }

        errors.value = checkTransactionState(itemState.value);

        if (errors.value.length) return;

        isLoading.value = true;
        try {
            const data = await createTransaction(itemState.value);

            showSuccess();

            await navigateTo(`/${module.urlName}/transactions/${data.item.id}`);
        } catch (e) {
            if (e instanceof ApiError) {
                errors.value = e.formHints();
            }
        } finally {
            isLoading.value = false;
        }
    };
</script>

<template>
    <div>
        <div class="form_title">
            <div class="title">Данные</div>
            <div class="buttons">
                <UButton
                    :disabled="isLoadingAnything"
                    :loading="isLoadingAnything"
                    @click="save"
                >
                    Сохранить
                </UButton>
            </div>
        </div>
        <div
            v-if="errors.length"
            class="mt-4"
        >
            <UAlert
                title="Возникли ошибки!"
                icon="i-lucide-ban"
            >
                <template #description>
                    <template
                        v-for="error in errors"
                        :key="error"
                    >
                        <div>– {{ error }}</div>
                    </template>
                </template>
            </UAlert>
        </div>
        <div class="mt-4">
            <BudgetWidgetTransactionForm
                v-if="itemState && categories"
                ref="form"
                v-model="itemState"
                v-model:data-item="itemObject"
                mode="new"
                :disabled="isLoadingAnything"
                :categories="categories.items"
            />
        </div>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';
</style>
