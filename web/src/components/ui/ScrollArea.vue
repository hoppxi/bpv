<script setup lang="ts">
import { cn } from "@/lib/utils";

defineProps<{
  class?: string;
}>();
</script>

<template>
  <div
    :class="cn(
      'overflow-auto scrollbar-thin',
      $props.class
    )"
  >
    <slot />
  </div>
</template>
