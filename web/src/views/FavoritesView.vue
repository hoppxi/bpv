<script setup lang="ts">
import { computed } from "vue";
import { Heart, Music2, Play } from "lucide-vue-next";
import type { AudioFile } from "@/types";
import { getCoverArtUrl } from "@/lib/api";
import { formatDuration } from "@/lib/utils";
import TopBar from "@/components/layout/TopBar.vue";
import { useFavorites } from "@/composables/useFavorites";

const props = defineProps<{
  tracks: AudioFile[];
  currentTrack: AudioFile | null;
  isPlaying: boolean;
}>();

const emit = defineEmits<{
  playTrack: [track: AudioFile, context?: AudioFile[]];
}>();

const { favorites, toggleFavorite } = useFavorites();

const favoriteTracks = computed(() => props.tracks.filter((t) => favorites.value.has(t.file_path)));
</script>

<template>
  <div class="flex flex-col h-full">
    <TopBar title="Favorites" />

    <div class="flex-1 overflow-y-auto p-6 scrollbar-thin">
      <div
        v-if="favoriteTracks.length === 0"
        class="flex flex-col items-center justify-center h-full text-muted-foreground"
      >
        <div class="w-20 h-20 rounded-full bg-secondary flex items-center justify-center mb-4">
          <Heart class="w-10 h-10 text-muted-foreground/50" />
        </div>
        <p class="text-xl font-semibold mb-2">No favorites yet</p>
        <p class="text-sm max-w-sm text-center">
          Click the heart icon on any track to add it to your favorites collection.
        </p>
      </div>

      <div v-else class="space-y-2">
        <div
          v-for="(track, i) in favoriteTracks"
          :key="track.file_path"
          @click="emit('playTrack', track, favoriteTracks)"
          class="group flex items-center gap-4 p-3 rounded-xl hover:bg-white/5 transition-all cursor-pointer border border-transparent hover:border-white/10"
          :class="currentTrack?.file_path === track.file_path ? 'bg-white/10 border-white/10' : ''"
        >
          <div
            class="w-8 text-center text-sm font-medium text-muted-foreground group-hover:text-foreground"
          >
            <span
              v-if="currentTrack?.file_path === track.file_path && isPlaying"
              class="flex justify-center"
            >
              <div
                class="w-1 h-1 bg-primary rounded-full animate-bounce mx-0.5"
                style="animation-delay: 0ms"
              ></div>
              <div
                class="w-1 h-1 bg-primary rounded-full animate-bounce mx-0.5"
                style="animation-delay: 150ms"
              ></div>
              <div
                class="w-1 h-1 bg-primary rounded-full animate-bounce mx-0.5"
                style="animation-delay: 300ms"
              ></div>
            </span>
            <span v-else>{{ i + 1 }}</span>
          </div>

          <div class="relative w-12 h-12 rounded-lg overflow-hidden flex-shrink-0 shadow-md">
            <img
              v-if="getCoverArtUrl(track)"
              :src="getCoverArtUrl(track)!"
              :alt="track.album"
              class="w-full h-full object-cover"
            />
            <div
              v-else
              class="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-700 to-gray-600"
            >
              <Music2 class="w-6 h-6 text-white/50" />
            </div>

            <div
              class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
            >
              <Play class="w-5 h-5 text-white fill-current" />
            </div>
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <p
                class="font-semibold truncate"
                :class="
                  currentTrack?.file_path === track.file_path ? 'text-primary' : 'text-foreground'
                "
              >
                {{ track.title }}
              </p>
            </div>
            <p class="text-xs text-muted-foreground truncate">
              {{ track.artist }} <span v-if="track.album">• {{ track.album }}</span>
            </p>
          </div>

          <span class="text-xs font-medium text-muted-foreground tabular-nums px-2">
            {{ formatDuration(track.duration) }}
          </span>

          <button
            @click.stop="toggleFavorite(track.file_path)"
            class="p-2 rounded-full hover:bg-white/10 text-red-500 transition-colors"
          >
            <Heart class="w-5 h-5" fill="currentColor" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
