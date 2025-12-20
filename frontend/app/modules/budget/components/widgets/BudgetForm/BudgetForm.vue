<script setup lang="ts">
    import { MONTHES } from '~/modules/budget/domain/model/const/const';
    import type { IBudgetItem, IBudgetItemState } from '~/modules/budget/domain/model/types/budget';
    import type { ICategory } from '~/modules/budget/domain/model/types/category';

    const props = defineProps<{
        disabled?: boolean;
        mode: 'new' | 'edit';
        categories: ICategory[];
    }>();

    const dataModel = defineModel<IBudgetItemState>({ required: true });

    const dataItem = defineModel<IBudgetItem | null>('dataItem', { required: true });

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

    const monthsListOptions = computed(() => {
        return Object.entries(MONTHES).map(([id, value]) => ({
            value: Number(id),
            label: value,
        }));
    });

    const amount = ref(Number(dataModel.value.amount));
    if (isNaN(amount.value)) {
        amount.value = 0;
    }

    watch(
        amount,
        (value) => {
            dataModel.value.amount = value.toString();
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
            <div class="title">Месяц:</div>
            <div class="value">
                <USelect
                    v-model="dataModel.period.month"
                    :items="monthsListOptions"
                    size="xl"
                    class="w-full"
                    :disabled="disabled"
                />
            </div>
        </div>
        <div>
            <div class="title">Год:</div>
            <div class="value">
                <UInputNumber
                    v-model="dataModel.period.year"
                    :min="2000"
                    :max="2100"
                    :step="1"
                    size="xl"
                    class="w-full"
                    :disabled="disabled"
                />
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
            <div class="title">Лимит:</div>
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
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    .linksBlockTitle {
        font-size: 20px;
        margin-bottom: 20px;
    }
</style>
