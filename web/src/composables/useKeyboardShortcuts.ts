import { onMounted, onUnmounted } from "vue";

interface ShortcutHandlers {
  togglePlayPause: () => void;
  nextTrack: () => void;
  previousTrack: () => void;
  volumeUp: () => void;
  volumeDown: () => void;
  toggleMute: () => void;
  seekForward: () => void;
  seekBackward: () => void;
  toggleShuffle: () => void;
  toggleRepeat: () => void;
  openSearch: () => void;
}

export function useKeyboardShortcuts(handlers: ShortcutHandlers) {
  function handleKeydown(e: KeyboardEvent) {
    if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
      return;
    }

    switch (e.code) {
      case "Space":
        e.preventDefault();
        handlers.togglePlayPause();
        break;
      case "ArrowRight":
        if (e.shiftKey) {
          handlers.nextTrack();
        } else {
          handlers.seekForward();
        }
        break;
      case "ArrowLeft":
        if (e.shiftKey) {
          handlers.previousTrack();
        } else {
          handlers.seekBackward();
        }
        break;
      case "ArrowUp":
        e.preventDefault();
        handlers.volumeUp();
        break;
      case "ArrowDown":
        e.preventDefault();
        handlers.volumeDown();
        break;
      case "KeyM":
        handlers.toggleMute();
        break;
      case "KeyS":
        if (!e.ctrlKey && !e.metaKey) {
          handlers.toggleShuffle();
        }
        break;
      case "KeyR":
        if (!e.ctrlKey && !e.metaKey) {
          handlers.toggleRepeat();
        }
        break;
      case "Slash":
        e.preventDefault();
        handlers.openSearch();
        break;
      case "KeyN":
        handlers.nextTrack();
        break;
      case "KeyP":
        handlers.previousTrack();
        break;
    }
  }

  onMounted(() => {
    window.addEventListener("keydown", handleKeydown);
  });

  onUnmounted(() => {
    window.removeEventListener("keydown", handleKeydown);
  });
}
