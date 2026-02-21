import type { AudioFile, LibraryResponse, HealthResponse } from "@/types";

const API_BASE = "/api";

async function fetchJSON<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, options);
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }
  return response.json();
}

export async function fetchHealth(): Promise<HealthResponse> {
  return fetchJSON<HealthResponse>(`${API_BASE}/health`);
}

export async function fetchLibrary(): Promise<LibraryResponse> {
  return fetchJSON<LibraryResponse>(`${API_BASE}/library`);
}

export async function fetchMetadata(filePath: string): Promise<AudioFile> {
  const data = await fetchJSON<{ file: AudioFile }>(
    `${API_BASE}/metadata/${encodeURIComponent(filePath)}`,
  );
  return data.file;
}

export async function fetchBasePath(): Promise<string> {
  const data = await fetchJSON<{ base_path: string }>(`${API_BASE}/base-path`);
  return data.base_path;
}

export async function searchLibrary(query: string): Promise<AudioFile[]> {
  const data = await fetchJSON<{ results: AudioFile[] }>(
    `${API_BASE}/search?q=${encodeURIComponent(query)}`,
  );
  return data.results || [];
}

export async function fetchArtists(): Promise<{ name: string; count: number }[]> {
  const data = await fetchJSON<{ artists: { name: string; count: number }[] }>(
    `${API_BASE}/artists`,
  );
  return data.artists || [];
}

export async function fetchArtistTracks(name: string): Promise<AudioFile[]> {
  const data = await fetchJSON<{ songs: AudioFile[] }>(
    `${API_BASE}/artist/${encodeURIComponent(name)}`,
  );
  return data.songs || [];
}

export async function fetchAlbums(): Promise<{ name: string; count: number }[]> {
  const data = await fetchJSON<{ albums: { name: string; count: number }[] }>(`${API_BASE}/albums`);
  return data.albums || [];
}

export async function fetchAlbumTracks(name: string): Promise<AudioFile[]> {
  const data = await fetchJSON<{ songs: AudioFile[] }>(
    `${API_BASE}/album/${encodeURIComponent(name)}`,
  );
  return data.songs || [];
}

export async function fetchGenres(): Promise<{ name: string; count: number }[]> {
  const data = await fetchJSON<{ genres: { name: string; count: number }[] }>(`${API_BASE}/genres`);
  return data.genres || [];
}

export async function fetchGenreTracks(name: string): Promise<AudioFile[]> {
  const data = await fetchJSON<{ songs: AudioFile[] }>(
    `${API_BASE}/genre/${encodeURIComponent(name)}`,
  );
  return data.songs || [];
}

export async function fetchComposers(): Promise<{ name: string; count: number }[]> {
  const data = await fetchJSON<{ composers: { name: string; count: number }[] }>(
    `${API_BASE}/composers`,
  );
  return data.composers || [];
}

export async function fetchComposerTracks(name: string): Promise<AudioFile[]> {
  const data = await fetchJSON<{ songs: AudioFile[] }>(
    `${API_BASE}/composer/${encodeURIComponent(name)}`,
  );
  return data.songs || [];
}

export async function triggerScan(): Promise<void> {
  await fetch(`${API_BASE}/scan-simple`, { method: "POST" });
}

export async function fetchFavorites(): Promise<string[]> {
  try {
    const data = await fetchJSON<{ favorites: string[] }>(`${API_BASE}/favorites`);
    return data.favorites || [];
  } catch {
    return [];
  }
}

export async function addFavorite(filePath: string): Promise<void> {
  const response = await fetch(`${API_BASE}/favorites`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ file_path: filePath }),
  });
  if (!response.ok) throw new Error("Failed to add favorite");
}

export async function removeFavorite(filePath: string): Promise<void> {
  const response = await fetch(`${API_BASE}/favorites`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ file_path: filePath }),
  });
  if (!response.ok) throw new Error("Failed to remove favorite");
}

export async function recordPlay(filePath: string): Promise<void> {
  try {
    await fetch(`${API_BASE}/stats/play`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ file_path: filePath }),
    });
  } catch {}
}

export async function fetchStats(): Promise<Record<string, number>> {
  try {
    const data = await fetchJSON<{ plays: Record<string, number> }>(`${API_BASE}/stats`);
    return data.plays || {};
  } catch {
    return {};
  }
}

export interface QueueState {
  file_paths: string[];
  current_index: number;
  shuffle: boolean;
  repeat: number; // 0=off, 1=all, 2=one
}

export async function fetchQueue(): Promise<QueueState> {
  const data = await fetchJSON<{ queue: QueueState }>(`${API_BASE}/queue`);
  return (
    data.queue || {
      file_paths: [],
      current_index: 0,
      shuffle: false,
      repeat: 0,
    }
  );
}

export async function saveQueue(queue: QueueState): Promise<void> {
  await fetch(`${API_BASE}/queue`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(queue),
  });
}

export interface SettingsState {
  visualizer_type?: string;
  show_visualizer?: boolean;
  auto_play?: boolean;
  crossfade?: boolean;
  gapless?: boolean;
  eq_bass?: number;
  eq_mid?: number;
  eq_treble?: number;
  eq_enabled?: boolean;
}

export async function fetchSettings(): Promise<SettingsState> {
  const data = await fetchJSON<{ settings: SettingsState }>(`${API_BASE}/settings`);
  return data.settings || {};
}

export async function saveSettings(settings: SettingsState): Promise<void> {
  await fetch(`${API_BASE}/settings`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(settings),
  });
}

export function getAudioUrl(filePath: string, basePath: string): string {
  const relativePath = filePath.replace(basePath, "");
  return `/files/${encodeURIComponent(relativePath)}`;
}

export function getCoverArtUrl(track: AudioFile): string | null {
  if (track.cover_art && track.cover_art_mime) {
    return `data:${track.cover_art_mime};base64,${track.cover_art}`;
  }
  return null;
}
