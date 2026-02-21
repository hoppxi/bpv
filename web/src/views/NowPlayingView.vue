<script setup lang="ts">
import { ref, computed } from "vue";
import {
  ChevronDown,
  Heart,
  Music2,
  Share2,
  ListMusic,
  Disc3,
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
  Maximize2,
  Minimize2,
  Activity,
} from "lucide-vue-next";
import type { AudioFile, VisualizerType, RepeatMode } from "@/types";
import { getCoverArtUrl } from "@/lib/api";
import {
  formatTime,
  formatDuration,
  formatFileSize,
  formatBitrate,
  formatSampleRate,
} from "@/lib/utils";
import Slider from "@/components/ui/Slider.vue";
import Badge from "@/components/ui/Badge.vue";
import Separator from "@/components/ui/Separator.vue";
import VisualizerCanvas from "@/components/visualizer/VisualizerCanvas.vue";
import { useFavorites } from "@/composables/useFavorites";

const { toggleFavorite, isFavorite } = useFavorites();

const props = defineProps<{
  currentTrack: AudioFile | null;
  isPlaying: boolean;
  isLoading?: boolean;
  currentTime: number;
  duration: number;
  volume?: number;
  muted?: boolean;
  shuffle?: boolean;
  repeat?: RepeatMode;
  frequencyData: Uint8Array;
  timeDomainData: Uint8Array;
  visualizerType: VisualizerType;
  colorPalette: string[];
  dominantColor: string;
}>();

const emit = defineEmits<{
  close: [];
  togglePlayPause: [];
  next: [];
  previous: [];
  seek: [time: number];
  volumeChange: [volume: number];
  toggleMute: [];
  toggleShuffle: [];
  cycleRepeat: [];
  changeVisualizer: [];
}>();

const coverUrl = computed(() => (props.currentTrack ? getCoverArtUrl(props.currentTrack) : null));

const progress = computed(() => {
  if (!props.duration) return 0;
  return (props.currentTime / props.duration) * 100;
});

const volumeIcon = computed(() => {
  if (props.muted || (props.volume || 1) === 0) return VolumeX;
  if ((props.volume || 1) < 0.5) return Volume1;
  return Volume2;
});

const repeatIcon = computed(() => {
  return props.repeat === "one" ? Repeat1 : Repeat;
});

const volumePercent = computed(() => Math.round((props.volume || 1) * 100));

function onSeek(val: number) {
  emit("seek", (val / 100) * props.duration);
}

function onVolumeChange(val: number) {
  emit("volumeChange", val / 100);
}

const isFullscreen = ref(false);

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen().catch((e) => {
      console.error(`Error attempting to enable fullscreen mode: ${e.message} (${e.name})`);
    });
    isFullscreen.value = true;
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen();
      isFullscreen.value = false;
    }
  }
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex flex-col">
    <VisualizerCanvas
      v-if="visualizerType !== 'none'"
      :frequency-data="frequencyData"
      :time-domain-data="timeDomainData"
      :visualizer-type="visualizerType"
      :color-palette="colorPalette"
      :dominant-color="dominantColor"
      :is-playing="isPlaying"
      fullscreen
    />

    <div
      class="absolute inset-0 pointer-events-none z-0 transition-all duration-1000"
      :style="{
        background: `radial-gradient(circle at 30% 70%, ${dominantColor}),
                     radial-gradient(circle at 70% 30%, ${colorPalette[2] || dominantColor})`,
      }"
    />
    <div
      class="absolute inset-0 z-0 bg-gradient-to-t from-background via-background/80 to-background/30 pointer-events-none"
    />

    <!-- Exit Button -->
    <div class="absolute top-6 left-6 z-50 flex gap-4">
      <button
        @click="emit('close')"
        class="w-10 h-10 rounded-full glass flex items-center justify-center hover:bg-white/10 transition-colors pointer-events-auto cursor-pointer"
      >
        <ChevronDown class="w-6 h-6" />
      </button>
    </div>

    <!-- Fullscreen Button -->
    <div class="absolute top-6 right-6 z-50">
      <button
        @click="toggleFullscreen"
        class="w-10 h-10 rounded-full glass flex items-center justify-center hover:bg-white/10 transition-colors pointer-events-auto cursor-pointer"
      >
        <Minimize2 v-if="isFullscreen" class="w-5 h-5" />
        <Maximize2 v-else class="w-5 h-5" />
      </button>
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col relative z-20 p-8 h-full">
      <!-- Top/Center area: Visualizer (implicit) & Lyrics -->
      <div class="flex-1 flex items-center justify-center relative min-h-0">
        <div
          v-if="currentTrack?.lyrics"
          class="max-w-2xl w-full text-center overflow-y-auto max-h-full scrollbar-none p-4 animate-in fade-in slide-in-from-bottom-4 duration-700"
        >
          <p
            class="text-lg md:text-2xl font-semibold leading-relaxed whitespace-pre-wrap text-foreground/90 drop-shadow-md"
          >
            {{ currentTrack.lyrics }}
          </p>
        </div>
        <div v-else class="flex flex-col items-center justify-center text-muted-foreground/50">
          <!-- Placeholder if no lyrics, visualizer takes stage -->
        </div>
      </div>

      <!-- Bottom Layout -->
      <div class="flex flex-col md:flex-row gap-8 items-end w-full mt-auto">
        <!-- Cover Art (Bottom Left) -->
        <div class="flex-shrink-0 animate-in fade-in zoom-in-95 duration-500">
          <div class="w-48 h-48 md:w-64 md:h-64 shadow-2xl relative group">
            <img
              v-if="coverUrl"
              :src="coverUrl"
              :alt="currentTrack?.album"
              class="w-full h-full object-cover rounded-lg shadow-2xl"
            />
            <div
              v-else
              class="w-full h-full flex items-center justify-center bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg"
            >
              <Music2 class="w-24 h-24 text-white/40" />
            </div>
          </div>
        </div>

        <!-- Track Info & Controls -->
        <div
          class="flex-1 flex flex-col w-full min-w-0 gap-6 pb-2 animate-in slide-in-from-bottom-8 fade-in duration-700 delay-100"
        >
          <!-- Track Info -->
          <div class="flex items-end justify-between gap-4">
            <div class="min-w-0">
              <h1 class="text-3xl md:text-5xl font-bold truncate mb-2 drop-shadow-lg">
                {{ currentTrack?.title }}
              </h1>
              <div class="flex items-center gap-3 text-xl md:text-2xl text-muted-foreground">
                <span class="truncate">{{ currentTrack?.artist }}</span>
                <span v-if="currentTrack?.album" class="text-muted-foreground/60">•</span>
                <span v-if="currentTrack?.album" class="truncate text-muted-foreground/80">{{
                  currentTrack?.album
                }}</span>
              </div>
            </div>
            <button
              @click="currentTrack && toggleFavorite(currentTrack.file_path)"
              :class="[
                'w-12 h-12 rounded-full glass flex items-center justify-center transition-colors flex-shrink-0',
                currentTrack && isFavorite(currentTrack.file_path)
                  ? 'text-red-500'
                  : 'text-muted-foreground hover:text-foreground',
              ]"
            >
              <Heart
                class="w-6 h-6"
                :fill="currentTrack && isFavorite(currentTrack.file_path) ? 'currentColor' : 'none'"
              />
            </button>
          </div>

          <!-- Progress Bar -->
          <div class="w-full space-y-2">
            <Slider
              :model-value="progress"
              :max="100"
              :step="0.1"
              class="w-full h-2"
              @update:model-value="onSeek"
            />
            <div
              class="flex justify-between text-sm text-muted-foreground font-medium tabular-nums"
            >
              <span>{{ formatTime(currentTime) }}</span>
              <span>{{ formatTime(duration) }}</span>
            </div>
          </div>

          <!-- Controls -->
          <div class="flex items-center justify-between">
            <!-- Left Controls (Shuffle/Repeat) -->
            <div class="flex items-center gap-4">
              <button
                @click="emit('toggleShuffle')"
                :class="[
                  'w-10 h-10 flex items-center justify-center rounded-full transition-colors',
                  shuffle ? 'text-primary' : 'text-muted-foreground hover:text-foreground',
                ]"
              >
                <Shuffle class="w-5 h-5" />
              </button>
              <button
                @click="emit('cycleRepeat')"
                :class="[
                  'w-10 h-10 flex items-center justify-center rounded-full transition-colors',
                  repeat !== 'off' ? 'text-primary' : 'text-muted-foreground hover:text-foreground',
                ]"
              >
                <component :is="repeatIcon" class="w-5 h-5" />
              </button>

              <button
                @click="emit('changeVisualizer')"
                class="w-10 h-10 flex items-center justify-center rounded-full text-muted-foreground hover:text-foreground transition-colors"
                title="Change Visualizer"
              >
                <Activity class="w-5 h-5" />
              </button>
            </div>

            <!-- Center Controls (Play/Prev/Next) -->
            <div class="flex items-center gap-6">
              <button
                @click="emit('previous')"
                class="w-12 h-12 flex items-center justify-center rounded-full text-foreground hover:bg-white/10 transition-colors"
              >
                <SkipBack class="w-6 h-6" fill="currentColor" />
              </button>

              <button
                @click="emit('togglePlayPause')"
                class="w-16 h-16 flex items-center justify-center rounded-full bg-primary text-primary-foreground hover:scale-105 active:scale-95 transition-all shadow-xl"
              >
                <Play v-if="!isPlaying" class="w-8 h-8 ml-1" fill="currentColor" />
                <Pause v-else class="w-8 h-8" fill="currentColor" />
              </button>

              <button
                @click="emit('next')"
                class="w-12 h-12 flex items-center justify-center rounded-full text-foreground hover:bg-white/10 transition-colors"
              >
                <SkipForward class="w-6 h-6" fill="currentColor" />
              </button>
            </div>

            <!-- Right Controls (Volume) -->
            <div class="flex items-center gap-3 w-48 justify-end">
              <button
                @click="emit('toggleMute')"
                class="w-10 h-10 flex items-center justify-center rounded-full text-muted-foreground hover:text-foreground transition-colors"
              >
                <component :is="volumeIcon" class="w-5 h-5" />
              </button>

              <Slider
                :model-value="volumePercent"
                :max="100"
                :step="1"
                class="flex-1"
                @update:model-value="onVolumeChange"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.glass {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
