import React from "react";
import VolumeControl from "./VolumeControl";
import {
  Play,
  Pause,
  SkipBack,
  SkipForward,
  Shuffle,
  Repeat,
  List,
  BarChart3,
  Waves,
  Circle,
  Sparkles,
} from "lucide-react";
import { PlayerControlsProps, VisualizerType } from "@/types";
import "@/styles/player-controls.scss";

const PlayerControls: React.FC<PlayerControlsProps> = ({
  isPlaying,
  shuffle,
  repeat,
  visualizerType,
  onPlayPause,
  onNext,
  onPrevious,
  onShuffleChange,
  onRepeatChange,
  onVisualizerChange,
  onOpenModal,
  volume,
  onVolumeChange,
}) => {
  const visualizerIcons: Record<VisualizerType, React.ReactNode> = {
    bars: <BarChart3 size={20} />,
    wave: <Waves size={20} />,
    particles: <Sparkles size={20} />,
    circle: <Circle size={20} />,
    sphere: <Circle size={20} />,
    lines: <BarChart3 size={20} />,
    mesh: <BarChart3 size={20} />,
    radial: <Circle size={20} />,
    spectrum: <BarChart3 size={20} />,
    orb: <Circle size={20} />,
  };

  const visualizerTypes: VisualizerType[] = [
    "bars",
    "wave",
    "particles",
    "circle",
    "sphere",
    "lines",
    "mesh",
    "radial",
    "spectrum",
    "orb",
  ];

  return (
    <div className="player-controls">
      <div className="player-controls__main">
        {/* Previous */}
        <button
          className="player-controls__btn player-controls__btn--previous"
          onClick={onPrevious}
          title="Previous"
        >
          <SkipBack size={24} />
        </button>

        {/* Play/Pause */}
        <button
          className="player-controls__btn player-controls__btn--play-pause"
          onClick={onPlayPause}
          title={isPlaying ? "Pause" : "Play"}
        >
          {isPlaying ? <Pause size={32} /> : <Play size={32} />}
        </button>

        {/* Next */}
        <button
          className="player-controls__btn player-controls__btn--next"
          onClick={onNext}
          title="Next"
        >
          <SkipForward size={24} />
        </button>

        {/* Shuffle */}
        <button
          className={`player-controls__btn player-controls__btn--shuffle ${
            shuffle ? "player-controls__btn--active" : ""
          }`}
          onClick={() => onShuffleChange(!shuffle)}
          title="Shuffle"
        >
          <Shuffle size={20} />
        </button>

        {/* Repeat */}
        <button
          className={`player-controls__btn player-controls__btn--repeat ${
            repeat ? "player-controls__btn--active" : ""
          }`}
          onClick={() => onRepeatChange(!repeat)}
          title="Repeat"
        >
          <Repeat size={20} />
        </button>
      </div>

      <div className="player-controls__secondary">
        {/* Volume Control */}
        <VolumeControl volume={volume} onVolumeChange={onVolumeChange} />

        {/* Visualizer Selector */}
        <div className="player-controls__visualizer-selector">
          <select
            value={visualizerType}
            onChange={(e) =>
              onVisualizerChange(e.target.value as VisualizerType)
            }
            className="player-controls__visualizer-select"
            title="Visualizer Type"
          >
            {visualizerTypes.map((type) => (
              <option key={type} value={type}>
                {type.charAt(0).toUpperCase() + type.slice(1)}
              </option>
            ))}
          </select>
          <div className="player-controls__visualizer-icon">
            {visualizerIcons[visualizerType]}
          </div>
        </div>

        {/* Library Button */}
        <button
          className="player-controls__btn player-controls__btn--library"
          onClick={onOpenModal}
          title="Open Library"
        >
          <List size={20} />
        </button>
      </div>
    </div>
  );
};

export default PlayerControls;
