import { DirectorySchema } from "@/lib/types/directory";
import { PlaylistSchema } from "@/lib/types/playlist";
import { getUuid } from "@/lib/utils";
import { useDroppable } from "@dnd-kit/core";
import { ActionIcon, Group, Stack, TextInput } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useState } from "react";
import { LuFolder, LuFolderOpen, LuFolderPlus, LuPencilLine, LuTrash2, LuTriangle } from "react-icons/lu";
import { ModalCenter } from "../atoms/ModalCenter";
import { DirectoryTreePlaylistEditable } from "./DirectoryTreePlaylistEditable";
import { Button } from "../atoms/Button";

type Props = {
  directory: DirectorySchema;
  level: number;
  onUpdate: (directory: DirectorySchema) => void;
  onDelete: (directory: DirectorySchema) => void;
}

const offSet = 32;

export const DirectoryTreeNodeEditable = ({ directory, level, onUpdate, onDelete }: Props) => {
  const [expanded, setExpanded] = useState(true)
  const [name, setName] = useState(directory.name)

  const [opened, { open, close }] = useDisclosure()

  const { isOver, setNodeRef } = useDroppable({
    id: directory.iid,
  })

  const handleCreate = () => {
    const newDirectory: DirectorySchema = {
      iid: getUuid(),
      name: "Subdirectory",
      children: [],
      playlists: [],
    }
    const updated = { ...directory, children: [...(directory.children ?? []), newDirectory] }

    onUpdate(updated)
    setExpanded(true)
  }

  const handleUpdate = () => {
    const updated = { ...directory, name }
    onUpdate(updated)
    close()
  }

  const handleDelete = () => {
    onDelete(directory)
  }

  const handleDeletePlaylist = (playlist: PlaylistSchema) => {
    const updated = { ...directory, playlists: directory.playlists.filter(p => p.id !== playlist.id) }
    onUpdate(updated)
  }

  return (
    <>
      <Stack ref={setNodeRef} gap={2}>
        <Group gap="xs" justify="space-between" style={{ marginLeft: level * offSet }}>
          <Group py={2} px="xs" gap="xs" onClick={() => setExpanded(prev => !prev)} className={`cursor-pointer hover:bg-[var(--mantine-color-background-1)] rounded-xl ${isOver && "bg-[var(--mantine-color-background-1)]"}`}>
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
          <Group gap={0}>
            <ActionIcon onClick={open} variant="subtle" className="text-muted">
              <LuPencilLine />
            </ActionIcon>
            <ActionIcon variant="subtle" onClick={handleCreate} className="text-muted">
              <LuFolderPlus />
            </ActionIcon>
            <ActionIcon variant="subtle" onClick={handleDelete} className="text-muted">
              <LuTrash2 />
            </ActionIcon>
          </Group>
        </Group>
        {expanded && (
          <>
            <Stack gap={2} style={{ marginLeft: (level + 1) * offSet + 5 }}>
              {directory.playlists.map(p => <DirectoryTreePlaylistEditable key={p.id} playlist={p} onDelete={handleDeletePlaylist} />)}
            </Stack>
            {directory.children?.map(d => (
              <DirectoryTreeNodeEditable
                key={d.iid}
                directory={d}
                level={level + 1}
                onUpdate={updatedChild => {
                  const updated = {
                    ...directory,
                    children: directory.children?.map(c => c.iid === updatedChild.iid ? updatedChild : c)
                  }
                  onUpdate(updated)
                }}
                onDelete={updatedChild => {
                  const updated = {
                    ...directory,
                    children: directory.children?.filter(c => c.iid !== updatedChild.iid)
                  }
                  onUpdate(updated)
                }}
              />
            ))}
          </>
        )}
      </Stack>
      <ModalCenter title="Rename directory" opened={opened} onClose={close} size="md">
        <Stack>
          <TextInput
            value={name}
            onChange={e => setName(e.target.value)}
          />
          <Group justify="end">
            <Button onClick={close} color="primary.6" variant="outline">
              Cancel
            </Button>
            <Button onClick={handleUpdate} color="primary.6">
              Save
            </Button>
          </Group>
        </Stack>
      </ModalCenter>
    </>
  )
}
