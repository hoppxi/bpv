import React from "react";
import { VisualizerType, SettingsTabProps } from "@/types";
import {
  BarChart3,
  Waves,
  Sparkles,
  Circle,
  ToggleLeft,
  ToggleRight,
  RefreshCw,
  Heart,
} from "lucide-react";
import "@/styles/modal-tabs.scss";

const SettingsTab: React.FC<SettingsTabProps> = ({
  visualizerType,
  shuffle,
  repeat,
  onVisualizerChange,
  onShuffleChange,
  onRepeatChange,
}) => {
  const visualizerOptions: {
    type: VisualizerType;
    label: string;
    icon: React.ReactNode;
  }[] = [
    { type: "bars", label: "Bars", icon: <BarChart3 size={20} /> },
    { type: "wave", label: "Wave", icon: <Waves size={20} /> },
    { type: "particles", label: "Particles", icon: <Sparkles size={20} /> },
    { type: "circle", label: "Circle", icon: <Circle size={20} /> },
    { type: "sphere", label: "Sphere", icon: <Circle size={20} /> },
    { type: "lines", label: "Lines", icon: <BarChart3 size={20} /> },
    { type: "mesh", label: "Mesh", icon: <BarChart3 size={20} /> },
    { type: "radial", label: "Radial", icon: <Circle size={20} /> },
    { type: "spectrum", label: "Spectrum", icon: <BarChart3 size={20} /> },
    { type: "orb", label: "Orb", icon: <Circle size={20} /> },
  ];

  const handleClearCache = () => {
    if (
      window.confirm(
        "Clear all cached data? This will remove your library cache and settings."
      )
    ) {
      localStorage.removeItem("musicLibrary");
      localStorage.removeItem("libraryLastUpdated");
      localStorage.removeItem("currentTrack");
      localStorage.removeItem("playbackState");
      window.location.reload();
    }
  };

  return (
    <div className="tab-content">
      <div className="tab-content__header">
        <div className="tab-content__stats">
          <h3>Settings</h3>
          <p>Customize your music player experience</p>
        </div>
      </div>

      <div className="tab-content__list">
        <div className="settings-list">
          {/* Visualizer Settings */}
          <div className="settings-group">
            <h4 className="settings-group__title">Visualizer</h4>
            <div className="visualizer-grid">
              {visualizerOptions.map((option) => (
                <button
                  key={option.type}
                  className={`visualizer-option ${
                    visualizerType === option.type
                      ? "visualizer-option--active"
                      : ""
                  }`}
                  onClick={() => onVisualizerChange(option.type)}
                >
                  <div className="visualizer-option__icon">{option.icon}</div>
                  <div className="visualizer-option__label">{option.label}</div>
                </button>
              ))}
            </div>
          </div>

          {/* Playback Settings */}
          <div className="settings-group">
            <h4 className="settings-group__title">Playback</h4>
            <div className="settings-item">
              <div className="settings-item__info">
                <div className="settings-item__label">Shuffle</div>
                <div className="settings-item__description">
                  Play songs in random order
                </div>
              </div>
              <button
                className={`settings-toggle ${
                  shuffle ? "settings-toggle--active" : ""
                }`}
                onClick={() => onShuffleChange(!shuffle)}
              >
                {shuffle ? <ToggleRight size={24} /> : <ToggleLeft size={24} />}
              </button>
            </div>

            <div className="settings-item">
              <div className="settings-item__info">
                <div className="settings-item__label">Repeat</div>
                <div className="settings-item__description">
                  Loop the current playlist
                </div>
              </div>
              <button
                className={`settings-toggle ${
                  repeat ? "settings-toggle--active" : ""
                }`}
                onClick={() => onRepeatChange(!repeat)}
              >
                {repeat ? <ToggleRight size={24} /> : <ToggleLeft size={24} />}
              </button>
            </div>
          </div>

          {/* Cache Settings */}
          <div className="settings-group">
            <h4 className="settings-group__title">Storage</h4>
            <div className="settings-item">
              <div className="settings-item__info">
                <div className="settings-item__label">Clear Cache</div>
                <div className="settings-item__description">
                  Remove cached library data and settings
                </div>
              </div>
              <button className="settings-button" onClick={handleClearCache}>
                <RefreshCw size={18} />
                Clear Data
              </button>
            </div>
          </div>

          {/* App Info */}
          <div className="settings-group">
            <h4 className="settings-group__title">Information</h4>
            <div className="settings-info">
              <div className="settings-info__item">
                <span className="settings-info__label">BPV:</span>
                <span className="settings-info__value">
                  Browser based local music player
                </span>
              </div>
              <div className="settings-info__item">
                <span className="settings-info__label">Version:</span>
                <span className="settings-info__value">0.1.0</span>
              </div>
              <div className="settings-info__item">
                <span className="settings-info__label">Built with:</span>
                <span className="settings-info__value">
                  React + TypeScript + Go and
                  <Heart />
                </span>
              </div>
              <div className="settings-info__item">
                <span className="settings-info__label">Audio formats:</span>
                <span className="settings-info__value">
                  MP3, FLAC, WAV, AAC, M4A
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SettingsTab;
