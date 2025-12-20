<script setup lang="ts">
    import type { BudgetWidgetBudgetForm } from '#components';
    import { showSuccess } from '~/core/components/shared/inform/toast';
    import { module } from '~/modules/budget/const';
    import { setModuleBreadcrums } from '~/modules/budget/domain/actions/setModuleBreadcrums';
    import { createBudget } from '~/modules/budget/domain/api/bugdet/createBudget';
    import { loadCategoriesList } from '~/modules/budget/domain/api/categories/fetchCategoriesList';
    import { checkBudgetState } from '~/modules/budget/domain/hooks/checkBudgetState';
    import type { IBudgetItemState } from '~/modules/budget/domain/model/types/budget';
    import { setMenu } from '~/plugins/app/model/actions/setMenu';
    import { ApiError } from '~/shared/errors/errors';

    useSeoMeta({
        title: 'Создание объекта',
    });

    setMenu(module.urlName, 'budgets');

    setModuleBreadcrums([
        {
            name: 'Список бюджетов',
            to: `/budgets`,
        },
        {
            name: 'Создание объекта',
        },
    ]);

    const form = ref<InstanceType<typeof BudgetWidgetBudgetForm> | null>(null);

    const initState: IBudgetItemState = {
        amount: '',
        period: {
            month: new Date().getMonth() + 1,
            year: new Date().getFullYear(),
        },
        categoryID: 1,
    };

    const { data: categories, status: categoriesStatus } = loadCategoriesList();

    const itemObject = ref<null>(null);

    const itemState = ref<IBudgetItemState>(initState);

    const isLoading = ref(false);

    const errors = ref<string[]>([]);

    const isLoadingAnything = computed(() => isLoading.value || categoriesStatus.value === 'pending');

    const save = async () => {
        if (isLoadingAnything.value || !itemState.value) return;

        if (form.value) {
            await form.value.syncAllData();
        }

        errors.value = checkBudgetState(itemState.value);

        if (errors.value.length) return;

        isLoading.value = true;
        try {
            const data = await createBudget(itemState.value);

            showSuccess();

            await navigateTo(`/${module.urlName}/budgets/${data.item.id}`);
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
            <BudgetWidgetBudgetForm
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
