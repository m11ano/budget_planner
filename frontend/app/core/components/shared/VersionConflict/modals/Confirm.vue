<script setup lang="ts">
defineProps<{
    text: string;
}>();

const emit = defineEmits<{ close: [boolean] }>();
</script>

<template>
    <UModal
        :title="`Уведомление`"
        @close="emit('close', false)"
    >
        <template #body>
            {{ text }}
        </template>
        <template #footer>
            <div class="flex justify-end gap-10">
                <UButton
                    label="OK"
                    color="primary"
                    @click="emit('close', true)"
                />
                <UButton
                    label="Отмена"
                    color="graylight"
                    @click="emit('close', false)"
                />
            </div>
        </template>
    </UModal>
</template>
