<script setup lang="ts">
import {
  Home,
  Library,
  Mic2,
  Disc3,
  Music4,
  Users,
  Heart,
  Search,
  Settings,
  Headphones,
} from "lucide-vue-next";
import type { ViewMode } from "@/types";

const props = defineProps<{
  currentView: ViewMode;
}>();

const emit = defineEmits<{
  navigate: [view: ViewMode];
}>();

const navItems: { icon: any; label: string; view: ViewMode }[] = [
  { icon: Home, label: "Home", view: "home" },
  { icon: Search, label: "Search", view: "search" },
  { icon: Library, label: "Library", view: "library" },
  { icon: Mic2, label: "Artists", view: "artists" },
  { icon: Disc3, label: "Albums", view: "albums" },
  { icon: Music4, label: "Genres", view: "genres" },
  { icon: Users, label: "Composers", view: "composers" },
  { icon: Heart, label: "Favorites", view: "favorites" },
];
</script>

<template>
  <aside
    class="fixed left-0 top-0 bottom-[var(--player-height)] w-[var(--sidebar-width)] bg-card/50 backdrop-blur-xl border-r border-border flex flex-col z-30"
  >
    <div class="p-5 flex items-center gap-3">
      <div
        class="w-9 h-9 rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center shadow-lg"
      >
        <Headphones class="w-5 h-5 text-white" />
      </div>
      <div>
        <h1 class="text-lg font-bold tracking-tight">BPV</h1>
        <p class="text-[10px] text-muted-foreground tracking-widest uppercase">Music Player</p>
      </div>
    </div>

    <nav class="flex-1 px-3 py-2 space-y-1 overflow-y-auto">
      <p class="px-3 py-2 text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
        Browse
      </p>
      <button
        v-for="item in navItems"
        :key="item.view"
        @click="emit('navigate', item.view)"
        :class="[
          'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200',
          currentView === item.view
            ? 'bg-primary/15 text-primary shadow-sm'
            : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground',
        ]"
      >
        <component :is="item.icon" class="w-4 h-4 shrink-0" />
        <span>{{ item.label }}</span>
        <div
          v-if="currentView === item.view"
          class="ml-auto w-1.5 h-1.5 rounded-full bg-primary animate-pulse-glow"
        />
      </button>

      <div class="h-px bg-border my-4" />

      <p class="px-3 py-2 text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
        System
      </p>
      <button
        @click="emit('navigate', 'settings')"
        :class="[
          'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200',
          currentView === 'settings'
            ? 'bg-primary/15 text-primary shadow-sm'
            : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground',
        ]"
      >
        <Settings class="w-4 h-4" />
        <span>Settings</span>
      </button>
    </nav>

    <div class="p-4 border-t border-border">
      <p class="text-[10px] text-muted-foreground text-center">BPV v0.2.0</p>
    </div>
  </aside>
</template>
