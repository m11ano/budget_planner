<script setup lang="ts">
    import EditorJS, { type ToolConstructable, type ToolSettings } from '@editorjs/editorjs';
    import Header from '@editorjs/header';
    import Underline from '@editorjs/underline';
    import ImageGallery from '~/shared/editorjs/modules/gallery2';
    import ImageTool from '~/shared/editorjs/modules/imageTool';
    import Hyperlink from '~/shared/editorjs/modules/link/Hyperlink';
    import CustomList from '~/shared/editorjs/modules/list/list';
    import Paragraph from '~/shared/editorjs/modules/paragraph';
    import Sortable from '~/shared/editorjs/modules/sortable/Sortable';
    import { normalizeParagraphEmpties } from './helpers';
    import { i18n } from './i18n';
    import { uploadToApi } from './uploader';

    interface Props {
        modelValue?: any;
        placeholder?: string;
        useLists?: boolean;
        useHeaders?: boolean;
        useImages?: boolean;
        useGalleries?: boolean;
        imageUploaderUrl?: string;
    }
    const props = defineProps<Props>();
    const emit = defineEmits<{
        (e: 'update:modelValue', value: any): void;
    }>();

    const holderRef = ref<HTMLDivElement | null>(null);
    let editor: EditorJS | null = null;

    const tools: Record<string, ToolConstructable | ToolSettings> = {
        paragraph: {
            class: Paragraph as any,
            inlineToolbar: true,
            config: {
                preserveBlank: true,
            },
        },
        link: {
            class: Hyperlink as any,
            config: {
                shortcut: 'CMD+L',
                availableTargets: ['_blank'],
                availableRels: ['nofollow'],
                validate: false,
            },
        },
        underline: Underline,
    };

    if (props.useLists) {
        tools.list = {
            class: CustomList as any,
            inlineToolbar: true,
            config: {
                defaultStyle: 'unordered',
                counterTypes: ['numeric'],
            },
        };
    }

    if (props.useHeaders) {
        tools.header = {
            class: Header as any,
            config: {
                placeholder: 'Введите заголовок',
                levels: [1, 2, 3],
                defaultLevel: 2,
            },
        };
    }

    if (props.useImages) {
        tools.image = {
            class: ImageTool as any,
            config: {
                captionPlaceholder: 'Ссылка при клике (необязательно)...',
                uploader: {
                    uploadByFile: (file: File) => {
                        return new Promise((ready, reject) => {
                            uploadToApi(props.imageUploaderUrl || '', file)
                                .then((data) => {
                                    ready(data);
                                })
                                .catch((e) => {
                                    reject(e);
                                });
                        });
                    },
                    uploadByUrl: (url: string) => {
                        return new Promise((ready) => {
                            ready(url);
                        })
                            .then((imageUrl) => {
                                return {
                                    success: 1,
                                    file: {
                                        url: imageUrl,
                                        type: 'url',
                                    },
                                };
                            })
                            .catch(() => {
                                return {
                                    success: 0,
                                };
                            });
                    },
                },
            },
        };
    }

    if (props.useGalleries) {
        tools.gallery = {
            class: ImageGallery as any,
            config: {
                sortableJs: Sortable,
                uploader: {
                    uploadByFile: (file: File) => {
                        return new Promise((ready, reject) => {
                            uploadToApi(props.imageUploaderUrl || '', file)
                                .then((data) => {
                                    ready(data);
                                })
                                .catch((e) => {
                                    reject(e);
                                });
                        });
                    },
                },
            },
        };
    }

    const build = () => {
        editor = new EditorJS({
            holder: holderRef.value!,
            placeholder: props.placeholder || 'Введите текст...',
            tools: tools,
            data: JSON.parse(JSON.stringify(props.modelValue || {})),
            i18n: i18n,
            async onChange(api) {
                const raw = await api.saver.save();
                const normalized = normalizeParagraphEmpties(raw);
                emit('update:modelValue', normalized);
            },
        });
    };

    const clear = () => {
        editor?.destroy();
        editor = null;
    };

    onMounted(() => {
        build();
    });

    onBeforeUnmount(() => {
        clear();
    });

    const save = async () => {
        if (!editor) return null;
        await editor.save();
    };

    const rebuild = async () => {
        clear();
        build();
    };

    defineExpose({ save, rebuild });
</script>

<template>
    <div>
        <div
            ref="holderRef"
            class="editorjs__holder"
        />
    </div>
</template>

<style lang="less">
    @import '@styles/includes';

    .editorjs__holder {
        width: 100%;
        border: 1px solid #ccc;
        border-radius: 6px;
        padding: 8px;
        background-color: #f8f8f8;

        .ce-block {
            background-color: #fff;
            margin: 10px 0;

            a {
                color: #4c5dff;
            }

            h1.ce-header {
                font-weight: normal;
                font-size: 28px;
            }

            h2.ce-header {
                font-weight: normal;
                font-size: 24px;
            }

            h3.ce-header {
                font-weight: normal;
                font-size: 20px;
            }
        }
    }
</style>
