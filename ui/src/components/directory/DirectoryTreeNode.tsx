import { DirectorySchema } from "@/lib/types/directory";
import { ActionIcon, Group, Stack } from "@mantine/core";
import { useState } from "react";
import { LuFolder, LuFolderOpen, LuTriangle } from "react-icons/lu";
import { DirectoryTreePlaylist } from "./DirectoryTreePlaylist";

type Props = {
  directory: DirectorySchema;
  level: number;
}

const offSet = 32;

export const DirectoryTreeNode = ({ directory, level }: Props) => {
  const [expanded, setExpanded] = useState(true)

  return (
    <Stack gap={2}>
      <Group style={{ marginLeft: level * offSet }}>
        <Group py={2} px="xs" gap="xs" onClick={() => setExpanded(prev => !prev)} className="cursor-pointer hover:bg-(--mantine-color-background-1) rounded-xl">
          <ActionIcon variant="transparent" color="black" size="sm">
            <LuTriangle className={`fill-black w-2 ${expanded && "rotate-180"}`} />
          </ActionIcon>
          {expanded
            ? <LuFolderOpen className="size-6" />
            : <LuFolder className="size-6" />
          }
          <p className="font-semibold">{directory.name}</p>
          <p className="text-muted text-sm">{directory.playlists.length + (directory.children?.length ?? 0)}</p>
        </Group>
      </Group>
      {expanded && (
        <>
          <Stack gap={2} style={{ marginLeft: (level + 1) * offSet + 5 }}>
            {directory.playlists.map(p => <DirectoryTreePlaylist key={p.id} playlist={p} />)}
          </Stack>
          {directory.children?.map(d => (
            <DirectoryTreeNode
              key={d.iid}
              directory={d}
              level={level + 1}
            />
          ))}
        </>
      )}
    </Stack>
  )
}
