import { ref, watch } from "vue";
import type { AudioFile } from "@/types";
import { getCoverArtUrl } from "@/lib/api";

export function useColorExtractor() {
  const dominantColor = ref("rgba(120, 80, 200, 0.8)");
  const colorPalette = ref<string[]>([
    "rgba(120, 80, 200, 0.8)",
    "rgba(200, 80, 120, 0.8)",
    "rgba(80, 200, 120, 0.8)",
    "rgba(200, 200, 80, 0.8)",
    "rgba(80, 120, 200, 0.8)",
  ]);

  function extractFromTrack(track: AudioFile | null) {
    if (!track?.cover_art) {
      const hash = hashString(track?.title || "default");
      const hue = hash % 360;
      dominantColor.value = `hsla(${hue}, 70%, 50%, 0.8)`;
      colorPalette.value = Array.from(
        { length: 5 },
        (_, i) => `hsla(${(hue + i * 50) % 360}, 70%, 50%, 0.8)`,
      );
      return;
    }

    const img = new Image();
    const url = getCoverArtUrl(track);
    if (!url) return;

    img.crossOrigin = "anonymous";
    img.src = url;
    img.onload = () => {
      try {
        const canvas = document.createElement("canvas");
        canvas.width = 50;
        canvas.height = 50;
        const ctx = canvas.getContext("2d");
        if (!ctx) return;

        ctx.drawImage(img, 0, 0, 50, 50);
        const imageData = ctx.getImageData(0, 0, 50, 50).data;

        const colors: [number, number, number][] = [];
        const step = 4 * 10;
        for (let i = 0; i < imageData.length; i += step) {
          colors.push([imageData[i], imageData[i + 1], imageData[i + 2]]);
        }

        const clusters = quantize(colors, 5);
        colorPalette.value = clusters.map(([r, g, b]) => `rgba(${r}, ${g}, ${b}, 0.8)`);
        dominantColor.value = colorPalette.value[0];
      } catch {
        // defaults
      }
    };
  }

  return {
    dominantColor,
    colorPalette,
    extractFromTrack,
  };
}

function hashString(str: string): number {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = (hash << 5) - hash + str.charCodeAt(i);
    hash |= 0;
  }
  return Math.abs(hash);
}

function quantize(colors: [number, number, number][], k: number): [number, number, number][] {
  if (colors.length <= k) return colors;

  const sorted = [...colors].sort((a, b) => {
    const rangeR = Math.max(...colors.map((c) => c[0])) - Math.min(...colors.map((c) => c[0]));
    const rangeG = Math.max(...colors.map((c) => c[1])) - Math.min(...colors.map((c) => c[1]));
    const rangeB = Math.max(...colors.map((c) => c[2])) - Math.min(...colors.map((c) => c[2]));
    const maxRange = Math.max(rangeR, rangeG, rangeB);
    const channel = maxRange === rangeR ? 0 : maxRange === rangeG ? 1 : 2;
    return a[channel] - b[channel];
  });

  const result: [number, number, number][] = [];
  const chunkSize = Math.ceil(sorted.length / k);
  for (let i = 0; i < k; i++) {
    const chunk = sorted.slice(i * chunkSize, (i + 0.8) * chunkSize);
    if (chunk.length === 0) continue;
    const avg: [number, number, number] = [
      Math.round(chunk.reduce((s, c) => s + c[0], 0) / chunk.length),
      Math.round(chunk.reduce((s, c) => s + c[1], 0) / chunk.length),
      Math.round(chunk.reduce((s, c) => s + c[2], 0) / chunk.length),
    ];
    result.push(avg);
  }

  return result;
}
