import React from "react";
import PlayerControls from "./PlayerControls";
import ProgressBar from "./ProgressBar";
import { Music2 } from "lucide-react";
import { PlayerProps } from "@/types";
import "@/styles/player.scss";

const Player: React.FC<PlayerProps> = ({
  currentTrack,
  isPlaying,
  currentTime,
  duration,
  volume,
  shuffle,
  repeat,
  visualizerType,
  onPlayPause,
  onNext,
  onPrevious,
  onSeek,
  onVolumeChange,
  onShuffleChange,
  onRepeatChange,
  onVisualizerChange,
  onOpenModal,
}) => {
  if (!currentTrack) {
    return (
      <div className="player player--empty">
        <div className="player__content">
          <div className="player__info">
            <h2 className="player__title">No track selected</h2>
            <p className="player__subtitle">Choose a song to start playing</p>
          </div>
          <button className="player__library-btn" onClick={onOpenModal}>
            Open Library
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="player">
      <div className="player__content">
        {/* Track Info */}
        <div className="player__info">
          <div className="player__cover-container">
            {currentTrack.cover_art ? (
              <img
                src={`data:${currentTrack.cover_art_mime};base64,${currentTrack.cover_art}`}
                alt={currentTrack.album}
                className="player__cover"
              />
            ) : (
              <div className="player__cover player__cover--placeholder">
                <span className="player__cover-icon">
                  <Music2 />
                </span>
              </div>
            )}
          </div>

          <div className="player__text">
            <h2 className="player__title" title={currentTrack.title}>
              {currentTrack.title}
            </h2>
            <p className="player__artist" title={currentTrack.artist}>
              {currentTrack.artist}
            </p>
            <p className="player__album" title={currentTrack.album}>
              {currentTrack.album}
            </p>
          </div>
        </div>

        {/* Progress Bar */}
        <ProgressBar
          currentTime={currentTime}
          duration={duration}
          onSeek={onSeek}
        />

        {/* Controls */}
        <PlayerControls
          isPlaying={isPlaying}
          shuffle={shuffle}
          repeat={repeat}
          visualizerType={visualizerType}
          onPlayPause={onPlayPause}
          onNext={onNext}
          onPrevious={onPrevious}
          onShuffleChange={onShuffleChange}
          onRepeatChange={onRepeatChange}
          onVisualizerChange={onVisualizerChange}
          onOpenModal={onOpenModal}
          volume={volume}
          onVolumeChange={onVolumeChange}
        />
      </div>
    </div>
  );
};

export default Player;
