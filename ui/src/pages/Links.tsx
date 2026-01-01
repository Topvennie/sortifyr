import { Button } from "@/components/atoms/Button"
import { Page, PageTitle, Section, SectionTitle } from "@/components/atoms/Page"
import { LinkConnections } from "@/components/link/LinkConnections"
import { LinkTree } from "@/components/link/LinkTree"
import { Confirm } from "@/components/molecules/Confirm"
import { useDirectoryGetAll } from "@/lib/api/directory"
import { usePlaylistGetAll } from "@/lib/api/playlist"
import { useLinkAnchor } from "@/lib/hooks/useLinkAnchor"
import { LinkAnchorProvider } from "@/lib/providers/LinkAnchorProvider"
import { Directory } from "@/lib/types/directory"
import { Playlist } from "@/lib/types/playlist"
import { Badge, Group, Stack } from "@mantine/core"
import { useDisclosure, useMediaQuery } from "@mantine/hooks"
import { useMemo } from "react"

export const Links = () => {
  return (
    <LinkAnchorProvider>
      <LinksInner />
    </LinkAnchorProvider>
  )
}

const explanation = `Drag connections between directories and playlists on the left and the right. Click on a connection to remove it.
Some connections are not shown because the source / destination is not visible. They are indicated with a red number.
Tracks from sources will periodically be added to the destination
`

const getPlaylists = (directory: Directory): Playlist[] => {
  return [
    ...directory.playlists,
    ...(directory.children?.flatMap(getPlaylists) ?? [])
  ]
}


const LinksInner = () => {
  const { data: playlistsAll, isLoading: isLoadingPlaylists } = usePlaylistGetAll()
  const { data: directories, isLoading: isLoadingDirectories } = useDirectoryGetAll()

  const unAssigned = useMemo(() => {
    const playlists = (directories?.flatMap(getPlaylists) ?? []).map(p => p.id)
    return playlistsAll?.filter(p => !playlists.includes(p.id)) ?? []
  }, [playlistsAll, directories])

  const { visibleAnchorsRef, connections, layoutVersion, resetConnections, saveConnections } = useLinkAnchor()
  const hidden = useMemo(() => {
    return connections.filter(({ from, to }) => !(visibleAnchorsRef.current[from] && visibleAnchorsRef.current[to])).length
  }, [connections, layoutVersion]) // eslint-disable-line react-hooks/exhaustive-deps

  const [openedReset, { open: openReset, close: closeReset }] = useDisclosure()
  const [openedSave, { open: openSave, close: closeSave }] = useDisclosure()

  const mdPoint = useMediaQuery('(width >= 64em)');

  if (!mdPoint) {
    return (
      <Page>
        <PageTitle
          title="Links"
          description="Connect directories and playlists"
        />
        <Section>
          <SectionTitle
            title="Visual linking"
            description="This screen is only available on big screens."
          />
        </Section>
      </Page>
    )
  }

  return (
    <>
      <Page className="select-none">
        <PageTitle
          title="Links"
          description="Connect directories and playlists"
        />
        <Section>
          <Group justify="space-between">
            <SectionTitle
              title="Visual linking"
              description={explanation}
            />
            <Stack gap="xs" align="end" className="flex-1">
              <Badge color="secondary.2">{`${connections.length} Connection${connections.length !== 1 ? "s" : ""}`}</Badge>
              <Badge color="red">{`${hidden} Hidden`}</Badge>
            </Stack>
          </Group>
          <div className="flex-1 flex gap-4 overflow-hidden">
            <LinkTree
              roots={directories ?? []}
              unAssigned={unAssigned}
              side="left"
              title="Source"
              isLoading={isLoadingPlaylists || isLoadingDirectories}
              className="flex-1"
            />
            <div className="h-full w-[20%]" />
            <LinkTree
              roots={directories ?? []}
              unAssigned={unAssigned}
              side="right"
              title="Target"
              isLoading={isLoadingPlaylists || isLoadingDirectories}
              className="flex-1"
            />
          </div>

          <Group justify="end">
            <Button onClick={openReset} c="" variant="default" className="text-muted">
              Reset changes
            </Button>
            <Button onClick={openSave}>
              Apply changes
            </Button>
          </Group>
        </Section>
      </Page>

      <LinkConnections />

      <Confirm
        opened={openedReset}
        onClose={closeReset}
        modalTitle="Reset"
        title="Reset links"
        description="Are you sure you want to discard all changes?"
        onConfirm={() => { resetConnections(); closeReset() }}
      />
      <Confirm
        opened={openedSave}
        onClose={closeSave}
        modalTitle="Save"
        title="Save links"
        description="Are you sure you want to save?"
        onConfirm={() => { saveConnections(); closeSave() }}
      />
    </>
  )
}
