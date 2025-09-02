// NOT IMPLMENTED OR TESTED

import { useState, useEffect, useCallback } from "react";
import { extractDominantColor, generateColorPalette } from "@/utils/colorUtils";
import { AudioFile } from "@/types";

interface UseColorExtractorProps {
  track: AudioFile | null;
}

export function useColorExtractor({ track }: UseColorExtractorProps) {
  const [dominantColor, setDominantColor] = useState<string>("#ff6b6b");
  const [colorPalette, setColorPalette] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const extractColors = useCallback(async (audioFile: AudioFile) => {
    if (!audioFile.cover_art) {
      // Generate random color if no cover art
      const randomColor = `hsl(${Math.random() * 360}, 70%, 60%)`;
      setDominantColor(randomColor);
      setColorPalette(generateColorPalette(randomColor));
      return;
    }

    setIsLoading(true);
    try {
      const color = await extractDominantColor(
        `data:${audioFile.cover_art_mime};base64,${audioFile.cover_art}`
      );
      setDominantColor(color);
      setColorPalette(generateColorPalette(color));
    } catch (error) {
      console.error("Failed to extract color:", error);
      // Fallback to random color
      const randomColor = `hsl(${Math.random() * 360}, 70%, 60%)`;
      setDominantColor(randomColor);
      setColorPalette(generateColorPalette(randomColor));
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    if (track) {
      extractColors(track);
    }
  }, [track, extractColors]);

  return {
    dominantColor,
    colorPalette,
    isLoading,
    extractColors,
  };
}
