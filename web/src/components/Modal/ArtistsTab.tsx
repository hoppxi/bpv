import React, { useState } from "react";
import { ArtistsTabProps } from "@/types";
import { Play, Users, Search, BarChart2, Music2 } from "lucide-react";
import { formatTime, formatFileSize } from "@/utils";
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

  if (selectedArtist) {
    return (
      <div className="tab-content">
        <div className="tab-content__header">
          <button
            className="tab-content__back-btn"
            onClick={() => setSelectedArtist(null)}
          >
            ← Back to Artists
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
            {artistSongs.map((track, index) => (
              <div
                key={track.file_path}
                className={`song-item ${
                  currentTrack?.file_path === track.file_path
                    ? "song-item--active"
                    : ""
                }`}
                onClick={() => onPlayTrack(track)}
              >
                <div className="song-item__index">
                  {currentTrack?.file_path === track.file_path ? (
                    <div className="song-item__playing-indicator">
                      <BarChart2 />
                    </div>
                  ) : (
                    <div className="song-item__number">{index + 1}</div>
                  )}

                  {track?.cover_art ? (
                    <img
                      src={`data:${track?.cover_art_mime};base64,${track?.cover_art}`}
                      alt={track?.album}
                      className="song-item__cover"
                    />
                  ) : (
                    <div className="song-item__cover-placeholder">
                      <Music2 />
                    </div>
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
        <div className="grid-list">
          {artists.map(([artist, count]) => (
            <div
              key={artist}
              className="grid-item"
              onClick={() => setSelectedArtist(artist)}
            >
              <div className="grid-item__icon">
                <Users size={32} />
              </div>
              <div className="grid-item__info">
                <div className="grid-item__title">{artist}</div>
                <div className="grid-item__subtitle">{count} songs</div>
              </div>
              <button
                className="grid-item__action"
                onClick={(e) => {
                  e.stopPropagation();
                  handlePlayArtist(artist);
                }}
                title="Play artist"
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

export default ArtistsTab;
