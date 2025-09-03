import { LibraryResponse, AudioFile } from "./audio";
import { VisualizerType } from "./visualizer";

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  library: LibraryResponse;
  currentTrack: AudioFile | null;
  visualizerType: VisualizerType;
  shuffle: boolean;
  repeat: boolean;
  onPlayTrack: (track: AudioFile) => void;
  onVisualizerChange: (type: VisualizerType) => void;
  onShuffleChange: (shuffle: boolean) => void;
  onRepeatChange: (repeat: boolean) => void;
  onRefreshLibrary: () => void;
}

export type TabType =
  | "library"
  | "artists"
  | "albums"
  | "genres"
  | "search"
  | "settings"
  | "storage";

export interface ArtistsTabProps {
  library: LibraryResponse;
  currentTrack: AudioFile | null;
  onPlayTrack: (track: AudioFile) => void;
}

export interface LibraryTabProps extends ArtistsTabProps {
  onRefreshLibrary: () => void;
}

export interface SearchTabProps {
  library: LibraryResponse;
  currentTrack: AudioFile | null;
  searchQuery: string;
  onSearchChange: (query: string) => void;
  onPlayTrack: (track: AudioFile) => void;
}

export interface SettingsTabProps {
  visualizerType: VisualizerType;
  shuffle: boolean;
  repeat: boolean;
  onVisualizerChange: (type: VisualizerType) => void;
  onShuffleChange: (shuffle: boolean) => void;
  onRepeatChange: (repeat: boolean) => void;
}

export type { ArtistsTabProps as AlbumsTabProps };
export type { ArtistsTabProps as GenresTabProps };
