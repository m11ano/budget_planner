<script setup lang="ts">
    import { ApiError } from '~/shared/errors/errors';
    import type { IFile } from '~/shared/types/files';
    import { fetchUpload } from './api/fetchUpload';
    import Error from './modals/Error.vue';
    import type { UploadedFile } from './model/types/types';

    const props = defineProps<{
        mode: 'solo' | 'multi';
        leadTarget?: string;
        uploadUrl: string;
        acceptTypes?: string;
        modelValue?: UploadedFile[];
        isNoneRemovable?: boolean;
        disabled?: boolean;
    }>();

    const leadFile = (file: UploadedFile): IFile => {
        if (props.leadTarget && file.targets[props.leadTarget] !== undefined) {
            return file.targets[props.leadTarget]!;
        }

        return Object.values(file.targets)[0]!;
    };

    const emit = defineEmits<{
        (e: 'update:modelValue', payload: UploadedFile[]): void;
    }>();

    const { mode, uploadUrl, acceptTypes = '', modelValue = [] } = props;

    const isLoading = ref(false);

    const files = ref<UploadedFile[]>([...modelValue]);

    const filesWithLeads = computed(() => {
        return files.value.map((file) => {
            return {
                ...file,
                lead: leadFile(file),
            };
        });
    });

    watch(
        () => modelValue,
        (newVal) => {
            files.value = newVal ? [...newVal] : [];
        },
        { deep: true },
    );

    const updateModelValue = () => {
        emit('update:modelValue', [...files.value]);
    };

    const isImage = (file: UploadedFile) => {
        return /\.(jpe?g|png|gif|webp)$/i.test(leadFile(file).url);
    };

    const getFileExtension = (url: string): string => {
        const clean = url.split('?')[0]!.split('#')[0]!;
        const lastDot = clean.lastIndexOf('.');
        if (lastDot === -1) return '';
        const lastSlash = clean.lastIndexOf('/');
        if (lastDot < lastSlash) return '';
        return clean.slice(lastDot + 1);
    };

    const overlay = useOverlay();
    const showErrors = (name: string, errors: string[]) => {
        const modal = overlay.create(Error, {
            props: {
                errors,
                name,
            },
            destroyOnClose: true,
        });
        modal.open();
    };

    const handleFiles = async (event: Event) => {
        const input = event.target as HTMLInputElement;
        if (!input.files || input.files.length === 0) return;
        if (isLoading.value) {
            input.value = '';
            return;
        }

        const selectedFiles = Array.from(input.files);
        isLoading.value = true;

        const errors: string[] = [];
        const limit = useNuxtApp().$config.public.maxFilesize;
        for (const file of selectedFiles) {
            if (file.size > limit) {
                errors.push(`Файл ${file.name} слишком большой, лимит ${limit / 1024 / 1024} МБ`);
            }
        }

        if (errors.length > 0) {
            showErrors('', errors);
            isLoading.value = false;
            input.value = '';
            return;
        }

        for (let i = 0; i < selectedFiles.length; i++) {
            const formData = new FormData();
            formData.append('file', selectedFiles[i]!);

            try {
                const uploaded = await fetchUpload(uploadUrl, selectedFiles[i]!);

                if (mode === 'solo') {
                    files.value = [];
                }

                const item: UploadedFile = {
                    targets: uploaded.files,
                };

                files.value.push(item);
                updateModelValue();
            } catch (e) {
                if (e instanceof ApiError) {
                    showErrors(selectedFiles[i]!.name, e.formHints());
                }
            }

            if (i < selectedFiles.length - 1) {
                await new Promise<void>((resolve) => setTimeout(resolve, 1000));
            }
        }

        isLoading.value = false;
        input.value = '';
    };

    const deleteFile = (index: number) => {
        files.value.splice(index, 1);
        updateModelValue();
    };

    const moveUp = (index: number) => {
        if (index <= 0) return;
        const temp = files.value[index - 1]!;
        files.value[index - 1] = files.value[index]!;
        files.value[index] = temp;
        updateModelValue();
    };

    const moveDown = (index: number) => {
        if (index >= files.value.length - 1) return;
        const temp = files.value[index + 1];
        files.value[index + 1] = files.value[index]!;
        files.value[index] = temp!;
        updateModelValue();
    };
</script>

<template>
    <div>
        <div>
            <div
                v-for="(file, index) in filesWithLeads"
                :key="file.lead.id"
                :class="$style.file_item"
            >
                <div :class="$style.icon">
                    <a
                        :href="file.lead.url"
                        target="_blank"
                    >
                        <template v-if="isImage(file)">
                            <img
                                v-if="isImage(file)"
                                :src="file.lead.url"
                                alt="preview"
                            />
                        </template>
                        <template v-else> {{ getFileExtension(file.lead.url) ? `.${getFileExtension(file.lead.url)}` : '.unknown' }} </template>
                    </a>
                </div>

                <div :class="$style.file_id">
                    <div>{{ file.lead.filename }}</div>
                    <template
                        v-for="(f, k) in file.targets"
                        :key="k"
                    >
                        <div>
                            <a
                                :href="f.url"
                                target="_blank"
                            >
                                {{ k }}
                            </a>
                        </div>
                    </template>
                </div>

                <div :class="$style.actions">
                    <template v-if="mode === 'multi'">
                        <UButton
                            icon="i-lucide-chevron-up"
                            size="xs"
                            color="graylight"
                            :disabled="index === 0"
                            @click="moveUp(index)"
                        />
                        <UButton
                            icon="i-lucide-chevron-down"
                            size="xs"
                            color="graylight"
                            :disabled="index === files.length - 1"
                            @click="moveDown(index)"
                        />
                    </template>
                    <template v-if="!isNoneRemovable">
                        <UButton
                            icon="i-lucide-trash"
                            size="xs"
                            color="info"
                            @click="deleteFile(index)"
                        />
                    </template>
                </div>
            </div>
        </div>

        <div>
            <UButton
                as="label"
                class="bg-blue-500 text-white px-4 py-2 rounded cursor-pointer hover:bg-blue-600"
                :loading="isLoading"
                :disabled="isLoading || disabled"
            >
                {{ mode === 'multi' ? 'Выбрать файлы' : 'Выбрать файл' }}
                <input
                    type="file"
                    :multiple="mode === 'multi'"
                    :accept="acceptTypes"
                    class="hidden"
                    @change="handleFiles"
                />
            </UButton>
        </div>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    .file_item {
        padding: 5px;
        border-radius: 10px;
        background-color: var(--color-neutral-200);
        display: flex;
        gap: 10px;
        align-items: center;
        margin-bottom: 10px;

        > .icon {
            flex-shrink: 0;
            > a {
                display: block;
                width: 50px;
                height: 50px;
                // background: #fff;
                display: flex;
                align-items: center;
                justify-content: center;
                overflow: hidden;

                > img {
                    max-width: 50px;
                    max-height: 50px;
                }
            }
        }

        > .file_id {
            display: flex;
            flex-direction: column;
            gap: 2px;
            font-size: 11px;
            overflow: hidden;
            word-break: break-all;

            > div {
                > a {
                    text-decoration: underline;
                    color: #4c5dff;
                }
            }

            .width-size-sm-less({
                font-size: 10px;
            });
        }

        > .actions {
            flex-shrink: 0;
            margin-left: auto;
            display: flex;
            gap: 5px;
        }
    }
</style>
