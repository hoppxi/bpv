export async function extractDominantColor(imageUrl: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = "Anonymous";

    img.onload = () => {
      try {
        const canvas = document.createElement("canvas");
        const ctx = canvas.getContext("2d");

        if (!ctx) {
          reject(new Error("Could not get canvas context"));
          return;
        }

        canvas.width = 100;
        canvas.height = 100;
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
        const data = imageData.data;
        const colorMap = new Map<string, number>();

        // Sample pixels (every 4th pixel for performance)
        for (let i = 0; i < data.length; i += 16) {
          const r = data[i];
          const g = data[i + 1];
          const b = data[i + 2];
          const a = data[i + 3];

          // Skip transparent pixels
          if (a < 128) continue;

          const color = `${r},${g},${b}`;
          colorMap.set(color, (colorMap.get(color) || 0) + 1);
        }

        // Find most frequent color
        let maxCount = 0;
        let dominantColor = "255,255,255"; // Default white

        colorMap.forEach((count, color) => {
          if (count > maxCount) {
            maxCount = count;
            dominantColor = color;
          }
        });

        const [r, g, b] = dominantColor.split(",").map(Number);
        resolve(`rgb(${r}, ${g}, ${b})`);
      } catch (error) {
        reject(error);
      }
    };

    img.onerror = () => {
      reject(new Error("Failed to load image"));
    };

    img.src = imageUrl;
  });
}

export function generateColorPalette(baseColor: string): string[] {
  const rgb = baseColor.match(/\d+/g)?.map(Number);
  if (!rgb || rgb.length < 3) {
    return [
      baseColor,
      adjustColorBrightness(baseColor, 0.8),
      adjustColorBrightness(baseColor, 0.6),
      adjustColorBrightness(baseColor, 0.4),
      adjustColorBrightness(baseColor, 0.2),
    ];
  }

  const [r, g, b] = rgb;
  const hsl = rgbToHsl(r, g, b);

  const palette: string[] = [];

  for (let i = 0; i < 5; i++) {
    const lightness = 0.9 - i * 0.15; // 0.9, 0.75, 0.6, 0.45, 0.3
    const color = hslToRgb(hsl[0], hsl[1], lightness);
    palette.push(`rgb(${color[0]}, ${color[1]}, ${color[2]})`);
  }

  return palette;
}

function adjustColorBrightness(color: string, factor: number): string {
  const rgb = color.match(/\d+/g)?.map(Number);
  if (!rgb || rgb.length < 3) return color;

  const [r, g, b] = rgb.map((c) => Math.round(c * factor));
  return `rgb(${r}, ${g}, ${b})`;
}

function rgbToHsl(r: number, g: number, b: number): [number, number, number] {
  r /= 255;
  g /= 255;
  b /= 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h = 0,
    s = 0,
    l = (max + min) / 2;

  if (max !== min) {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }

    h /= 6;
  }

  return [h, s, l];
}

function hslToRgb(h: number, s: number, l: number): [number, number, number] {
  let r, g, b;

  if (s === 0) {
    r = g = b = l; // achromatic
  } else {
    const hue2rgb = (p: number, q: number, t: number) => {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1 / 6) return p + (q - p) * 6 * t;
      if (t < 1 / 2) return q;
      if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
      return p;
    };

    const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    const p = 2 * l - q;

    r = hue2rgb(p, q, h + 1 / 3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1 / 3);
  }

  return [Math.round(r * 255), Math.round(g * 255), Math.round(b * 255)];
}

export function isDarkColor(color: string): boolean {
  const rgb = color.match(/\d+/g)?.map(Number);
  if (!rgb || rgb.length < 3) return false;

  const [r, g, b] = rgb;
  // Calculate luminance
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
  return luminance < 0.5;
}

export function getContrastColor(color: string): string {
  return isDarkColor(color) ? "#ffffff" : "#000000";
}
