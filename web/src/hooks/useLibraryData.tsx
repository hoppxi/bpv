import { useState, useEffect } from "react";
import { LibraryResponse } from "@/types";
import { IDB } from "@/utils/indexedDB";

export function useLibraryData() {
  const [library, setLibrary] = useState<LibraryResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchLibrary = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch("/api/library");

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: LibraryResponse = await response.json();

      if (data.status === "ok") {
        setLibrary(data);
        IDB.setItem("musicLibrary", JSON.stringify(data));
        IDB.setItem("libraryLastUpdated", new Date().toISOString());
      } else {
        throw new Error("Invalid response from server");
      }
    } catch (err) {
      console.error("Failed to fetch library:", err);
      setError(err instanceof Error ? err.message : "Unknown error occurred");

      // Try to load from cache if available
      const cachedLibrary = localStorage.getItem("musicLibrary");
      if (cachedLibrary) {
        try {
          setLibrary(JSON.parse(cachedLibrary));
          setError("Using cached library data - server unavailable");
        } catch (parseError) {
          console.error("Failed to parse cached library:", parseError);
        }
      }
    } finally {
      setLoading(false);
    }
  };

  const refreshLibrary = async () => {
    try {
      setLoading(true);
      const response = await fetch("/api/scan-simple", { method: "POST" });

      if (response.ok) {
        setTimeout(fetchLibrary, 1000);
      } else {
        throw new Error("Scan request failed");
      }
    } catch (err) {
      console.error("Failed to initiate scan:", err);
      setError(err instanceof Error ? err.message : "Scan failed");
    }
  };

  useEffect(() => {
    const lastUpdated = localStorage.getItem("libraryLastUpdated");
    const cachedLibrary = localStorage.getItem("musicLibrary");

    if (cachedLibrary && lastUpdated) {
      const lastUpdatedDate = new Date(lastUpdated);
      const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000);

      if (lastUpdatedDate > oneHourAgo) {
        try {
          setLibrary(JSON.parse(cachedLibrary));
          setLoading(false);
          // Still fetch fresh data in background
          fetchLibrary();
          return;
        } catch (error) {
          console.error("Failed to parse cached library:", error);
        }
      }
    }

    fetchLibrary();
  }, []);

  return { library, loading, error, refreshLibrary, refetch: fetchLibrary };
}
