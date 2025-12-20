<script setup lang="ts">
    import type { FormSubmitEvent } from '@nuxt/ui';
    import * as v from 'valibot';
    import { reactive, ref } from 'vue';
    import { doLogin } from '~/plugins/auth/model';
    import { ApiError } from '~/shared/errors/errors';

    const isLoading = ref(false);
    const showPassword = ref(false);
    const toast = useToast();

    const schema = v.object({
        email: v.pipe(v.string(), v.trim(), v.email()),
        password: v.string(),
    });
    type Schema = v.InferOutput<typeof schema>;

    const formState = reactive<Schema>({
        email: '',
        password: '',
    });

    async function onSubmit(e: FormSubmitEvent<Schema>) {
        if (isLoading.value) return;

        isLoading.value = true;
        try {
            const { email, password } = e.data;

            await doLogin(email, password);

            toast.add({
                title: 'Успех',
                description: 'Вы успешно авторизовались',
                color: 'success',
                icon: 'i-lucide-check-circle',
            });
        } catch (err) {
            const errors: string[] = [];
            const hints: Record<string, string> = {
                ACCOUNT_NOT_CONFIRMED: 'Аккаунт не подтвержден',
                ACCOUNT_BLOCKED: 'Аккаунт заблокирован',
            };

            if (err instanceof ApiError) {
                if (err.code === 401) {
                    if (hints[err.textCode]) {
                        errors.push(hints[err.textCode]!);
                    } else {
                        errors.push('Неверный email или пароль');
                    }
                } else {
                    errors.push(...err.formHints());
                }
            } else {
                errors.push('Неизвестная ошибка');
            }

            toast.add({
                title: 'Возникли ошибки',
                description: errors.join('\n'),
                color: 'error',
                icon: 'i-lucide-alert-triangle',
            });
        } finally {
            isLoading.value = false;
        }
    }
</script>

<template>
    <div :class="$style.wrapper">
        <div :class="$style.title">Вход в приложение</div>

        <UForm
            :schema="schema"
            :state="formState"
            class="space-y-4"
            :class="$style.form"
            @submit.prevent="onSubmit"
        >
            <UFormField
                label="Email"
                name="email"
            >
                <UInput
                    v-model="formState.email"
                    class="w-full"
                    size="xl"
                />
            </UFormField>

            <UFormField
                label="Пароль"
                name="password"
            >
                <UInput
                    v-model="formState.password"
                    class="w-full"
                    size="xl"
                    :type="showPassword ? 'text' : 'password'"
                    :ui="{ trailing: 'pe-1' }"
                >
                    <template #trailing>
                        <UButton
                            color="neutral"
                            variant="link"
                            size="sm"
                            :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                            :aria-label="showPassword ? 'Скрыть пароль' : 'Показать пароль'"
                            :aria-pressed="showPassword"
                            aria-controls="password"
                            @click="showPassword = !showPassword"
                        />
                    </template>
                </UInput>
            </UFormField>

            <div class="flex justify-center">
                <UButton
                    type="submit"
                    size="xl"
                    trailing-icon="i-lucide-arrow-right"
                    :loading="isLoading"
                >
                    Войти
                </UButton>
            </div>

            <div class="flex justify-center mt-8">
                <NuxtLink
                    to="/registration"
                    :class="$style.link"
                    >Регистрация</NuxtLink
                >
            </div>
        </UForm>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    .wrapper {
        width: 100%;
        max-width: 360px;
    }

    .title {
        font-size: 32px;
        text-align: center;
        margin-bottom: 20px;
        font-family: 'Strong';

        .width-size-sm-less({
            font-size: 24px;
            margin-bottom: 15px;
        });
    }

    .form {
        background-color: var(--color-neutral-200);
        border-radius: 7px;
        padding: 30px;

        .width-size-sm-less({
            padding: 20px;
        });
    }

    .link {
        color: var(--color-blue-500);
        text-decoration: underline;
    }
</style>
