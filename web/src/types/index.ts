export interface AudioFile {
  file_path: string;
  file_name: string;
  file_size: number;
  file_type: string;
  modified: string;
  title: string;
  artist: string;
  album: string;
  album_artist?: string;
  composer?: string;
  genre?: string;
  year?: number;
  track?: number;
  total_tracks?: number;
  disc?: number;
  total_discs?: number;
  duration?: number;
  bitrate?: number;
  sample_rate?: number;
  channels?: number;
  comment?: string;
  lyrics?: string;
  bpm?: number;
  cover_art?: string;
  cover_art_mime?: string;
}

export interface LibraryResponse {
  status: string;
  total_files: number;
  audio_files: number;
  artists: Record<string, number>;
  albums: Record<string, number>;
  genres: Record<string, number>;
  files: AudioFile[];
  scan_time: string;
  errors?: string[];
}

export interface PlaybackState {
  isPlaying: boolean;
  currentTime: number;
  duration: number;
  volume: number;
  muted: boolean;
  repeat: "off" | "one" | "all";
  shuffle: boolean;
}

export interface QueueItem {
  file: AudioFile;
  index: number;
}
