import { useLinkAnchor } from "@/lib/hooks/useLinkAnchor";
import { Side } from "@/lib/types/general";
import { Playlist } from "@/lib/types/playlist";
import { useMemo, useState } from "react";
import { LuCircle } from "react-icons/lu";
import { PlaylistCover } from "../playlist/PlaylistCover";
import { getLinkPlaylistId } from "./util";

type Props = {
  playlist: Playlist;
  side: Side;
}

export const LinkTreePlaylist = ({ playlist, side }: Props) => {
  const { visibleAnchorsRef, registerAnchor, startConnection, finishConnection, hoveredConnection, connections, layoutVersion } = useLinkAnchor()

  const id = getLinkPlaylistId(playlist, side)

  const isHoveredConnection = side === "left" ? hoveredConnection?.from === id : hoveredConnection?.to === id
  const [isHovered, setIsHovered] = useState(false)

  const hidden = useMemo(() => {
    let cons: string[]
    if (side === "left") cons = connections.filter(c => c.from === id).map(c => c.to)
    else cons = connections.filter(c => c.to === id).map(c => c.from)

    return cons.filter(c => !visibleAnchorsRef.current[c]).length
  }, [connections, layoutVersion, side]) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div
      onMouseDown={() => startConnection(id)}
      onMouseUp={() => finishConnection(id)}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      className={`flex items-center gap-2 rounded-xl ${isHovered && "bg-(--mantine-color-background-1)"}`}
    >
      <PlaylistCover playlist={playlist} />
      <p className="line-clamp-1 overflow-hidden text-ellipsis wrap-break-word">{playlist.name}</p>
      <p className="text-muted text-sm">{playlist.trackAmount}</p>
      {hidden > 0 && <p className="text-red-500 text-sm">{`(${hidden})`}</p>}
      <div ref={el => registerAnchor(id, { el, side, playlist })} className={`absolute ${side === "left" ? "right-1" : "left-1"}`}>
        <LuCircle className={`${isHoveredConnection && "text-red-500"}`} />
      </div>
    </div>
  )
}
