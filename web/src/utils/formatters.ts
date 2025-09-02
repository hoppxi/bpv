export function formatTime(seconds: number): string {
  if (isNaN(seconds) || !isFinite(seconds)) return "0:00";

  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, "0")}`;
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 B";

  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
}

export function formatBitrate(bitrate: number): string {
  if (bitrate < 1000) return `${bitrate} kbps`;
  return `${(bitrate / 1000).toFixed(1)} Mbps`;
}

export function formatDuration(duration: string): string {
  try {
    const seconds = parseInt(duration);
    if (isNaN(seconds)) return "0:00";
    return formatTime(seconds);
  } catch {
    return "0:00";
  }
}

export function getFileExtension(filename: string): string {
  return filename.split(".").pop()?.toLowerCase() || "";
}

export function formatTrackNumber(track: number, totalTracks: number): string {
  if (!track) return "";
  if (!totalTracks || totalTracks <= 1) return track.toString();
  return `${track}/${totalTracks}`;
}

export function formatDiscNumber(disc: number, totalDiscs: number): string {
  if (!disc) return "";
  if (!totalDiscs || totalDiscs <= 1) return disc.toString();
  return `${disc}/${totalDiscs}`;
}

export function truncateText(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength - 3) + "...";
}
