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
  | "orb";

export interface VisualizerConfig {
  type: VisualizerType;
  color: string;
  intensity: number;
  sensitivity: number;
  complexity: number;
}

export interface AudioAnalysis {
  frequencyData: Uint8Array;
  timeDomainData: Uint8Array;
  waveform: number[];
  beat: boolean;
  energy: number;
}
