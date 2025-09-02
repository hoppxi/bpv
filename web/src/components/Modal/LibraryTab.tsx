import React from "react";
import { Play, RefreshCw, ListMusic } from "lucide-react";
import { LibraryTabProps } from "@/types";
import { formatTime, formatFileSize } from "@/utils/formatters";
import "@/styles/modal-tabs.scss";

const LibraryTab: React.FC<LibraryTabProps> = ({
  library,
  currentTrack,
  onPlayTrack,
  onRefreshLibrary,
}) => {
  const handlePlayAll = () => {
    if (library.files.length > 0) {
      onPlayTrack(library.files[0]);
    }
  };

  return (
    <div className="tab-content">
      <div className="tab-content__header">
        <div className="tab-content__stats">
          <h3>Your Library</h3>
          <p>
            {library.audio_files} songs • {Object.keys(library.artists).length}{" "}
            artists • {Object.keys(library.albums).length} albums
          </p>
        </div>
        <div className="tab-content__actions">
          <button
            className="tab-content__action-btn"
            onClick={handlePlayAll}
            disabled={library.files.length === 0}
          >
            <Play size={16} />
            Play All
          </button>
          <button
            className="tab-content__action-btn"
            onClick={onRefreshLibrary}
          >
            <RefreshCw size={16} />
            Rescan
          </button>
        </div>
      </div>

      <div className="tab-content__list">
        {library.files.length === 0 ? (
          <div className="tab-content__empty">
            <ListMusic size={48} />
            <p>No music found</p>
            <button
              className="tab-content__action-btn"
              onClick={onRefreshLibrary}
            >
              Scan Library
            </button>
          </div>
        ) : (
          <div className="songs-list">
            {library.files.map((track, index) => (
              <div
                key={track.file_path}
                className={`song-item ${
                  currentTrack?.file_path === track.file_path
                    ? "song-item--active"
                    : ""
                }`}
                onClick={() => onPlayTrack(track)}
              >
                <div className="song-item__number">
                  {currentTrack?.file_path === track.file_path ? (
                    <div className="song-item__playing-indicator">▶</div>
                  ) : (
                    index + 1
                  )}
                </div>
                <div className="song-item__info">
                  <div className="song-item__title">{track.title}</div>
                  <div className="song-item__artist">{track.artist}</div>
                </div>
                <div className="song-item__album">{track.album}</div>
                <div className="song-item__duration">
                  {formatTime(parseInt(track.duration) || 0)}
                </div>
                <div className="song-item__size">
                  {formatFileSize(track.file_size)}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default LibraryTab;
