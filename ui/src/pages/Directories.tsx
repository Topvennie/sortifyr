import { Button } from "@/components/atoms/Button"
import { Page, PageTitle, Section, SectionTitle } from "@/components/atoms/Page"
import { DirectoryPlaylistDraggable, DirectoryPlaylistSelector } from "@/components/directory/DirectoryPlaylistSelector"
import { DirectoryTree } from "@/components/directory/DirectoryTree"
import { Confirm } from "@/components/molecules/Confirm"
import { Search } from "@/components/molecules/Search"
import { useDirectoryGetAll, useDirectorySync } from "@/lib/api/directory"
import { usePlaylistGetAll } from "@/lib/api/playlist"
import { convertDirectorySchema, DirectorySchema } from "@/lib/types/directory"
import { Playlist, PlaylistSchema } from "@/lib/types/playlist"
import { getUuid } from "@/lib/utils"
import { DndContext, DragEndEvent, DragOverlay, DragStartEvent } from '@dnd-kit/core'
import { restrictToWindowEdges } from '@dnd-kit/modifiers'
import { Group } from "@mantine/core"
import { useDebouncedValue, useDisclosure } from "@mantine/hooks"
import { notifications } from "@mantine/notifications"
import { useEffect, useMemo, useState } from "react"
import { LuHand, LuPlus } from "react-icons/lu"

export const Directories = () => {
  return (
    <DndContext>
      <DirectoriesInner />
    </DndContext>
  )
}

const getPlaylists = (directory: DirectorySchema): PlaylistSchema[] => {
  return [
    ...directory.playlists,
    ...(directory.children?.flatMap(getPlaylists) ?? [])
  ]
}

const getDirectories = (directories: DirectorySchema[]): DirectorySchema[] => {
  return [...directories, ...(directories.flatMap(d => getDirectories(d?.children ?? [])))]
}

const DirectoriesInner = () => {
  const [filter, setFilter] = useState("")
  const [debounced] = useDebouncedValue(filter, 200);
  const [activePlaylist, setActivePlaylist] = useState<Playlist | null>(null);

  const [openedReset, { open: openReset, close: closeReset }] = useDisclosure()
  const [openedSave, { open: openSave, close: closeSave }] = useDisclosure()

  const save = useDirectorySync()

  const { data: playlistsAll, isLoading: isLoadingPlaylists } = usePlaylistGetAll()
  const { data: directoriesAll, isLoading: isLoadingDirectories } = useDirectoryGetAll()

  const [roots, setRoots] = useState<DirectorySchema[]>([])
  useEffect(() => {
    setRoots(convertDirectorySchema(directoriesAll ?? []))
  }, [directoriesAll])

  const playlists = useMemo(() => {
    const playlistsInDirectories = new Set(
      roots.flatMap(getPlaylists).map(p => p.id)
    )

    return playlistsAll?.filter(p =>
      p.name.toLowerCase().includes(debounced.toLowerCase()) && !playlistsInDirectories.has(p.id)
    ) ?? []
  }, [playlistsAll, debounced, roots])

  const handleDirectoryCreate = () => {
    const newRoot: DirectorySchema = {
      iid: getUuid(),
      name: "Directory",
      playlists: [],
      children: []
    }

    const updated = [...roots, newRoot]
    setRoots(updated)
  }

  const handleDirectoryUpdate = (directory: DirectorySchema) => {
    const oldRoots = roots.filter(r => r.iid !== directory.iid)
    const newRoots = [...oldRoots, directory].sort((a, b) => a.name > b.name ? 1 : -1)

    setRoots(newRoots)
  }

  const handleDirectoryDelete = (directory: DirectorySchema) => {
    const newRoots = roots.filter(r => r.iid !== directory.iid)
    setRoots(newRoots)
  }

  const handleDragStart = (event: DragStartEvent) => {
    const playlist = playlistsAll?.find(p => p.id === event.active.id) ?? null
    setActivePlaylist(playlist)
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

  const handleReset = () => {
    setRoots(convertDirectorySchema(directoriesAll ?? []))
    closeReset()
  }

  const handleSave = () => {
    save.mutate(roots, {
      onSuccess: () => notifications.show({ message: "Directories synced" }),
      onSettled: () => closeSave()
    })
  }

  return (
    <>
      <DndContext onDragStart={handleDragStart} onDragEnd={handleDragEnd} modifiers={[restrictToWindowEdges]}>
        <Page>
          <PageTitle
            title="Directory configurator"
            description="Keep playlists tidy by stacking them into simple folders."
          />

          <div className="flex-1 flex flex-col md:flex-row gap-4 md:overflow-hidden">
            <Section className="flex-none md:w-[70%]">
              <Group justify="space-between">
                <SectionTitle
                  title="Directories"
                  description="Drag playlists here to group them."
                />
                <Button onClick={handleDirectoryCreate} color="secondary.1" leftSection={<LuPlus />}>
                  New root
                </Button>
              </Group>
              <DirectoryTree
                roots={roots}
                isLoading={isLoadingDirectories}
                title="Root directories"
                editable
                onUpdate={handleDirectoryUpdate}
                onDelete={handleDirectoryDelete}
                className="flex-1"
              />
              <Group justify="end">
                <Button onClick={openReset} c="" variant="default" className="text-muted">
                  Reset
                </Button>
                <Button onClick={openSave}>
                  Save
                </Button>
              </Group>
            </Section>

            <Section>
              <SectionTitle
                title="Unassigned playlists"
                description="Anything you add to a directory disappears from here"
              />
              <Search
                placeholder="Filter by playlist name..."
                value={filter}
                onChange={e => setFilter(e.target.value)}
              />
              <DirectoryPlaylistSelector
                playlists={playlists}
                isLoading={isLoadingPlaylists}
              />
              <Group wrap="nowrap">
                <LuHand className="text-muted size-6" />
                <p className="text-muted">Drag any playlist row onto a directory to add it</p>
              </Group>
            </Section>
          </div>
        </Page>
        <DragOverlay>
          {activePlaylist && <DirectoryPlaylistDraggable playlist={activePlaylist} />}
        </DragOverlay>
      </DndContext>
      <Confirm
        opened={openedReset}
        onClose={closeReset}
        modalTitle="Reset"
        title="Reset directory structure"
        description="Are you sure you want to discard all changes?"
        onConfirm={handleReset}
      />
      <Confirm
        opened={openedSave}
        onClose={closeSave}
        modalTitle="Save"
        title="Save directory structure"
        description="Are you sure you want to save?"
        onConfirm={handleSave}
      />
    </>
  )
}

