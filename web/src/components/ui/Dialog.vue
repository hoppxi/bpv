<script setup lang="ts">
import { cn } from "@/lib/utils";

defineProps<{
  class?: string;
  open?: boolean;
}>();

defineEmits<{
  "update:open": [value: boolean];
}>();
</script>

<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
        <div
          class="fixed inset-0 bg-black/80 backdrop-blur-sm"
          @click="$emit('update:open', false)"
        />
        <div
          :class="
            cn(
              'relative z-50 w-full max-w-lg max-h-[85vh] overflow-auto rounded-xl border bg-card p-6 shadow-2xl animate-fade-in',
              $props.class,
            )
          "
        >
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.2s ease;
}
.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}
</style>
