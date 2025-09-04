import { formatTime, formatFileSize } from "@/utils";
import { BarChart2, Music2 } from "lucide-react";
import { LibraryResponse, AudioFile, TrackRowProps } from "@/types";

const TrackRow: React.FC<TrackRowProps> = ({
  library,
  currentTrack,
  onPlayTrack,
  index,
}: TrackRowProps) => {
  const track =
    "files" in library
      ? (library as LibraryResponse).files[index]
      : (library as AudioFile[])[index];
  const isActive = currentTrack?.file_path === track.file_path;

  return (
    <div
      key={track.file_path}
      className={`song-item ${isActive ? "song-item--active" : ""}`}
      onClick={() => onPlayTrack(track)}
    >
      <div className="song-item__index">
        {isActive ? (
          <div className="song-item__playing-indicator">
            <BarChart2 />
          </div>
        ) : (
          <div className="song-item__number">{index + 1}</div>
        )}

        {track?.cover_art ? (
          <img
            src={`data:${track?.cover_art_mime};base64,${track?.cover_art}`}
            alt={track?.album}
            className="song-item__cover"
          />
        ) : (
          <div className="song-item__cover-placeholder">
            <Music2 />
          </div>
        )}
      </div>
      <div className="song-item__info">
        <div className="song-item__title">{track.title}</div>
        <div className="song-item__artist">{track.artist}</div>
      </div>
      <div className="song-item__album">{track.album}</div>
      <div className="song-item__duration">
        {formatTime(parseInt(track.duration) || 0)}
      </div>
      <div className="song-item__size">{formatFileSize(track.file_size)}</div>
    </div>
  );
};
export default TrackRow;
