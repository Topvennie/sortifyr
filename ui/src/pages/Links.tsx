import { ModalCenter } from "@/components/atoms/ModalCenter"
import { LinkConnections } from "@/components/link/LinkConnections"
import { LinkTree } from "@/components/link/LinkTree"
import { Confirm } from "@/components/molecules/Confirm"
import { LoadingSpinner } from "@/components/molecules/LoadingSpinner"
import { useDirectoryGetAll } from "@/lib/api/directory"
import { useLinkAnchor } from "@/lib/hooks/useLinkAnchor"
import { LinkAnchorProvider } from "@/lib/providers/LinkAnchorProvider"
import { Button, Stack, Title, Text } from "@mantine/core"
import { useDisclosure, useMediaQuery } from "@mantine/hooks"

export const Links = () => {
  return (
    <LinkAnchorProvider>
      <LinksInner />
    </LinkAnchorProvider>
  )
}

const explanation = `A link is a connection between 2 items and represents a one way synchronization of songs.
If the left item is a directory then it will synchronize every song from nested playlist in that directory.
If the right item is a directory then it will synchronize to any nested playlist in that directory.

If a link is not visibile because the playlist is collapsed then a red number will appear with the amount of non visible links.
`

const LinksInner = () => {
  const { data: directories, isLoading } = useDirectoryGetAll()

  const { resetConnections, saveConnections } = useLinkAnchor()

  const [openedInfo, { open: openInfo, close: closeInfo }] = useDisclosure()
  const [openedReset, { open: openReset, close: closeReset }] = useDisclosure()
  const [openedSave, { open: openSave, close: closeSave }] = useDisclosure()

  const mdPoint = useMediaQuery('(width >= 64em)');

  if (isLoading) return <LoadingSpinner />
  if (!mdPoint) return (
    <div className="flex flex-col justify-center gap-4">
      <Title order={1} className="text-center">Links</Title>
      <span className="text-center">This page only supports big screen, apologies!</span>
    </div>
  )

  const handleInfo = () => {
    openInfo()
  }

  const handleResetInit = () => {
    openReset()
  }

  const handleReset = () => {
    resetConnections()
    closeReset()
  }

  const handleSaveInit = () => {
    openSave()
  }

  const handleSave = async () => {
    await saveConnections()
    closeSave()
  }

  return (
    <>
      <div className="grid grid-cols-3 gap-8">
        <Title order={1} className="col-span-full text-center">Links</Title>

        <div className="col-span-full">
          <div className="flex items-center justify-end gap-2">
            <Button onClick={handleInfo} variant="outline" className="mr-8">Info</Button>
            <Button onClick={handleResetInit} color="red">Reset</Button>
            <Button onClick={handleSaveInit}>Save</Button>
          </div>
        </div>

        <div className="col-span-1 col-start-1 space-y-1">
          {directories?.map(d => <LinkTree key={d.id} side="left" directory={d} />)}
        </div>

        <div className="col-span-1 col-start-3 space-y-1">
          {directories?.map(d => <LinkTree key={d.id} side="right" directory={d} />)}
        </div>
      </div>

      <LinkConnections />
      <ModalCenter opened={openedInfo} onClose={closeInfo} title="Info">
        <Stack>
          <Text fw="bold">Explanation</Text>
          <div className="whitespace-pre-wrap">{explanation}</div>
        </Stack>
      </ModalCenter>
      <Confirm
        opened={openedReset}
        onClose={closeReset}
        modalTitle="Reset"
        title="Reset links"
        description="Are you sure you want to discard all changes?"
        onConfirm={handleReset}
      />
      <Confirm
        opened={openedSave}
        onClose={closeSave}
        modalTitle="Save"
        title="Save links"
        description="Are you sure you want to save?"
        onConfirm={handleSave}
      />
    </>
  )
}

