import { LinkButton } from "@/components/atoms/LinkButton"
import { Page, PageTitle, Section, SectionTitle } from "@/components/atoms/Page"
import { Confirm } from "@/components/molecules/Confirm"
import { Table } from "@/components/molecules/Table"
import { useGeneratorDelete, useGeneratorGetAll, useGeneratorRefresh } from "@/lib/api/generator"
import { Generator } from "@/lib/types/generator"
import { getErrorMessage } from "@/lib/utils"
import { ActionIcon, Badge, Checkbox, Group, Stack } from "@mantine/core"
import { useDisclosure } from "@mantine/hooks"
import { notifications } from "@mantine/notifications"
import { useNavigate } from "@tanstack/react-router"
import { useState } from "react"
import { LuCheck, LuListRestart, LuPencil, LuSparkles, LuTrash2, LuTriangle } from "react-icons/lu"

export const GeneratorOverview = () => {
  const { data: generators, isLoading } = useGeneratorGetAll()
  const generatorRefresh = useGeneratorRefresh()
  const generatorDelete = useGeneratorDelete()

  const [generatorToDelete, setGeneratorToDelete] = useState<Generator | null>(null)
  const [checkedPlaylist, setCheckedPlaylist] = useState(false)
  const [opened, { open, close }] = useDisclosure()

  const navigate = useNavigate()

  const [expandedRecordIds, setExpandedRecordIds] = useState<number[]>([]);

  const handleRefresh = (gen: Generator) => {
    generatorRefresh.mutateAsync(gen, {
      onSuccess: () => notifications.show({ message: "Updating generator" }),
      onError: async error => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      }
    })
  }

  const handleEdit = (gen: Generator) => {
    navigate({ to: "/generator/edit/$generatorId", params: { generatorId: gen.id.toString() } })
  }

  const handleDeleteInit = (gen: Generator) => {
    setGeneratorToDelete(gen)
    setCheckedPlaylist(false)
    open()
  }

  const handleDelete = () => {
    if (!generatorToDelete) return

    generatorDelete.mutateAsync({ generator: generatorToDelete, deletePlaylist: checkedPlaylist }, {
      onSuccess: () => notifications.show({ message: "Generated deleted" }),
      onError: async error => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      },
      onSettled: () => close(),
    })
  }

  return (
    <>
      <Page>
        <Group justify="space-between">
          <PageTitle
            title="Generate new playlists"
            description="Create playlists from presets and fine-tune them before saving."
          />
          <LinkButton to={"/generator/create"} leftSection={<LuSparkles />} radius="lg">New Generator</LinkButton>
        </Group>

        <Section>
          <SectionTitle
            title="Generated playlists"
          />

          <Table
            columns={[
              {
                accessor: "expanded",
                title: "",
                width: 24,
                render: ({ id }) => <LuTriangle className={`fill-black w-2 transform duration-300 ${expandedRecordIds.includes(id) ? "" : "rotate-180"}`} />
              },
              {
                accessor: "name",
                title: "Name & Description",
                render: ({ name, description }) => (
                  <div>
                    <p className="font-semibold">{name}</p>
                    <p className="text-muted text-sm">{description}</p>
                  </div>
                )
              },
              {
                accessor: "playlist",
                title: "Playlist",
                textAlign: "right",
                render: ({ playlistId }) => playlistId && <LuCheck className="ml-auto text-green-500 size-6" />,
              },
              {
                accessor: "interval_days",
                title: "Refresh",
                textAlign: "right",
                render: ({ intervalDays }) => intervalDays > 0 && <p>{`Every${intervalDays !== 1 ? ' ' + intervalDays : ''} day${intervalDays !== 1 ? 's' : ''}`}</p>
              },
              {
                accessor: "spotifyOutdated",
                title: "Spotify Status",
                textAlign: "right",
                render: ({ playlistId, spotifyOutdated }) => playlistId && <Badge color={spotifyOutdated ? "red" : "gray"}>{spotifyOutdated ? "Outdated" : "Up to date"}</Badge>,
              },
              {
                accessor: "actions",
                title: "",
                width: 120,
                render: gen => (
                  <div className="flex gap-0 flex-nowrap">
                    <div className="flex-1" />
                    <ActionIcon onClick={() => handleRefresh(gen)} variant="subtle" color="black" hidden={!gen.playlistId}><LuListRestart className="size-5" /></ActionIcon>
                    <ActionIcon onClick={() => handleEdit(gen)} variant="subtle" color="black"><LuPencil /></ActionIcon>
                    <ActionIcon onClick={() => handleDeleteInit(gen)} variant="subtle" color="red"><LuTrash2 /></ActionIcon>
                  </div>
                )
              },
            ]}
            rowExpansion={{
              content: ({ record: { tracks } }) => (
                <Table
                  noHeader
                  backgroundColor="background.1"
                  columns={[
                    { accessor: "name" },
                  ]}
                  records={tracks}
                  height={180}
                  className="m-4"
                />
              ),
              expanded: {
                recordIds: expandedRecordIds,
                onRecordIdsChange: setExpandedRecordIds,
              }
            }}
            records={generators ?? []}
            fetching={isLoading}
            noRecordsText="No generators"
          />
        </Section>
      </Page>
      <Confirm
        opened={opened}
        onClose={close}
        modalTitle="Delete"
        title="Delete Generator"
        description="Are you sure you want to delete the generator"
        content={generatorToDelete?.playlistId && (
          <Stack>
            <Checkbox checked={checkedPlaylist} onChange={e => setCheckedPlaylist(e.target.checked)} label="Delete Spotify playlist" />
          </Stack>
        )}
        onConfirm={handleDelete}
      />

    </>
  )
}
