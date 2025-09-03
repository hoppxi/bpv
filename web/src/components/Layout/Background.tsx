import React, { useEffect, useRef, useState } from "react";
import { BackgroundProps } from "@/types";
import { useColorExtractor, useAudioVisualizer } from "@/hooks";
import "@/styles/background.scss";

const Background: React.FC<BackgroundProps> = ({
  track,
  visualizerType,
  isPlaying,
  audioElement,
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  const { dominantColor, colorPalette } = useColorExtractor({ track });

  const { frequencyData, timeDomainData } = useAudioVisualizer(
    audioElement,
    isPlaying
  );

  // Update canvas dimensions on resize
  useEffect(() => {
    const updateDimensions = () => {
      setDimensions({
        width: window.innerWidth,
        height: window.innerHeight,
      });
    };

    updateDimensions();
    window.addEventListener("resize", updateDimensions);

    return () => window.removeEventListener("resize", updateDimensions);
  }, []);

  // Main animation loop
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas || !isPlaying) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    let animationFrameId: number;

    const renderVisualizer = () => {
      try {
        ctx.clearRect(0, 0, canvas.width, canvas.height);

        switch (visualizerType) {
          case "bars":
            renderBars(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "wave":
            renderWave(
              ctx,
              timeDomainData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "particles":
            renderParticles(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "circle":
            renderCircle(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "sphere":
            renderSphere(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "lines":
            renderLines(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "mesh":
            renderMesh(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "radial":
            renderRadial(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "spectrum":
            renderSpectrum(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          case "orb":
            renderOrb(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
            break;
          default:
            renderBars(
              ctx,
              frequencyData,
              canvas.width,
              canvas.height,
              colorPalette
            );
        }

        animationFrameId = requestAnimationFrame(renderVisualizer);
      } catch (error) {
        console.warn("Error in visualizer rendering:", error);
      }
    };

    renderVisualizer();

    return () => {
      if (animationFrameId) {
        cancelAnimationFrame(animationFrameId);
      }
    };
  }, [
    visualizerType,
    isPlaying,
    frequencyData,
    timeDomainData,
    colorPalette,
    dimensions,
  ]);

  const renderBars = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
    const barWidth = (width / data.length) * 7.5;
    let x = 0;

    for (let i = 0; i < data.length; i++) {
      const barHeight = (data[i] / 255) * height * 0.8;
      const colorIndex = Math.floor((i / data.length) * palette.length);

      const gradient = ctx.createLinearGradient(
        x,
        height - barHeight,
        x,
        height
      );
      gradient.addColorStop(0, palette[colorIndex] || dominantColor);
      gradient.addColorStop(1, "rgba(255, 255, 255, 0.3)");

      ctx.fillStyle = gradient;
      ctx.fillRect(x, height - barHeight, barWidth, barHeight);

      x += barWidth + 1;
    }
  };

  const renderWave = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
    ctx.beginPath();
    ctx.lineWidth = 4;
    ctx.lineCap = "round";
    ctx.lineJoin = "round";

    // Create gradient for wave
    const gradient = ctx.createLinearGradient(0, 0, width, 0);
    gradient.addColorStop(0, palette[0] || dominantColor);
    gradient.addColorStop(
      0.5,
      palette[Math.floor(palette.length / 2)] || dominantColor
    );
    gradient.addColorStop(1, palette[palette.length - 1] || dominantColor);

    ctx.strokeStyle = gradient;
    ctx.shadowBlur = 15;
    ctx.shadowColor = dominantColor;

    const sliceWidth = width / data.length;
    let x = 0;

    for (let i = 0; i < data.length; i++) {
      // Convert byte data to waveform values (-1 to 1)
      const v = data[i] / 128.0 - 1;
      const y = (v * height) / 3 + height / 2;

      if (i === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }

      x += sliceWidth;
    }

    ctx.stroke();
    ctx.shadowBlur = 0;
  };

  const renderParticles = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
    const centerX = width / 2;
    const centerY = height / 2;
    const maxRadius = Math.min(width, height) * 0.45;

    for (let i = 0; i < data.length; i += 2) {
      const angle = (i / data.length) * Math.PI * 2;
      const radius =
        ((data[i] + data[i + 16] + data[i + 32]) / (3 * 255)) * maxRadius + 30;
      const x = centerX + Math.cos(angle) * radius;
      const y = centerY + Math.sin(angle) * radius;

      const size = ((data[i] + data[data.length - i]) / (2 * 255)) * 10 + 2;
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

  const renderCircle = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
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

      // Change color based on frequency
      ctx.strokeStyle = palette[colorIndex] || dominantColor;

      if (i === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    }

    ctx.closePath();
    ctx.stroke();
  };

  const renderSphere = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
    const centerX = width / 2;
    const centerY = height / 2;
    const baseRadius = Math.min(width, height) * 0.3;

    const pulse = (data[0] / 255) * 0.5 + 0.5;
    const radius = baseRadius * pulse;

    // Create gradient
    const gradient = ctx.createRadialGradient(
      centerX,
      centerY,
      0,
      centerX,
      centerY,
      radius * 1.5
    );
    gradient.addColorStop(0, palette[0] || dominantColor);
    gradient.addColorStop(1, "rgba(0,0,0,0)");

    ctx.beginPath();
    ctx.fillStyle = gradient;
    ctx.arc(centerX, centerY, radius, 0, Math.PI * 2);
    ctx.fill();
  };

  const renderLines = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
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

  const renderMesh = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
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
          size
        );
      }
    }
  };

  const renderRadial = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
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
      ctx.lineTo(
        centerX + Math.cos(angle) * radius,
        centerY + Math.sin(angle) * radius
      );
      ctx.stroke();
    }
  };

  const renderSpectrum = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
    const barCount = 100;
    const barWidth = width / barCount;

    for (let i = 0; i < barCount; i++) {
      const value = data[Math.floor((i / barCount) * data.length)];
      const barHeight = (value / 255) * height * 0.8;
      const x = i * barWidth;
      const colorIndex = Math.floor((i / barCount) * palette.length);

      // Create gradient for each bar
      const gradient = ctx.createLinearGradient(
        x,
        height - barHeight,
        x,
        height
      );
      gradient.addColorStop(0, palette[colorIndex] || dominantColor);
      gradient.addColorStop(1, "rgba(255,255,255,0.3)");

      ctx.fillStyle = gradient;
      ctx.fillRect(x, height - barHeight, barWidth - 1, barHeight);
    }
  };

  const renderOrb = (
    ctx: CanvasRenderingContext2D,
    data: Uint8Array,
    width: number,
    height: number,
    palette: string[]
  ) => {
    const centerX = width / 2;
    const centerY = height / 2;
    const baseRadius = Math.min(width, height) * 0.2;

    const pulse = (data[0] / 255) * 0.5 + 0.5;
    const radius = baseRadius * pulse;

    // Create gradient
    const gradient = ctx.createRadialGradient(
      centerX,
      centerY,
      0,
      centerX,
      centerY,
      radius
    );
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

  return (
    <div className="background">
      {/* Background overlay with dominant color */}
      <div
        className="background__overlay"
        style={{
          background: `linear-gradient(135deg, ${dominantColor.replace(
            ")",
            ", 30%)"
          )} 0%, ${dominantColor.replace(")", ", 70%)")} 100%)`,
        }}
      >
        {track?.cover_art && (
          <img
            src={`data:${track?.cover_art_mime};base64,${track?.cover_art}`}
            alt=""
            style={{
              width: "100svw",
              height: "100svh",
              objectFit: "cover",
              opacity: 0.4,
              filter: "blur(17px) brightness(0.7)",
            }}
          />
        )}
      </div>

      {/* Canvas for visualizers */}
      <canvas
        ref={canvasRef}
        className="background__canvas"
        width={dimensions.width}
        height={dimensions.height}
      />

      {/* Additional background elements */}
      <div className="background__particles">
        {Array.from({ length: 50 }).map((_, i) => (
          <div
            key={i}
            className="background__particle"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 5}s`,
              animationDuration: `${5 + Math.random() * 10}s`,
              background: colorPalette[i % colorPalette.length],
            }}
          />
        ))}
      </div>
    </div>
  );
};

export default Background;
