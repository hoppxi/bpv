import React, { useState } from "react";
import { AlbumsTabProps } from "@/types";
import { Play, Disc, Search, ArrowLeft } from "lucide-react";
import TrackRow from "./TrackRow";
import { Virtuoso } from "react-virtuoso";
import "@/styles/modal-tabs.scss";

const AlbumsTab: React.FC<AlbumsTabProps> = ({
  library,
  currentTrack,
  onPlayTrack,
}) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedAlbum, setSelectedAlbum] = useState<string | null>(null);

  const albums = Object.entries(library.albums)
    .sort(([a], [b]) => a.localeCompare(b))
    .filter(([album]) =>
      album.toLowerCase().includes(searchQuery.toLowerCase())
    );

  const albumSongs = selectedAlbum
    ? library.files.filter((file) => file.album === selectedAlbum)
    : [];

  const handlePlayAlbum = (album: string) => {
    const albumTracks = library.files.filter((file) => file.album === album);
    if (albumTracks.length > 0) {
      onPlayTrack(albumTracks[0]);
    }
  };

  const Row = React.useCallback(
    (index: number) => {
      return (
        <TrackRow
          library={albumSongs}
          currentTrack={currentTrack}
          onPlayTrack={onPlayTrack}
          index={index}
        />
      );
    },
    [albumSongs, currentTrack, onPlayTrack, library]
  );

  if (selectedAlbum) {
    return (
      <div className="tab-content">
        <div className="tab-content__header">
          <button
            className="tab-content__back-btn"
            onClick={() => setSelectedAlbum(null)}
          >
            <ArrowLeft size={14} /> Back to Albums
          </button>
          <div className="tab-content__stats">
            <h3>{selectedAlbum}</h3>
            <p>{albumSongs.length} songs</p>
          </div>
          <button
            className="tab-content__action-btn"
            onClick={() => handlePlayAlbum(selectedAlbum)}
          >
            <Play size={16} />
            Play All
          </button>
        </div>

        <div className="tab-content__list">
          <div className="songs-list">
            <Virtuoso
              totalCount={albumSongs.length}
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
          <h3>Albums</h3>
          <p>{albums.length} albums</p>
        </div>
        <div className="tab-content__search">
          <Search size={18} />
          <input
            type="text"
            placeholder="Search albums..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="tab-content__search-input"
          />
        </div>
      </div>

      <div className="tab-content__list">
        <div className="grid-list">
          {albums.map(([album, count]) => (
            <div
              key={album}
              className="grid-item"
              onClick={() => setSelectedAlbum(album)}
            >
              <div className="grid-item__icon">
                <Disc size={32} />
              </div>
              <div className="grid-item__info">
                <div className="grid-item__title">{album}</div>
                <div className="grid-item__subtitle">{count} songs</div>
              </div>
              <button
                className="grid-item__action"
                onClick={(e) => {
                  e.stopPropagation();
                  handlePlayAlbum(album);
                }}
                title="Play album"
              >
                <Play size={16} />
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default AlbumsTab;
