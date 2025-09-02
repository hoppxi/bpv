import { useState, useEffect, useRef, useCallback } from "react";

export function useAudioVisualizer(
  audioElement: HTMLAudioElement | null,
  isPlaying: boolean
) {
  const [frequencyData, setFrequencyData] = useState<Uint8Array>(
    new Uint8Array(64)
  );
  const [timeDomainData, setTimeDomainData] = useState<Uint8Array>(
    new Uint8Array(64)
  );
  const animationRef = useRef<number | null>(null);

  // Web Audio API references
  const audioContextRef = useRef<AudioContext | null>(null);
  const analyserRef = useRef<AnalyserNode | null>(null);
  const sourceRef = useRef<MediaElementAudioSourceNode | null>(null);

  // Initialize Web Audio API
  const initAudioAnalysis = useCallback(() => {
    if (!audioElement || !window.AudioContext) {
      console.log("Web Audio API not supported");
      return false;
    }

    try {
      // Clean up previous context if exists
      if (
        audioContextRef.current &&
        audioContextRef.current.state !== "closed"
      ) {
        audioContextRef.current.close();
      }

      // Create new audio context
      const context = new (window.AudioContext ||
        (window as any).webkitAudioContext)();
      audioContextRef.current = context;

      // Create analyser node
      const analyser = context.createAnalyser();
      analyser.fftSize = 512; // Increased for better frequency resolution
      analyser.smoothingTimeConstant = 0.6; // Smoother transitions
      analyserRef.current = analyser;

      // Create source from audio element
      const source = context.createMediaElementSource(audioElement);
      sourceRef.current = source;

      // Connect nodes: source -> analyser -> destination
      // This is the critical part - we MUST keep the connection to destination
      source.connect(analyser);
      analyser.connect(context.destination); // This ensures audio continues to play

      console.log("Audio analysis initialized successfully");
      return true;
    } catch (error) {
      console.error("Failed to initialize audio analysis:", error);
      return false;
    }
  }, [audioElement]);

  // Update visualizer data
  const updateVisualizerData = useCallback(() => {
    if (!analyserRef.current) return;

    try {
      const analyser = analyserRef.current;

      // Get frequency data
      const freqData = new Uint8Array(analyser.frequencyBinCount);
      analyser.getByteFrequencyData(freqData);
      setFrequencyData(freqData);

      // Get time domain data
      const timeData = new Uint8Array(analyser.frequencyBinCount);
      analyser.getByteTimeDomainData(timeData);
      setTimeDomainData(timeData);

      animationRef.current = requestAnimationFrame(updateVisualizerData);
    } catch (error) {
      console.warn("Error updating visualizer data:", error);
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
        animationRef.current = null;
      }
    }
  }, []);

  // Handle audio element changes and playback state
  useEffect(() => {
    if (!audioElement || !isPlaying) {
      // Stop animation when not playing
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
        animationRef.current = null;
      }
      return;
    }

    // Initialize audio analysis when we have an audio element and it's playing
    const initialize = async () => {
      try {
        const success = initAudioAnalysis();
        if (success) {
          // Wait a bit for the audio context to be ready
          await new Promise((resolve) => setTimeout(resolve, 100));
          animationRef.current = requestAnimationFrame(updateVisualizerData);
        }
      } catch (error) {
        console.error("Failed to initialize visualizer:", error);
      }
    };

    initialize();

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
        animationRef.current = null;
      }
    };
  }, [audioElement, isPlaying, initAudioAnalysis, updateVisualizerData]);

  // Clean up on unmount
  useEffect(() => {
    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
        animationRef.current = null;
      }

      // Clean up audio context
      if (
        audioContextRef.current &&
        audioContextRef.current.state !== "closed"
      ) {
        audioContextRef.current.close().catch(console.error);
      }
    };
  }, []);

  return { frequencyData, timeDomainData };
}
