import { TrackGroupGridProps } from "@/types";
import { Play } from "lucide-react";
import React from "react";

const TrackGroupGrid: React.FC<TrackGroupGridProps> = ({
  metadata,
  group,
  handlePlayGroup,
  onclick,
}: TrackGroupGridProps) => {
  return (
    <div className="grid-list">
      {group.map(([track, count]) => (
        <div key={count} className="grid-item" onClick={() => onclick(track)}>
          <div className="grid-item__icon">{metadata.icon}</div>
          <div className="grid-item__info">
            <div className="grid-item__title">{track}</div>
            <div className="grid-item__subtitle">
              {count + "song" + (count == 1 ? "" : "s")}
            </div>
          </div>
          <button
            className="grid-item__action"
            onClick={(e) => {
              e.stopPropagation();
              handlePlayGroup(track);
            }}
            title={`Play ${metadata.name}`}
          >
            <Play size={16} />
          </button>
        </div>
      ))}
    </div>
  );
};

export default TrackGroupGrid;
