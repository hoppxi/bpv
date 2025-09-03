import React, { useState, useRef, useEffect } from "react";
import { Volume2, VolumeX, Volume1, Volume } from "lucide-react";
import { VolumeControlProps } from "@/types";
import "@/styles/volume-control.scss";

const VolumeControl: React.FC<VolumeControlProps> = ({
  volume,
  onVolumeChange,
}) => {
  const [isDragging, setIsDragging] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [previousVolume, setPreviousVolume] = useState(volume);
  const [showSlider, setShowSlider] = useState(false);
  const volumeSliderRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const handleVolumeClick = () => {
    if (isMuted) {
      // Unmute and restore previous volume
      setIsMuted(false);
      onVolumeChange(previousVolume);
    } else {
      // Mute and remember current volume
      setIsMuted(true);
      setPreviousVolume(volume);
      onVolumeChange(0);
    }
  };

  const handleSliderClick = (e: React.MouseEvent) => {
    if (!volumeSliderRef.current) return;

    const rect = volumeSliderRef.current.getBoundingClientRect();
    const y = rect.bottom - e.clientY;
    const height = rect.height;
    const newVolume = Math.max(0, Math.min(1, y / height));

    onVolumeChange(newVolume);
    if (newVolume > 0) {
      setIsMuted(false);
    }
  };

  const handleDragStart = (e: React.MouseEvent) => {
    setIsDragging(true);
    handleSliderClick(e);
  };

  const handleDrag = (e: React.MouseEvent) => {
    if (!isDragging) return;
    handleSliderClick(e);
  };

  const handleDragEnd = () => {
    setIsDragging(false);
  };

  // Add global mouse event listeners for dragging
  useEffect(() => {
    const handleGlobalMouseMove = (e: MouseEvent) => {
      if (isDragging && volumeSliderRef.current) {
        const rect = volumeSliderRef.current.getBoundingClientRect();
        const y = rect.bottom - e.clientY;
        const height = rect.height;
        const newVolume = Math.max(0, Math.min(1, y / height));

        onVolumeChange(newVolume);
        if (newVolume > 0) {
          setIsMuted(false);
        }
      }
    };

    const handleGlobalMouseUp = () => {
      if (isDragging) {
        setIsDragging(false);
      }
    };

    if (isDragging) {
      document.addEventListener("mousemove", handleGlobalMouseMove);
      document.addEventListener("mouseup", handleGlobalMouseUp);
    }

    return () => {
      document.removeEventListener("mousemove", handleGlobalMouseMove);
      document.removeEventListener("mouseup", handleGlobalMouseUp);
    };
  }, [isDragging, onVolumeChange]);

  // Handle clicks outside to close slider
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setShowSlider(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const getVolumeIcon = () => {
    if (isMuted || volume === 0) {
      return <VolumeX size={20} />;
    } else if (volume < 0.33) {
      return <Volume size={20} />;
    } else if (volume < 0.66) {
      return <Volume1 size={20} />;
    } else {
      return <Volume2 size={20} />;
    }
  };

  const volumePercentage = isMuted ? 0 : volume * 100;

  return (
    <div className="volume-control" ref={containerRef}>
      <button
        className="volume-control__button"
        onClick={handleVolumeClick}
        onMouseEnter={() => setShowSlider(true)}
        title={isMuted ? "Unmute" : "Mute"}
      >
        {getVolumeIcon()}
      </button>

      {showSlider && (
        <div className="volume-control__slider-container">
          <div
            ref={volumeSliderRef}
            className="volume-control__slider"
            onClick={handleSliderClick}
            onMouseDown={handleDragStart}
            onMouseMove={handleDrag}
            onMouseUp={handleDragEnd}
          >
            <div className="volume-control__slider-background">
              <div
                className="volume-control__slider-fill"
                style={{ height: `${volumePercentage}%` }}
              />
            </div>

            <div
              className="volume-control__slider-handle"
              style={{ bottom: `${volumePercentage}%` }}
              onMouseDown={handleDragStart}
            />
          </div>

          <div className="volume-control__tooltip">
            {Math.round(volumePercentage)}%
          </div>
        </div>
      )}
    </div>
  );
};

export default VolumeControl;
