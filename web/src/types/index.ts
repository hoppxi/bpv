export interface AudioFile {
  file_path: string;
  file_name: string;
  file_size: number;
  file_type: string;
  modified: string;
  title: string;
  artist: string;
  album: string;
  album_artist: string;
  composer: string;
  genre: string;
  year: number;
  track: number;
  total_tracks: number;
  disc: number;
  total_discs: number;
  duration: number;
  bitrate: number;
  sample_rate: number;
  channels: number;
  comment: string;
  lyrics: string;
  bpm: number;
  cover_art: string;
  cover_art_mime: string;
  raw_metadata: Record<string, any>;
  error: string;
}

export interface LibraryResponse {
  status: string;
  music_dir: string;
  total_files: number;
  audio_files: number;
  artists: Record<string, number>;
  albums: Record<string, number>;
  genres: Record<string, number>;
  composers: Record<string, number>;
  files: AudioFile[];
  scan_time: string;
  errors: string[];
}

export interface HealthResponse {
  status: string;
  port: number;
  music_dir: string;
  version: string;
}

export interface PlaybackState {
  isPlaying: boolean;
  currentTime: number;
  duration: number;
  volume: number;
  muted: boolean;
  shuffle: boolean;
  repeat: RepeatMode;
  speed: number;
}

export type RepeatMode = "off" | "all" | "one";

export type VisualizerType =
  | "none"
  | "bars"
  | "wave"
  | "particles"
  | "circle"
  | "sphere"
  | "lines"
  | "mesh"
  | "radial"
  | "spectrum"
  | "orb"
  | "galaxy"
  | "dna"
  | "aurora"
  | "terrain"
  | "retroBars"
  | "sunburst"
  | "hexagons"
  | "blocks"
  | "spiral"
  | "tunnel"
  | "flower"
  | "neonGrid"
  | "kaleidoscope"
  | "drops"
  | "rings"
  | "segmentedBars"
  | "seismic"
  | "pixels"
  | "lightning"
  | "polarWave"
  | "confetti"
  | "glitch"
  | "infinity"
  | "rain";

export interface VisualizerConfig {
  type: VisualizerType;
  color: string;
  intensity: number;
  sensitivity: number;
}

export interface QueueItem {
  track: AudioFile;
  id: string;
}

export type SortField = "title" | "artist" | "album" | "duration" | "year" | "genre";
export type SortDirection = "asc" | "desc";

export interface SortConfig {
  field: SortField;
  direction: SortDirection;
}

export interface ThemeConfig {
  mode: "dark" | "light";
  accentColor: string;
}

export interface EqSettings {
  bass: number;
  mid: number;
  treble: number;
  enabled: boolean;
}

export type ViewMode =
  | "home"
  | "library"
  | "artists"
  | "albums"
  | "genres"
  | "composers"
  | "favorites"
  | "search"
  | "settings"
  | "now-playing";
