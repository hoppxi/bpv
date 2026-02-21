import { ref, watch, onUnmounted } from "vue";
import type { AudioFile, RepeatMode } from "@/types";
import { getAudioUrl, recordPlay } from "@/lib/api";

export function useAudioPlayer() {
  const audio = new Audio();
  audio.preload = "metadata";

  const isPlaying = ref(false);
  const currentTime = ref(0);
  const duration = ref(0);
  const volume = ref(0.7);
  const muted = ref(false);
  const isLoading = ref(false);
  const currentTrack = ref<AudioFile | null>(null);
  const playbackSpeed = ref(1);
  const errorMsg = ref<string | null>(null);

  let audioContext: AudioContext | null = null;
  let analyser: AnalyserNode | null = null;
  let source: MediaElementAudioSourceNode | null = null;
  let bassFilter: BiquadFilterNode | null = null;
  let midFilter: BiquadFilterNode | null = null;
  let trebleFilter: BiquadFilterNode | null = null;
  let sourceConnected = false;

  function initAudioContext() {
    if (audioContext) return;

    try {
      audioContext = new AudioContext();
      analyser = audioContext.createAnalyser();
      analyser.fftSize = 512;
      analyser.smoothingTimeConstant = 0.6;

      bassFilter = audioContext.createBiquadFilter();
      bassFilter.type = "lowshelf";
      bassFilter.frequency.value = 200;
      bassFilter.gain.value = 0;

      midFilter = audioContext.createBiquadFilter();
      midFilter.type = "peaking";
      midFilter.frequency.value = 1000;
      midFilter.Q.value = 1;
      midFilter.gain.value = 0;

      trebleFilter = audioContext.createBiquadFilter();
      trebleFilter.type = "highshelf";
      trebleFilter.frequency.value = 4000;
      trebleFilter.gain.value = 0;

      source = audioContext.createMediaElementSource(audio);
      sourceConnected = true;

      source.connect(bassFilter);
      bassFilter.connect(midFilter);
      midFilter.connect(trebleFilter);
      trebleFilter.connect(analyser);
      analyser.connect(audioContext.destination);
    } catch (e) {
      console.warn("Failed to initialize audio context:", e);
    }
  }

  audio.addEventListener("loadedmetadata", () => {
    duration.value = audio.duration || 0;
    isLoading.value = false;
  });

  audio.addEventListener("timeupdate", () => {
    currentTime.value = audio.currentTime;
  });

  audio.addEventListener("play", () => {
    isPlaying.value = true;
    errorMsg.value = null;
  });

  audio.addEventListener("pause", () => {
    isPlaying.value = false;
  });

  audio.addEventListener("waiting", () => {
    isLoading.value = true;
  });

  audio.addEventListener("canplay", () => {
    isLoading.value = false;
  });

  audio.addEventListener("error", () => {
    isPlaying.value = false;
    isLoading.value = false;
    errorMsg.value = "Failed to play audio";
  });

  let onTrackEndCallback: (() => void) | null = null;
  audio.addEventListener("ended", () => {
    isPlaying.value = false;
    if (onTrackEndCallback) onTrackEndCallback();
  });

  function onTrackEnd(cb: () => void) {
    onTrackEndCallback = cb;
  }

  async function playTrack(
    track: AudioFile,
    basePath: string,
    options: { startAt?: number; crossfade?: boolean } = {}
  ) {
    try {
      if (options.crossfade && isPlaying.value) {
        // Fade out
        const startVolume = volume.value;
        const fadeOutDuration = 1000; // 1s
        const steps = 20;
        const stepTime = fadeOutDuration / steps;
        const volStep = startVolume / steps;

        for (let i = 0; i < steps; i++) {
          audio.volume = Math.max(0, audio.volume - volStep);
          await new Promise((r) => setTimeout(r, stepTime));
        }
        audio.volume = 0;
      }

      isLoading.value = true;
      errorMsg.value = null;

      initAudioContext();
      if (audioContext?.state === "suspended") {
        await audioContext.resume();
      }

      const url = getAudioUrl(track.file_path, basePath);
      audio.src = url;
      audio.currentTime = options.startAt || 0;
      audio.volume = options.crossfade ? 0 : volume.value; // Start silent if fading in
      audio.load();

      currentTrack.value = track;

      await audio.play();
      isPlaying.value = true;

      recordPlay(track.file_path);

      if (options.crossfade) {
        // Fade in
        const targetVolume = volume.value;
        const fadeInDuration = 1000;
        const steps = 20;
        const stepTime = fadeInDuration / steps;
        const volStep = targetVolume / steps;

        for (let i = 0; i < steps; i++) {
          const newVol = Math.min(targetVolume, audio.volume + volStep);
          audio.volume = newVol;
          await new Promise((r) => setTimeout(r, stepTime));
        }
        audio.volume = targetVolume;
      } else {
        audio.volume = volume.value;
      }

    } catch (e: any) {
      console.error("Error playing track:", e);
      errorMsg.value = "Failed to play audio";
      isPlaying.value = false;
      isLoading.value = false;
      audio.volume = volume.value; // Restore volume on error
    }
  }

  function togglePlayPause() {
    if (!audio.src) return;

    initAudioContext();
    if (audioContext?.state === "suspended") {
      audioContext.resume();
    }

    if (isPlaying.value) {
      audio.pause();
    } else {
      audio.play().catch((e) => {
        console.error("Play failed:", e);
        errorMsg.value = "Failed to play";
      });
    }
  }

  function seek(time: number) {
    if (!audio.src) return;
    audio.currentTime = Math.max(0, Math.min(time, audio.duration || 0));
    currentTime.value = audio.currentTime;
  }

  function setVolume(val: number) {
    volume.value = val;
    audio.volume = val;
  }

  function toggleMute() {
    muted.value = !muted.value;
    audio.muted = muted.value;
  }

  function setSpeed(speed: number) {
    playbackSpeed.value = speed;
    audio.playbackRate = speed;
  }

  function skipForward(seconds = 10) {
    seek(audio.currentTime + seconds);
  }

  function skipBackward(seconds = 10) {
    seek(audio.currentTime - seconds);
  }

  function setEQ(bass: number, mid: number, treble: number) {
    if (bassFilter) bassFilter.gain.value = bass;
    if (midFilter) midFilter.gain.value = mid;
    if (trebleFilter) trebleFilter.gain.value = treble;
  }

  function getAnalyser(): AnalyserNode | null {
    return analyser;
  }

  function getAudioElement(): HTMLAudioElement {
    return audio;
  }

  watch(volume, (v) => {
    audio.volume = v;
  });

  watch(muted, (m) => {
    audio.muted = m;
  });

  onUnmounted(() => {
    audio.pause();
    audio.src = "";
    if (audioContext && audioContext.state !== "closed") {
      audioContext.close().catch(() => { });
    }
  });

  return {
    isPlaying,
    currentTime,
    duration,
    volume,
    muted,
    isLoading,
    currentTrack,
    playbackSpeed,
    errorMsg,

    playTrack,
    togglePlayPause,
    seek,
    setVolume,
    toggleMute,
    setSpeed,
    skipForward,
    skipBackward,
    setEQ,
    onTrackEnd,
    getAnalyser,
    getAudioElement,
  };
}
