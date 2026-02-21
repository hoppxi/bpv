<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { Music4, Play, Music2 } from "lucide-vue-next";
import type { AudioFile } from "@/types";
import { fetchGenres, fetchGenreTracks, getCoverArtUrl } from "@/lib/api";
import { formatDuration } from "@/lib/utils";
import TopBar from "@/components/layout/TopBar.vue";
import ScrollArea from "@/components/ui/ScrollArea.vue";

const props = defineProps<{
  tracks: AudioFile[];
  currentTrack: AudioFile | null;
  isPlaying: boolean;
}>();

const emit = defineEmits<{
  playTrack: [track: AudioFile, context?: AudioFile[]];
}>();

const genres = ref<{ name: string; count: number }[]>([]);
const selectedGenre = ref<string | null>(null);
const genreTracks = ref<AudioFile[]>([]);
const searchQuery = ref("");

const genreColors = [
  "from-red-600/40 to-orange-600/40",
  "from-blue-600/40 to-cyan-600/40",
  "from-green-600/40 to-emerald-600/40",
  "from-purple-600/40 to-pink-600/40",
  "from-yellow-600/40 to-amber-600/40",
  "from-indigo-600/40 to-violet-600/40",
  "from-rose-600/40 to-fuchsia-600/40",
  "from-teal-600/40 to-lime-600/40",
];

onMounted(async () => {
  try {
    genres.value = await fetchGenres();
    genres.value.sort((a, b) => b.count - a.count);
  } catch {
    const genreMap = new Map<string, number>();
    props.tracks.forEach((t) => {
      if (t.genre && t.genre !== "Unknown Genre") {
        genreMap.set(t.genre, (genreMap.get(t.genre) || 0) + 1);
      }
    });
    genres.value = Array.from(genreMap.entries())
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => b.count - a.count);
  }
});

async function selectGenre(name: string) {
  selectedGenre.value = name;
  try {
    genreTracks.value = await fetchGenreTracks(name);
  } catch {
    genreTracks.value = props.tracks.filter((t) => t.genre === name);
  }
}

const filteredGenres = computed(() => {
  if (!searchQuery.value) return genres.value;
  const q = searchQuery.value.toLowerCase();
  return genres.value.filter((g) => g.name.toLowerCase().includes(q));
});
</script>

<template>
  <div class="flex flex-col h-full">
    <TopBar
      title="Genres"
      :show-search="true"
      :search-query="searchQuery"
      @search="searchQuery = $event"
    />

    <ScrollArea class="flex-1" v-if="!selectedGenre">
      <div class="p-6">
        <div class="grid grid-cols-3 gap-4">
          <button
            v-for="(genre, i) in filteredGenres"
            :key="genre.name"
            @click="selectGenre(genre.name)"
            :class="[
              'group relative h-32 rounded-xl overflow-hidden text-left p-5 transition-all duration-200 hover:scale-[1.02] bg-gradient-to-br',
              genreColors[i % genreColors.length],
            ]"
          >
            <Music4 class="absolute right-4 bottom-4 w-12 h-12 text-white/10" />
            <p class="text-lg font-bold text-white">{{ genre.name }}</p>
            <p class="text-sm text-white/70 mt-1">{{ genre.count }} tracks</p>
          </button>
        </div>
      </div>
    </ScrollArea>

    <ScrollArea v-else class="flex-1">
      <div class="p-6">
        <button
          @click="selectedGenre = null"
          class="text-sm text-primary hover:underline mb-4 inline-block"
        >
          ← Back to Genres
        </button>
        <h3 class="text-2xl font-bold mb-6">{{ selectedGenre }}</h3>
        <div class="space-y-0.5">
          <div
            v-for="(track, i) in genreTracks"
            :key="track.file_path"
            @click="emit('playTrack', track, genreTracks)"
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
      </div>
    </ScrollArea>
  </div>
</template>
