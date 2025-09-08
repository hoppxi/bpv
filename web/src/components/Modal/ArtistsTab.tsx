import React, { useState } from "react";
import { ArtistsTabProps } from "@/types";
import { Play, Users, Search, ArrowLeft } from "lucide-react";
import { Virtuoso } from "react-virtuoso";
import TrackRow from "./TrackRow";
import TrackGroupGrid from "./TrackGroupGrid";
import "@/styles/modal-tabs.scss";

const ArtistsTab: React.FC<ArtistsTabProps> = ({
  library,
  currentTrack,
  onPlayTrack,
}) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedArtist, setSelectedArtist] = useState<string | null>(null);

  const artists = Object.entries(library.artists)
    .sort(([a], [b]) => a.localeCompare(b))
    .filter(([artist]) =>
      artist.toLowerCase().includes(searchQuery.toLowerCase())
    );

  const artistSongs = selectedArtist
    ? library.files.filter((file) => file.artist === selectedArtist)
    : [];

  const handlePlayArtist = (artist: string) => {
    const artistTracks = library.files.filter((file) => file.artist === artist);
    if (artistTracks.length > 0) {
      onPlayTrack(artistTracks[0]);
    }
  };

  const Row = React.useCallback(
    (index: number) => {
      return (
        <TrackRow
          library={artistSongs}
          currentTrack={currentTrack}
          onPlayTrack={onPlayTrack}
          index={index}
        />
      );
    },
    [artistSongs, currentTrack, onPlayTrack, library]
  );

  if (selectedArtist) {
    return (
      <div className="tab-content">
        <div className="tab-content__header">
          <button
            className="tab-content__back-btn"
            onClick={() => setSelectedArtist(null)}
          >
            <ArrowLeft size={14} /> Back to Artists
          </button>
          <div className="tab-content__stats">
            <h3>{selectedArtist}</h3>
            <p>{artistSongs.length} songs</p>
          </div>
          <button
            className="tab-content__action-btn"
            onClick={() => handlePlayArtist(selectedArtist)}
          >
            <Play size={16} />
            Play All
          </button>
        </div>

        <div className="tab-content__list">
          <div className="songs-list">
            <Virtuoso
              totalCount={artistSongs.length}
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
          <h3>Artists</h3>
          <p>{artists.length} artists</p>
        </div>
        <div className="tab-content__search">
          <Search size={18} />
          <input
            type="text"
            placeholder="Search artists..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="tab-content__search-input"
          />
        </div>
      </div>

      <div className="tab-content__list">
        <TrackGroupGrid
          metadata={{ name: "artist", icon: <Users size={32} /> }}
          group={artists}
          handlePlayGroup={handlePlayArtist}
          onclick={(track) => setSelectedArtist(track)}
        />
      </div>
    </div>
  );
};

export default ArtistsTab;
