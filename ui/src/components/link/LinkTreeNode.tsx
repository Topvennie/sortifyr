import { useLinkAnchor } from "@/lib/hooks/useLinkAnchor";
import { Directory } from "@/lib/types/directory";
import { Side } from "@/lib/types/general";
import { Group, Stack } from "@mantine/core";
import { useMemo, useState } from "react";
import { LuCircle, LuFolderOpen } from "react-icons/lu";
import { LinkTreePlaylist } from "./LinkTreePlaylist";
import { getLinkDirectoryId } from "./util";

type Props = {
  directory: Directory;
  side: Side;
  level: number;
}

const offSet = 32;

export const LinkTreeNode = ({ directory, side, level }: Props) => {
  const { visibleAnchorsRef, registerAnchor, startConnection, finishConnection, connections, layoutVersion, hoveredConnection } = useLinkAnchor()

  const id = getLinkDirectoryId(directory, side)

  const isHoveredConnection = side === "left" ? hoveredConnection?.from === id : hoveredConnection?.to === id
  const [isHovered, setIsHovered] = useState(false)

  const hidden = useMemo(() => {
    let cons: string[]
    if (side === "left") cons = connections.filter(c => c.from === id).map(c => c.to)
    else cons = connections.filter(c => c.to === id).map(c => c.from)

    return cons.filter(c => !visibleAnchorsRef.current[c]).length
  }, [connections, layoutVersion, side]) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Stack
      gap={2}
      className="cursor-pointer"
    >
      <Group
        py={2}
        px="xs"
        gap="xs"
        style={{ marginLeft: level * offSet }}
        onMouseDown={() => startConnection(id)}
        onMouseUp={() => finishConnection(id)}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        className={`rounded-xl ${isHovered && "bg-(--mantine-color-background-1)"}`}
      >
        <LuFolderOpen className="size-6" />
        <p className="font-semibold">{directory.name}</p>
        <p className="text-muted text-sm">{directory.playlists.length + (directory.children?.length ?? 0)}</p>
        {hidden > 0 && <p className="text-red-500 text-sm">{`(${hidden})`}</p>}
        <div ref={el => registerAnchor(id, { el, side, directory })} className={`absolute ${side === "left" ? "right-1" : "left-1"}`}>
          <LuCircle className={`${isHoveredConnection && "text-red-500"}`} />
        </div>
      </Group>
      <Stack gap={2} style={{ marginLeft: (level + 1) * offSet + 5 }}>
        {directory.playlists.map(p => <LinkTreePlaylist key={p.id} playlist={p} side={side} />)}
      </Stack>
      {directory.children?.map(d => (
        <LinkTreeNode
          key={d.id}
          directory={d}
          side={side}
          level={level + 1}
        />
      ))}
    </Stack>
  )
}
