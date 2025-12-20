<script setup lang="ts">
    import { showSuccess } from '~/core/components/shared/inform/toast';
    import { module } from '~/modules/budget/const';
    import { setModuleBreadcrums } from '~/modules/budget/domain/actions/setModuleBreadcrums';
    import { setMenu } from '~/plugins/app/model/actions/setMenu';
    import { ApiError } from '~/shared/errors/errors';

    useSeoMeta({ title: 'Импорт CSV' });

    setMenu(module.urlName, '');
    setModuleBreadcrums([{ name: 'Импорт CSV' }]);

    const isLoading = ref(false);
    const errors = ref<string[]>([]);

    const file = ref<File | null>(null);

    const uploadCsv = async (csvFile: File) => {
        const form = new FormData();
        form.append('file', csvFile, csvFile.name);

        return await useNuxtApp().$apiFetch('v1/ledger/transactions/import', {
            method: 'POST',
            body: form,
        });
    };

    const isLoadingAnything = computed(() => isLoading.value);

    const onFileChange = (payload: Event | FileList | File[] | File | null) => {
        errors.value = [];

        let f: File | null = null;

        if (payload instanceof File) {
            f = payload;
        } else if (payload instanceof Event) {
            const input = payload.target as HTMLInputElement | null;
            f = input?.files?.[0] ?? null;
        } else if (payload && typeof payload === 'object' && 'length' in payload) {
            f = payload[0] ?? null;
        }

        if (!f) {
            file.value = null;
            return;
        }

        const nameOk = f.name.toLowerCase().endsWith('.csv');
        const typeOk = !f.type || f.type.includes('csv') || f.type.includes('comma-separated') || f.type === 'text/plain';

        if (!nameOk && !typeOk) {
            file.value = null;
            errors.value = ['Пожалуйста, выберите CSV файл.'];
            return;
        }

        file.value = f;
    };

    const save = async () => {
        if (isLoadingAnything.value) return;

        errors.value = [];

        if (!file.value) {
            errors.value = ['Выберите файл для загрузки.'];
            return;
        }

        isLoading.value = true;
        try {
            await uploadCsv(file.value);
            showSuccess();

            file.value = null;
        } catch (e: unknown) {
            if (e instanceof ApiError) {
                errors.value = e.formHints();
            } else {
                errors.value = ['Не удалось загрузить файл. Попробуйте ещё раз.'];
            }
        } finally {
            isLoading.value = false;
        }
    };
</script>

<template>
    <div>
        <div class="form_title">
            <div class="title">Импорт CSV</div>
            <div class="buttons">
                <UButton
                    :disabled="isLoadingAnything"
                    :loading="isLoadingAnything"
                    @click="save"
                >
                    Загрузить
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
            <div class="form-table">
                <div>
                    <div class="title">Файл:</div>
                    <div class="value">
                        <UInput
                            type="file"
                            accept=".csv,text/csv"
                            :disabled="isLoadingAnything"
                            @change="onFileChange"
                        />
                        <div
                            v-if="file"
                            class="mt-2 text-sm opacity-70"
                        >
                            Выбрано: <b>{{ file.name }}</b> ({{ Math.ceil(file.size / 1024) }} KB)
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';
</style>
