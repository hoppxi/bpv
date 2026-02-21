<script setup lang="ts">
import { cn } from "@/lib/utils";

defineProps<{
  class?: string;
  value?: number;
  max?: number;
}>();
</script>

<template>
  <div :class="cn('relative h-1.5 w-full overflow-hidden rounded-full bg-secondary', $props.class)">
    <div
      class="h-full rounded-full bg-primary transition-all duration-150"
      :style="{ width: `${Math.min(((value || 0) / (max || 100)) * 100, 100)}%` }"
    />
  </div>
</template>
