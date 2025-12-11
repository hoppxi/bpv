import { useState, useEffect } from "react";
import { AudioFile, VisualizerType } from "@/types";
import {
  useAudioPlayer,
  useLibraryData,
  // useIDB,
  useIndexedDB,
} from "@/hooks";
import { Player, Modal, Background } from "@/components";
import { IDB } from "./utils";

function App() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentTrack, setCurrentTrack] = useState<AudioFile | null>(null);

  const { library, loading, error, refreshLibrary } = useLibraryData();

  const [visualizerType, setVisualizerType] = useIndexedDB<VisualizerType>(
    "visualizerType",
    "bars"
  );
  const [shuffle, setShuffle] = useIndexedDB<boolean>("shuffle", false);
  const [repeat, setRepeat] = useIndexedDB<boolean>("repeat", false);
  const [volume, setVolume] = useIndexedDB<number>("volume", 0.7);
  const [seekPosition, setSeekPosition] = useIndexedDB<number | undefined>(
    "seek-position",
    0
  );

  const audioPlayer = useAudioPlayer({
    volume,
    onTrackEnd: async () => await handleTrackEnd(),
    onTrackChange: setCurrentTrack,
    onTimeUpdate: handleTimeUpdate,
  });

  function handleTimeUpdate(currentTime: number | undefined) {
    setSeekPosition(currentTime);
  }

  // Set initial track if none is set but library exists
  useEffect(() => {
    (async () => {
      if (!currentTrack && library && library.files.length > 0) {
        const lastPlayedTrackId = await IDB.getItem("lastPlayedTrack");
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
    })();
  }, [library, currentTrack]);

  useEffect(() => {
    if (currentTrack && audioPlayer.currentTime > 0) {
      // Save both global and track-specific seek position
      IDB.setItem(
        `seek-position-${currentTrack.file_path}`,
        audioPlayer.currentTime
      );
      setSeekPosition(audioPlayer.currentTime);
    }
  }, [audioPlayer.currentTime, currentTrack]);

  useEffect(() => {
    (async () => {
      if (currentTrack) {
        const trackSeekPosition: number | undefined | null = await IDB.getItem(
          `seek-position-${currentTrack.file_path}`
        );
        const positionToUse = trackSeekPosition || seekPosition;

        audioPlayer.playTrack(currentTrack, positionToUse);
        await IDB.setItem("lastPlayedTrack", currentTrack.file_path);

        if (trackSeekPosition) {
          await IDB.setItem(
            `seek-position-${currentTrack.file_path}`,
            trackSeekPosition
          );
        }
      }
    })();
  }, [currentTrack]);

  async function handleTrackEnd() {
    if (!library || !currentTrack) return;

    await IDB.removeItem(`seek-position-${currentTrack.file_path}`);

    const currentIndex = library.files.findIndex(
      (file) => file.file_path === currentTrack.file_path
    );

    if (shuffle) {
      const randomIndex = Math.floor(Math.random() * library.files.length);
      setCurrentTrack(library.files[randomIndex]);
    } else if (repeat) {
      audioPlayer.playTrack(currentTrack);
    } else {
      setCurrentTrack(library.files[currentIndex + 1]);
    }
  }

  function handlePlayTrack(track: AudioFile) {
    setCurrentTrack(track);
  }

  async function handleNextTrack() {
    if (!library || !currentTrack) return;

    await IDB.removeItem(`seek-position-${currentTrack.file_path}`);

    const currentIndex = library.files.findIndex(
      (file) => file.file_path === currentTrack.file_path
    );

    if (shuffle) {
      const randomIndex = Math.floor(Math.random() * library.files.length);
      setCurrentTrack(library.files[randomIndex]);
    } else if (currentIndex < library.files.length - 1) {
      setCurrentTrack(library.files[currentIndex + 1]);
    }
  }

  async function handlePreviousTrack() {
    if (!library || !currentTrack) return;

    await IDB.removeItem(`seek-position-${currentTrack.file_path}`);

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
        onNext={async () => await handleNextTrack()}
        onPrevious={async () => await handlePreviousTrack()}
        onSeek={(position) => {
          audioPlayer.seek(position);
          setSeekPosition(position);
          if (currentTrack) {
            IDB.setItem(`seek-position-${currentTrack.file_path}`, position);
          }
        }}
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
