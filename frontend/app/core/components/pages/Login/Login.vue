<script setup lang="ts">
    import Alert from '../../shared/Alert/modals/Alert.vue';

    useSeoMeta({
        title: 'Вход в панель управления',
    });

    const route = useRoute();
    const router = useRouter();

    const showAccountVerified = async (afterClose?: () => void) => {
        const modal = useOverlay().create(Alert, {
            destroyOnClose: true,
            props: {
                text: 'Ваш аккаунт успешно активирован!',
            },
        });

        const instance = modal.open();

        await instance.result;

        if (afterClose) {
            afterClose();
        }
    };

    onMounted(() => {
        if (route.query['account-verified'] === 'true') {
            showAccountVerified(() => {
                router.replace({
                    query: {
                        ...route.query,
                        'account-verified': undefined,
                    },
                });
            });
        }
    });
</script>

<template>
    <div :id="$style.form"><FeatureLogin /></div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    #form {
        display: flex;
        justify-content: center;
        .std-wrapper(@std-x-padding, @std-x-padding, @std-x-padding-sm, @std-x-padding-sm);
    }
</style>
