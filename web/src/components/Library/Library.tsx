import { useState } from "react";
import { useLibrary } from "@/hooks/useLIbrary";
import { usePlayer } from "@/hooks/usePlayer";
import { RefreshCw, Play, MoreHorizontal } from "lucide-react";
import "@/styles/Library/Library.css";

export default function Library() {
  const { library, loading, error, refreshLibrary, scanLibrary } = useLibrary();
  const { play } = usePlayer();
  const [view, setView] = useState<"songs" | "artists" | "albums">("songs");

  const handlePlay = (file: any) => {
    play(file, library?.files || []);
  };

  const handlePlayAll = () => {
    if (library?.files && library.files.length > 0) {
      play(library.files[0], library.files);
    }
  };

  if (loading) {
    return (
      <div className="library-loading">
        <RefreshCw className="spinner" size={32} />
        <p>Loading your music library...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="library-error">
        <h2>Error Loading Library</h2>
        <p>{error}</p>
        <button onClick={refreshLibrary} className="retry-button">
          Try Again
        </button>
        <button onClick={scanLibrary} className="scan-button">
          Scan Library
        </button>
      </div>
    );
  }

  if (!library || library.audio_files === 0) {
    return (
      <div className="library-empty">
        <h2>No Music Found</h2>
        <p>Your music library is empty. Scan your library to add music.</p>
        <button onClick={scanLibrary} className="scan-button">
          Scan Library
        </button>
      </div>
    );
  }

  return (
    <div className="library">
      <div className="library-header">
        <h1>Your Library</h1>
        <div className="library-actions">
          <button onClick={handlePlayAll} className="play-all-button">
            <Play size={16} />
            Play All
          </button>
          <button onClick={refreshLibrary} className="refresh-button">
            <RefreshCw size={16} />
          </button>
          <div className="view-toggle">
            <button
              onClick={() => setView("songs")}
              className={view === "songs" ? "active" : ""}
            >
              Songs
            </button>
            <button
              onClick={() => setView("artists")}
              className={view === "artists" ? "active" : ""}
            >
              Artists
            </button>
            <button
              onClick={() => setView("albums")}
              className={view === "albums" ? "active" : ""}
            >
              Albums
            </button>
          </div>
        </div>
      </div>

      <div className="library-stats">
        <div className="stat">
          <span className="stat-number">{library.audio_files}</span>
          <span className="stat-label">Songs</span>
        </div>
        <div className="stat">
          <span className="stat-number">
            {Object.keys(library.artists).length}
          </span>
          <span className="stat-label">Artists</span>
        </div>
        <div className="stat">
          <span className="stat-number">
            {Object.keys(library.albums).length}
          </span>
          <span className="stat-label">Albums</span>
        </div>
        <div className="stat">
          <span className="stat-number">
            {Object.keys(library.genres).length}
          </span>
          <span className="stat-label">Genres</span>
        </div>
      </div>

      {view === "songs" && (
        <div className="songs-list">
          <div className="songs-header">
            <span className="header-number">#</span>
            <span className="header-title">Title</span>
            <span className="header-artist">Artist</span>
            <span className="header-album">Album</span>
            <span className="header-duration">Duration</span>
            <span className="header-actions"></span>
          </div>
          {library.files.map((file, index) => (
            <div key={file.file_path} className="song-row">
              <span className="song-number">{index + 1}</span>
              <div className="song-info">
                {file.cover_art ? (
                  <img
                    src={`data:${file.cover_art_mime};base64,${file.cover_art}`}
                    alt={file.album}
                    className="song-cover"
                  />
                ) : (
                  <div className="song-cover-placeholder">
                    <Music size={16} />
                  </div>
                )}
                <div className="song-details">
                  <span className="song-title">{file.title}</span>
                  {file.album && (
                    <span className="song-album-mobile">{file.album}</span>
                  )}
                </div>
              </div>
              <span className="song-artist">{file.artist}</span>
              <span className="song-album">{file.album}</span>
              <span className="song-duration">
                {file.duration ? formatDuration(file.duration) : "--:--"}
              </span>
              <div className="song-actions">
                <button
                  onClick={() => handlePlay(file)}
                  className="play-button"
                >
                  <Play size={16} />
                </button>
                <button className="more-button">
                  <MoreHorizontal size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {view === "artists" && (
        <div className="artists-grid">
          {Object.entries(library.artists).map(([artist, count]) => (
            <div key={artist} className="artist-card">
              <div className="artist-avatar">
                {artist.charAt(0).toUpperCase()}
              </div>
              <div className="artist-info">
                <h3 className="artist-name">{artist}</h3>
                <p className="artist-songs">
                  {count} song{count !== 1 ? "s" : ""}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}

      {view === "albums" && (
        <div className="albums-grid">
          {Object.entries(library.albums).map(([album, count]) => (
            <div key={album} className="album-card">
              <div className="album-cover">
                <Music size={32} />
              </div>
              <div className="album-info">
                <h3 className="album-name">{album}</h3>
                <p className="album-songs">
                  {count} song{count !== 1 ? "s" : ""}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, "0")}`;
}

function Music({ size }: { size: number }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
    >
      <path d="M9 18V5l12-2v13" />
      <circle cx="6" cy="18" r="3" />
      <circle cx="18" cy="16" r="3" />
    </svg>
  );
}
