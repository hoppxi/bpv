<script setup lang="ts">
import { ref, computed } from "vue";
import { Play, Pause, Clock, Music2, ArrowUpDown, Hash } from "lucide-vue-next";
import type { AudioFile, SortConfig } from "@/types";
import { getCoverArtUrl, fetchBasePath } from "@/lib/api";
import { formatDuration } from "@/lib/utils";
import TopBar from "@/components/layout/TopBar.vue";
import ScrollArea from "@/components/ui/ScrollArea.vue";
import Badge from "@/components/ui/Badge.vue";

const props = defineProps<{
  tracks: AudioFile[];
  currentTrack: AudioFile | null;
  isPlaying: boolean;
  favorites: Set<string>;
}>();

const emit = defineEmits<{
  playTrack: [track: AudioFile, context?: AudioFile[]];
}>();

const searchQuery = ref("");
const sort = ref<SortConfig>({ field: "title", direction: "asc" });

function toggleSort(field: SortConfig["field"]) {
  if (sort.value.field === field) {
    sort.value.direction = sort.value.direction === "asc" ? "desc" : "asc";
  } else {
    sort.value = { field, direction: "asc" };
  }
}

const filteredTracks = computed(() => {
  let tracks = [...props.tracks];

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    tracks = tracks.filter(
      (t) =>
        t.title.toLowerCase().includes(q) ||
        t.artist.toLowerCase().includes(q) ||
        t.album.toLowerCase().includes(q),
    );
  }

  tracks.sort((a, b) => {
    const dir = sort.value.direction === "asc" ? 1 : -1;
    switch (sort.value.field) {
      case "title":
        return a.title.localeCompare(b.title) * dir;
      case "artist":
        return a.artist.localeCompare(b.artist) * dir;
      case "album":
        return a.album.localeCompare(b.album) * dir;
      case "duration":
        return (a.duration - b.duration) * dir;
      case "year":
        return ((a.year || 0) - (b.year || 0)) * dir;
      case "genre":
        return a.genre.localeCompare(b.genre) * dir;
      default:
        return 0;
    }
  });

  return tracks;
});

const isCurrentTrack = (track: AudioFile) => props.currentTrack?.file_path === track.file_path;
</script>

<template>
  <div class="flex flex-col h-full">
    <TopBar
      title="Library"
      :show-search="true"
      :search-query="searchQuery"
      @search="searchQuery = $event"
    />

    <div
      class="grid grid-cols-[2.5rem_1fr_1fr_1fr_5rem] gap-4 px-6 py-2 text-xs font-medium text-muted-foreground uppercase tracking-wider border-b border-border/50 sticky top-14 bg-background/80 backdrop-blur-md z-10"
    >
      <span class="text-center">#</span>
      <button
        @click="toggleSort('title')"
        class="flex items-center gap-1 hover:text-foreground transition-colors"
      >
        Title
        <ArrowUpDown v-if="sort.field === 'title'" class="w-3 h-3" />
      </button>
      <button
        @click="toggleSort('artist')"
        class="flex items-center gap-1 hover:text-foreground transition-colors"
      >
        Artist
        <ArrowUpDown v-if="sort.field === 'artist'" class="w-3 h-3" />
      </button>
      <button
        @click="toggleSort('album')"
        class="flex items-center gap-1 hover:text-foreground transition-colors"
      >
        Album
        <ArrowUpDown v-if="sort.field === 'album'" class="w-3 h-3" />
      </button>
      <button
        @click="toggleSort('duration')"
        class="flex items-center gap-1 justify-end hover:text-foreground transition-colors"
      >
        <Clock class="w-3 h-3" />
        <ArrowUpDown v-if="sort.field === 'duration'" class="w-3 h-3" />
      </button>
    </div>

    <ScrollArea class="flex-1">
      <div class="pb-4">
        <div
          v-for="(track, i) in filteredTracks"
          :key="track.file_path"
          @click="emit('playTrack', track, filteredTracks)"
          :class="[
            'group items-center grid grid-cols-[2.5rem_3rem_1fr_1fr_1fr_5rem] gap-4 gap-y-4 px-6 p-3 rounded-xl hover:bg-white/5 transition-all cursor-pointer border border-transparent hover:border-white/10',
            isCurrentTrack(track) ? 'bg-white/10 border-white/10' : '',
          ]"
        >
          <div
            class="w-8 text-center text-sm font-medium text-muted-foreground group-hover:text-foreground"
          >
            <span v-if="isCurrentTrack(track) && isPlaying" class="flex justify-center">
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
            <p
              class="font-semibold truncate"
              :class="isCurrentTrack(track) ? 'text-primary' : 'text-foreground'"
            >
              {{ track.title }}
            </p>
          </div>

          <span class="text-muted-foreground truncate">{{ track.artist }}</span>
          <span class="text-muted-foreground truncate">{{ track.album }}</span>

          <span class="text-muted-foreground text-xs text-right tabular-nums">{{
            formatDuration(track.duration)
          }}</span>
        </div>
        <div
          v-if="filteredTracks.length === 0"
          class="flex flex-col items-center justify-center py-20 text-muted-foreground"
        >
          <Music2 class="w-12 h-12 mb-4 opacity-50" />
          <p class="text-lg font-medium">No tracks found</p>
          <p class="text-sm">Try adjusting your search query</p>
        </div>
      </div>
    </ScrollArea>
  </div>
</template>
