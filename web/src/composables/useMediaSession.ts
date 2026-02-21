import { watch } from "vue";
import type { Ref } from "vue";
import type { AudioFile } from "@/types";
import { getCoverArtUrl } from "@/lib/api";

export function useMediaSession(
  currentTrack: Ref<AudioFile | null>,
  handlers: {
    play: () => void;
    pause: () => void;
    nextTrack: () => void;
    previousTrack: () => void;
    seekTo: (time: number) => void;
  },
) {
  if (!("mediaSession" in navigator)) return;

  navigator.mediaSession.setActionHandler("play", handlers.play);
  navigator.mediaSession.setActionHandler("pause", handlers.pause);
  navigator.mediaSession.setActionHandler("previoustrack", handlers.previousTrack);
  navigator.mediaSession.setActionHandler("nexttrack", handlers.nextTrack);
  navigator.mediaSession.setActionHandler("seekto", (details) => {
    if (details.seekTime !== undefined) {
      handlers.seekTo(details.seekTime);
    }
  });

  watch(
    currentTrack,
    (track) => {
      if (!track) {
        navigator.mediaSession.metadata = null;
        return;
      }

      const artwork: MediaImage[] = [];
      const coverUrl = getCoverArtUrl(track);
      if (coverUrl) {
        artwork.push({
          src: coverUrl,
          sizes: "300x300",
          type: track.cover_art_mime || "image/jpeg",
        });
      }

      navigator.mediaSession.metadata = new MediaMetadata({
        title: track.title,
        artist: track.artist,
        album: track.album,
        artwork,
      });
    },
    { immediate: true },
  );
}
