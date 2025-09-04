import React from "react";
import { Play, RefreshCw, ListMusic } from "lucide-react";
import { Virtuoso } from "react-virtuoso";
import { LibraryTabProps } from "@/types";
import TrackRow from "./TrackRow";
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

  const Row = React.useCallback(
    (index: number) => {
      return (
        <TrackRow
          library={library}
          currentTrack={currentTrack}
          onPlayTrack={onPlayTrack}
          index={index}
        />
      );
    },
    [library, currentTrack, onPlayTrack]
  );

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
          <div className="song-list">
            <Virtuoso
              totalCount={library.files.length}
              itemContent={Row}
              style={{ height: "calc(100vh - 400px)" }}
              overscan={200}
              increaseViewportBy={{ top: 200, bottom: 200 }}
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default LibraryTab;
