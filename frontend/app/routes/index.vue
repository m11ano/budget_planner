<script setup lang="ts">
    definePageMeta({});
</script>

<template>
    <div>
        <template v-if="$authData.isAuth"> <LazyPageMain /></template>
        <template v-else> <LazyPageLogin /></template>
    </div>
</template>
