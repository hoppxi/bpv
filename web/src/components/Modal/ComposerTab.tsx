import React, { useState } from "react";
import { ComposersTabProps } from "@/types";
import { Play, Users, Search, ArrowLeft } from "lucide-react";
import { Virtuoso } from "react-virtuoso";
import TrackRow from "./TrackRow";
import TrackGroupGrid from "./TrackGroupGrid";
import "@/styles/modal-tabs.scss";

const ComposersTab: React.FC<ComposersTabProps> = ({
  library,
  currentTrack,
  onPlayTrack,
}) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedComposer, setSelectedComposer] = useState<string | null>(null);

  const composers = Object.entries(library.composers)
    .sort(([a], [b]) => a.localeCompare(b))
    .filter(([composer]) =>
      composer.toLowerCase().includes(searchQuery.toLowerCase())
    );

  const composerSongs = selectedComposer
    ? library.files.filter((file) => file.composer === selectedComposer)
    : [];

  const handlePlayComposer = (composer: string) => {
    const composerTracks = library.files.filter(
      (file) => file.composer === composer
    );
    if (composerTracks.length > 0) {
      onPlayTrack(composerTracks[0]);
    }
  };

  const Row = React.useCallback(
    (index: number) => {
      return (
        <TrackRow
          library={composerSongs}
          currentTrack={currentTrack}
          onPlayTrack={onPlayTrack}
          index={index}
        />
      );
    },
    [composerSongs, currentTrack, onPlayTrack, library]
  );

  if (selectedComposer) {
    return (
      <div className="tab-content">
        <div className="tab-content__header">
          <button
            className="tab-content__back-btn"
            onClick={() => setSelectedComposer(null)}
          >
            <ArrowLeft size={14} /> Back to Composers
          </button>
          <div className="tab-content__stats">
            <h3>{selectedComposer}</h3>
            <p>{composerSongs.length} songs</p>
          </div>
          <button
            className="tab-content__action-btn"
            onClick={() => handlePlayComposer(selectedComposer)}
          >
            <Play size={16} />
            Play All
          </button>
        </div>

        <div className="tab-content__list">
          <div className="songs-list">
            <Virtuoso
              totalCount={composerSongs.length}
              itemContent={Row}
              style={{ height: "calc(100vh - 400px)" }}
              overscan={200}
              increaseViewportBy={{ top: 200, bottom: 200 }}
            />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="tab-content">
      <div className="tab-content__header">
        <div className="tab-content__stats">
          <h3>Composers</h3>
          <p>{composers.length} composers</p>
        </div>
        <div className="tab-content__search">
          <Search size={18} />
          <input
            type="text"
            placeholder="Search composers..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="tab-content__search-input"
          />
        </div>
      </div>

      <div className="tab-content__list">
        <TrackGroupGrid
          metadata={{ name: "composer", icon: <Users size={32} /> }}
          group={composers}
          handlePlayGroup={handlePlayComposer}
          onclick={(track: string) => setSelectedComposer(track)}
        />
      </div>
    </div>
  );
};

export default ComposersTab;
