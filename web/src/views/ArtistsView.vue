<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { Mic2, Play, Music2, ChevronRight } from "lucide-vue-next";
import type { AudioFile } from "@/types";
import { fetchArtists, fetchArtistTracks, getCoverArtUrl } from "@/lib/api";
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

const artists = ref<{ name: string; count: number }[]>([]);
const selectedArtist = ref<string | null>(null);
const artistTracks = ref<AudioFile[]>([]);
const searchQuery = ref("");
const loading = ref(false);

onMounted(async () => {
  try {
    artists.value = await fetchArtists();
    artists.value.sort((a, b) => a.name.localeCompare(b.name));
  } catch {
    const artistMap = new Map<string, number>();
    props.tracks.forEach((t) => {
      if (t.artist && t.artist !== "Unknown Artist") {
        artistMap.set(t.artist, (artistMap.get(t.artist) || 0) + 1);
      }
    });
    artists.value = Array.from(artistMap.entries())
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => a.name.localeCompare(b.name));
  }
});

async function selectArtist(name: string) {
  selectedArtist.value = name;
  loading.value = true;
  try {
    artistTracks.value = await fetchArtistTracks(name);
  } catch {
    artistTracks.value = props.tracks.filter((t) => t.artist === name);
  }
  loading.value = false;
}

const filteredArtists = computed(() => {
  if (!searchQuery.value) return artists.value;
  const q = searchQuery.value.toLowerCase();
  return artists.value.filter((a) => a.name.toLowerCase().includes(q));
});

function getArtistCover(artistName: string): string | null {
  const track = props.tracks.find((t) => t.artist === artistName && t.cover_art);
  return track ? getCoverArtUrl(track) : null;
}
</script>

<template>
  <div class="flex flex-col h-full">
    <TopBar
      title="Artists"
      :show-search="true"
      :search-query="searchQuery"
      @search="searchQuery = $event"
    />

    <div class="flex flex-1 min-h-0">
      <ScrollArea class="w-72 border-r border-border/50" v-if="!selectedArtist">
        <div class="p-3 space-y-1">
          <button
            v-for="artist in filteredArtists"
            :key="artist.name"
            @click="selectArtist(artist.name)"
            class="w-full flex items-center gap-3 p-2.5 rounded-lg hover:bg-accent/50 transition-colors text-left group"
          >
            <div class="w-10 h-10 rounded-full overflow-hidden bg-secondary flex-shrink-0">
              <img
                v-if="getArtistCover(artist.name)"
                :src="getArtistCover(artist.name)!"
                :alt="artist.name"
                class="w-full h-full object-cover"
              />
              <div
                v-else
                class="w-full h-full flex items-center justify-center bg-gradient-to-br from-pink-600/50 to-purple-600/50"
              >
                <Mic2 class="w-4 h-4 text-white/70" />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-sm font-medium truncate">{{ artist.name }}</p>
              <p class="text-xs text-muted-foreground">{{ artist.count }} tracks</p>
            </div>
            <ChevronRight
              class="w-4 h-4 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity"
            />
          </button>
        </div>
      </ScrollArea>

      <ScrollArea v-if="!selectedArtist" class="flex-1">
        <div class="p-6">
          <div class="grid grid-cols-4 gap-4">
            <button
              v-for="artist in filteredArtists"
              :key="artist.name"
              @click="selectArtist(artist.name)"
              class="group glass rounded-xl p-4 text-center hover:bg-white/10 transition-all duration-200"
            >
              <div class="w-20 h-20 mx-auto rounded-full overflow-hidden mb-3 bg-secondary">
                <img
                  v-if="getArtistCover(artist.name)"
                  :src="getArtistCover(artist.name)!"
                  :alt="artist.name"
                  class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                />
                <div
                  v-else
                  class="w-full h-full flex items-center justify-center bg-gradient-to-br from-pink-600/50 to-purple-600/50"
                >
                  <Mic2 class="w-8 h-8 text-white/60" />
                </div>
              </div>
              <p class="text-sm font-medium truncate">{{ artist.name }}</p>
              <p class="text-xs text-muted-foreground">{{ artist.count }} tracks</p>
            </button>
          </div>
        </div>
      </ScrollArea>

      <ScrollArea v-else class="flex-1">
        <div class="p-6">
          <button
            @click="selectedArtist = null"
            class="text-sm text-primary hover:underline mb-4 inline-block"
          >
            ← Back to Artists
          </button>
          <h3 class="text-2xl font-bold mb-6">{{ selectedArtist }}</h3>
          <div class="space-y-0.5">
            <div
              v-for="(track, i) in artistTracks"
              :key="track.file_path"
              @click="emit('playTrack', track, artistTracks)"
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
                      currentTrack?.file_path === track.file_path
                        ? 'text-primary'
                        : 'text-foreground'
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
  </div>
</template>
