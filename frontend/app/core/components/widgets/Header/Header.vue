<script setup lang="ts">
    import { doLogout } from '~/plugins/auth/model';

    const { smLess } = useScreenSizes();

    const { $authData } = useNuxtApp();

    const drawlerOpen = ref(false);

    useNuxtApp().hook('page:finish', () => {
        drawlerOpen.value = false;
    });
</script>

<template>
    <header id="auth_header">
        <div>
            <div class="left">
                <div class="menu">
                    <UDrawer
                        v-model:open="drawlerOpen"
                        direction="left"
                    >
                        <UButton
                            label="Меню"
                            :size="smLess ? 'sm' : 'md'"
                            icon="i-lucide-menu"
                        />

                        <template #content>
                            <div :id="$style.draw_menu"><WidgetMenu /></div>
                        </template>
                    </UDrawer>
                </div>
            </div>
            <div class="right">
                <div class="greet">Привет, {{ $authData.userData?.account.profileName }}!</div>
                <div class="logout">
                    <UButton
                        size="sm"
                        icon="i-lucide-log-out"
                        variant="subtle"
                        color="neutral"
                        @click="doLogout"
                    >
                        <template v-if="!smLess">Выход</template>
                    </UButton>
                </div>
            </div>
        </div>
    </header>
</template>

<style lang="less" module>
    @import '@styles/includes';

    #draw_menu {
        padding: 15px 15px 15px 15px;
        width: 250px;
        height: 100dvh;
        overflow-y: scroll;
    }
</style>
