import React, { useState } from "react";
import { TabType, ModalProps } from "@/types";
import LibraryTab from "./LibraryTab";
import ArtistsTab from "./ArtistsTab";
import AlbumsTab from "./AlbumsTab";
import GenresTab from "./GenresTab";
import SearchTab from "./SearchTab";
import SettingsTab from "./SettingsTab";
import StorageTab from "./StorageTab";
import {
  X,
  Library,
  Users,
  Disc,
  Tag,
  Search,
  Settings,
  Database,
} from "lucide-react";
import "@/styles/modal.scss";

const Modal: React.FC<ModalProps> = ({
  isOpen,
  onClose,
  library,
  currentTrack,
  visualizerType,
  shuffle,
  repeat,
  onPlayTrack,
  onVisualizerChange,
  onShuffleChange,
  onRepeatChange,
  onRefreshLibrary,
}) => {
  const [activeTab, setActiveTab] = useState<TabType>("artists");
  const [searchQuery, setSearchQuery] = useState("");

  if (!isOpen) return null;

  const tabs = [
    { id: "library" as TabType, label: "Library", icon: Library },
    { id: "artists" as TabType, label: "Artists", icon: Users },
    { id: "albums" as TabType, label: "Albums", icon: Disc },
    { id: "genres" as TabType, label: "Genres", icon: Tag },
    { id: "search" as TabType, label: "Search", icon: Search },
    { id: "settings" as TabType, label: "Settings", icon: Settings },
    { id: "storage" as TabType, label: "Storage", icon: Database },
  ];

  const renderTabContent = () => {
    switch (activeTab) {
      case "library":
        return (
          <LibraryTab
            library={library}
            currentTrack={currentTrack}
            onPlayTrack={onPlayTrack}
            onRefreshLibrary={onRefreshLibrary}
          />
        );
      case "artists":
        return (
          <ArtistsTab
            library={library}
            currentTrack={currentTrack}
            onPlayTrack={onPlayTrack}
          />
        );
      case "albums":
        return (
          <AlbumsTab
            library={library}
            currentTrack={currentTrack}
            onPlayTrack={onPlayTrack}
          />
        );
      case "genres":
        return (
          <GenresTab
            library={library}
            currentTrack={currentTrack}
            onPlayTrack={onPlayTrack}
          />
        );
      case "search":
        return (
          <SearchTab
            library={library}
            currentTrack={currentTrack}
            searchQuery={searchQuery}
            onSearchChange={setSearchQuery}
            onPlayTrack={onPlayTrack}
          />
        );
      case "settings":
        return (
          <SettingsTab
            visualizerType={visualizerType}
            shuffle={shuffle}
            repeat={repeat}
            onVisualizerChange={onVisualizerChange}
            onShuffleChange={onShuffleChange}
            onRepeatChange={onRepeatChange}
          />
        );
      case "storage":
        return <StorageTab />;
      default:
        return null;
    }
  };

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div className="modal-overlay" onClick={handleOverlayClick}>
      <div className="modal">
        {/* Header */}
        <div className="modal__header">
          <h2 className="modal__title">Music Library</h2>
          <button className="modal__close-btn" onClick={onClose}>
            <X size={24} />
          </button>
        </div>

        {/* Tabs */}
        <div className="modal__tabs">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                className={`modal__tab ${
                  activeTab === tab.id ? "modal__tab--active" : ""
                }`}
                onClick={() => setActiveTab(tab.id)}
              >
                <Icon size={18} />
                <span>{tab.label}</span>
              </button>
            );
          })}
        </div>

        {/* Content */}
        <div className="modal__content">{renderTabContent()}</div>
      </div>
    </div>
  );
};

export default Modal;
