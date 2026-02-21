<script setup lang="ts">
import { cn } from "@/lib/utils";

defineProps<{
  class?: string;
  modelValue?: number;
  max?: number;
  min?: number;
  step?: number;
}>();

defineEmits<{
  "update:modelValue": [value: number];
}>();
</script>

<template>
  <input
    type="range"
    :class="cn(
      'w-full h-1.5 rounded-full appearance-none cursor-pointer bg-secondary',
      '[&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-3 [&::-webkit-slider-thumb]:h-3 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-primary [&::-webkit-slider-thumb]:shadow-md [&::-webkit-slider-thumb]:transition-transform [&::-webkit-slider-thumb]:hover:scale-125',
      $props.class
    )"
    :value="modelValue"
    :min="min ?? 0"
    :max="max ?? 100"
    :step="step ?? 1"
    @input="$emit('update:modelValue', Number(($event.target as HTMLInputElement).value))"
  />
</template>
