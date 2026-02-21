<script setup lang="ts">
import { Search, Bell } from "lucide-vue-next";
import Input from "@/components/ui/Input.vue";

const props = defineProps<{
  title: string;
  searchQuery?: string;
  showSearch?: boolean;
}>();

const emit = defineEmits<{
  search: [query: string];
}>();
</script>

<template>
  <header
    class="h-14 flex items-center justify-between px-6 border-b border-border/50 bg-background/50 backdrop-blur-md sticky top-0 z-20"
  >
    <h2 class="text-lg font-semibold">{{ title }}</h2>
    <div class="flex items-center gap-3">
      <div v-if="showSearch" class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <Input
          :model-value="searchQuery"
          @update:model-value="emit('search', $event)"
          placeholder="Search..."
          class="pl-9 w-64 bg-secondary/50 border-0"
        />
      </div>
    </div>
  </header>
</template>
