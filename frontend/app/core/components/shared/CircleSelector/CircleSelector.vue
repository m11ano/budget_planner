<script setup lang="ts">
    import type { ICicleMapItem } from './model/types/types';

    defineProps<{
        map: ICicleMapItem[];
    }>();

    const vModel = defineModel<number>();

    const colors: Record<string, string> = {
        '1': '#ff0000',
        '2': '#00ff00',
        '3': '#0000ff',
    };
</script>

<template>
    <div>
        <div :class="$style.wrapper">
            <template
                v-for="item in map"
                :key="item.value"
            >
                <button
                    type="button"
                    :class="vModel === item.value ? $style.active : false"
                    :title="item.title"
                    @click="vModel = item.value"
                >
                    <span :style="{ backgroundColor: item.color || colors[item.value.toString()] || undefined }" />
                </button>
            </template>
        </div>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    .wrapper {
        display: flex;
        gap: 10px;

        > button {
            display: block;
            padding: 3px;
            border: 3px solid transparent;
            border-radius: 100%;

            > span {
                display: block;
                width: 30px;
                height: 30px;
                border-radius: 100%;
            }

            &.active {
                border-color: #000;
            }
        }
    }
</style>
