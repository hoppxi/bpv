import React, { useState } from "react";
import { GenresTabProps } from "@/types";
import { Play, Tag, Search, ArrowLeft } from "lucide-react";
import { Virtuoso } from "react-virtuoso";
import TrackRow from "./TrackRow";
import "@/styles/modal-tabs.scss";

const GenresTab: React.FC<GenresTabProps> = ({
  library,
  currentTrack,
  onPlayTrack,
}) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedGenre, setSelectedGenre] = useState<string | null>(null);

  const genres = Object.entries(library.genres)
    .sort(([a], [b]) => a.localeCompare(b))
    .filter(([genre]) =>
      genre.toLowerCase().includes(searchQuery.toLowerCase())
    );

  const genreSongs = selectedGenre
    ? library.files.filter((file) => file.genre === selectedGenre)
    : [];

  const handlePlayGenre = (genre: string) => {
    const genreTracks = library.files.filter((file) => file.genre === genre);
    if (genreTracks.length > 0) {
      onPlayTrack(genreTracks[0]);
    }
  };

  const Row = React.useCallback(
    (index: number) => {
      return (
        <TrackRow
          library={genreSongs}
          currentTrack={currentTrack}
          onPlayTrack={onPlayTrack}
          index={index}
        />
      );
    },
    [genreSongs, currentTrack, onPlayTrack, library]
  );

  if (selectedGenre) {
    return (
      <div className="tab-content">
        <div className="tab-content__header">
          <button
            className="tab-content__back-btn"
            onClick={() => setSelectedGenre(null)}
          >
            <ArrowLeft size={14} /> Back to Genres
          </button>
          <div className="tab-content__stats">
            <h3>{selectedGenre}</h3>
            <p>{genreSongs.length} songs</p>
          </div>
          <button
            className="tab-content__action-btn"
            onClick={() => handlePlayGenre(selectedGenre)}
          >
            <Play size={16} />
            Play All
          </button>
        </div>

        <div className="tab-content__list">
          <div className="songs-list">
            <Virtuoso
              totalCount={genreSongs.length}
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
          <h3>Genres</h3>
          <p>{genres.length} genres</p>
        </div>
        <div className="tab-content__search">
          <Search size={18} />
          <input
            type="text"
            placeholder="Search genres..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="tab-content__search-input"
          />
        </div>
      </div>

      <div className="tab-content__list">
        <div className="grid-list">
          {genres.map(([genre, count]) => (
            <div
              key={genre}
              className="grid-item"
              onClick={() => setSelectedGenre(genre)}
            >
              <div className="grid-item__icon">
                <Tag size={32} />
              </div>
              <div className="grid-item__info">
                <div className="grid-item__title">{genre}</div>
                <div className="grid-item__subtitle">{count} songs</div>
              </div>
              <button
                className="grid-item__action"
                onClick={(e) => {
                  e.stopPropagation();
                  handlePlayGenre(genre);
                }}
                title="Play genre"
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

export default GenresTab;
