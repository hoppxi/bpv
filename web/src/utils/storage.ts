export const STORAGE_KEYS = {
  LIBRARY: "musicLibrary",
  LAST_UPDATED: "libraryLastUpdated",
  CURRENT_TRACK: "currentTrack",
  PLAYBACK_STATE: "playbackState",
  VOLUME: "volume",
  SHUFFLE: "shuffle",
  REPEAT: "repeat",
  VISUALIZER: "visualizerType",
  SETTINGS: "playerSettings",
} as const;

export function getStoredData<T>(key: string, defaultValue: T): T {
  try {
    const item = localStorage.getItem(key);
    return item ? JSON.parse(item) : defaultValue;
  } catch (error) {
    console.error(`Error reading from localStorage key "${key}":`, error);
    return defaultValue;
  }
}

export function setStoredData<T>(key: string, value: T): void {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (error) {
    console.error(`Error writing to localStorage key "${key}":`, error);
  }
}

export function clearStoredData(key: string): void {
  try {
    localStorage.removeItem(key);
  } catch (error) {
    console.error(`Error removing localStorage key "${key}":`, error);
  }
}

export function clearAllStoredData(): void {
  try {
    localStorage.clear();
  } catch (error) {
    console.error("Error clearing localStorage:", error);
  }
}

export function getStoredPlaybackState() {
  return getStoredData(STORAGE_KEYS.PLAYBACK_STATE, {
    currentTime: 0,
    isPlaying: false,
    volume: 0.7,
  });
}

export function setStoredPlaybackState(state: {
  currentTime: number;
  isPlaying: boolean;
  volume: number;
}) {
  setStoredData(STORAGE_KEYS.PLAYBACK_STATE, state);
}

export function getStoredSettings() {
  return getStoredData(STORAGE_KEYS.SETTINGS, {
    visualizerType: "bars",
    shuffle: false,
    repeat: false,
    theme: "dark",
  });
}

export function setStoredSettings(settings: {
  visualizerType: string;
  shuffle: boolean;
  repeat: boolean;
  theme: string;
}) {
  setStoredData(STORAGE_KEYS.SETTINGS, settings);
}

export function migrateStorageData(): void {
  // Migration logic for future updates
  const oldKeys = ["music-player-settings", "audio-player-state"];

  oldKeys.forEach((oldKey) => {
    if (localStorage.getItem(oldKey)) {
      try {
        // Migrate to new format if needed
        localStorage.removeItem(oldKey);
      } catch {
        localStorage.removeItem(oldKey);
      }
    }
  });
}

// Initialize storage migration
if (typeof window !== "undefined") {
  migrateStorageData();
}
