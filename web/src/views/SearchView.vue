<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { Search, Music2, X } from "lucide-vue-next";
import type { AudioFile } from "@/types";
import { searchLibrary, getCoverArtUrl } from "@/lib/api";
import { formatDuration, debounce } from "@/lib/utils";
import Input from "@/components/ui/Input.vue";
import ScrollArea from "@/components/ui/ScrollArea.vue";

const props = defineProps<{
  tracks: AudioFile[];
  currentTrack: AudioFile | null;
  isPlaying: boolean;
}>();

const emit = defineEmits<{
  playTrack: [track: AudioFile, context?: AudioFile[]];
}>();

const query = ref("");
const results = ref<AudioFile[]>([]);
const loading = ref(false);
const hasSearched = ref(false);

const doSearch = debounce(async (q: string) => {
  if (!q.trim()) {
    results.value = [];
    hasSearched.value = false;
    return;
  }
  loading.value = true;
  hasSearched.value = true;
  try {
    results.value = await searchLibrary(q);
  } catch {
    const lower = q.toLowerCase();
    results.value = props.tracks.filter(
      (t) =>
        t.title.toLowerCase().includes(lower) ||
        t.artist.toLowerCase().includes(lower) ||
        t.album.toLowerCase().includes(lower) ||
        t.genre.toLowerCase().includes(lower),
    );
  }
  loading.value = false;
}, 300);

watch(query, (val) => doSearch(val));
</script>

<template>
  <div class="flex flex-col h-full">
    <div class="p-6 pb-0">
      <div class="relative max-w-xl">
        <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
        <input
          v-model="query"
          type="text"
          placeholder="Search for songs, artists, albums..."
          class="w-full h-12 pl-12 pr-10 rounded-xl bg-secondary/50 border border-border/50 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50 transition-all"
          autofocus
        />
        <button
          v-if="query"
          @click="query = ''"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </div>

    <ScrollArea class="flex-1 mt-4">
      <div class="px-6 pb-6">
        <div v-if="loading" class="flex items-center justify-center py-20">
          <div
            class="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin"
          />
        </div>

        <div
          v-else-if="hasSearched && results.length === 0"
          class="flex flex-col items-center justify-center py-20 text-muted-foreground"
        >
          <Search class="w-12 h-12 mb-4 opacity-50" />
          <p class="text-lg font-medium">No results found</p>
          <p class="text-sm">Try a different search term</p>
        </div>

        <div v-else-if="results.length > 0" class="space-y-0.5">
          <p class="text-sm text-muted-foreground mb-4">{{ results.length }} results</p>
          <div
            v-for="(track, i) in results"
            :key="track.file_path"
            @click="emit('playTrack', track, results)"
            :class="[
              'group flex items-center gap-4 p-3 rounded-xl hover:bg-white/5 transition-all cursor-pointer border border-transparent hover:border-white/10',
              currentTrack?.file_path === track.file_path ? 'bg-white/10 border-white/10' : '',
            ]"
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

            <span class="text-xs text-muted-foreground tabular-nums">{{
              formatDuration(track.duration)
            }}</span>
          </div>
        </div>

        <div v-else class="flex flex-col items-center justify-center py-20 text-muted-foreground">
          <Search class="w-16 h-16 mb-4 opacity-30" />
          <p class="text-lg font-medium">Search your music</p>
          <p class="text-sm mt-1">Find songs, artists, albums, and genres</p>
          <p class="text-xs mt-4 text-muted-foreground/60">
            Press <kbd class="px-1.5 py-0.5 rounded bg-secondary text-xs">/</kbd> anywhere to search
          </p>
        </div>
      </div>
    </ScrollArea>
  </div>
</template>
