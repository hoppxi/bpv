<script setup lang="ts">
import { cn } from "@/lib/utils";

defineProps<{
  class?: string;
  variant?: "default" | "secondary" | "destructive" | "outline";
}>();
</script>

<template>
  <span
    :class="cn(
      'inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
      {
        'border-transparent bg-primary text-primary-foreground shadow': variant === 'default' || !variant,
        'border-transparent bg-secondary text-secondary-foreground': variant === 'secondary',
        'border-transparent bg-destructive text-destructive-foreground shadow': variant === 'destructive',
        'text-foreground': variant === 'outline',
      },
      $props.class
    )"
  >
    <slot />
  </span>
</template>
