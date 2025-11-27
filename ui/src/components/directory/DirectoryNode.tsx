import { DirectorySchema } from "@/lib/types/directory";
import { PlaylistSchema } from "@/lib/types/playlist";
import { getUuid } from "@/lib/utils";
import { useDroppable } from "@dnd-kit/core";
import { ActionIcon, TextInput } from "@mantine/core";
import { MouseEvent, useState } from "react";
import { FaCheck, FaPencil, FaPlus, FaRegFolder, FaRegFolderOpen, FaTrashCan, FaX } from "react-icons/fa6";
import { DirectoryPlaylist } from "./DirectoryPlaylist";

type Props = {
  directory: DirectorySchema;
  onUpdate: (directory: DirectorySchema) => void;
  onDelete: (directory: DirectorySchema) => void;
  level: number;
}

export const DirectoryNode = ({ directory, onUpdate, onDelete, level }: Props) => {
  const [expanded, setExpanded] = useState(false)
  const [editing, setEditing] = useState(false)
  const [name, setName] = useState(directory.name)

  console.log(directory)

  const { isOver, setNodeRef } = useDroppable({
    id: directory.iid,
  })

  const handleCreate = (e: MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()

    const newDir: DirectorySchema = {
      iid: getUuid(),
      name: "New subdirectory",
      children: [],
      playlists: [],
    }
    const updated = { ...directory, children: [...(directory.children ?? []), newDir] }
    onUpdate(updated)
    setExpanded(true)
  }

  const handleDelete = (directoryDelete: DirectorySchema) => {
    const updated = { ...directory, children: directory.children?.filter(c => c.iid !== directoryDelete.iid) }
    onUpdate(updated)
  }

  const handleDeletePlaylist = (playlist: PlaylistSchema) => {
    const updated = { ...directory, playlists: directory.playlists.filter(p => p.id !== playlist.id) }
    onUpdate(updated)
  }

  const handleExpand = () => {
    if (editing) return

    setExpanded(prev => !prev)
  }

  const handleChangeName = (e: MouseEvent<HTMLButtonElement>, save: boolean) => {
    e.stopPropagation()

    if (!editing) {
      setName(directory.name)
      setEditing(true)
      return
    }

    if (save) {
      const updated = { ...directory, name }
      onUpdate(updated)
    }

    setEditing(false)
  }

  return (
    <div ref={setNodeRef} className="flex flex-col gap-1">
      <div
        onClick={handleExpand}
        style={{ marginLeft: level * 16 }}
        className={`flex items-center justify-between rounded-md bg-gray-200 p-4 cursor-pointer ${isOver ? "brightness-75" : ""}`}
      >
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <ActionIcon color="black" variant="subtle" disabled={editing} >
              {expanded
                ? <FaRegFolderOpen />
                : <FaRegFolder />
              }
            </ActionIcon>
            {editing ? (
              <TextInput
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            ) : (
              <span className="font-semibold">{directory.name}</span>
            )}
          </div>
          {editing ? (
            <>
              <ActionIcon onClick={e => handleChangeName(e, true)} size="xs" variant="subtle">
                <FaCheck color="green" />
              </ActionIcon>
              <ActionIcon onClick={e => handleChangeName(e, false)} size="xs" variant="subtle">
                <FaX color="red" />
              </ActionIcon>
            </>
          ) : (
            <ActionIcon onClick={e => handleChangeName(e, false)} size="xs" variant="subtle">
              <FaPencil color="black" />
            </ActionIcon>
          )}
        </div>
        <div className="flex items-center gap-2">
          <ActionIcon onClick={handleCreate} variant="subtle">
            <FaPlus />
          </ActionIcon>
          <ActionIcon onClick={e => {
            e.stopPropagation()
            onDelete(directory)
          }} color="red" variant="subtle">
            <FaTrashCan />
          </ActionIcon>
        </div>
      </div>

      {expanded && (
        <>
          <div style={{ marginLeft: (level + 1) * 16 }} className="flex flex-col gap-1">
            {directory.playlists?.map((p: PlaylistSchema) => (
              <DirectoryPlaylist key={p.id} playlist={p} onDelete={handleDeletePlaylist} className={isOver ? "brightness-75" : ""} />
            ))}
          </div>
          <div className="flex flex-col gap-1">
            {directory.children?.map((child: DirectorySchema) => (
              <DirectoryNode
                key={child.iid}
                directory={child}
                level={level + 1}
                onUpdate={updatedChild => {
                  const updated = {
                    ...directory,
                    children: directory.children?.map((c: DirectorySchema) =>
                      c.iid === updatedChild.iid ? updatedChild : c
                    )
                  }
                  onUpdate(updated)
                }}
                onDelete={handleDelete}
              />
            ))}
          </div>
        </>
      )}

    </div>
  )
}
