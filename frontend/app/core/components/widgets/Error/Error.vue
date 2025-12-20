<script setup lang="ts">
    defineProps<{ code?: number }>();

    const errorDesc = ref<Record<string, string>>({
        '403': 'Доступ запрещен',
    });
</script>

<template>
    <div>
        <UAlert icon="i-lucide-terminal">
            <template #title><b>Ошибка!</b></template>
            <template #description>
                <template v-if="code && errorDesc[code.toString()] !== undefined">
                    {{ errorDesc[code.toString()] }}
                </template>
                <template v-else>
                    <slot />
                </template>
            </template>
        </UAlert>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';
</style>
