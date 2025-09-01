import {
  createContext,
  useContext,
  useRef,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { AudioFile, PlaybackState, QueueItem } from "@/types";

interface PlayerContextType {
  currentTrack: AudioFile | null;
  queue: QueueItem[];
  currentIndex: number;
  playbackState: PlaybackState;
  play: (track: AudioFile, queue?: AudioFile[]) => void;
  pause: () => void;
  resume: () => void;
  stop: () => void;
  seek: (time: number) => void;
  next: () => void;
  previous: () => void;
  setVolume: (volume: number) => void;
  toggleMute: () => void;
  toggleRepeat: () => void;
  toggleShuffle: () => void;
  addToQueue: (tracks: AudioFile[]) => void;
  clearQueue: () => void;
}

const PlayerContext = createContext<PlayerContextType | null>(null);

export function PlayerProvider({ children }: { children: ReactNode }) {
  const audioRef = useRef<HTMLAudioElement>(null);
  const [currentTrack, setCurrentTrack] = useState<AudioFile | null>(null);
  const [queue, setQueue] = useState<QueueItem[]>([]);
  const [currentIndex, setCurrentIndex] = useState(-1);
  const [playbackState, setPlaybackState] = useState<PlaybackState>({
    isPlaying: false,
    currentTime: 0,
    duration: 0,
    volume: 1,
    muted: false,
    repeat: "off",
    shuffle: false,
  });

  // Initialize audio element
  useEffect(() => {
    const audio = new Audio();
    audioRef.current = audio;

    const updateTime = () => {
      setPlaybackState((prev) => ({
        ...prev,
        currentTime: audio.currentTime,
        duration: audio.duration || 0,
      }));
    };

    const handleEnded = () => {
      if (playbackState.repeat === "one") {
        audio.currentTime = 0;
        audio.play();
      } else {
        next();
      }
    };

    audio.addEventListener("timeupdate", updateTime);
    audio.addEventListener("ended", handleEnded);
    audio.addEventListener("loadedmetadata", updateTime);

    return () => {
      audio.removeEventListener("timeupdate", updateTime);
      audio.removeEventListener("ended", handleEnded);
      audio.removeEventListener("loadedmetadata", updateTime);
      audio.pause();
    };
  }, [playbackState.repeat]);

  const play = (track: AudioFile, newQueue?: AudioFile[]) => {
    if (!audioRef.current) return;

    if (newQueue) {
      const queueItems = newQueue.map((file, index) => ({ file, index }));
      setQueue(queueItems);
      setCurrentIndex(
        newQueue.findIndex((f) => f.file_path === track.file_path)
      );
    }

    setCurrentTrack(track);
    audioRef.current.src = `/files/${encodeURIComponent(track.file_name)}`;
    audioRef.current.currentTime = 0;
    audioRef.current.play();
    setPlaybackState((prev) => ({ ...prev, isPlaying: true }));
  };

  const pause = () => {
    audioRef.current?.pause();
    setPlaybackState((prev) => ({ ...prev, isPlaying: false }));
  };

  const resume = () => {
    audioRef.current?.play();
    setPlaybackState((prev) => ({ ...prev, isPlaying: true }));
  };

  const stop = () => {
    audioRef.current?.pause();
    audioRef.current!.currentTime = 0;
    setPlaybackState((prev) => ({ ...prev, isPlaying: false, currentTime: 0 }));
  };

  const seek = (time: number) => {
    if (audioRef.current) {
      audioRef.current.currentTime = time;
      setPlaybackState((prev) => ({ ...prev, currentTime: time }));
    }
  };

  const next = () => {
    if (queue.length === 0) return;

    let nextIndex = currentIndex + 1;
    if (nextIndex >= queue.length) {
      if (playbackState.repeat === "all") {
        nextIndex = 0;
      } else {
        stop();
        return;
      }
    }

    setCurrentIndex(nextIndex);
    setCurrentTrack(queue[nextIndex].file);
    play(queue[nextIndex].file);
  };

  const previous = () => {
    if (queue.length === 0 || currentIndex <= 0) return;

    const prevIndex = currentIndex - 1;
    setCurrentIndex(prevIndex);
    setCurrentTrack(queue[prevIndex].file);
    play(queue[prevIndex].file);
  };

  const setVolume = (volume: number) => {
    if (audioRef.current) {
      audioRef.current.volume = volume;
      setPlaybackState((prev) => ({ ...prev, volume }));
    }
  };

  const toggleMute = () => {
    if (audioRef.current) {
      audioRef.current.muted = !audioRef.current.muted;
      setPlaybackState((prev) => ({ ...prev, muted: audioRef.current!.muted }));
    }
  };

  const toggleRepeat = () => {
    setPlaybackState((prev) => ({
      ...prev,
      repeat:
        prev.repeat === "off" ? "all" : prev.repeat === "all" ? "one" : "off",
    }));
  };

  const toggleShuffle = () => {
    setPlaybackState((prev) => ({ ...prev, shuffle: !prev.shuffle }));
  };

  const addToQueue = (tracks: AudioFile[]) => {
    const newItems = tracks.map((file, index) => ({
      file,
      index: queue.length + index,
    }));
    setQueue((prev) => [...prev, ...newItems]);
  };

  const clearQueue = () => {
    setQueue([]);
    setCurrentIndex(-1);
    stop();
  };

  const value: PlayerContextType = {
    currentTrack,
    queue,
    currentIndex,
    playbackState,
    play,
    pause,
    resume,
    stop,
    seek,
    next,
    previous,
    setVolume,
    toggleMute,
    toggleRepeat,
    toggleShuffle,
    addToQueue,
    clearQueue,
  };

  return (
    <PlayerContext.Provider value={value}>
      {children}
      <audio ref={audioRef} style={{ display: "none" }} />
    </PlayerContext.Provider>
  );
}

export function usePlayer() {
  const context = useContext(PlayerContext);
  if (!context) {
    throw new Error("usePlayer must be used within a PlayerProvider");
  }
  return context;
}
