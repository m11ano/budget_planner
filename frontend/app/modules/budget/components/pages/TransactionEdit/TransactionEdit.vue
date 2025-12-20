<script setup lang="ts">
    import type { BudgetWidgetTransactionForm } from '#components';
    import { showErrors, showSuccess } from '~/core/components/shared/inform/toast';
    import { module } from '~/modules/budget/const';
    import { setModuleBreadcrums } from '~/modules/budget/domain/actions/setModuleBreadcrums';
    import { loadCategoriesList } from '~/modules/budget/domain/api/categories/fetchCategoriesList';
    import { fetchTransaction } from '~/modules/budget/domain/api/transaction/fetchTransaction';
    import { updateTransaction } from '~/modules/budget/domain/api/transaction/updateTransaction';
    import { checkTransactionState } from '~/modules/budget/domain/hooks/checkTransactionState';
    import type { ITransactionItem, ITransactionItemState } from '~/modules/budget/domain/model/types/transaction';
    import { setMenu } from '~/plugins/app/model/actions/setMenu';
    import { ApiError } from '~/shared/errors/errors';

    const props = defineProps<{
        id: string;
    }>();

    useSeoMeta({
        title: 'Редактирование объекта',
    });

    setMenu(module.urlName, 'transactions');

    setModuleBreadcrums([
        {
            name: 'Список транзакций',
            to: `/transactions`,
        },
        {
            name: 'Редактирование объекта',
        },
    ]);

    const form = ref<InstanceType<typeof BudgetWidgetTransactionForm> | null>(null);

    const itemState = ref<ITransactionItemState | null>(null);
    const itemObject = ref<ITransactionItem | null>(null);

    const isLoading = ref(false);

    const errors = ref<string[]>([]);

    const { data: categories, status: categoriesStatus } = loadCategoriesList();

    const isLoadingAnything = computed(() => isLoading.value || categoriesStatus.value === 'pending');

    const fetchItem = async (): Promise<ITransactionItem | null> => {
        isLoading.value = true;
        try {
            const data = await fetchTransaction(props.id);
            return data.item;
        } catch (e) {
            if (e instanceof ApiError) {
                if (e.code === 404) {
                    showError({
                        statusCode: e.code,
                        statusMessage: 'Объект не найден',
                    });
                } else {
                    showErrors(e.formHints());
                }
            }
        } finally {
            isLoading.value = false;
        }

        return null;
    };

    const updateItemState = (item: ITransactionItem) => {
        itemObject.value = item;

        const stateValue: ITransactionItemState = {
            amount: item.amount,
            categoryID: item.categoryID,
            isIncome: item.isIncome,
            occurredOn: item.occurredOn,
            description: item.description,
        };

        itemState.value = stateValue;
    };

    watch(
        () => props.id,
        async () => {
            const data = await fetchItem();
            if (data) {
                updateItemState(data);
            }
            if (form.value) {
                form.value.rebuild();
            }
        },
        {
            immediate: true,
        },
    );

    const save = async () => {
        if (isLoadingAnything.value || !itemState.value || !itemObject.value) return;

        if (form.value) {
            await form.value.syncAllData();
        }

        errors.value = checkTransactionState(itemState.value);

        if (errors.value.length) return;

        isLoading.value = true;
        try {
            await updateTransaction(itemObject.value.id, {
                ...itemState.value,
            });

            const data = await fetchItem();
            if (data) {
                updateItemState(data);
                if (form.value) {
                    form.value.rebuild();
                }
            }

            showSuccess();
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
        <div>
            <div class="form-table">
                <div>
                    <div class="title">ID:</div>
                    <div class="value">
                        {{ itemObject?.id }}
                    </div>
                </div>
            </div>
        </div>

        <div class="form_title mt-10">
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
                v-if="itemState && itemObject && categories"
                ref="form"
                v-model="itemState"
                v-model:data-item="itemObject"
                mode="edit"
                :disabled="isLoadingAnything"
                :categories="categories.items"
            />
        </div>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';
</style>
