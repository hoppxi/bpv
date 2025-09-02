import React, { useState, useEffect } from "react";
import { Trash2, Download, Upload, PieChart } from "lucide-react";
import { formatFileSize } from "@/utils";
import "@/styles/modal-tabs.scss";

const StorageTab: React.FC = () => {
  const [storageData, setStorageData] = useState<{
    total: number;
    used: number;
    free: number;
    items: Array<{ key: string; size: number; lastModified: string }>;
  }>({
    total: 0,
    used: 0,
    free: 0,
    items: [],
  });

  useEffect(() => {
    calculateStorageUsage();
  }, []);

  const calculateStorageUsage = () => {
    let totalSize = 0;
    const items: Array<{ key: string; size: number; lastModified: string }> =
      [];

    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key) {
        const value = localStorage.getItem(key) || "";
        const size = new Blob([value]).size;
        totalSize += size;

        items.push({
          key,
          size,
          lastModified: new Date().toLocaleDateString(),
        });
      }
    }

    // Estimate total storage (usually 5MB for most browsers)
    const totalStorage = 5 * 1024 * 1024; // 5MB in bytes
    const freeStorage = totalStorage - totalSize;

    setStorageData({
      total: totalStorage,
      used: totalSize,
      free: freeStorage,
      items: items.sort((a, b) => b.size - a.size),
    });
  };

  const handleExportData = () => {
    const data: Record<string, string> = {};
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key) {
        data[key] = localStorage.getItem(key) || "";
      }
    }

    const blob = new Blob([JSON.stringify(data, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "bpv-music-player-backup.json";
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleImportData = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const data = JSON.parse(e.target?.result as string);
        if (
          window.confirm("This will overwrite your current data. Continue?")
        ) {
          localStorage.clear();
          Object.entries(data).forEach(([key, value]) => {
            localStorage.setItem(key, value as string);
          });
          calculateStorageUsage();
          alert("Data imported successfully!");
        }
      } catch (error) {
        alert("Failed to import data: Invalid file format");
      }
    };
    reader.readAsText(file);
  };

  const handleClearItem = (key: string) => {
    if (window.confirm(`Delete "${key}"?`)) {
      localStorage.removeItem(key);
      calculateStorageUsage();
    }
  };

  const handleClearAll = () => {
    if (window.confirm("Clear all stored data? This cannot be undone.")) {
      localStorage.clear();
      calculateStorageUsage();
    }
  };

  const usagePercentage = (storageData.used / storageData.total) * 100;

  return (
    <div className="tab-content">
      <div className="tab-content__header">
        <div className="tab-content__stats">
          <h3>Storage</h3>
          <p>Manage your local storage data</p>
        </div>
        <div className="tab-content__actions">
          <button
            className="tab-content__action-btn"
            onClick={handleExportData}
          >
            <Download size={16} />
            Export
          </button>
          <label className="tab-content__action-btn">
            <Upload size={16} />
            Import
            <input
              type="file"
              accept=".json"
              onChange={handleImportData}
              style={{ display: "none" }}
            />
          </label>
        </div>
      </div>

      <div className="tab-content__list">
        <div className="storage-content">
          {/* Storage Overview */}
          <div className="storage-overview">
            <div className="storage-overview__chart">
              <div className="storage-chart">
                <div
                  className="storage-chart__fill"
                  style={{ height: `${usagePercentage}%` }}
                />
                <PieChart size={48} />
              </div>
            </div>
            <div className="storage-overview__stats">
              <div className="storage-stat">
                <span className="storage-stat__label">Used:</span>
                <span className="storage-stat__value">
                  {formatFileSize(storageData.used)}
                </span>
              </div>
              <div className="storage-stat">
                <span className="storage-stat__label">Free:</span>
                <span className="storage-stat__value">
                  {formatFileSize(storageData.free)}
                </span>
              </div>
              <div className="storage-stat">
                <span className="storage-stat__label">Total:</span>
                <span className="storage-stat__value">
                  {formatFileSize(storageData.total)}
                </span>
              </div>
              <div className="storage-stat">
                <span className="storage-stat__label">Usage:</span>
                <span className="storage-stat__value">
                  {usagePercentage.toFixed(1)}%
                </span>
              </div>
            </div>
          </div>

          {/* Storage Items */}
          <div className="storage-items">
            <h4 className="storage-items__title">Stored Items</h4>
            <div className="storage-items__list">
              {storageData.items.map((item) => (
                <div key={item.key} className="storage-item">
                  <div className="storage-item__info">
                    <div className="storage-item__name">{item.key}</div>
                    <div className="storage-item__details">
                      {formatFileSize(item.size)} • {item.lastModified}
                    </div>
                  </div>
                  <button
                    className="storage-item__delete"
                    onClick={() => handleClearItem(item.key)}
                    title="Delete item"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              ))}
            </div>
          </div>

          {/* Clear All Button */}
          {storageData.items.length > 0 && (
            <div className="storage-actions">
              <button className="storage-clear-btn" onClick={handleClearAll}>
                <Trash2 size={16} />
                Clear All Data
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default StorageTab;
