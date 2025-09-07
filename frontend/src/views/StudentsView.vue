<script setup lang="ts">

const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

</script>

<template>
    <RouterView @loading="(value: boolean) => emit('loading', value)" />
</template>