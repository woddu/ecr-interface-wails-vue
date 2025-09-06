<!-- SearchInput.vue -->
<script setup lang="ts">

const props = defineProps<{
  modelValue: number,
  size?: 'sm' | 'md' | 'lg' | 'xl',
  label: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
}>()


</script>

<template>
  <UFormField :label="props.label" class="inline-block w-auto" :size="props.size">
    <UInputNumber
      class="max-w-fit"
      :model-value="props.modelValue"
      @update:model-value="emit('update:modelValue', $event)"
      color="primary"
      variant="soft"
      :min="0"
      :max="99"
      orientation="vertical" />
  </UFormField>
</template>
