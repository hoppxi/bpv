<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from "vue";
import type { VisualizerType } from "@/types";
import { renderers } from "./renderers";

const props = defineProps<{
  frequencyData: Uint8Array;
  timeDomainData: Uint8Array;
  visualizerType: VisualizerType;
  colorPalette: string[];
  dominantColor: string;
  isPlaying: boolean;
  fullscreen?: boolean;
}>();

const canvasRef = ref<HTMLCanvasElement | null>(null);
const dimensions = ref({ width: 0, height: 0 });

function updateDimensions() {
  dimensions.value = {
    width: window.innerWidth,
    height: window.innerHeight,
  };
}

onMounted(() => {
  updateDimensions();
  window.addEventListener("resize", updateDimensions);
});

onUnmounted(() => {
  window.removeEventListener("resize", updateDimensions);
});

let animationId: number | null = null;

function render() {
  const canvas = canvasRef.value;
  if (!canvas) return;

  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  ctx.clearRect(0, 0, canvas.width, canvas.height);

  if (props.visualizerType === "none" || !props.isPlaying) return;

  const renderFn = renderers[props.visualizerType];
  if (!renderFn) return;

  const data = props.visualizerType === "wave" ? props.timeDomainData : props.frequencyData;
  renderFn(ctx, data, canvas.width, canvas.height, props.colorPalette, props.dominantColor);

  animationId = requestAnimationFrame(render);
}

watch(
  () => [props.isPlaying, props.visualizerType, props.frequencyData],
  () => {
    if (animationId) cancelAnimationFrame(animationId);
    if (props.isPlaying && props.visualizerType !== "none") {
      render();
    }
  }
);

onUnmounted(() => {
  if (animationId) cancelAnimationFrame(animationId);
});
</script>

<template>
  <canvas
    ref="canvasRef"
    :class="[
      'fixed inset-0 pointer-events-none z-0',
      fullscreen ? 'z-10' : ''
    ]"
    :width="dimensions.width"
    :height="dimensions.height"
  />
</template>
