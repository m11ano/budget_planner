<script setup lang="ts">
    useSeoMeta({
        title: 'Регистрация',
    });
</script>

<template>
    <div :id="$style.form"><FeatureRegistration /></div>
</template>

<style lang="less" module>
    @import '@styles/includes';

    #form {
        display: flex;
        justify-content: center;
        .std-wrapper(@std-x-padding, @std-x-padding, @std-x-padding-sm, @std-x-padding-sm);
    }
</style>
