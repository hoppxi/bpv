import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { AudioFile, LibraryResponse } from "@/types";
import { api } from "@/utils/api";

interface LibraryContextType {
  library: LibraryResponse | null;
  loading: boolean;
  error: string | null;
  refreshLibrary: () => Promise<void>;
  scanLibrary: () => Promise<void>;
  getArtistSongs: (artist: string) => AudioFile[];
  getAlbumSongs: (album: string) => AudioFile[];
  searchSongs: (query: string) => AudioFile[];
}

const LibraryContext = createContext<LibraryContextType | null>(null);

export function LibraryProvider({ children }: { children: ReactNode }) {
  const [library, setLibrary] = useState<LibraryResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    refreshLibrary();
  }, []);

  const refreshLibrary = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getLibrary();
      setLibrary(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load library");
      console.error("Error loading library:", err);
    } finally {
      setLoading(false);
    }
  };

  const scanLibrary = async () => {
    try {
      setLoading(true);
      await api.scanLibrary();
      // Wait a bit before refreshing to allow scan to start
      setTimeout(refreshLibrary, 1000);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to scan library");
      console.error("Error scanning library:", err);
    }
  };

  const getArtistSongs = (artist: string): AudioFile[] => {
    if (!library) return [];
    return library.files.filter((file) => file.artist === artist);
  };

  const getAlbumSongs = (album: string): AudioFile[] => {
    if (!library) return [];
    return library.files.filter((file) => file.album === album);
  };

  const searchSongs = (query: string): AudioFile[] => {
    if (!library) return [];
    const lowerQuery = query.toLowerCase();
    return library.files.filter(
      (file) =>
        file.title.toLowerCase().includes(lowerQuery) ||
        file.artist.toLowerCase().includes(lowerQuery) ||
        file.album.toLowerCase().includes(lowerQuery) ||
        (file.genre && file.genre.toLowerCase().includes(lowerQuery))
    );
  };

  const value: LibraryContextType = {
    library,
    loading,
    error,
    refreshLibrary,
    scanLibrary,
    getArtistSongs,
    getAlbumSongs,
    searchSongs,
  };

  return (
    <LibraryContext.Provider value={value}>{children}</LibraryContext.Provider>
  );
}

export function useLibrary() {
  const context = useContext(LibraryContext);
  if (!context) {
    throw new Error("useLibrary must be used within a LibraryProvider");
  }
  return context;
}
