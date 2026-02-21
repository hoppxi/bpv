<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { Disc3, Play, Music2 } from "lucide-vue-next";
import type { AudioFile } from "@/types";
import { fetchAlbums, fetchAlbumTracks, getCoverArtUrl } from "@/lib/api";
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

const albums = ref<{ name: string; count: number }[]>([]);
const selectedAlbum = ref<string | null>(null);
const albumTracks = ref<AudioFile[]>([]);
const searchQuery = ref("");

onMounted(async () => {
  try {
    albums.value = await fetchAlbums();
    albums.value.sort((a, b) => a.name.localeCompare(b.name));
  } catch {
    const albumMap = new Map<string, number>();
    props.tracks.forEach((t) => {
      if (t.album && t.album !== "Unknown Album") {
        albumMap.set(t.album, (albumMap.get(t.album) || 0) + 1);
      }
    });
    albums.value = Array.from(albumMap.entries())
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => a.name.localeCompare(b.name));
  }
});

async function selectAlbum(name: string) {
  selectedAlbum.value = name;
  try {
    albumTracks.value = await fetchAlbumTracks(name);
  } catch {
    albumTracks.value = props.tracks.filter((t) => t.album === name);
  }
}

const filteredAlbums = computed(() => {
  if (!searchQuery.value) return albums.value;
  const q = searchQuery.value.toLowerCase();
  return albums.value.filter((a) => a.name.toLowerCase().includes(q));
});

function getAlbumCover(albumName: string): string | null {
  const track = props.tracks.find((t) => t.album === albumName && t.cover_art);
  return track ? getCoverArtUrl(track) : null;
}

function getAlbumArtist(albumName: string): string {
  const track = props.tracks.find((t) => t.album === albumName);
  return track?.artist || "Various Artists";
}
</script>

<template>
  <div class="flex flex-col h-full">
    <TopBar
      title="Albums"
      :show-search="true"
      :search-query="searchQuery"
      @search="searchQuery = $event"
    />

    <ScrollArea class="flex-1" v-if="!selectedAlbum">
      <div class="p-6">
        <div class="grid grid-cols-5 gap-5">
          <button
            v-for="album in filteredAlbums"
            :key="album.name"
            @click="selectAlbum(album.name)"
            class="group text-left"
          >
            <div
              class="relative aspect-square rounded-xl overflow-hidden mb-3 bg-secondary shadow-lg"
            >
              <img
                v-if="getAlbumCover(album.name)"
                :src="getAlbumCover(album.name)!"
                :alt="album.name"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              />
              <div
                v-else
                class="w-full h-full flex items-center justify-center bg-gradient-to-br from-blue-600/50 to-purple-600/50"
              >
                <Disc3 class="w-10 h-10 text-white/60" />
              </div>
              <div
                class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors flex items-center justify-center"
              >
                <div
                  class="w-12 h-12 rounded-full bg-primary flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all transform translate-y-2 group-hover:translate-y-0 shadow-lg"
                >
                  <Play class="w-6 h-6 text-white ml-0.5" fill="currentColor" />
                </div>
              </div>
            </div>
            <p class="text-sm font-medium line-clamp-1">{{ album.name }}</p>
            <p class="text-xs text-muted-foreground">
              {{ getAlbumArtist(album.name) }} · {{ album.count }} tracks
            </p>
          </button>
        </div>
      </div>
    </ScrollArea>

    <ScrollArea v-else class="flex-1">
      <div class="p-6">
        <button
          @click="selectedAlbum = null"
          class="text-sm text-primary hover:underline mb-6 inline-block"
        >
          ← Back to Albums
        </button>
        <div class="flex gap-6 mb-8">
          <div class="w-48 h-48 rounded-xl overflow-hidden bg-secondary shadow-2xl flex-shrink-0">
            <img
              v-if="getAlbumCover(selectedAlbum)"
              :src="getAlbumCover(selectedAlbum)!"
              class="w-full h-full object-cover"
            />
            <div
              v-else
              class="w-full h-full flex items-center justify-center bg-gradient-to-br from-blue-600/50 to-purple-600/50"
            >
              <Disc3 class="w-16 h-16 text-white/60" />
            </div>
          </div>
          <div class="flex flex-col justify-end">
            <p class="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">
              Album
            </p>
            <h3 class="text-3xl font-bold mb-2">{{ selectedAlbum }}</h3>
            <p class="text-sm text-muted-foreground">
              {{ getAlbumArtist(selectedAlbum) }} · {{ albumTracks.length }} tracks
            </p>
          </div>
        </div>

        <div class="space-y-0.5">
          <div
            v-for="(track, i) in albumTracks"
            :key="track.file_path"
            @click="emit('playTrack', track, albumTracks)"
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
