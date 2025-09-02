import React, { useState, useRef, useEffect } from "react";
import { formatTime } from "@/utils";
import "@/styles/progress-bar.scss";

interface ProgressBarProps {
  currentTime: number;
  duration: number;
  onSeek: (time: number) => void;
}

const ProgressBar: React.FC<ProgressBarProps> = ({
  currentTime,
  duration,
  onSeek,
}) => {
  const [isDragging, setIsDragging] = useState(false);
  const [hoverTime, setHoverTime] = useState<number | null>(null);
  const progressBarRef = useRef<HTMLDivElement>(null);

  const percentage = duration > 0 ? (currentTime / duration) * 100 : 0;
  const formattedCurrentTime = formatTime(currentTime);
  const formattedDuration = formatTime(duration);

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!progressBarRef.current) return;

    const rect = progressBarRef.current.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const width = rect.width;
    const hoverPercentage = Math.max(0, Math.min(100, (x / width) * 100));
    const hoverTimeValue = (hoverPercentage / 100) * duration;

    setHoverTime(hoverTimeValue);
  };

  const handleMouseLeave = () => {
    setHoverTime(null);
  };

  const handleClick = (e: React.MouseEvent) => {
    if (!progressBarRef.current) return;

    const rect = progressBarRef.current.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const width = rect.width;
    const newPercentage = Math.max(0, Math.min(100, (x / width) * 100));
    const newTime = (newPercentage / 100) * duration;

    onSeek(newTime);
  };

  const handleDragStart = (e: React.MouseEvent) => {
    setIsDragging(true);
    handleClick(e);
  };

  const handleDrag = (e: React.MouseEvent) => {
    if (!isDragging) return;
    handleClick(e);
  };

  const handleDragEnd = () => {
    setIsDragging(false);
  };

  // Add global mouse event listeners for dragging
  useEffect(() => {
    const handleGlobalMouseMove = (e: MouseEvent) => {
      if (isDragging && progressBarRef.current) {
        const rect = progressBarRef.current.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const width = rect.width;
        const newPercentage = Math.max(0, Math.min(100, (x / width) * 100));
        const newTime = (newPercentage / 100) * duration;

        onSeek(newTime);
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
  }, [isDragging, duration, onSeek]);

  return (
    <div className="progress-bar">
      <span className="progress-bar__time progress-bar__time--current">
        {formattedCurrentTime}
      </span>

      <div
        ref={progressBarRef}
        className="progress-bar__container"
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
        onClick={handleClick}
        onMouseDown={handleDragStart}
        onMouseMoveCapture={handleDrag}
        onMouseUp={handleDragEnd}
      >
        <div className="progress-bar__background">
          <div
            className="progress-bar__progress"
            style={{ width: `${percentage}%` }}
          />

          {hoverTime !== null && !isDragging && (
            <div
              className="progress-bar__hover-indicator"
              style={{ left: `${(hoverTime / duration) * 100}%` }}
            >
              <div className="progress-bar__tooltip">
                {formatTime(hoverTime)}
              </div>
            </div>
          )}

          <div
            className="progress-bar__handle"
            style={{ left: `${percentage}%` }}
            onMouseDown={handleDragStart}
          />
        </div>
      </div>

      <span className="progress-bar__time progress-bar__time--duration">
        {formattedDuration}
      </span>
    </div>
  );
};

export default ProgressBar;
