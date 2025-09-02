import { AudioFile } from "./audio";
import { VisualizerType } from "./visualizer";

export interface BackgroundProps {
  track: AudioFile | null;
  visualizerType: VisualizerType;
  isPlaying: boolean;
  audioElement: HTMLAudioElement | null;
}
