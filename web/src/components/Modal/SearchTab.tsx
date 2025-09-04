import React, { useState, useEffect } from "react";
import { SearchTabProps, AudioFile } from "@/types";
import { Play, Search, Music } from "lucide-react";
import { searchLibrary } from "@/utils";
import { Virtuoso } from "react-virtuoso";
import TrackRow from "./TrackRow";
import "@/styles/modal-tabs.scss";

const SearchTab: React.FC<SearchTabProps> = ({
  library,
  currentTrack,
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

  const Row = React.useCallback(
    (index: number) => {
      return (
        <TrackRow
          library={searchResults}
          currentTrack={currentTrack}
          onPlayTrack={onPlayTrack}
          index={index}
        />
      );
    },
    [searchResults, currentTrack, onPlayTrack, library]
  );

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
            <Virtuoso
              totalCount={searchResults.length}
              itemContent={Row}
              style={{ height: "calc(90vh - 260px)" }}
              overscan={200}
              increaseViewportBy={{ top: 200, bottom: 200 }}
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default SearchTab;
