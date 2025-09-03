import { useState, useRef, useCallback, useEffect } from "react";
import { baseFilePath } from "@/utils/api";
import { UseAudioPlayerProps, AudioFile } from "@/types";

export function useAudioPlayer({
  volume,
  onTrackEnd,
  onTrackChange,
  onTimeUpdate,
}: UseAudioPlayerProps) {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [currentTrack, setCurrentTrack] = useState<AudioFile | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const audioContextRef = useRef<AudioContext | null>(null);
  const analyserRef = useRef<AnalyserNode | null>(null);
  const audioSourceRef = useRef<MediaElementAudioSourceNode | null>(null);
  const [frequencyData, setFrequencyData] = useState<Uint8Array>(
    new Uint8Array(0)
  );

  useEffect(() => {
    if (!audioRef.current) {
      audioRef.current = new Audio();
      audioRef.current.preload = "metadata";
    }

    const audio = audioRef.current;

    const handleLoadedMetadata = () => {
      setDuration(audio.duration || 0);
      setIsLoading(false);
    };

    const handleTimeUpdate = () => {
      setCurrentTime(audio.currentTime);
    };

    const handleEnded = () => {
      setIsPlaying(false);
      onTrackEnd();
    };

    const handlePlay = () => {
      setIsPlaying(true);
      setError(null);
    };

    const handlePause = () => {
      setIsPlaying(false);
    };

    const handleError = (e: Event) => {
      setIsPlaying(false);
      setIsLoading(false);
      setError("Failed to play audio");
      console.error("Audio error:", e);
    };

    const handleWaiting = () => {
      setIsLoading(true);
    };

    const handleCanPlay = () => {
      setIsLoading(false);
    };

    // Add event listeners
    audio.addEventListener("loadedmetadata", handleLoadedMetadata);
    audio.addEventListener("timeupdate", handleTimeUpdate);
    audio.addEventListener("ended", handleEnded);
    audio.addEventListener("play", handlePlay);
    audio.addEventListener("pause", handlePause);
    audio.addEventListener("error", handleError);
    audio.addEventListener("waiting", handleWaiting);
    audio.addEventListener("canplay", handleCanPlay);

    // Set initial volume
    audio.volume = volume;

    return () => {
      audio.removeEventListener("loadedmetadata", handleLoadedMetadata);
      audio.removeEventListener("timeupdate", handleTimeUpdate);
      audio.removeEventListener("ended", handleEnded);
      audio.removeEventListener("play", handlePlay);
      audio.removeEventListener("pause", handlePause);
      audio.removeEventListener("error", handleError);
      audio.removeEventListener("waiting", handleWaiting);
      audio.removeEventListener("canplay", handleCanPlay);

      cleanupAudioContext();
    };
  }, [volume, onTrackEnd]);

  // Clean up audio context safely
  const cleanupAudioContext = useCallback(() => {
    if (audioSourceRef.current) {
      try {
        audioSourceRef.current.disconnect();
      } catch (e) {
        console.warn("Error disconnecting audio source:", e);
      }
      audioSourceRef.current = null;
    }

    if (analyserRef.current) {
      try {
        analyserRef.current.disconnect();
      } catch (e) {
        console.warn("Error disconnecting analyser:", e);
      }
      analyserRef.current = null;
    }

    if (audioContextRef.current) {
      if (audioContextRef.current.state !== "closed") {
        try {
          audioContextRef.current.close();
        } catch (e) {
          console.warn("Error closing audio context:", e);
        }
      }
      audioContextRef.current = null;
    }
  }, []);

  useEffect(() => {
    if (!analyserRef.current || !isPlaying) return;

    let animationFrameId: number;
    const updateFrequencyData = () => {
      try {
        const dataArray = new Uint8Array(
          analyserRef.current!.frequencyBinCount
        );
        analyserRef.current!.getByteFrequencyData(dataArray);
        setFrequencyData(dataArray);
        animationFrameId = requestAnimationFrame(updateFrequencyData);
      } catch (error) {
        console.warn("Error updating frequency data:", error);
        if (animationFrameId) {
          cancelAnimationFrame(animationFrameId);
        }
      }
    };

    animationFrameId = requestAnimationFrame(updateFrequencyData);

    return () => {
      if (animationFrameId) {
        cancelAnimationFrame(animationFrameId);
      }
    };
  }, [isPlaying]);

  const playTrack = useCallback(
    async (track: AudioFile, initialPosition: number = 0) => {
      if (!audioRef.current) return;

      try {
        setIsLoading(true);
        setError(null);

        // Stop current playback
        audioRef.current.pause();
        audioRef.current.currentTime = initialPosition;

        const audioPath = `/files/${encodeURIComponent(
          track.file_path.replace(await baseFilePath(), "")
        )}`;
        audioRef.current.src = audioPath;

        audioRef.current.load();
        setCurrentTrack(track);
        onTrackChange(track);

        const playPromise = audioRef.current.play();
        if (playPromise !== undefined) {
          playPromise
            .then(() => {
              setIsPlaying(true);
            })
            .catch((err) => {
              console.error("Play failed:", err);
              setError("Failed to play audio");
              setIsPlaying(false);
              setIsLoading(false);
            });
        }
      } catch (err) {
        console.error("Error playing track:", err);
        setError("Failed to play audio");
        setIsPlaying(false);
        setIsLoading(false);
      }
    },
    [onTrackChange, onTimeUpdate]
  );

  useEffect(() => {
    if (!audioRef.current) return;

    const handleTimeUpdate = () => {
      if (onTimeUpdate) {
        onTimeUpdate(audioRef.current?.currentTime);
      }
    };

    audioRef.current?.addEventListener("timeupdate", handleTimeUpdate);
    return () =>
      audioRef.current?.removeEventListener("timeupdate", handleTimeUpdate);
  }, [onTimeUpdate]);

  const togglePlayPause = useCallback(() => {
    if (!audioRef.current) return;

    try {
      if (isPlaying) {
        audioRef.current.pause();
        cleanupAudioContext();
      } else {
        const playPromise = audioRef.current.play();
        if (playPromise !== undefined) {
          playPromise
            .then(() => {
              setIsPlaying(true);
            })
            .catch((err) => {
              console.error("Play failed:", err);
              setError("Failed to play audio");
            });
        }
      }
    } catch (err) {
      console.error("Error toggling play/pause:", err);
      setError("Failed to toggle playback");
    }
  }, [isPlaying, cleanupAudioContext]);

  const seek = useCallback((time: number) => {
    if (!audioRef.current) return;

    try {
      audioRef.current.currentTime = time;
      setCurrentTime(time);
    } catch (err) {
      console.error("Error seeking:", err);
      setError("Failed to seek");
    }
  }, []);

  const setVolume = useCallback((newVolume: number) => {
    if (!audioRef.current) return;

    try {
      audioRef.current.volume = newVolume;
    } catch (err) {
      console.error("Error setting volume:", err);
      setError("Failed to set volume");
    }
  }, []);

  const mute = useCallback(() => {
    if (!audioRef.current) return;

    try {
      audioRef.current.muted = !audioRef.current.muted;
    } catch (err) {
      console.error("Error muting:", err);
      setError("Failed to mute");
    }
  }, []);

  const skipForward = useCallback((seconds: number = 10) => {
    if (!audioRef.current) return;

    try {
      audioRef.current.currentTime += seconds;
    } catch (err) {
      console.error("Error skipping forward:", err);
      setError("Failed to skip forward");
    }
  }, []);

  const skipBackward = useCallback((seconds: number = 10) => {
    if (!audioRef.current) return;

    try {
      audioRef.current.currentTime -= seconds;
    } catch (err) {
      console.error("Error skipping backward:", err);
      setError("Failed to skip backward");
    }
  }, []);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (audioRef.current) {
        audioRef.current.pause();
        audioRef.current.src = "";
      }
      cleanupAudioContext();
    };
  }, [cleanupAudioContext]);

  return {
    // State
    isPlaying,
    currentTime,
    duration,
    currentTrack,
    isLoading,
    error,
    frequencyData,

    audioElement: audioRef.current,

    // Methods
    playTrack,
    togglePlayPause,
    seek,
    setVolume,
    mute,
    skipForward,
    skipBackward,

    // Getters
    get volume() {
      return audioRef.current?.volume || 0;
    },
    get muted() {
      return audioRef.current?.muted || false;
    },
    get playbackRate() {
      return audioRef.current?.playbackRate || 1;
    },
  };
}
