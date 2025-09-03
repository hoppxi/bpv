import { AudioFile } from "./audio";
import { VisualizerType } from "./visualizer";

export interface PlayerProps {
  currentTrack: AudioFile | null;
  isPlaying: boolean;
  currentTime: number;
  duration: number;
  volume: number;
  shuffle: boolean;
  repeat: boolean;
  visualizerType: VisualizerType;
  onPlayPause: () => void;
  onNext: () => void;
  onPrevious: () => void;
  onSeek: (time: number) => void;
  onVolumeChange: (volume: number) => void;
  onShuffleChange: (shuffle: boolean) => void;
  onRepeatChange: (repeat: boolean) => void;
  onVisualizerChange: (type: VisualizerType) => void;
  onOpenModal: () => void;
}

export interface PlayerControlsProps {
  isPlaying: boolean;
  shuffle: boolean;
  repeat: boolean;
  visualizerType: VisualizerType;
  onPlayPause: () => void;
  onNext: () => void;
  onPrevious: () => void;
  onShuffleChange: (shuffle: boolean) => void;
  onRepeatChange: (repeat: boolean) => void;
  onVisualizerChange: (type: VisualizerType) => void;
  onOpenModal: () => void;
  volume: number;
  onVolumeChange: (volume: number) => void;
}

export interface ProgressBarProps {
  currentTime: number;
  duration: number;
  onSeek: (time: number) => void;
}

export interface VolumeControlProps {
  volume: number;
  onVolumeChange: (volume: number) => void;
}
