<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from "vue";
import { X, Music2, ListMusic, Play, GripVertical } from "lucide-vue-next";
import type { AudioFile } from "@/types";
import { getCoverArtUrl } from "@/lib/api";
import { formatDuration } from "@/lib/utils";

const props = defineProps<{
  queue: AudioFile[];
  currentIndex: number;
  currentTrack: AudioFile | null;
  isOpen: boolean;
}>();

const emit = defineEmits<{
  close: [];
  playValues: [index: number];
  removeTrack: [index: number];
  clearQueue: [];
}>();

const containerRef = ref<HTMLElement | null>(null);

function scrollToCurrent() {
    if (!containerRef.value) return;
    const currentItem = containerRef.value.querySelector('.current-track');
    if (currentItem) {
        currentItem.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
}

// Watch for changes in current index to auto-scroll
// In a real implementation we would watch props.currentIndex

</script>

<template>
  <div
    class="fixed inset-y-0 right-0 z-50 w-full md:w-[400px] bg-background border-l border-border shadow-2xl transform transition-transform duration-300 ease-in-out flex flex-col"
    :class="isOpen ? 'translate-x-0' : 'translate-x-full'"
  >
    <!-- Header -->
    <div class="flex items-center justify-between p-4 border-b border-border bg-card/50 backdrop-blur-md">
      <div class="flex items-center gap-2">
        <ListMusic class="w-5 h-5 text-primary" />
        <h2 class="font-semibold text-lg">Play Queue</h2>
        <span class="text-xs text-muted-foreground ml-2 badge bg-secondary px-2 py-0.5 rounded-full">
            {{ queue.length }} tracks
        </span>
      </div>
      <div class="flex items-center gap-2">
          <button 
            @click="emit('clearQueue')"
            class="text-xs text-muted-foreground hover:text-red-500 transition-colors px-2 py-1"
            v-if="queue.length > 0"
          >
            Clear
          </button>
          <button
            @click="emit('close')"
            class="p-2 rounded-full hover:bg-secondary transition-colors"
          >
            <X class="w-5 h-5" />
          </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="queue.length === 0" class="flex-1 flex flex-col items-center justify-center text-muted-foreground p-8">
      <ListMusic class="w-16 h-16 mb-4 opacity-20" />
      <p class="text-center font-medium">Items you play will appear here</p>
      <p class="text-xs text-center mt-2 opacity-60">Add some music to start listening</p>
    </div>

    <!-- Queue List -->
    <div v-else class="flex-1 overflow-y-auto p-2 space-y-1" ref="containerRef">
        
        <!-- Now Playing Section -->
        <div v-if="currentTrack" class="mb-4 sticky top-0 z-10 bg-background/95 backdrop-blur pb-2 pt-2 px-2 border-b border-white/5">
             <h3 class="text-xs font-bold uppercase tracking-wider text-muted-foreground mb-2 px-2">Now Playing</h3>
            <div class="current-track flex items-center gap-3 p-3 rounded-xl bg-primary/10 border border-primary/20">
                <div class="relative w-12 h-12 rounded-lg overflow-hidden flex-shrink-0 shadow-sm">
                    <img
                        v-if="getCoverArtUrl(currentTrack)"
                        :src="getCoverArtUrl(currentTrack)!"
                        class="w-full h-full object-cover"
                    />
                     <div v-else class="w-full h-full flex items-center justify-center bg-secondary">
                        <Music2 class="w-5 h-5 text-muted-foreground" />
                    </div>
                     <div class="absolute inset-0 flex items-center justify-center bg-black/20">
                        <div class="flex space-x-1">
                            <div class="w-1 h-3 bg-white rounded-full animate-bounce" style="animation-delay: 0s"></div>
                            <div class="w-1 h-3 bg-white rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
                            <div class="w-1 h-3 bg-white rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
                        </div>
                    </div>
                </div>
                <div class="flex-1 min-w-0">
                    <p class="font-bold text-primary truncate">{{ currentTrack.title }}</p>
                    <p class="text-xs text-muted-foreground truncate">{{ currentTrack.artist }}</p>
                </div>
                 <span class="text-xs font-medium tabular-nums text-muted-foreground">
                    {{ formatDuration(currentTrack.duration) }}
                </span>
            </div>
        </div>

        <!-- Next Up Section -->
         <h3 class="text-xs font-bold uppercase tracking-wider text-muted-foreground mb-2 px-4 mt-6">Next Up</h3>
         
         <div 
            v-for="(track, index) in queue" 
            :key="`${track.file_path}-${index}`"
            v-show="index > currentIndex"
            @click="emit('playValues', index)"
            class="group flex items-center gap-3 p-2 rounded-lg hover:bg-white/5 transition-colors cursor-pointer border border-transparent"
             :class="index === currentIndex ? 'opacity-50 pointer-events-none hidden' : ''"
         >
            <div class="w-6 flex items-center justify-center text-xs text-muted-foreground opacity-50 group-hover:opacity-100">
                <span class="group-hover:hidden">{{ index + 1 }}</span>
                <Play class="w-3 h-3 hidden group-hover:block fill-current" />
            </div>

             <div class="w-10 h-10 rounded-md overflow-hidden bg-secondary flex-shrink-0">
                 <img
                    v-if="getCoverArtUrl(track)"
                    :src="getCoverArtUrl(track)!"
                    class="w-full h-full object-cover opacity-80 group-hover:opacity-100 transition-opacity"
                />
                <div v-else class="w-full h-full flex items-center justify-center">
                    <Music2 class="w-4 h-4 text-muted-foreground" />
                </div>
             </div>

             <div class="flex-1 min-w-0">
                 <p class="text-sm font-medium truncate text-foreground/90 group-hover:text-foreground">{{ track.title }}</p>
                 <p class="text-xs text-muted-foreground truncate group-hover:text-muted-foreground/80">{{ track.artist }}</p>
             </div>

             <span class="text-xs text-muted-foreground tabular-nums opacity-0 group-hover:opacity-100 transition-opacity">
                {{ formatDuration(track.duration) }}
             </span>
             
             <!-- Could add remove button here later -->
         </div>

    </div>
  </div>

  <!-- Backdrop -->
  <div
    v-if="isOpen"
    @click="emit('close')"
    class="fixed inset-0 bg-black/50 backdrop-blur-sm z-40 transition-opacity duration-300"
  />
</template>
