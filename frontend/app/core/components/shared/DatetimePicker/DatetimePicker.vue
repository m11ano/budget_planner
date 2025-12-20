<script setup lang="ts">
    import type theme from '#build/ui/input';
    import type { AppConfig } from '@nuxt/schema';
    import type { ComponentConfig } from '@nuxt/ui';
    import { useVModel } from '@vueuse/core';

    type Input = ComponentConfig<typeof theme, AppConfig, 'input'>;

    type Model = string | null;

    const props = defineProps<{
        modelValue: Model;
        time?: boolean;
        seconds?: boolean;
        label?: string;
        disabled?: boolean;
        placeholder?: string;
        size?: Input['variants']['size'];
    }>();

    const emit = defineEmits<{
        'update:modelValue': [Model];
    }>();

    const model = useVModel(props, 'modelValue', emit);

    function rfc3339ToLocalInput(v: Model): string {
        if (!v) return '';
        const d = new Date(v);
        if (Number.isNaN(d.getTime())) return '';

        const Y = d.getFullYear();
        const M = String(d.getMonth() + 1).padStart(2, '0');
        const D = String(d.getDate()).padStart(2, '0');
        const h = String(d.getHours()).padStart(2, '0');
        const m = String(d.getMinutes()).padStart(2, '0');
        const s = String(d.getSeconds()).padStart(2, '0');
        if (!props.time) return `${Y}-${M}-${D}`;

        return props.seconds ? `${Y}-${M}-${D}T${h}:${m}:${s}` : `${Y}-${M}-${D}T${h}:${m}`;
    }

    function localInputToRFC3339(v: string): Model {
        if (!v) return null;
        const d = new Date(v);
        if (Number.isNaN(d.getTime())) return null;
        return d.toISOString();
    }

    const inputValue = computed<string>({
        get() {
            return rfc3339ToLocalInput(model.value);
        },
        set(v: string) {
            const next = localInputToRFC3339(v);
            if (next !== model.value) model.value = next;
        },
    });
</script>

<template>
    <UInput
        v-model="inputValue"
        :type="time ? 'datetime-local' : 'date'"
        :placeholder="placeholder"
        :disabled="disabled"
        :step="seconds ? 1 : undefined"
        :size="size"
    />
</template>
