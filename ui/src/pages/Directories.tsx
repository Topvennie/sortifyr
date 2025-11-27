
import { DirectoryToolbar } from "@/components/directory/DirectoryToolbar"
import { DirectoryTree } from "@/components/directory/DirectoryTree"
import { LoadingSpinner } from "@/components/molecules/LoadingSpinner"
import { useDirectoryGetAll } from "@/lib/api/directory"
import { usePlaylistGetAll } from "@/lib/api/playlist"
import { DndContext, DragEndEvent } from '@dnd-kit/core'
import {
  restrictToWindowEdges,
} from '@dnd-kit/modifiers'

import {
  convertDirectorySchema,
  DirectorySchema
} from "@/lib/types/directory"

import { DirectoryPlaylistSelector } from "@/components/directory/DirectoryPlaylistSelector"
import { PlaylistSchema } from "@/lib/types/playlist"
import { useEffect, useMemo, useState } from "react"
import { Title } from "@mantine/core"

const getPlaylists = (directory: DirectorySchema): PlaylistSchema[] => {
  return [
    ...directory.playlists,
    ...(directory.children?.flatMap(getPlaylists) ?? [])
  ]
}

const getDirectories = (directories: DirectorySchema[]): DirectorySchema[] => {
  return [...directories, ...(directories.flatMap(d => getDirectories(d?.children ?? [])))]
}

export const Directories = () => {
  const { data: playlistsAll, isLoading: isLoadingPlaylists } = usePlaylistGetAll()
  const { data: directories, isLoading: isLoadingDirectories } = useDirectoryGetAll()

  const [roots, setRoots] = useState<DirectorySchema[]>([])

  useEffect(() => {
    const converted = convertDirectorySchema(directories ?? [])
    setRoots(converted)
  }, [directories])

  const handleDirectoryUpdate = (updatedRoot: DirectorySchema) => {
    const oldRoots = roots.filter(r => r.iid !== updatedRoot.iid)
    const newRoots = [...oldRoots, updatedRoot].sort((a, b) => a.name > b.name ? 1 : -1)

    setRoots(newRoots)
  }

  const handleDirectoryDelete = (root: DirectorySchema) => {
    const newRoots = roots.filter(r => r.iid !== root.iid)
    setRoots(newRoots)
  }

  const playlistsAvailable = useMemo(() => {
    const playlistsInDirectories = new Set(
      roots?.flatMap(getPlaylists).map(p => p.id)
    )

    return playlistsAll?.filter(p => !playlistsInDirectories.has(p.id))
  }, [playlistsAll, roots])

  if (isLoadingPlaylists || isLoadingDirectories) {
    return <LoadingSpinner />
  }

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event

    if (!over?.id) return

    const newRoots = [...roots]

    const directories = getDirectories(newRoots)
    const directory = directories.find(d => d.iid === over.id)
    if (!directory) return

    const playlist = playlistsAll?.find(p => p.id === active.id)
    if (!playlist) return

    directory.playlists.push(playlist)
    setRoots(newRoots)
  }

  return (
    <DndContext onDragEnd={handleDragEnd} modifiers={[restrictToWindowEdges]}>
      <div className="grid grid-cols-4 gap-8">
        <Title order={1} className="col-span-full">Directories</Title>

        <div className="order-1 lg:order-2 col-span-full lg:col-span-1">
          <DirectoryPlaylistSelector playlists={playlistsAvailable ?? []} />
        </div>

        <div className="order-2 lg:order-1 col-span-full lg:col-span-3 space-y-4">
          <DirectoryToolbar roots={roots} setRoots={setRoots} />
          {roots.map(r => <DirectoryTree key={r.iid} root={r} onUpdate={handleDirectoryUpdate} onDelete={handleDirectoryDelete} />)}
        </div>

      </div>
    </DndContext>
  )
}

