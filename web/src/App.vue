<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";
import type { AudioFile, ViewMode, VisualizerType, RepeatMode, EqSettings } from "@/types";
import { triggerScan, fetchQueue, saveQueue, type QueueState } from "@/lib/api";
import { generateId } from "@/lib/utils";

import { useAudioPlayer } from "@/composables/useAudioPlayer";
import { useAudioVisualizer } from "@/composables/useAudioVisualizer";
import { useColorExtractor } from "@/composables/useColorExtractor";
import { useLibraryData } from "@/composables/useLibraryData";
import { useStorage } from "@/composables/useStorage";
import { useKeyboardShortcuts } from "@/composables/useKeyboardShortcuts";
import { useMediaSession } from "@/composables/useMediaSession";

import AppSidebar from "@/components/layout/AppSidebar.vue";
import PlayerBar from "@/components/layout/PlayerBar.vue";
import QueueSidebar from "@/components/layout/QueueSidebar.vue";

import VisualizerCanvas from "@/components/visualizer/VisualizerCanvas.vue";

import HomeView from "@/views/HomeView.vue";
import LibraryView from "@/views/LibraryView.vue";
import ArtistsView from "@/views/ArtistsView.vue";
import AlbumsView from "@/views/AlbumsView.vue";
import GenresView from "@/views/GenresView.vue";
import ComposersView from "@/views/ComposersView.vue";
import SearchView from "@/views/SearchView.vue";
import FavoritesView from "@/views/FavoritesView.vue";
import SettingsView from "@/views/SettingsView.vue";
import NowPlayingView from "@/views/NowPlayingView.vue";

const currentView = ref<ViewMode>("home");
const showNowPlaying = ref(false);
const showQueuePanel = ref(false);

const queue = ref<AudioFile[]>([]);
const queueIndex = ref(-1);
const shuffle = ref(false);
const repeat = ref<RepeatMode>("off");

const visualizerType = useStorage<VisualizerType>("bpv-visualizer-type", "bars");
const showVisualizer = useStorage<boolean>("bpv-show-visualizer", true);
const autoPlay = useStorage<boolean>("bpv-auto-play", false);
const crossfade = useStorage<boolean>("bpv-crossfade", false);
const gapless = useStorage<boolean>("bpv-gapless", true);
const eqSettings = useStorage<EqSettings>("bpv-eq", { bass: 0, mid: 0, treble: 0, enabled: false });
import { useFavorites } from "@/composables/useFavorites";

const { favorites, toggleFavorite, isFavorite } = useFavorites();

function removeFromQueue(index: number) {
  queue.value.splice(index, 1);
  if (index < queueIndex.value) {
    queueIndex.value--;
  } else if (index === queueIndex.value) {
    if (queue.value.length > 0) {
      if (queueIndex.value >= queue.value.length) {
        queueIndex.value = queue.value.length - 1;
      }
      playAudio(queue.value[queueIndex.value], basePath.value);
    } else {
      queueIndex.value = -1;
    }
  }
}

function clearQueue() {
  queue.value = [];
  queueIndex.value = -1;
}

function playQueueIndex(index: number) {
  if (index >= 0 && index < queue.value.length) {
    queueIndex.value = index;
    playAudio(queue.value[index], basePath.value);
  }
}

const {
  isPlaying,
  currentTime,
  duration,
  volume,
  muted,
  isLoading,
  currentTrack,
  playbackSpeed,
  errorMsg,
  playTrack: playAudio,
  togglePlayPause,
  seek,
  setVolume,
  toggleMute,
  setSpeed,
  skipForward,
  skipBackward,
  setEQ,
  onTrackEnd,
  getAnalyser,
  getAudioElement,
} = useAudioPlayer();

const { library, loading: libraryLoading, basePath, refreshLibrary } = useLibraryData();
const allTracks = computed(() => library.value?.files || []);

const {
  frequencyData,
  timeDomainData,
  start: startVisualizer,
  stop: stopVisualizer,
} = useAudioVisualizer(getAnalyser);

const { dominantColor, colorPalette, extractFromTrack } = useColorExtractor();

function buildQueue(tracks: AudioFile[], startTrack?: AudioFile) {
  if (shuffle.value) {
    const shuffled = [...tracks].sort(() => Math.random() - 0.5);
    if (startTrack) {
      const idx = shuffled.findIndex((t) => t.file_path === startTrack.file_path);
      if (idx > 0) {
        shuffled.splice(idx, 1);
        shuffled.unshift(startTrack);
      }
    }
    queue.value = shuffled;
  } else {
    queue.value = [...tracks];
  }

  if (startTrack) {
    queueIndex.value = queue.value.findIndex((t) => t.file_path === startTrack.file_path);
  } else {
    queueIndex.value = 0;
  }
}

async function handlePlayTrack(track: AudioFile, context?: AudioFile[]) {
  let tracks: AudioFile[] = [];

  if (context) {
    tracks = context;
  } else {
    switch (currentView.value) {
      case "library":
        tracks = allTracks.value;
        break;
      case "artists":
      case "albums":
      case "genres":
      case "composers":
      case "favorites":
      case "home":
      case "search":
      default:
        tracks = allTracks.value;
        break;
    }
  }

  buildQueue(tracks, track);
  await playAudio(track, basePath.value);
  extractFromTrack(track);
  startVisualizer();
  persistQueue();
}

async function playNext() {
  if (queue.value.length === 0) return;

  let nextIdx = queueIndex.value + 1;

  if (nextIdx >= queue.value.length) {
    if (repeat.value === "all") {
      nextIdx = 0;
    } else {
      return;
    }
  }

  queueIndex.value = nextIdx;
  const track = queue.value[nextIdx];
  await playAudio(track, basePath.value);
  extractFromTrack(track);
  persistQueue();
}

async function playPrevious() {
  if (queue.value.length === 0) return;

  if (currentTime.value > 3) {
    seek(0);
    return;
  }

  let prevIdx = queueIndex.value - 1;
  if (prevIdx < 0) {
    if (repeat.value === "all") {
      prevIdx = queue.value.length - 1;
    } else {
      prevIdx = 0;
    }
  }

  queueIndex.value = prevIdx;
  const track = queue.value[prevIdx];
  await playAudio(track, basePath.value);
  extractFromTrack(track);
  persistQueue();
}

onTrackEnd(() => {
  if (repeat.value === "one") {
    seek(0);
    const audio = getAudioElement();
    audio.play();
  } else {
    playNext();
  }
});

const isCurrentFavorite = computed(() => isFavorite(currentTrack.value?.file_path));

function toggleShuffle() {
  shuffle.value = !shuffle.value;
  if (shuffle.value && queue.value.length > 0) {
    buildQueue(queue.value, currentTrack.value || undefined);
  }
  persistQueue();
}

function cycleRepeat() {
  const modes: RepeatMode[] = ["off", "all", "one"];
  const idx = modes.indexOf(repeat.value);
  repeat.value = modes[(idx + 1) % modes.length];
  persistQueue();
}

watch(
  eqSettings,
  (settings) => {
    if (settings.enabled) {
      setEQ(settings.bass, settings.mid, settings.treble);
    } else {
      setEQ(0, 0, 0);
    }
  },
  { deep: true },
);

watch(playbackSpeed, (speed) => {
  setSpeed(speed);
});

useKeyboardShortcuts({
  togglePlayPause,
  nextTrack: playNext,
  previousTrack: playPrevious,
  volumeUp: () => setVolume(Math.min(1, volume.value + 0.05)),
  volumeDown: () => setVolume(Math.max(0, volume.value - 0.05)),
  toggleMute,
  seekForward: () => skipForward(10),
  seekBackward: () => skipBackward(10),
  toggleShuffle,
  toggleRepeat: cycleRepeat,
  openSearch: () => (currentView.value = "search"),
});

useMediaSession(currentTrack, {
  play: () => togglePlayPause(),
  pause: () => togglePlayPause(),
  nextTrack: playNext,
  previousTrack: playPrevious,
  seekTo: seek,
});

async function rescanLibrary() {
  await triggerScan();
  await refreshLibrary();
}

function persistQueue() {
  const state: QueueState = {
    file_paths: queue.value.map((t) => t.file_path),
    current_index: Math.max(0, Math.min(queueIndex.value, queue.value.length - 1)),
    shuffle: shuffle.value,
    repeat: repeat.value === "all" ? 1 : repeat.value === "one" ? 2 : 0,
  };
  saveQueue(state).catch(() => {});
}

async function restoreQueueFromDaemon() {
  const q = await fetchQueue().catch<QueueState | null>(() => null);
  if (!q || !q.file_paths || q.file_paths.length === 0) return;

  const byPath = new Map<string, AudioFile>();
  for (const t of allTracks.value) {
    byPath.set(t.file_path, t);
  }

  const tracks: AudioFile[] = [];
  for (const p of q.file_paths) {
    const t = byPath.get(p);
    if (t) tracks.push(t);
  }
  if (tracks.length === 0) return;

  queue.value = tracks;
  shuffle.value = q.shuffle;
  repeat.value = q.repeat === 1 ? "all" : q.repeat === 2 ? "one" : "off";

  const idx = q.current_index >= 0 && q.current_index < tracks.length ? q.current_index : 0;
  queueIndex.value = idx;

  if (autoPlay.value) {
    const track = queue.value[idx];
    await playAudio(track, basePath.value);
    extractFromTrack(track);
    startVisualizer();
  }
}

onMounted(async () => {
  watch(
    () => allTracks.value.length,
    async (len, oldLen) => {
      if (len > 0 && oldLen === 0) {
        await restoreQueueFromDaemon();
      }
    },
    { immediate: true },
  );
});

const viewTitles: Record<ViewMode, string> = {
  home: "Home",
  library: "Library",
  artists: "Artists",
  albums: "Albums",
  genres: "Genres",
  composers: "Composers",
  favorites: "Favorites",
  search: "Search",
  settings: "Settings",
  "now-playing": "Now Playing",
};

const visualizerTypes: VisualizerType[] = [
  "bars",
  "wave",
  "particles",
  "circle",
  "sphere",
  "lines",
  "mesh",
  "radial",
  "spectrum",
  "orb",
  "galaxy",
  "dna",
  "aurora",
  "terrain",
  "retroBars",
  "sunburst",
  "hexagons",
  "blocks",
  "spiral",
  "tunnel",
  "flower",
  "neonGrid",
  "kaleidoscope",
  "drops",
  "rings",
  "segmentedBars",
  "seismic",
  "pixels",
  "lightning",
  "polarWave",
  "confetti",
  "glitch",
  "infinity",
  "rain",
  "none",
];

function cycleVisualizer() {
  const currentIndex = visualizerTypes.indexOf(visualizerType.value);
  const nextIndex = (currentIndex + 1) % visualizerTypes.length;
  visualizerType.value = visualizerTypes[nextIndex];
}
</script>

<template>
  <div class="h-screen flex flex-col relative overflow-hidden">
    <VisualizerCanvas
      v-if="showVisualizer && visualizerType !== 'none'"
      :frequency-data="frequencyData"
      :time-domain-data="timeDomainData"
      :visualizer-type="visualizerType"
      :color-palette="colorPalette"
      :dominant-color="dominantColor"
      :is-playing="isPlaying"
    />

    <div
      class="fixed inset-0 pointer-events-none z-0 opacity-20 transition-all duration-1000"
      :style="{
        background: `radial-gradient(circle at 30% 70%, ${dominantColor}, transparent 70%),
                     radial-gradient(circle at 70% 30%, ${colorPalette[2] || dominantColor}, transparent 60%)`,
      }"
    />

    <QueueSidebar
      :is-open="showQueuePanel"
      :queue="queue"
      :current-index="queueIndex"
      :current-track="currentTrack"
      @close="showQueuePanel = false"
      @play-values="playQueueIndex"
      @remove-track="removeFromQueue"
      @clear-queue="clearQueue"
    />

    <AppSidebar :current-view="currentView" @navigate="currentView = $event" />

    <main
      class="ml-[var(--sidebar-width)] mb-[var(--player-height)] h-full overflow-auto relative z-10"
    >
      <div v-if="libraryLoading" class="flex items-center justify-center h-full">
        <div class="text-center space-y-4">
          <div
            class="w-12 h-12 border-3 border-primary border-t-transparent rounded-full animate-spin mx-auto"
          />
          <p class="text-muted-foreground">Loading your music library...</p>
        </div>
      </div>

      <template v-else>
        <HomeView
          v-if="currentView === 'home'"
          :library="library"
          :current-track="currentTrack"
          @play-track="handlePlayTrack"
          @navigate="currentView = $event as ViewMode"
        />

        <LibraryView
          v-else-if="currentView === 'library'"
          :tracks="allTracks"
          :current-track="currentTrack"
          :is-playing="isPlaying"
          :favorites="favorites"
          @play-track="handlePlayTrack"
        />

        <ArtistsView
          v-else-if="currentView === 'artists'"
          :tracks="allTracks"
          :current-track="currentTrack"
          :is-playing="isPlaying"
          @play-track="handlePlayTrack"
        />

        <AlbumsView
          v-else-if="currentView === 'albums'"
          :tracks="allTracks"
          :current-track="currentTrack"
          :is-playing="isPlaying"
          @play-track="handlePlayTrack"
        />

        <GenresView
          v-else-if="currentView === 'genres'"
          :tracks="allTracks"
          :current-track="currentTrack"
          :is-playing="isPlaying"
          @play-track="handlePlayTrack"
        />

        <ComposersView
          v-else-if="currentView === 'composers'"
          :tracks="allTracks"
          :current-track="currentTrack"
          :is-playing="isPlaying"
          @play-track="handlePlayTrack"
        />

        <SearchView
          v-else-if="currentView === 'search'"
          :tracks="allTracks"
          :current-track="currentTrack"
          :is-playing="isPlaying"
          @play-track="handlePlayTrack"
        />

        <FavoritesView
          v-else-if="currentView === 'favorites'"
          :tracks="allTracks"
          :favorites="favorites"
          :current-track="currentTrack"
          :is-playing="isPlaying"
          @play-track="handlePlayTrack"
          @toggle-favorite="toggleFavorite"
        />

        <SettingsView
          v-else-if="currentView === 'settings'"
          :visualizer-type="visualizerType"
          :eq-settings="eqSettings"
          :playback-speed="playbackSpeed"
          :crossfade="crossfade"
          :gapless="gapless"
          :show-visualizer="showVisualizer"
          :auto-play="autoPlay"
          :library-path="basePath"
          :total-tracks="allTracks.length"
          @update:visualizer-type="visualizerType = $event"
          @update:eq-settings="eqSettings = $event"
          @update:playback-speed="
            playbackSpeed = $event;
            setSpeed($event);
          "
          @update:crossfade="crossfade = $event"
          @update:gapless="gapless = $event"
          @update:show-visualizer="showVisualizer = $event"
          @update:auto-play="autoPlay = $event"
          @rescan-library="rescanLibrary"
        />
      </template>
    </main>

    <PlayerBar
      :current-track="currentTrack"
      :is-playing="isPlaying"
      :is-loading="isLoading"
      :current-time="currentTime"
      :duration="duration"
      :volume="volume"
      :muted="muted"
      :shuffle="shuffle"
      :repeat="repeat"
      :is-favorite="isCurrentFavorite"
      @toggle-play-pause="togglePlayPause"
      @next="playNext"
      @previous="playPrevious"
      @seek="seek"
      @volume-change="setVolume"
      @toggle-mute="toggleMute"
      @toggle-shuffle="toggleShuffle"
      @cycle-repeat="cycleRepeat"
      @open-queue="showQueuePanel = !showQueuePanel"
      @open-now-playing="showNowPlaying = true"
      @change-visualizer="cycleVisualizer"
    />

    <Transition name="slide-up">
      <NowPlayingView
        v-if="showNowPlaying"
        :current-track="currentTrack"
        :is-playing="isPlaying"
        :is-loading="isLoading"
        :current-time="currentTime"
        :duration="duration"
        :volume="volume"
        :muted="muted"
        :shuffle="shuffle"
        :repeat="repeat"
        :is-favorite="isCurrentFavorite"
        :frequency-data="frequencyData"
        :time-domain-data="timeDomainData"
        :visualizer-type="visualizerType"
        :color-palette="colorPalette"
        :dominant-color="dominantColor"
        @close="showNowPlaying = false"
        @toggle-play-pause="togglePlayPause"
        @next="playNext"
        @previous="playPrevious"
        @seek="seek"
        @volume-change="setVolume"
        @toggle-mute="toggleMute"
        @toggle-shuffle="toggleShuffle"
        @cycle-repeat="cycleRepeat"
        @change-visualizer="cycleVisualizer"
      />
    </Transition>
  </div>
</template>
