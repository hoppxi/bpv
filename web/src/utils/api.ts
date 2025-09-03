import { AudioFile } from "@/types";

export async function fetchMetadata(filePath: string): Promise<AudioFile> {
  const response = await fetch(`/api/metadata/${encodeURIComponent(filePath)}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch metadata: ${response.status}`);
  }
  const data = await response.json();
  return data.file;
}

export async function fetchCoverArt(filePath: string): Promise<string> {
  const response = await fetch(`/api/cover/${encodeURIComponent(filePath)}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch cover art: ${response.status}`);
  }
  const blob = await response.blob();
  return URL.createObjectURL(blob);
}

export async function searchLibrary(query: string): Promise<AudioFile[]> {
  const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
  if (!response.ok) {
    throw new Error(`Search failed: ${response.status}`);
  }
  const data = await response.json();
  return data.results || [];
}

export async function checkHealth(): Promise<boolean> {
  try {
    const response = await fetch("/api/health");
    return response.ok;
  } catch {
    return false;
  }
}

export async function baseFilePath(): Promise<string> {
  const response = await fetch("/api/base-path");
  if (!response.ok) {
    throw new Error(`Failed to fetch base path: ${response.status}`);
  }
  const data = await response.json();
  return data.base_path;
}
