import { Playlist as PlaylistType } from "@/lib/types/playlist";
import { cn } from "@/lib/utils";
import { useDraggable } from "@dnd-kit/core";
import { useAutoAnimate } from "@formkit/auto-animate/react";
import { Stack, StackProps } from "@mantine/core";
import { LuGripVertical } from "react-icons/lu";
import { LoadingSpinner } from "../molecules/LoadingSpinner";
import { PlaylistCover } from "../playlist/PlaylistCover";

type Props = {
  playlists: PlaylistType[];
  isLoading: boolean;
} & StackProps

export const DirectoryPlaylistSelector = ({ playlists, isLoading, className, ...props }: Props) => {
  const [bodyRef] = useAutoAnimate<HTMLTableSectionElement>();

  return (
    <Stack ref={bodyRef} p="md" gap="lg" justify="start" bg="background.0" className={cn("rounded-xl overflow-auto select-none", className)} {...props}>
      {isLoading
        ? <LoadingSpinner />
        : playlists.map(p => <Playlist key={p.id} playlist={p} />)
      }
    </Stack>
  )
}

type PlaylistProps = {
  playlist: PlaylistType;
}

const Playlist = ({ playlist }: PlaylistProps) => {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({ id: playlist.id })

  return (
    <div className={`flex items-center justify-between gap-2 ${isDragging && "invisible"}`}>
      <div ref={setNodeRef} className={`cursor-grab flex items-center gap-2`} {...listeners} {...attributes} >
        <LuGripVertical className="text-muted shrink-0" />
        <PlaylistCover playlist={playlist} />
        <p className="flex-1 line-clamp-2 overflow-hidden text-ellipsis break-words">
          {playlist.name}
        </p>
      </div>
      <p className={`ml-auto text-muted shrink-0 whitespace-nowrap`}>
        {playlist.trackAmount}
      </p>
    </div>
  )
}

type DirectoryPlaylistDraggableProps = {
  playlist: PlaylistType;
}

export const DirectoryPlaylistDraggable = ({ playlist }: DirectoryPlaylistDraggableProps) => {
  return (
    <div className="flex items-center gap-2">
      <LuGripVertical className="text-muted flex-shrink-0" />
      <PlaylistCover playlist={playlist} />
      <p className="flex-1 line-clamp-2 break-words">
        {playlist.name}
      </p>
    </div>
  )
}
