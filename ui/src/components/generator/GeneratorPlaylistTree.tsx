import { useDirectoryGetAll } from "@/lib/api/directory";
import { Playlist } from "@/lib/types/playlist"
import { ActionIcon, Checkbox, Group, Stack } from "@mantine/core";
import { LoadingSpinner } from "../molecules/LoadingSpinner";
import { Directory } from "@/lib/types/directory";
import { useMemo, useState } from "react";
import { LuFolder, LuFolderOpen, LuTriangle } from "react-icons/lu";
import { PlaylistCover } from "../playlist/PlaylistCover";

type Props = {
  playlists: Playlist[];
  excluded: number[];
  isLoading: boolean;
  onToggle: (playlistIds: number[], isExcluded: boolean) => void;
}

const getPlaylists = (directory: Directory): Playlist[] => {
  return [
    ...directory.playlists,
    ...(directory.children?.flatMap(getPlaylists) ?? [])
  ]
}

export const GeneratorPlaylistTree = ({ playlists, excluded, isLoading: isLoadingPlaylists, onToggle }: Props) => {
  const { data: directories, isLoading: isLoadingDirectories } = useDirectoryGetAll()

  const unAssigned = useMemo(() => {
    const directoryPlaylists = (directories?.flatMap(getPlaylists) ?? []).map(p => p.id)
    return playlists?.filter(p => !directoryPlaylists.includes(p.id)) ?? []
  }, [playlists, directories])

  const child = () => {
    if (isLoadingDirectories || isLoadingPlaylists) return <LoadingSpinner />

    return (
      <>
        {directories?.length !== 0 && (
          <>
            <p className="text-muted">Directories</p>
            {directories?.map(d => <TreeNode key={d.id} directory={d} level={0} excluded={excluded} onToggle={onToggle} />)}
            <p className="text-muted">Unassigned Playlists</p>
          </>
        )}
        <Stack gap={2} style={{ marginLeft: 6 }}>
          {unAssigned.map(p => <TreePlaylist key={p.id} playlist={p} excluded={excluded.includes(p.id)} onToggle={onToggle} />)}
        </Stack>
      </>
    )
  }

  return (
    <Stack p="md" gap="md" justify="start" bg="background.0" className="overflow-auto rounded-xl">
      {child()}
    </Stack>
  )
}

type TreeNodeProps = {
  directory: Directory;
  level: number;
  excluded: number[];
  onToggle: (playlistIds: number[], isExcluded: boolean) => void;
}

const offset = 32;

const TreeNode = ({ directory, level, excluded, onToggle }: TreeNodeProps) => {
  const [expanded, setExpanded] = useState(true)

  const playlistsIds = useMemo(() => getPlaylists(directory).map(p => p.id), [directory])

  const selected: boolean | null = useMemo(() => {
    const filtered = excluded.filter(e => playlistsIds.includes(e))
    if (filtered.length === 0) return true
    if (filtered.length === playlistsIds.length) return false
    return null
  }, [excluded, playlistsIds])

  return (
    <Stack gap={2}>
      <Group py={2} px="xs" gap="xs" style={{ marginLeft: level * offset }}>
        <ActionIcon onClick={() => setExpanded(prev => !prev)} variant="transparent" color="black" size="sm">
          <LuTriangle className={`fill-black w-2 ${expanded && "rotate-180"}`} />
        </ActionIcon>
        <Checkbox checked={!!selected} indeterminate={selected === null} onChange={() => onToggle(playlistsIds, selected === true || selected === null)} color="secondary.2" />
        {expanded
          ? <LuFolderOpen className="size-6" />
          : <LuFolder className="size-6" />
        }
        <p className="font-semibold">{directory.name}</p>
        <p className="text-muted text-sm">{directory.playlists.length + (directory.children?.length ?? 0)}</p>
      </Group>
      {expanded && (
        <>
          <Stack gap={2} style={{ marginLeft: (level + 1) * offset + 2 }}>
            {directory.playlists.map(p => <TreePlaylist key={p.id} playlist={p} excluded={excluded.includes(p.id)} onToggle={onToggle} />)}
          </Stack>
          {directory.children?.map(d => (
            <TreeNode key={d.id} directory={d} level={level + 1} excluded={excluded} onToggle={onToggle} />
          ))}
        </>
      )}
    </Stack>
  )
}

type TreePlaylistProps = {
  playlist: Playlist;
  excluded: boolean;
  onToggle: (playlistIds: number[], isExcluded: boolean) => void;
}

const TreePlaylist = ({ playlist, excluded, onToggle }: TreePlaylistProps) => {
  return (
    <div onClick={() => onToggle([playlist.id], !excluded)} className={`flex items-center gap-2 px-2 py-0.5 rounded-xl cursor-pointer ${!excluded && "bg-(--mantine-color-background-2)"}`}>
      <Checkbox color="secondary.2" checked={!excluded} readOnly />
      <PlaylistCover playlist={playlist} />
      <p className="line-clamp-1 overflow-hidden text-ellipsis wrap-break-word">{playlist.name}</p>
      <p className="text-muted text-sm">{playlist.trackAmount}</p>
    </div>
  )
}
