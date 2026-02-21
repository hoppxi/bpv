<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { Users, Music2 } from "lucide-vue-next";
import type { AudioFile } from "@/types";
import { fetchComposers, fetchComposerTracks, getCoverArtUrl } from "@/lib/api";
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

const composers = ref<{ name: string; count: number }[]>([]);
const selectedComposer = ref<string | null>(null);
const composerTracks = ref<AudioFile[]>([]);
const searchQuery = ref("");

onMounted(async () => {
  try {
    composers.value = await fetchComposers();
    composers.value.sort((a, b) => a.name.localeCompare(b.name));
  } catch {
    const composerMap = new Map<string, number>();
    props.tracks.forEach((t) => {
      if (t.composer && t.composer !== "Unknown Composer" && t.composer !== "") {
        composerMap.set(t.composer, (composerMap.get(t.composer) || 0) + 1);
      }
    });
    composers.value = Array.from(composerMap.entries())
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => a.name.localeCompare(b.name));
  }
});

async function selectComposer(name: string) {
  selectedComposer.value = name;
  try {
    composerTracks.value = await fetchComposerTracks(name);
  } catch {
    composerTracks.value = props.tracks.filter((t) => t.composer === name);
  }
}

const filteredComposers = computed(() => {
  if (!searchQuery.value) return composers.value;
  const q = searchQuery.value.toLowerCase();
  return composers.value.filter((c) => c.name.toLowerCase().includes(q));
});
</script>

<template>
  <div class="flex flex-col h-full">
    <TopBar
      title="Composers"
      :show-search="true"
      :search-query="searchQuery"
      @search="searchQuery = $event"
    />

    <ScrollArea class="flex-1" v-if="!selectedComposer">
      <div class="p-6 space-y-1">
        <button
          v-for="composer in filteredComposers"
          :key="composer.name"
          @click="selectComposer(composer.name)"
          class="w-full flex items-center gap-4 p-3 rounded-lg hover:bg-accent/50 transition-colors text-left"
        >
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-amber-600/50 to-orange-600/50 flex items-center justify-center"
          >
            <Users class="w-4 h-4 text-white/70" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium">{{ composer.name }}</p>
            <p class="text-xs text-muted-foreground">{{ composer.count }} tracks</p>
          </div>
        </button>
        <div
          v-if="filteredComposers.length === 0"
          class="flex flex-col items-center justify-center py-20 text-muted-foreground"
        >
          <Users class="w-12 h-12 mb-4 opacity-50" />
          <p>No composers found</p>
        </div>
      </div>
    </ScrollArea>

    <ScrollArea v-else class="flex-1">
      <div class="p-6">
        <button
          @click="selectedComposer = null"
          class="text-sm text-primary hover:underline mb-4 inline-block"
        >
          ← Back to Composers
        </button>
        <h3 class="text-2xl font-bold mb-6">{{ selectedComposer }}</h3>
        <div class="space-y-0.5">
          <div
            v-for="(track, i) in composerTracks"
            :key="track.file_path"
            @click="emit('playTrack', track, composerTracks)"
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
