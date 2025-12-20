<script setup lang="ts">
defineProps<{
    name: string;
    errors: string[];
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <UModal
        :title="`Возникли ошибки при загрузке ${name}`"
        @close="emit('close')"
    >
        <template #body>
            <template
                v-for="error in errors"
                :key="error"
            >
                <div>– {{ error }}</div>
            </template>
        </template>
        <template #footer>
            <div class="flex justify-end">
                <UButton
                    label="OK"
                    color="primary"
                    @click="emit('close')"
                />
            </div>
        </template>
    </UModal>
</template>
