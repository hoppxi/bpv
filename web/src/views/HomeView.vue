<script setup lang="ts">
import { computed } from "vue";
import { Play, Music2, Clock, Disc3, Mic2, TrendingUp } from "lucide-vue-next";
import type { AudioFile, LibraryResponse } from "@/types";
import { getCoverArtUrl } from "@/lib/api";
import { formatDuration } from "@/lib/utils";

const props = defineProps<{
  library: LibraryResponse | null;
  currentTrack: AudioFile | null;
}>();

const emit = defineEmits<{
  playTrack: [track: AudioFile];
  navigate: [view: string];
}>();

const recentTracks = computed(() => {
  if (!props.library) return [];
  return [...props.library.files]
    .sort((a, b) => new Date(b.modified).getTime() - new Date(a.modified).getTime())
    .slice(0, 8);
});

const stats = computed(() => {
  if (!props.library) return { tracks: 0, artists: 0, albums: 0, genres: 0 };
  return {
    tracks: props.library.audio_files,
    artists: Object.keys(props.library.artists || {}).length,
    albums: Object.keys(props.library.albums || {}).length,
    genres: Object.keys(props.library.genres || {}).length,
  };
});
</script>

<template>
  <div class="p-6 space-y-8">
    <div>
      <h1 class="text-3xl font-bold mb-2">
        Good
        {{
          new Date().getHours() < 12
            ? "Morning"
            : new Date().getHours() < 18
              ? "Afternoon"
              : "Evening"
        }}
        🎶
      </h1>
      <p class="text-muted-foreground">Enjoy your music collection</p>
    </div>

    <div class="grid grid-cols-4 gap-4">
      <div
        class="glass rounded-xl p-4 cursor-pointer hover:bg-white/10 transition-colors"
        @click="emit('navigate', 'library')"
      >
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-purple-500/20 flex items-center justify-center">
            <Music2 class="w-5 h-5 text-purple-400" />
          </div>
          <div>
            <p class="text-2xl font-bold">{{ stats.tracks }}</p>
            <p class="text-xs text-muted-foreground">Tracks</p>
          </div>
        </div>
      </div>
      <div
        class="glass rounded-xl p-4 cursor-pointer hover:bg-white/10 transition-colors"
        @click="emit('navigate', 'artists')"
      >
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-pink-500/20 flex items-center justify-center">
            <Mic2 class="w-5 h-5 text-pink-400" />
          </div>
          <div>
            <p class="text-2xl font-bold">{{ stats.artists }}</p>
            <p class="text-xs text-muted-foreground">Artists</p>
          </div>
        </div>
      </div>
      <div
        class="glass rounded-xl p-4 cursor-pointer hover:bg-white/10 transition-colors"
        @click="emit('navigate', 'albums')"
      >
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-blue-500/20 flex items-center justify-center">
            <Disc3 class="w-5 h-5 text-blue-400" />
          </div>
          <div>
            <p class="text-2xl font-bold">{{ stats.albums }}</p>
            <p class="text-xs text-muted-foreground">Albums</p>
          </div>
        </div>
      </div>
      <div
        class="glass rounded-xl p-4 cursor-pointer hover:bg-white/10 transition-colors"
        @click="emit('navigate', 'genres')"
      >
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-green-500/20 flex items-center justify-center">
            <TrendingUp class="w-5 h-5 text-green-400" />
          </div>
          <div>
            <p class="text-2xl font-bold">{{ stats.genres }}</p>
            <p class="text-xs text-muted-foreground">Genres</p>
          </div>
        </div>
      </div>
    </div>

    <div>
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold">Recently Added</h2>
        <button @click="emit('navigate', 'library')" class="text-sm text-primary hover:underline">
          See all
        </button>
      </div>
      <div class="grid grid-cols-4 gap-4">
        <button
          v-for="track in recentTracks"
          :key="track.file_path"
          @click="emit('playTrack', track)"
          class="group glass rounded-xl p-3 text-left hover:bg-white/10 transition-all duration-200"
        >
          <div class="relative aspect-square rounded-lg overflow-hidden mb-3 bg-secondary">
            <img
              v-if="getCoverArtUrl(track)"
              :src="getCoverArtUrl(track)!"
              :alt="track.album"
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
            />
            <div
              v-else
              class="w-full h-full flex items-center justify-center bg-gradient-to-br from-purple-600/50 to-pink-600/50"
            >
              <Music2 class="w-8 h-8 text-white/60" />
            </div>

            <div
              class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors flex items-center justify-center"
            >
              <div
                class="w-10 h-10 rounded-full bg-primary flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all transform translate-y-2 group-hover:translate-y-0 shadow-lg"
              >
                <Play class="w-5 h-5 text-white ml-0.5" fill="currentColor" />
              </div>
            </div>
          </div>
          <p class="text-sm font-medium line-clamp-1">{{ track.title }}</p>
          <p class="text-xs text-muted-foreground line-clamp-1">{{ track.artist }}</p>
        </button>
      </div>
    </div>
  </div>
</template>
