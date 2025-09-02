import { useState, useEffect } from "react";
import { AudioFile, VisualizerType } from "@/types";
import { useAudioPlayer, useLibraryData, useLocalStorage } from "@/hooks";
import { Player, Modal, Background } from "@/components";

function App() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentTrack, setCurrentTrack] = useState<AudioFile | null>(null);

  // Get library data
  const { library, loading, error, refreshLibrary } = useLibraryData();

  // Get user preferences from localStorage
  const [visualizerType, setVisualizerType] = useLocalStorage<VisualizerType>(
    "visualizerType",
    "bars"
  );
  const [shuffle, setShuffle] = useLocalStorage("shuffle", false);
  const [repeat, setRepeat] = useLocalStorage("repeat", false);
  const [volume, setVolume] = useLocalStorage("volume", 0.7);

  // Initialize audio player
  const audioPlayer = useAudioPlayer({
    volume,
    onTrackEnd: handleTrackEnd,
    onTrackChange: setCurrentTrack,
  });

  // Set initial track if none is set but library exists
  useEffect(() => {
    if (!currentTrack && library && library.files.length > 0) {
      const lastPlayedTrackId = localStorage.getItem("lastPlayedTrack");
      if (lastPlayedTrackId) {
        const track = library.files.find(
          (f) => f.file_path === lastPlayedTrackId
        );
        if (track) {
          setCurrentTrack(track);
          return;
        }
      }
      setCurrentTrack(library.files[0]);
    }
  }, [library, currentTrack]);

  // Play track when currentTrack changes
  useEffect(() => {
    if (currentTrack) {
      audioPlayer.playTrack(currentTrack);
      localStorage.setItem("lastPlayedTrack", currentTrack.file_path);
    }
  }, [currentTrack]);

  function handleTrackEnd() {
    if (!library || !currentTrack) return;

    const currentIndex = library.files.findIndex(
      (file) => file.file_path === currentTrack.file_path
    );

    if (shuffle) {
      // Play random track
      const randomIndex = Math.floor(Math.random() * library.files.length);
      setCurrentTrack(library.files[randomIndex]);
    } else if (repeat) {
      // Repeat current track
      audioPlayer.playTrack(currentTrack);
    } else if (currentIndex < library.files.length - 1) {
      // Play next track
      setCurrentTrack(library.files[currentIndex + 1]);
    }
    // Else: stop at end of playlist
  }

  function handlePlayTrack(track: AudioFile) {
    setCurrentTrack(track);
  }

  function handleNextTrack() {
    if (!library || !currentTrack) return;

    const currentIndex = library.files.findIndex(
      (file) => file.file_path === currentTrack.file_path
    );

    if (currentIndex < library.files.length - 1) {
      setCurrentTrack(library.files[currentIndex + 1]);
    } else if (repeat) {
      setCurrentTrack(library.files[0]);
    }
  }

  function handlePreviousTrack() {
    if (!library || !currentTrack) return;

    const currentIndex = library.files.findIndex(
      (file) => file.file_path === currentTrack.file_path
    );

    if (currentIndex > 0) {
      setCurrentTrack(library.files[currentIndex - 1]);
    } else if (repeat) {
      setCurrentTrack(library.files[library.files.length - 1]);
    }
  }

  if (loading) {
    return (
      <div className="app-loading">
        <div className="loading-spinner"></div>
        <p>Loading your music library...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="app-error">
        <h2>Error Loading Library</h2>
        <p>{error}</p>
        <button onClick={refreshLibrary}>Retry</button>
      </div>
    );
  }

  if (!library || library.files.length === 0) {
    return (
      <div className="app-empty">
        <h2>No Music Found</h2>
        <p>Scan your library to add music</p>
        <button onClick={() => setIsModalOpen(true)}>Open Library</button>
      </div>
    );
  }

  return (
    <div className="app">
      <Background
        track={currentTrack}
        visualizerType={visualizerType}
        isPlaying={audioPlayer.isPlaying}
        audioElement={audioPlayer.audioElement}
      />

      <Player
        currentTrack={currentTrack}
        isPlaying={audioPlayer.isPlaying}
        currentTime={audioPlayer.currentTime}
        duration={audioPlayer.duration}
        volume={audioPlayer.volume}
        shuffle={shuffle}
        repeat={repeat}
        visualizerType={visualizerType}
        onPlayPause={audioPlayer.togglePlayPause}
        onNext={handleNextTrack}
        onPrevious={handlePreviousTrack}
        onSeek={audioPlayer.seek}
        onVolumeChange={(vol) => {
          setVolume(vol);
          audioPlayer.setVolume(vol);
        }}
        onShuffleChange={setShuffle}
        onRepeatChange={setRepeat}
        onVisualizerChange={setVisualizerType}
        onOpenModal={() => setIsModalOpen(true)}
      />

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        library={library}
        currentTrack={currentTrack}
        visualizerType={visualizerType}
        shuffle={shuffle}
        repeat={repeat}
        onPlayTrack={handlePlayTrack}
        onVisualizerChange={setVisualizerType}
        onShuffleChange={setShuffle}
        onRepeatChange={setRepeat}
        onRefreshLibrary={refreshLibrary}
      />
    </div>
  );
}

export default App;
