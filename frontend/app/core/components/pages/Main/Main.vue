<script setup lang="ts">
    import { setAppBreadcrumbs } from '~/plugins/app/model/actions/setAppBreadcrumbs';
    import { setMenu } from '~/plugins/app/model/actions/setMenu';

    useSeoMeta({
        title: '',
    });

    setMenu('', '');

    setAppBreadcrumbs([]);
</script>

<template>
    <div>
        <div :id="$style.desc">Добро пожаловать в budget planner</div>
    </div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    #desc {
        font-family: 'Strong';
        font-size: 32px;

        .width-size-sm-less({
            font-size: 24px;
        });
    }
</style>
