<script setup lang="ts">
    import type { ICategory } from '~/modules/budget/domain/model/types/category';
    import type { ITransactionItem, ITransactionItemState } from '~/modules/budget/domain/model/types/transaction';

    const props = defineProps<{
        disabled?: boolean;
        mode: 'new' | 'edit';
        categories: ICategory[];
    }>();

    const dataModel = defineModel<ITransactionItemState>({ required: true });

    const dataItem = defineModel<ITransactionItem | null>('dataItem', { required: true });

    const show = ref(true);

    const rebuild = async () => {
        await syncAllData();

        show.value = false;
        await nextTick();
        show.value = true;
    };

    const syncAllData = async () => {};

    defineExpose({
        syncAllData,
        rebuild,
    });

    const categoriesListOptions = computed(() => {
        return props.categories.map((value) => ({
            value: value.id,
            label: value.title,
        }));
    });

    const opType = ref<'outcome' | 'income'>(dataModel.value.isIncome ? 'income' : 'outcome');

    watch(
        opType,
        (value) => {
            dataModel.value.isIncome = opType.value === 'income';
        },
        { immediate: true },
    );

    const typesListOptions = computed(() => {
        return [
            { value: 'outcome', label: 'Расход' },
            { value: 'income', label: 'Доход' },
        ];
    });

    const amount = ref<number>(Number(dataModel.value.amount));
    if (isNaN(amount.value)) {
        amount.value = 0;
    }

    if (!dataModel.value.isIncome && amount.value < 0) {
        amount.value = Math.abs(amount.value);
    }

    watch(
        [amount, () => dataModel.value.isIncome],
        () => {
            dataModel.value.amount = dataModel.value.isIncome ? amount.value.toString() : `-${amount.value.toString()}`;
        },
        { immediate: true },
    );
</script>

<template>
    <div
        v-if="show"
        class="form-table"
    >
        <div>
            <div class="title">Тип:</div>
            <div class="value">
                <template v-if="mode === 'edit'">
                    {{ dataModel.isIncome ? 'Доход' : 'Расход' }}
                </template>
                <template v-else>
                <USelect
                    v-model="opType"
                    :items="typesListOptions"
                    size="xl"
                    class="w-full"
                    :disabled="disabled"
                />
                </template>
            </div>
        </div>

        <div>
            <div class="title">Категория:</div>
            <div class="value">
                <USelect
                    v-model="dataModel.categoryID"
                    :items="categoriesListOptions"
                    size="xl"
                    class="w-full"
                    :disabled="disabled"
                />
            </div>
        </div>
        <div>
            <div class="title">Дата:</div>
            <div class="value">
                <UInput
                    v-model="dataModel.occurredOn"
                    type="date"
                    size="xl"
                    class="w-full"
                    :disabled="disabled"
                />
            </div>
        </div>
        <div>
            <div class="title">Сумма:</div>
            <div class="value">
                <UInputNumber
                    v-model="amount"
                    :min="0"
                    :step="0.01"
                    size="xl"
                    class="w-full"
                    :disabled="disabled"
                />
            </div>
        </div>
        <div>
            <div class="title">Описание:</div>
            <div class="value">
                <UInput
                    v-model="dataModel.description"
                    size="xl"
                    class="w-full"
                    :disabled="disabled"
                />
            </div>
        </div>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    .linksBlockTitle {
        font-size: 20px;
        margin-bottom: 20px;
    }
</style>
