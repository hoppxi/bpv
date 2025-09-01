import { AudioFile, LibraryResponse } from "@/types";

const API_BASE = "/api";

export const api = {
  // Health check
  async health() {
    const response = await fetch(`${API_BASE}/health`);
    return response.json();
  },

  // Library endpoints
  async getLibrary(): Promise<LibraryResponse> {
    const response = await fetch(`${API_BASE}/library`);
    if (!response.ok) throw new Error("Failed to fetch library");
    return response.json();
  },

  async scanLibrary(): Promise<void> {
    const response = await fetch(`${API_BASE}/scan`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ full_scan: true, extract_covers: true }),
    });
    if (!response.ok) throw new Error("Failed to start scan");
  },

  // Metadata endpoints
  async getMetadata(filePath: string): Promise<AudioFile> {
    const response = await fetch(
      `${API_BASE}/metadata/${encodeURIComponent(filePath)}`
    );
    if (!response.ok) throw new Error("Failed to fetch metadata");
    return response.json();
  },

  // Browse endpoints
  async getArtists() {
    const response = await fetch(`${API_BASE}/artists`);
    if (!response.ok) throw new Error("Failed to fetch artists");
    return response.json();
  },

  async getArtistSongs(artist: string) {
    const response = await fetch(
      `${API_BASE}/artist/${encodeURIComponent(artist)}`
    );
    if (!response.ok) throw new Error("Failed to fetch artist songs");
    return response.json();
  },

  async getAlbums() {
    const response = await fetch(`${API_BASE}/albums`);
    if (!response.ok) throw new Error("Failed to fetch albums");
    return response.json();
  },

  async getAlbumSongs(album: string) {
    const response = await fetch(
      `${API_BASE}/album/${encodeURIComponent(album)}`
    );
    if (!response.ok) throw new Error("Failed to fetch album songs");
    return response.json();
  },

  // Search
  async search(query: string) {
    const response = await fetch(
      `${API_BASE}/search?q=${encodeURIComponent(query)}`
    );
    if (!response.ok) throw new Error("Failed to search");
    return response.json();
  },

  // Cover art
  getCoverUrl(filePath: string): string {
    return `${API_BASE}/cover/${encodeURIComponent(filePath)}`;
  },

  // Music file URL
  getMusicUrl(filePath: string): string {
    return `/files/${encodeURIComponent(filePath)}`;
  },
};
