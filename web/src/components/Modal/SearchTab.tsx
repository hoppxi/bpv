import React, { useState, useEffect } from "react";
import { SearchTabProps, AudioFile } from "@/types";
import { Play, Search, Music } from "lucide-react";
import { formatTime, formatFileSize, searchLibrary } from "@/utils";
import "@/styles/modal-tabs.scss";

const SearchTab: React.FC<SearchTabProps> = ({
  library,
  searchQuery,
  onSearchChange,
  onPlayTrack,
}) => {
  const [searchResults, setSearchResults] = useState<AudioFile[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [localQuery, setLocalQuery] = useState(searchQuery);

  useEffect(() => {
    const delayDebounce = setTimeout(() => {
      onSearchChange(localQuery);
    }, 300);

    return () => clearTimeout(delayDebounce);
  }, [localQuery, onSearchChange]);

  useEffect(() => {
    const performSearch = async () => {
      if (!searchQuery.trim()) {
        setSearchResults([]);
        return;
      }

      setIsSearching(true);
      try {
        // First try client-side search for instant results
        const clientResults = library.files.filter(
          (file) =>
            file.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
            file.artist.toLowerCase().includes(searchQuery.toLowerCase()) ||
            file.album.toLowerCase().includes(searchQuery.toLowerCase()) ||
            file.genre.toLowerCase().includes(searchQuery.toLowerCase())
        );

        setSearchResults(clientResults);

        // Then try server-side search for better results
        if (clientResults.length === 0) {
          const serverResults = await searchLibrary(searchQuery);
          setSearchResults(serverResults);
        }
      } catch (error) {
        console.error("Search failed:", error);
        // Fallback to client-side search
        const fallbackResults = library.files.filter(
          (file) =>
            file.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
            file.artist.toLowerCase().includes(searchQuery.toLowerCase()) ||
            file.album.toLowerCase().includes(searchQuery.toLowerCase()) ||
            file.genre.toLowerCase().includes(searchQuery.toLowerCase())
        );
        setSearchResults(fallbackResults);
      } finally {
        setIsSearching(false);
      }
    };

    performSearch();
  }, [searchQuery, library.files]);

  const handlePlayAll = () => {
    if (searchResults.length > 0) {
      onPlayTrack(searchResults[0]);
    }
  };

  return (
    <div className="tab-content">
      <div className="tab-content__header">
        <div className="tab-content__stats">
          <h3>Search</h3>
          <p>
            {searchResults.length} results
            {isSearching && " • Searching..."}
          </p>
        </div>
        <div className="tab-content__search">
          <Search size={18} />
          <input
            type="text"
            placeholder="Search songs, artists, albums..."
            value={localQuery}
            onChange={(e) => setLocalQuery(e.target.value)}
            className="tab-content__search-input"
          />
        </div>
        {searchResults.length > 0 && (
          <button className="tab-content__action-btn" onClick={handlePlayAll}>
            <Play size={16} />
            Play All
          </button>
        )}
      </div>

      <div className="tab-content__list">
        {!searchQuery ? (
          <div className="tab-content__empty">
            <Search size={48} />
            <p>Enter a search term to find music</p>
          </div>
        ) : isSearching ? (
          <div className="tab-content__loading">
            <div className="loading-spinner"></div>
            <p>Searching...</p>
          </div>
        ) : searchResults.length === 0 ? (
          <div className="tab-content__empty">
            <Music size={48} />
            <p>No results found for "{searchQuery}"</p>
            <p className="tab-content__hint">
              Try searching by title, artist, album, or genre
            </p>
          </div>
        ) : (
          <div className="songs-list">
            {searchResults.map((track) => (
              <div
                key={track.file_path}
                className={`song-item ${
                  track.file_path === searchResults[0]?.file_path
                    ? "song-item--active"
                    : ""
                }`}
                onClick={() => onPlayTrack(track)}
              >
                <div className="song-item__number">
                  <Music size={16} />
                </div>
                <div className="song-item__info">
                  <div className="song-item__title">{track.title}</div>
                  <div className="song-item__artist">{track.artist}</div>
                  <div className="song-item__album">{track.album}</div>
                </div>
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

export default SearchTab;
