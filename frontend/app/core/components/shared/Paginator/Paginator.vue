<script setup lang="ts">
const limit = defineModel<number>({ required: true });

defineProps<{ disabled?: boolean }>();
</script>

<template>
    <div :class="$style.wrapper">
        <div :class="$style.paginator">
            <slot />
        </div>
        <div :class="$style.limit">
            <div :class="$style.label">Лимит:</div>
            <UInputNumber
                v-model="limit"
                style="width: 100px"
                :min="1"
                :max="100"
                :step="1"
                :disabled="disabled"
                orientation="vertical"
            />
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

.wrapper {
    display: flex;
    align-items: center;
    border-top: 1px solid var(--ui-border);
    padding-top: 15px;

    .width-size-less(1000px, {
       flex-wrap: wrap;
       gap:20px;
       justify-content: center
    });

    > .limit {
        flex-shrink: 0;
        width: 160px;
        display: flex;
        align-items: center;
        justify-content: flex-end;
        gap: 10px;
        flex-wrap: wrap;
        color: var(--ui-text-muted);

        > .label {
            color: var(--ui-text-muted);
        }

        .width-size-less(1000px, {
            width: auto;
            justify-content: center;
            gap: 5px;

            > .label {
                width: 100%;
                text-align: center
            }
        });
    }

    > .paginator {
        width: 100%;
        padding-left: 160px;
        display: flex;
        justify-content: center;

        .width-size-less(1000px, {
            padding-left: 0;
        });
    }
}
</style>
