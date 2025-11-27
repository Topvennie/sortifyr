import { PlaylistSchema } from "@/lib/types/playlist";
import { cn } from "@/lib/utils";
import { ActionIcon } from "@mantine/core";
import { ComponentProps } from "react";
import { FaRegCirclePlay, FaTrashCan } from "react-icons/fa6";

type Props = {
  playlist: PlaylistSchema;
  onDelete: (playlist: PlaylistSchema) => void;
} & ComponentProps<"div">

export const DirectoryPlaylist = ({ playlist, onDelete, className, ...props }: Props) => {
  return (
    <div className={cn("flex justify-between rounded-md bg-white p-4", className)} {...props}>
      <div className="flex items-center gap-2">
        <FaRegCirclePlay className="text-green-500 w-8" />
        <span>{playlist.name}</span>
        <span className="text-muted text-sm">{playlist.tracks}</span>
      </div>
      <ActionIcon onClick={() => onDelete(playlist)} color="red" variant="subtle">
        <FaTrashCan />
      </ActionIcon>
    </div>
  )
}
