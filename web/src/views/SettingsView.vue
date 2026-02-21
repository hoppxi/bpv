<script setup lang="ts">
import { ref, computed } from "vue";
import {
  Settings,
  Palette,
  Music2,
  Volume2,
  Keyboard,
  HardDrive,
  RefreshCw,
  Info,
  Sliders,
  Eye,
  Gauge,
} from "lucide-vue-next";
import type { VisualizerType, EqSettings } from "@/types";
import Switch from "@/components/ui/Switch.vue";
import Slider from "@/components/ui/Slider.vue";
import Button from "@/components/ui/Button.vue";
import Separator from "@/components/ui/Separator.vue";
import ScrollArea from "@/components/ui/ScrollArea.vue";
import Badge from "@/components/ui/Badge.vue";

const props = defineProps<{
  visualizerType: VisualizerType;
  eqSettings: EqSettings;
  playbackSpeed: number;
  crossfade: boolean;
  gapless: boolean;
  showVisualizer: boolean;
  autoPlay: boolean;
  libraryPath: string;
  totalTracks: number;
}>();

const emit = defineEmits<{
  "update:visualizerType": [type: VisualizerType];
  "update:eqSettings": [settings: EqSettings];
  "update:playbackSpeed": [speed: number];
  "update:crossfade": [enabled: boolean];
  "update:gapless": [enabled: boolean];
  "update:showVisualizer": [show: boolean];
  "update:autoPlay": [auto: boolean];
  rescanLibrary: [];
}>();

const visualizerTypes: { value: VisualizerType; label: string }[] = [
  { value: "none", label: "None" },
  { value: "bars", label: "Bars" },
  { value: "wave", label: "Waveform" },
  { value: "particles", label: "Particles" },
  { value: "circle", label: "Circle" },
  { value: "sphere", label: "Sphere" },
  { value: "lines", label: "Lines" },
  { value: "mesh", label: "Mesh" },
  { value: "radial", label: "Radial" },
  { value: "spectrum", label: "Spectrum" },
  { value: "orb", label: "Orb" },
  { value: "galaxy", label: "Galaxy" },
  { value: "dna", label: "DNA" },
  { value: "aurora", label: "Aurora" },
  { value: "terrain", label: "Terrain" },
  { value: "retroBars", label: "retroBars" },
  { value: "sunburst", label: "sunburst" },
  { value: "hexagons", label: "hexagons" },
  { value: "blocks", label: "blocks" },
  { value: "spiral", label: "spiral" },
  { value: "tunnel", label: "tunnel" },
  { value: "flower", label: "flower" },
  { value: "neonGrid", label: "neonGrid" },
  { value: "kaleidoscope", label: "kaleidoscope" },
  { value: "drops", label: "drops" },
  { value: "rings", label: "rings" },
  { value: "segmentedBars", label: "segmentedBars" },
  { value: "seismic", label: "seismic" },
  { value: "pixels", label: "pixels" },
  { value: "lightning", label: "lightning" },
  { value: "polarWave", label: "polarWave" },
  { value: "confetti", label: "confetti" },
  { value: "glitch", label: "glitch" },
  { value: "infinity", label: "infinity" },
  { value: "rain", label: "rain" },
];

const speedPresets = [0.5, 0.75, 1, 1.25, 1.5, 2];

const activeSection = ref("general");

const sections = [
  { id: "general", label: "General", icon: Settings },
  { id: "audio", label: "Audio & EQ", icon: Sliders },
  { id: "visualizer", label: "Visualizer", icon: Eye },
  { id: "playback", label: "Playback", icon: Gauge },
  { id: "library", label: "Library", icon: HardDrive },
  { id: "shortcuts", label: "Shortcuts", icon: Keyboard },
  { id: "about", label: "About", icon: Info },
];

const shortcuts = [
  { key: "Space", action: "Play / Pause" },
  { key: "→", action: "Seek forward 10s" },
  { key: "←", action: "Seek backward 10s" },
  { key: "Shift+→", action: "Next track" },
  { key: "Shift+←", action: "Previous track" },
  { key: "↑", action: "Volume up" },
  { key: "↓", action: "Volume down" },
  { key: "M", action: "Toggle mute" },
  { key: "S", action: "Toggle shuffle" },
  { key: "R", action: "Toggle repeat" },
  { key: "/", action: "Open search" },
  { key: "N", action: "Next track" },
  { key: "P", action: "Previous track" },
];

function updateEq(field: keyof EqSettings, value: number | boolean) {
  emit("update:eqSettings", { ...props.eqSettings, [field]: value });
}
</script>

<template>
  <div class="flex h-full">
    <nav class="w-56 border-r border-border/50 p-3 space-y-1">
      <button
        v-for="section in sections"
        :key="section.id"
        @click="activeSection = section.id"
        :class="[
          'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all',
          activeSection === section.id
            ? 'bg-primary/15 text-primary'
            : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground',
        ]"
      >
        <component :is="section.icon" class="w-4 h-4" />
        {{ section.label }}
      </button>
    </nav>

    <ScrollArea class="flex-1">
      <div class="p-6 max-w-2xl">
        <div v-if="activeSection === 'general'" class="space-y-6">
          <h3 class="text-xl font-semibold">General Settings</h3>

          <div class="glass rounded-xl p-4 space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium">Auto-play on startup</p>
                <p class="text-xs text-muted-foreground">
                  Automatically play last track when app loads
                </p>
              </div>
              <Switch
                :model-value="autoPlay"
                @update:model-value="emit('update:autoPlay', $event)"
              />
            </div>
            <Separator />
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium">Show Visualizer</p>
                <p class="text-xs text-muted-foreground">Display audio visualizer in background</p>
              </div>
              <Switch
                :model-value="showVisualizer"
                @update:model-value="emit('update:showVisualizer', $event)"
              />
            </div>
          </div>
        </div>

        <div v-if="activeSection === 'audio'" class="space-y-6">
          <h3 class="text-xl font-semibold">Audio & Equalizer</h3>

          <div class="glass rounded-xl p-4 space-y-4">
            <div class="flex items-center justify-between mb-2">
              <p class="text-sm font-medium">Equalizer</p>
              <Switch
                :model-value="eqSettings.enabled"
                @update:model-value="updateEq('enabled', $event)"
              />
            </div>

            <div
              :class="[
                'space-y-4 transition-opacity',
                !eqSettings.enabled && 'opacity-40 pointer-events-none',
              ]"
            >
              <div>
                <div class="flex justify-between text-xs text-muted-foreground mb-1">
                  <span>Bass</span>
                  <span>{{ eqSettings.bass > 0 ? "+" : "" }}{{ eqSettings.bass }} dB</span>
                </div>
                <Slider
                  :model-value="eqSettings.bass + 12"
                  :max="24"
                  :step="1"
                  @update:model-value="updateEq('bass', $event - 12)"
                />
              </div>
              <div>
                <div class="flex justify-between text-xs text-muted-foreground mb-1">
                  <span>Mid</span>
                  <span>{{ eqSettings.mid > 0 ? "+" : "" }}{{ eqSettings.mid }} dB</span>
                </div>
                <Slider
                  :model-value="eqSettings.mid + 12"
                  :max="24"
                  :step="1"
                  @update:model-value="updateEq('mid', $event - 12)"
                />
              </div>
              <div>
                <div class="flex justify-between text-xs text-muted-foreground mb-1">
                  <span>Treble</span>
                  <span>{{ eqSettings.treble > 0 ? "+" : "" }}{{ eqSettings.treble }} dB</span>
                </div>
                <Slider
                  :model-value="eqSettings.treble + 12"
                  :max="24"
                  :step="1"
                  @update:model-value="updateEq('treble', $event - 12)"
                />
              </div>
              <Button
                variant="outline"
                size="sm"
                @click="emit('update:eqSettings', { bass: 0, mid: 0, treble: 0, enabled: true })"
              >
                Reset EQ
              </Button>
            </div>
          </div>
        </div>

        <div v-if="activeSection === 'visualizer'" class="space-y-6">
          <h3 class="text-xl font-semibold">Visualizer</h3>

          <div class="glass rounded-xl p-4">
            <p class="text-sm font-medium mb-3">Visualizer Style</p>
            <div class="grid grid-cols-3 gap-2">
              <button
                v-for="viz in visualizerTypes"
                :key="viz.value"
                @click="emit('update:visualizerType', viz.value)"
                :class="[
                  'px-3 py-2 rounded-lg text-sm font-medium transition-all border',
                  visualizerType === viz.value
                    ? 'bg-primary text-primary-foreground border-primary shadow-md'
                    : 'bg-secondary/50 text-muted-foreground border-transparent hover:bg-secondary hover:text-foreground',
                ]"
              >
                {{ viz.label }}
              </button>
            </div>
          </div>
        </div>

        <div v-if="activeSection === 'playback'" class="space-y-6">
          <h3 class="text-xl font-semibold">Playback</h3>

          <div class="glass rounded-xl p-4 space-y-4">
            <div>
              <p class="text-sm font-medium mb-3">Playback Speed</p>
              <div class="flex gap-2">
                <button
                  v-for="speed in speedPresets"
                  :key="speed"
                  @click="emit('update:playbackSpeed', speed)"
                  :class="[
                    'px-3 py-1.5 rounded-lg text-sm font-medium transition-all',
                    playbackSpeed === speed
                      ? 'bg-primary text-primary-foreground'
                      : 'bg-secondary/50 text-muted-foreground hover:bg-secondary',
                  ]"
                >
                  {{ speed }}x
                </button>
              </div>
            </div>

            <Separator />

            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium">Crossfade</p>
                <p class="text-xs text-muted-foreground">Smooth transition between tracks</p>
              </div>
              <Switch
                :model-value="crossfade"
                @update:model-value="emit('update:crossfade', $event)"
              />
            </div>

            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium">Gapless Playback</p>
                <p class="text-xs text-muted-foreground">Remove silence between tracks</p>
              </div>
              <Switch :model-value="gapless" @update:model-value="emit('update:gapless', $event)" />
            </div>
          </div>
        </div>

        <div v-if="activeSection === 'library'" class="space-y-6">
          <h3 class="text-xl font-semibold">Library</h3>

          <div class="glass rounded-xl p-4 space-y-4">
            <div>
              <p class="text-sm font-medium">Music Directory</p>
              <p class="text-xs text-muted-foreground mt-1 font-mono bg-secondary/50 p-2 rounded">
                {{ libraryPath || "Not set" }}
              </p>
            </div>
            <Separator />
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium">Total Tracks</p>
                <p class="text-xs text-muted-foreground">Files indexed in library</p>
              </div>
              <Badge variant="secondary">{{ totalTracks }}</Badge>
            </div>
            <Separator />
            <Button variant="outline" @click="emit('rescanLibrary')">
              <RefreshCw class="w-4 h-4 mr-2" />
              Rescan Library
            </Button>
          </div>
        </div>

        <div v-if="activeSection === 'shortcuts'" class="space-y-6">
          <h3 class="text-xl font-semibold">Keyboard Shortcuts</h3>

          <div class="glass rounded-xl p-4">
            <div class="space-y-2">
              <div
                v-for="shortcut in shortcuts"
                :key="shortcut.key"
                class="flex items-center justify-between py-2"
              >
                <span class="text-sm">{{ shortcut.action }}</span>
                <kbd
                  class="px-2 py-1 rounded bg-secondary text-xs font-mono text-muted-foreground"
                  >{{ shortcut.key }}</kbd
                >
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeSection === 'about'" class="space-y-6">
          <h3 class="text-xl font-semibold">About BPV</h3>

          <div class="glass rounded-xl p-6 text-center space-y-3">
            <div
              class="w-16 h-16 mx-auto rounded-2xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center shadow-lg"
            >
              <Music2 class="w-8 h-8 text-white" />
            </div>
            <h4 class="text-lg font-bold">BPV Music Player</h4>
            <p class="text-sm text-muted-foreground">Version 0.2.0</p>
            <p class="text-xs text-muted-foreground max-w-sm mx-auto">
              A beautiful, browser-based music player built with Vue 3 + Go. Inspired by mpv,
              designed for your local music library.
            </p>
            <div class="flex gap-2 justify-center mt-4">
              <Badge>Vue 3</Badge>
              <Badge>Go</Badge>
              <Badge>Tailwind CSS</Badge>
              <Badge>Web Audio API</Badge>
            </div>
          </div>
        </div>
      </div>
    </ScrollArea>
  </div>
</template>
