import { usePlayer } from "@/hooks/usePlayer";
import {
  Play,
  Pause,
  SkipBack,
  SkipForward,
  Volume2,
  Repeat,
  Shuffle,
} from "lucide-react";
import "@/styles/Player/Player.css";

export default function Player() {
  const {
    currentTrack,
    playbackState,
    pause,
    resume,
    next,
    previous,
    seek,
    setVolume,
    toggleRepeat,
    toggleShuffle,
  } = usePlayer();

  if (!currentTrack) {
    return null;
  }

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, "0")}`;
  };

  const handleProgressClick = (e: React.MouseEvent<HTMLDivElement>) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const percent = (e.clientX - rect.left) / rect.width;
    seek(percent * playbackState.duration);
  };

  return (
    <div className="player">
      <div className="player-info">
        {currentTrack.cover_art ? (
          <img
            src={`data:${currentTrack.cover_art_mime};base64,${currentTrack.cover_art}`}
            alt={currentTrack.album}
            className="player-cover"
          />
        ) : (
          <div className="player-cover-placeholder">
            <Music size={24} />
          </div>
        )}
        <div className="player-track-info">
          <div className="player-title">{currentTrack.title}</div>
          <div className="player-artist">{currentTrack.artist}</div>
        </div>
      </div>

      <div className="player-controls">
        <div className="player-buttons">
          <button
            onClick={toggleShuffle}
            className={`player-button ${playbackState.shuffle ? "active" : ""}`}
          >
            <Shuffle size={16} />
          </button>
          <button onClick={previous} className="player-button">
            <SkipBack size={20} />
          </button>
          <button
            onClick={playbackState.isPlaying ? pause : resume}
            className="player-button play-button"
          >
            {playbackState.isPlaying ? <Pause size={24} /> : <Play size={24} />}
          </button>
          <button onClick={next} className="player-button">
            <SkipForward size={20} />
          </button>
          <button
            onClick={toggleRepeat}
            className={`player-button ${
              playbackState.repeat !== "off" ? "active" : ""
            }`}
          >
            <Repeat size={16} />
          </button>
        </div>

        <div className="player-progress">
          <span className="player-time">
            {formatTime(playbackState.currentTime)}
          </span>
          <div className="progress-bar" onClick={handleProgressClick}>
            <div
              className="progress-fill"
              style={{
                width: `${
                  (playbackState.currentTime / playbackState.duration) * 100
                }%`,
              }}
            />
          </div>
          <span className="player-time">
            {formatTime(playbackState.duration)}
          </span>
        </div>
      </div>

      <div className="player-volume">
        <Volume2 size={18} />
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          value={playbackState.volume}
          onChange={(e) => setVolume(parseFloat(e.target.value))}
          className="volume-slider"
        />
      </div>
    </div>
  );
}

// Add missing Music icon component
function Music({ size }: { size: number }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
    >
      <path d="M9 18V5l12-2v13" />
      <circle cx="6" cy="18" r="3" />
      <circle cx="18" cy="16" r="3" />
    </svg>
  );
}
