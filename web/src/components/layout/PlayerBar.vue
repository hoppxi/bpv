<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  Play,
  Pause,
  SkipBack,
  SkipForward,
  Shuffle,
  Repeat,
  Repeat1,
  Volume2,
  VolumeX,
  Volume1,
  Heart,
  ListMusic,
  Maximize2,
  Music2,
  Loader2,
  Activity,
} from "lucide-vue-next";
import Slider from "@/components/ui/Slider.vue";
import type { AudioFile, RepeatMode } from "@/types";
import { getCoverArtUrl } from "@/lib/api";
import { formatTime } from "@/lib/utils";
import { useFavorites } from "@/composables/useFavorites";

const { toggleFavorite } = useFavorites();

const props = defineProps<{
  currentTrack: AudioFile | null;
  isPlaying: boolean;
  isLoading: boolean;
  currentTime: number;
  duration: number;
  volume: number;
  muted: boolean;
  shuffle: boolean;
  repeat: RepeatMode;
  isFavorite: boolean;
}>();

const emit = defineEmits<{
  togglePlayPause: [];
  next: [];
  previous: [];
  seek: [time: number];
  volumeChange: [volume: number];
  toggleMute: [];
  toggleShuffle: [];
  cycleRepeat: [];
  openQueue: [];
  openNowPlaying: [];
  changeVisualizer: [];
}>();

const coverUrl = computed(() => (props.currentTrack ? getCoverArtUrl(props.currentTrack) : null));

const volumeIcon = computed(() => {
  if (props.muted || props.volume === 0) return VolumeX;
  if (props.volume < 0.5) return Volume1;
  return Volume2;
});

const repeatIcon = computed(() => {
  return props.repeat === "one" ? Repeat1 : Repeat;
});

const volumePercent = computed(() => Math.round(props.volume * 100));

function onSeek(val: number) {
  emit("seek", (val / 100) * props.duration);
}

function onVolumeChange(val: number) {
  emit("volumeChange", val / 100);
}

const progressPercent = computed(() => {
  if (!props.duration) return 0;
  return (props.currentTime / props.duration) * 100;
});
</script>

<template>
  <div
    class="fixed bottom-0 left-0 right-0 h-[var(--player-height)] bg-card/80 backdrop-blur-2xl border-t border-border z-40"
  >
    <div class="h-full flex items-center px-4 gap-4">
      <!-- left -->
      <div class="flex items-center gap-3 w-[300px] min-w-[200px]">
        <template v-if="currentTrack">
          <button
            @click="emit('openNowPlaying')"
            class="relative w-12 h-12 rounded-lg overflow-hidden flex-shrink-0 group shadow-lg"
          >
            <img
              v-if="coverUrl"
              :src="coverUrl"
              :alt="currentTrack.album"
              class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
            />
            <div
              v-else
              class="w-full h-full bg-gradient-to-br from-purple-600 to-pink-600 flex items-center justify-center"
            >
              <Music2 class="w-5 h-5 text-white/80" />
            </div>
            <div
              class="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-colors flex items-center justify-center"
            >
              <Maximize2
                class="w-4 h-4 text-white opacity-0 group-hover:opacity-100 transition-opacity"
              />
            </div>
          </button>

          <div class="min-w-0">
            <p class="text-sm font-medium truncate">{{ currentTrack.title }}</p>
            <p class="text-xs text-muted-foreground truncate">{{ currentTrack.artist }}</p>
          </div>

          <button
            @click="toggleFavorite(currentTrack.file_path)"
            :class="[
              'transition-colors',
              isFavorite ? 'text-red-500' : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            <Heart class="w-4 h-4" :fill="isFavorite ? 'currentColor' : 'none'" />
          </button>
        </template>
        <template v-else>
          <div class="w-12 h-12 rounded-lg bg-secondary flex items-center justify-center">
            <Music2 class="w-5 h-5 text-muted-foreground" />
          </div>
          <div>
            <p class="text-sm text-muted-foreground">No track selected</p>
          </div>
        </template>
      </div>

      <!-- center -->
      <div class="flex-1 flex flex-col items-center gap-1 max-w-[600px] mx-auto">
        <div class="flex items-center gap-3">
          <button
            @click="emit('toggleShuffle')"
            :class="[
              'w-8 h-8 flex items-center justify-center rounded-full transition-colors',
              shuffle ? 'text-primary' : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            <Shuffle class="w-4 h-4" />
          </button>

          <button
            @click="emit('previous')"
            class="w-8 h-8 flex items-center justify-center rounded-full text-foreground hover:bg-accent transition-colors"
          >
            <SkipBack class="w-4 h-4" fill="currentColor" />
          </button>

          <button
            @click="emit('togglePlayPause')"
            class="w-10 h-10 flex items-center justify-center rounded-full bg-white text-black hover:scale-105 active:scale-95 transition-all shadow-lg"
            :disabled="!currentTrack"
          >
            <Loader2 v-if="isLoading" class="w-5 h-5 animate-spin" />
            <Pause v-else-if="isPlaying" class="w-5 h-5" fill="currentColor" />
            <Play v-else class="w-5 h-5 ml-0.5" fill="currentColor" />
          </button>

          <button
            @click="emit('next')"
            class="w-8 h-8 flex items-center justify-center rounded-full text-foreground hover:bg-accent transition-colors"
          >
            <SkipForward class="w-4 h-4" fill="currentColor" />
          </button>

          <button
            @click="emit('cycleRepeat')"
            :class="[
              'w-8 h-8 flex items-center justify-center rounded-full transition-colors',
              repeat !== 'off' ? 'text-primary' : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            <component :is="repeatIcon" class="w-4 h-4" />
          </button>
        </div>

        <div class="w-full flex items-center gap-2">
          <span class="text-[11px] text-muted-foreground tabular-nums w-10 text-right">
            {{ formatTime(currentTime) }}
          </span>
          <Slider
            :model-value="progressPercent"
            :max="100"
            :step="0.1"
            class="flex-1"
            @update:model-value="onSeek"
          />
          <span class="text-[11px] text-muted-foreground tabular-nums w-10">
            {{ formatTime(duration) }}
          </span>
        </div>
      </div>

      <!-- right -->
      <div class="flex items-center gap-2 w-[200px] justify-end">
        <button
          @click="emit('changeVisualizer')"
          class="w-8 h-8 flex items-center justify-center rounded-full text-muted-foreground hover:text-foreground transition-colors"
          title="Change Visualizer"
        >
          <Activity class="w-4 h-4" />
        </button>

        <button
          @click="emit('openQueue')"
          class="w-8 h-8 flex items-center justify-center rounded-full text-muted-foreground hover:text-foreground transition-colors"
        >
          <ListMusic class="w-4 h-4" />
        </button>

        <button
          @click="emit('toggleMute')"
          class="w-8 h-8 flex items-center justify-center rounded-full text-muted-foreground hover:text-foreground transition-colors"
        >
          <component :is="volumeIcon" class="w-4 h-4" />
        </button>

        <Slider
          :model-value="volumePercent"
          :max="100"
          :step="1"
          class="w-24"
          @update:model-value="onVolumeChange"
        />
      </div>
    </div>
  </div>
</template>
