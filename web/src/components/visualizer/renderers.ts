import type { VisualizerType } from "@/types";

export type RendererFn = (
  ctx: CanvasRenderingContext2D,
  data: Uint8Array,
  width: number,
  height: number,
  palette: string[],
  dominantColor: string,
) => void;

export const renderBars: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const barWidth = (width / data.length) * 7.5;
  let x = 0;
  for (let i = 0; i < data.length; i++) {
    const barHeight = (data[i] / 255) * height * 0.8;
    const colorIndex = Math.floor((i / data.length) * palette.length);
    const gradient = ctx.createLinearGradient(x, height - barHeight, x, height);
    gradient.addColorStop(0, palette[colorIndex] || dominantColor);
    gradient.addColorStop(1, "rgba(255, 255, 255, 0.1)");
    ctx.fillStyle = gradient;
    ctx.fillRect(x, height - barHeight, barWidth, barHeight);
    x += barWidth + 1;
  }
};

export const renderWave: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  ctx.beginPath();
  ctx.lineWidth = 3;
  ctx.lineCap = "round";
  ctx.lineJoin = "round";
  const gradient = ctx.createLinearGradient(0, 0, width, 0);
  gradient.addColorStop(0, palette[0] || dominantColor);
  gradient.addColorStop(0.5, palette[Math.floor(palette.length / 2)] || dominantColor);
  gradient.addColorStop(1, palette[palette.length - 1] || dominantColor);
  ctx.strokeStyle = gradient;
  ctx.shadowBlur = 15;
  ctx.shadowColor = dominantColor;
  const sliceWidth = width / data.length;
  let x = 0;
  for (let i = 0; i < data.length; i++) {
    const v = data[i] / 128.0 - 1;
    const y = (v * height) / 3 + height / 2;
    if (i === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
    x += sliceWidth;
  }
  ctx.stroke();
  ctx.shadowBlur = 0;
};

export const renderParticles: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const maxRadius = Math.min(width, height) * 0.45;
  for (let i = 0; i < data.length; i += 2) {
    const angle = (i / data.length) * Math.PI * 2;
    const idx = Math.min(i + 16, data.length - 1);
    const idx2 = Math.min(i + 32, data.length - 1);
    const radius = ((data[i] + data[idx] + data[idx2]) / (3 * 255)) * maxRadius + 30;
    const x = centerX + Math.cos(angle) * radius;
    const y = centerY + Math.sin(angle) * radius;
    const ri = Math.abs(data.length - 1 - i);
    const size = ((data[i] + data[ri]) / (2 * 255)) * 10 + 2;
    const colorIndex = i % palette.length;
    ctx.shadowBlur = 15;
    ctx.shadowColor = palette[colorIndex] || dominantColor;
    ctx.beginPath();
    ctx.fillStyle = palette[colorIndex] || dominantColor;
    ctx.arc(x, y, size, 0, Math.PI * 2);
    ctx.fill();
    ctx.shadowBlur = 0;
  }
};

export const renderCircle: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const maxRadius = Math.min(width, height) * 0.4;
  ctx.beginPath();
  ctx.lineWidth = 3;
  for (let i = 0; i < data.length; i++) {
    const angle = (i / data.length) * Math.PI * 2;
    const radius = (data[i] / 255) * maxRadius + 20;
    const x = centerX + Math.cos(angle) * radius;
    const y = centerY + Math.sin(angle) * radius;
    const colorIndex = Math.floor((i / data.length) * palette.length);
    ctx.strokeStyle = palette[colorIndex] || dominantColor;
    if (i === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
  }
  ctx.closePath();
  ctx.stroke();
};

export const renderSphere: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const baseRadius = Math.min(width, height) * 0.3;
  const pulse = (data[0] / 255) * 0.5 + 0.5;
  const radius = baseRadius * pulse;
  const gradient = ctx.createRadialGradient(centerX, centerY, 0, centerX, centerY, radius * 1.5);
  gradient.addColorStop(0, palette[0] || dominantColor);
  gradient.addColorStop(1, "rgba(0,0,0,0)");
  ctx.beginPath();
  ctx.fillStyle = gradient;
  ctx.arc(centerX, centerY, radius, 0, Math.PI * 2);
  ctx.fill();
};

export const renderLines: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const lineCount = 50;
  const segmentWidth = width / lineCount;
  for (let i = 0; i < lineCount; i++) {
    const amplitude = (data[i % data.length] / 255) * height * 0.4;
    const x = i * segmentWidth;
    const colorIndex = i % palette.length;
    ctx.beginPath();
    ctx.strokeStyle = palette[colorIndex] || dominantColor;
    ctx.lineWidth = 2;
    ctx.moveTo(x, height / 2 - amplitude);
    ctx.lineTo(x, height / 2 + amplitude);
    ctx.stroke();
  }
};

export const renderMesh: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const gridSize = 20;
  const cellWidth = width / gridSize;
  const cellHeight = height / gridSize;
  for (let x = 0; x < gridSize; x++) {
    for (let y = 0; y < gridSize; y++) {
      const index = (x + y) % data.length;
      const size = (data[index] / 255) * cellWidth * 0.8;
      const colorIndex = (x + y) % palette.length;
      ctx.fillStyle = palette[colorIndex] || dominantColor;
      ctx.fillRect(
        x * cellWidth + (cellWidth - size) / 2,
        y * cellHeight + (cellHeight - size) / 2,
        size,
        size,
      );
    }
  }
};

export const renderRadial: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const maxRadius = Math.min(width, height) * 0.4;
  for (let i = 0; i < data.length; i++) {
    const angle = (i / data.length) * Math.PI * 2;
    const radius = (data[i] / 255) * maxRadius;
    const colorIndex = Math.floor((i / data.length) * palette.length);
    ctx.beginPath();
    ctx.strokeStyle = palette[colorIndex] || dominantColor;
    ctx.lineWidth = 2;
    ctx.moveTo(centerX, centerY);
    ctx.lineTo(centerX + Math.cos(angle) * radius, centerY + Math.sin(angle) * radius);
    ctx.stroke();
  }
};

export const renderSpectrum: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const barCount = 100;
  const barWidth = width / barCount;
  for (let i = 0; i < barCount; i++) {
    const value = data[Math.floor((i / barCount) * data.length)];
    const barHeight = (value / 255) * height * 0.8;
    const x = i * barWidth;
    const colorIndex = Math.floor((i / barCount) * palette.length);
    const gradient = ctx.createLinearGradient(x, height - barHeight, x, height);
    gradient.addColorStop(0, palette[colorIndex] || dominantColor);
    gradient.addColorStop(1, "rgba(255,255,255,0.1)");
    ctx.fillStyle = gradient;
    ctx.fillRect(x, height - barHeight, barWidth - 1, barHeight);
  }
};

export const renderOrb: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const baseRadius = Math.min(width, height) * 0.2;
  const pulse = (data[0] / 255) * 0.5 + 0.5;
  const radius = baseRadius * pulse;
  const gradient = ctx.createRadialGradient(centerX, centerY, 0, centerX, centerY, radius);
  gradient.addColorStop(0, palette[0] || dominantColor);
  gradient.addColorStop(1, "rgba(255,255,255,0.1)");
  ctx.beginPath();
  ctx.fillStyle = gradient;
  ctx.arc(centerX, centerY, radius, 0, Math.PI * 2);
  ctx.fill();
  for (let i = 0; i < 50; i++) {
    const angle = (i / 50) * Math.PI * 2;
    const distance = radius + 30 + (data[i % data.length] / 255) * 100;
    const particleSize = (data[i % data.length] / 255) * 3 + 1;
    const x = centerX + Math.cos(angle) * distance;
    const y = centerY + Math.sin(angle) * distance;
    const colorIndex = i % palette.length;
    ctx.beginPath();
    ctx.fillStyle = palette[colorIndex] || dominantColor;
    ctx.arc(x, y, particleSize, 0, Math.PI * 2);
    ctx.fill();
  }
};

export const renderGalaxy: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const time = Date.now() * 0.001;
  const arms = 3;

  for (let arm = 0; arm < arms; arm++) {
    const armAngle = (arm / arms) * Math.PI * 2;
    for (let i = 0; i < 80; i++) {
      const t = i / 80;
      const spiralRadius = t * Math.min(width, height) * 0.4;
      const angle = armAngle + t * Math.PI * 4 + time * 0.5;
      const freq = data[Math.floor(t * data.length)] / 255;
      const wobble = freq * 30;
      const x = centerX + Math.cos(angle) * (spiralRadius + wobble);
      const y = centerY + Math.sin(angle) * (spiralRadius + wobble);
      const size = freq * 4 + 1;
      const alpha = (1 - t) * 0.8 + 0.2;
      const colorIndex = Math.floor(t * palette.length);
      ctx.beginPath();
      ctx.fillStyle = (palette[colorIndex] || dominantColor).replace(/[\d.]+\)$/, `${alpha})`);
      ctx.shadowBlur = 10;
      ctx.shadowColor = palette[colorIndex] || dominantColor;
      ctx.arc(x, y, size, 0, Math.PI * 2);
      ctx.fill();
      ctx.shadowBlur = 0;
    }
  }
};

export const renderDna: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const time = Date.now() * 0.002;
  const points = 60;

  for (let i = 0; i < points; i++) {
    const t = i / points;
    const x = t * width;
    const freq = data[Math.floor(t * data.length)] / 255;
    const amplitude = height * 0.25 * (0.5 + freq * 0.5);
    const y1 = height / 2 + Math.sin(t * Math.PI * 4 + time) * amplitude;
    const y2 = height / 2 - Math.sin(t * Math.PI * 4 + time) * amplitude;

    if (i % 4 === 0) {
      ctx.beginPath();
      ctx.strokeStyle = (palette[0] || dominantColor).replace(/[\d.]+\)$/, "0.3)");
      ctx.lineWidth = 1;
      ctx.moveTo(x, y1);
      ctx.lineTo(x, y2);
      ctx.stroke();
    }

    ctx.beginPath();
    ctx.fillStyle = palette[0] || dominantColor;
    ctx.shadowBlur = 8;
    ctx.shadowColor = palette[0] || dominantColor;
    ctx.arc(x, y1, 3 + freq * 3, 0, Math.PI * 2);
    ctx.fill();

    ctx.beginPath();
    ctx.fillStyle = palette[Math.min(2, palette.length - 1)] || dominantColor;
    ctx.shadowColor = palette[Math.min(2, palette.length - 1)] || dominantColor;
    ctx.arc(x, y2, 3 + freq * 3, 0, Math.PI * 2);
    ctx.fill();
    ctx.shadowBlur = 0;
  }
};

export const renderAurora: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const time = Date.now() * 0.001;
  const layers = 5;

  for (let layer = 0; layer < layers; layer++) {
    ctx.beginPath();
    const baseY = height * 0.3 + layer * 30;
    const points: [number, number][] = [];

    for (let x = 0; x <= width; x += 5) {
      const t = x / width;
      const freq = data[Math.floor(t * data.length)] / 255;
      const wave1 = Math.sin(t * 3 + time + layer) * 40 * freq;
      const wave2 = Math.sin(t * 5 + time * 1.5 + layer * 0.5) * 20 * freq;
      const y = baseY + wave1 + wave2;
      points.push([x, y]);
      if (x === 0) ctx.moveTo(x, y);
      else ctx.lineTo(x, y);
    }

    ctx.lineTo(width, height);
    ctx.lineTo(0, height);
    ctx.closePath();

    const colorIndex = layer % palette.length;
    const gradient = ctx.createLinearGradient(0, baseY - 50, 0, height);
    gradient.addColorStop(0, (palette[colorIndex] || dominantColor).replace(/[\d.]+\)$/, "0.4)"));
    gradient.addColorStop(1, "rgba(0,0,0,0)");
    ctx.fillStyle = gradient;
    ctx.fill();
  }
};

export const renderTerrain: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const layers = 8;
  for (let layer = 0; layer < layers; layer++) {
    const baseY = height * 0.4 + layer * (height * 0.08);
    const colorIndex = Math.min(layer, palette.length - 1);
    const alpha = 0.3 + (layer / layers) * 0.5;

    ctx.beginPath();
    ctx.moveTo(0, height);

    for (let x = 0; x <= width; x += 3) {
      const t = x / width;
      const dataIdx = Math.floor((t + layer * 0.1) * data.length) % data.length;
      const freq = data[dataIdx] / 255;
      const mountainHeight = freq * height * 0.3 * (1 - (layer / layers) * 0.5);
      const y = baseY - mountainHeight;
      ctx.lineTo(x, y);
    }

    ctx.lineTo(width, height);
    ctx.closePath();

    ctx.fillStyle = (palette[colorIndex] || dominantColor).replace(/[\d.]+\)$/, `${alpha})`);
    ctx.fill();
  }
};

// --- NEW VISUALIZERS ---

export const renderRetroBars: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const barCount = 32;
  const barWidth = (width / barCount) * 0.8;
  const spacing = (width / barCount) * 0.2;
  const segmentHeight = 8;
  const gap = 2;

  for (let i = 0; i < barCount; i++) {
    const value = data[Math.floor((i / barCount) * data.length)];
    const amplitude = (value / 255) * height;
    const segments = Math.floor(amplitude / (segmentHeight + gap));
    const x = i * (barWidth + spacing) + spacing / 2;
    const colorIndex = i % palette.length;

    ctx.fillStyle = palette[colorIndex] || dominantColor;
    ctx.shadowBlur = 10;
    ctx.shadowColor = palette[colorIndex] || dominantColor;

    for (let j = 0; j < segments; j++) {
      const y = height - j * (segmentHeight + gap) - segmentHeight;
      ctx.fillRect(x, y, barWidth, segmentHeight);
    }
    ctx.shadowBlur = 0;
  }
};

export const renderSunburst: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const radius = Math.min(width, height) * 0.45;
  const rayCount = 120;

  for (let i = 0; i < rayCount; i++) {
    const angle = (i / rayCount) * Math.PI * 2;
    const value = data[Math.floor((i / rayCount) * data.length)];
    const rayLength = (value / 255) * radius * 0.8 + radius * 0.2;
    const colorIndex = Math.floor((i / rayCount) * palette.length);

    ctx.beginPath();
    ctx.strokeStyle = palette[colorIndex] || dominantColor;
    ctx.lineWidth = 2;
    ctx.moveTo(centerX, centerY);
    ctx.lineTo(centerX + Math.cos(angle) * rayLength, centerY + Math.sin(angle) * rayLength);
    ctx.stroke();
  }
};

export const renderHexagons: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const size = 30;
  const cols = Math.ceil(width / (size * 1.5));
  const rows = Math.ceil(height / (size * Math.sqrt(3)));

  const drawHex = (x: number, y: number, s: number) => {
    ctx.beginPath();
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 3) * i;
      const px = x + s * Math.cos(angle);
      const py = y + s * Math.sin(angle);
      if (i === 0) ctx.moveTo(px, py);
      else ctx.lineTo(px, py);
    }
    ctx.closePath();
  };

  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      const x = c * size * 1.5;
      const y = r * size * Math.sqrt(3) + (c % 2 === 0 ? 0 : (size * Math.sqrt(3)) / 2);

      const dataIndex = (c + r * cols) % data.length;
      const value = data[dataIndex] / 255;

      if (value > 0.3) {
        const hexSize = size * value;
        const colorIndex = (c + r) % palette.length;
        ctx.fillStyle = (palette[colorIndex] || dominantColor).replace(/[\d.]+\)$/, `${value})`);
        drawHex(x, y, hexSize);
        ctx.fill();
      }
    }
  }
};

export const renderBlocks: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const cols = 20;
  const rows = 15;
  const cellW = width / cols;
  const cellH = height / rows;

  for (let x = 0; x < cols; x++) {
    for (let y = 0; y < rows; y++) {
      const idx = Math.floor(((x + y) / (cols + rows)) * data.length);
      const val = data[idx] / 255;
      const size = val * Math.min(cellW, cellH) * 0.9;

      if (val > 0.1) {
        const cx = x * cellW + cellW / 2;
        const cy = y * cellH + cellH / 2;
        const colorIndex = (x * y) % palette.length;

        ctx.fillStyle = palette[colorIndex] || dominantColor;
        ctx.globalAlpha = val;
        ctx.fillRect(cx - size / 2, cy - size / 2, size, size);
        ctx.globalAlpha = 1.0;
      }
    }
  }
};

export const renderSpiral: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const maxRadius = Math.min(width, height) * 0.45;
  const points = data.length;

  ctx.beginPath();
  ctx.lineWidth = 2;
  ctx.strokeStyle = palette[0] || dominantColor;

  for (let i = 0; i < points; i++) {
    const angle = i * 0.1;
    const baseRadius = (i / points) * maxRadius;
    const wobble = (data[i] / 255) * 40;
    const radius = baseRadius + wobble;

    const x = centerX + Math.cos(angle) * radius;
    const y = centerY + Math.sin(angle) * radius;

    if (i === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
  }
  ctx.stroke();
};

export const renderTunnel: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const rectCount = 20;

  ctx.lineWidth = 2;

  for (let i = 0; i < rectCount; i++) {
    const t = i / rectCount;
    const val = data[Math.floor(t * data.length)] / 255;
    const size = Math.pow(t, 2) * Math.min(width, height);
    const offset = val * 50;

    const w = size + offset;
    const h = size + offset;

    const colorIndex = i % palette.length;
    ctx.strokeStyle = palette[colorIndex] || dominantColor;
    ctx.globalAlpha = t; // Fade out towards center
    ctx.strokeRect(centerX - w / 2, centerY - h / 2, w, h);
  }
  ctx.globalAlpha = 1.0;
};

export const renderFlower: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const petals = 12;
  const maxRadius = Math.min(width, height) * 0.35;

  for (let i = 0; i < petals; i++) {
    const angle = (i / petals) * Math.PI * 2;
    const val = data[Math.floor((i / petals) * data.length)] / 255;
    const petalLen = maxRadius * (0.5 + val);

    const colorIndex = i % palette.length;
    ctx.fillStyle = (palette[colorIndex] || dominantColor).replace(/[\d.]+\)$/, "0.6)");

    ctx.save();
    ctx.translate(centerX, centerY);
    ctx.rotate(angle);
    ctx.beginPath();
    ctx.ellipse(petalLen / 2, 0, petalLen / 2, 20 + val * 30, 0, 0, Math.PI * 2);
    ctx.fill();
    ctx.restore();
  }
};

export const renderNeonGrid: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const horizonY = height * 0.4;
  const time = Date.now() * 0.1;
  const spacing = 40;

  // Vertical lines (perspective)
  ctx.strokeStyle = palette[0] || dominantColor;
  ctx.lineWidth = 1;

  const centerX = width / 2;
  for (let x = -width; x < width * 2; x += spacing * 2) {
    ctx.beginPath();
    ctx.moveTo(centerX, horizonY);
    // Displace x based on freq
    const freq = data[Math.abs(Math.floor(x)) % data.length] / 255;
    ctx.lineTo(x + (x - centerX) * 2, height);
    ctx.stroke();
  }

  // Horizontal lines (moving)
  const offset = time % spacing;
  for (let y = horizonY; y < height; y += (y - horizonY) * 0.1 + 2) {
    const freq = data[Math.floor((y / height) * data.length)] / 255;
    const yPos = y + offset;
    if (yPos > height) continue;

    ctx.beginPath();
    ctx.lineWidth = 1 + freq * 3;
    ctx.strokeStyle = palette[1] || dominantColor;
    ctx.moveTo(0, yPos);
    ctx.lineTo(width, yPos);
    ctx.stroke();
  }
};

export const renderKaleidoscope: RendererFn = (
  ctx,
  data,
  width,
  height,
  palette,
  dominantColor,
) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const slices = 8;
  const sliceAngle = (Math.PI * 2) / slices;
  const radius = Math.min(width, height) * 0.45;

  // Draw one slice and rotate it
  for (let s = 0; s < slices; s++) {
    ctx.save();
    ctx.translate(centerX, centerY);
    ctx.rotate(s * sliceAngle);

    ctx.beginPath();
    ctx.moveTo(0, 0);
    for (let i = 0; i < data.length / 4; i++) {
      const t = i / (data.length / 4);
      const r = t * radius;
      const val = data[i] / 255;
      const angle = (val - 0.5) * sliceAngle; // Wiggle within slice

      const x = Math.cos(angle) * r;
      const y = Math.sin(angle) * r;
      ctx.lineTo(x, y);
    }

    const colorIndex = s % palette.length;
    ctx.strokeStyle = palette[colorIndex] || dominantColor;
    ctx.lineWidth = 2;
    ctx.stroke();

    // Mirror
    ctx.scale(1, -1);
    ctx.beginPath();
    ctx.moveTo(0, 0);
    for (let i = 0; i < data.length / 4; i++) {
      const t = i / (data.length / 4);
      const r = t * radius;
      const val = data[i] / 255;
      const angle = (val - 0.5) * sliceAngle;

      const x = Math.cos(angle) * r;
      const y = Math.sin(angle) * r;
      ctx.lineTo(x, y);
    }
    ctx.stroke();

    ctx.restore();
  }
};

export const renderDrops: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const cols = 50;
  const colWidth = width / cols;
  const time = Date.now() * 0.005;

  for (let i = 0; i < cols; i++) {
    const val = data[Math.floor((i / cols) * data.length)] / 255;
    // Pseudorandom speed based on column index
    const speed = 100 + (i % 7) * 50;
    const y = ((time * speed) % (height + 200)) - 200;
    const len = 50 + val * 200;

    const colorIndex = i % palette.length;
    const gradient = ctx.createLinearGradient(0, y, 0, y + len);
    gradient.addColorStop(0, "transparent");
    gradient.addColorStop(1, palette[colorIndex] || dominantColor);

    ctx.fillStyle = gradient;
    ctx.fillRect(i * colWidth, y, colWidth - 2, len);

    // Drop splash
    if (y + len > height - 10) {
      ctx.beginPath();
      ctx.fillStyle = palette[colorIndex] || dominantColor;
      ctx.arc(i * colWidth + colWidth / 2, height, val * 10, 0, Math.PI * 2);
      ctx.fill();
    }
  }
};

export const renderRings: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  // Concentric rings pulsing
  const centerX = width / 2;
  const centerY = height / 2;
  const maxRadius = Math.min(width, height) * 0.45;
  const ringCount = 20;

  for (let i = 0; i < ringCount; i++) {
    const dataIndex = Math.floor((i / ringCount) * data.length);
    const value = data[dataIndex] / 255;
    const radius = (i / ringCount) * maxRadius;
    const thickness = 2 + value * 10;

    ctx.beginPath();
    ctx.strokeStyle = palette[i % palette.length] || dominantColor;
    ctx.lineWidth = thickness;
    ctx.arc(centerX, centerY, radius + value * 20, 0, Math.PI * 2);
    ctx.stroke();
  }
};

export const renderSegmentedBars: RendererFn = (
  ctx,
  data,
  width,
  height,
  palette,
  dominantColor,
) => {
  // LED style bars
  const barCount = 32;
  const barWidth = (width / barCount) * 0.8;
  const spacing = (width / barCount) * 0.2;

  for (let i = 0; i < barCount; i++) {
    const val = data[Math.floor((i / barCount) * data.length)] / 255;
    const totalHeight = height * 0.8;
    const segments = 20;
    const activeSegments = Math.floor(val * segments);
    const segHeight = totalHeight / segments;

    const x = i * (barWidth + spacing) + spacing / 2;

    for (let j = 0; j < segments; j++) {
      if (j > activeSegments) break;
      const y = height - j * segHeight - segHeight;
      const colorIndex = Math.floor((j / segments) * palette.length);
      ctx.fillStyle = palette[colorIndex] || dominantColor;
      ctx.fillRect(x, y + 1, barWidth, segHeight - 2);
    }
  }
};

export const renderSeismic: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  // Joy division style
  const lines = 20;
  const lineSpacing = height / (lines + 4);

  ctx.lineWidth = 1.5;

  for (let i = 0; i < lines; i++) {
    ctx.beginPath();
    const yBase = (i + 2) * lineSpacing;
    const colorIndex = i % palette.length;
    ctx.strokeStyle = palette[colorIndex] || dominantColor;

    let started = false;
    for (let x = 0; x < width; x += 5) {
      // Map x to data index, but offset index by line number to create moving pattern feeling
      const idx = Math.floor((x / width) * data.length);
      // Add localized peak
      const val = data[idx] / 255;
      // Envelope to keep edges flat
      const centerDist = Math.abs(x - width / 2) / (width / 2);
      const envelope = Math.max(0, 1 - centerDist);

      const offset = val * envelope * -50;

      if (!started) {
        ctx.moveTo(x, yBase + offset);
        started = true;
      } else {
        ctx.lineTo(x, yBase + offset);
      }
    }
    ctx.stroke();
  }
};

export const renderPixels: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  // Grid of glowing squares
  const cols = 32;
  const rows = 16;
  const cellW = width / cols;
  const cellH = height / rows;

  for (let y = 0; y < rows; y++) {
    for (let x = 0; x < cols; x++) {
      const idx = Math.floor(((x + y * cols) / (cols * rows)) * data.length);
      const val = data[idx] / 255;

      if (val > 0.05) {
        ctx.fillStyle = (palette[(x + y) % palette.length] || dominantColor).replace(
          /[\d.]+\)$/,
          `${val})`,
        );
        const size = Math.max(1, val * Math.min(cellW, cellH) * 0.9);
        const cx = x * cellW + cellW / 2;
        const cy = y * cellH + cellH / 2;
        ctx.fillRect(cx - size / 2, cy - size / 2, size, size);
      }
    }
  }
};

export const renderLightning: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const bolts = 10;

  for (let i = 0; i < bolts; i++) {
    const val = data[Math.floor((i / bolts) * data.length)] / 255;
    if (val < 0.2) continue;

    const angle = (i / bolts) * Math.PI * 2;
    const len = (Math.min(width, height) / 2) * (0.5 + val);

    ctx.beginPath();
    ctx.strokeStyle = palette[i % palette.length] || dominantColor;
    ctx.lineWidth = 2;
    ctx.moveTo(centerX, centerY);

    let curX = centerX;
    let curY = centerY;
    const segments = 10;

    for (let j = 0; j < segments; j++) {
      const t = (j + 1) / segments;
      const tx = centerX + Math.cos(angle) * len * t;
      const ty = centerY + Math.sin(angle) * len * t;

      // Jitter
      const jitter = val * 20;
      const nx = tx + (Math.random() - 0.5) * jitter;
      const ny = ty + (Math.random() - 0.5) * jitter;

      ctx.lineTo(nx, ny);
      curX = nx;
      curY = ny;
    }
    ctx.stroke();
  }
};

export const renderPolarWave: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const radius = Math.min(width, height) * 0.3;

  ctx.beginPath();
  ctx.strokeStyle = palette[0] || dominantColor;
  ctx.lineWidth = 3;

  const len = data.length;
  for (let i = 0; i < len; i++) {
    const angle = (i / len) * Math.PI * 2;
    const val = data[i] / 255;
    const r = radius + val * 60;

    const x = centerX + Math.cos(angle) * r;
    const y = centerY + Math.sin(angle) * r;

    if (i === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
  }
  ctx.closePath();
  ctx.stroke();

  // Inner glow
  ctx.fillStyle = (palette[1] || dominantColor).replace(/[\d.]+\)$/, "0.1)");
  ctx.fill();
};

export const renderConfetti: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const count = 50;
  const time = Date.now() * 0.001;

  for (let i = 0; i < count; i++) {
    const val = data[Math.floor((i / count) * data.length)] / 255;
    const x = (i * 12345 + time * 50) % width;
    const y = (i * 67890 + time * 100) % height;
    const size = 5 + val * 15;

    ctx.save();
    ctx.translate(x, y);
    ctx.rotate(time * (i % 2 === 0 ? 1 : -1) + val * 5);
    ctx.fillStyle = palette[i % palette.length] || dominantColor;
    ctx.fillRect(-size / 2, -size / 2, size, size);
    ctx.restore();
  }
};

export const renderGlitch: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const bars = 20;
  for (let i = 0; i < bars; i++) {
    const val = data[Math.floor((i / bars) * data.length)] / 255;
    if (val < 0.1) continue;

    const h = Math.random() * 50 * val;
    const y = Math.random() * height;
    const w = width * 0.8 + Math.random() * width * 0.2;
    const x = (Math.random() - 0.5) * 50 * val; // Shake x

    ctx.fillStyle = (palette[i % palette.length] || dominantColor).replace(
      /[\d.]+\)$/,
      `${val * 0.7})`,
    );
    ctx.fillRect(x, y, w, h);
  }
};

export const renderInfinity: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const scale = Math.min(width, height) * 0.3;
  const points = 100;

  ctx.beginPath();
  ctx.lineWidth = 4;
  ctx.lineCap = "round";
  const gradient = ctx.createLinearGradient(0, 0, width, height);
  gradient.addColorStop(0, palette[0] || dominantColor);
  gradient.addColorStop(1, palette[palette.length - 1] || dominantColor);
  ctx.strokeStyle = gradient;

  for (let i = 0; i <= points; i++) {
    const t = (i / points) * Math.PI * 2;
    const val = data[Math.floor((i / points) * data.length)] / 255;

    // Lemniscate of Bernoulli
    const denom = 1 + Math.sin(t) * Math.sin(t);
    const xBase = (scale * Math.cos(t)) / denom;
    const yBase = (scale * Math.cos(t) * Math.sin(t)) / denom;

    // Extrude with audio
    const extrude = 1 + val * 0.5;

    ctx.lineTo(centerX + xBase * extrude, centerY + yBase * extrude);
  }
  ctx.stroke();
};

export const renderRain: RendererFn = (ctx, data, width, height, palette, dominantColor) => {
  const drops = 60;
  const time = Date.now() * 0.002;

  ctx.lineWidth = 2;

  for (let i = 0; i < drops; i++) {
    const val = data[Math.floor((i / drops) * data.length)] / 255;
    const x = (i / drops) * width;
    const speed = 200 + val * 300;
    const y = (time * speed + i * 100) % height;
    const len = 20 + val * 50;

    ctx.beginPath();
    ctx.strokeStyle = (palette[i % palette.length] || dominantColor).replace(
      /[\d.]+\)$/,
      `${0.2 + val * 0.8})`,
    );
    ctx.moveTo(x, y);
    ctx.lineTo(x, y + len);
    ctx.stroke();
  }
};

export const renderers: Record<VisualizerType, RendererFn | null> = {
  none: null,
  bars: renderBars,
  wave: renderWave,
  particles: renderParticles,
  circle: renderCircle,
  sphere: renderSphere,
  lines: renderLines,
  mesh: renderMesh,
  radial: renderRadial,
  spectrum: renderSpectrum,
  orb: renderOrb,
  galaxy: renderGalaxy,
  dna: renderDna,
  aurora: renderAurora,
  terrain: renderTerrain,
  // New
  retroBars: renderRetroBars,
  sunburst: renderSunburst,
  hexagons: renderHexagons,
  blocks: renderBlocks,
  spiral: renderSpiral,
  tunnel: renderTunnel,
  flower: renderFlower,
  neonGrid: renderNeonGrid,
  kaleidoscope: renderKaleidoscope,
  drops: renderDrops,
  // Newer
  rings: renderRings,
  segmentedBars: renderSegmentedBars,
  seismic: renderSeismic,
  pixels: renderPixels,
  lightning: renderLightning,
  polarWave: renderPolarWave,
  confetti: renderConfetti,
  glitch: renderGlitch,
  infinity: renderInfinity,
  rain: renderRain,
};
