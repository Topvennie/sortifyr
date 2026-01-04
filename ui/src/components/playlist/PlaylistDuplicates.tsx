import { usePlaylistGetDuplicates, usePlaylistRemoveDuplicates } from "@/lib/api/playlist"
import { Playlist } from "@/lib/types/playlist"
import { Group } from "@mantine/core"
import { useDisclosure } from "@mantine/hooks"
import { notifications } from "@mantine/notifications"
import { sortBy } from "lodash"
import { type DataTableSortStatus } from "mantine-datatable"
import { useMemo, useState } from "react"
import { SectionTitle } from "../atoms/Page"
import { Confirm } from "../molecules/Confirm"
import { Table } from "../molecules/Table"
import { PlaylistCover } from "./PlaylistCover"
import { getErrorMessage } from "@/lib/utils"
import { Button } from "../atoms/Button"
import { useTaskGetAll } from "@/lib/api/task"
import { TaskStatus } from "@/lib/types/task"

export const PlaylistDuplicates = () => {
  const { data: playlists, isLoading } = usePlaylistGetDuplicates()
  const removeDuplicates = usePlaylistRemoveDuplicates()

  const { data: tasks } = useTaskGetAll()
  const task = tasks?.find(t => t.uid === "task-playlist-duplicate")

  const [sortStatus, setSortStatus] = useState<DataTableSortStatus<Playlist>>({
    columnAccessor: "name",
    direction: "asc",
  })
  const records = useMemo(() => {
    const data = sortBy(playlists, sortStatus.columnAccessor);
    return sortStatus.direction === "desc" ? data.reverse() : data;
  }, [playlists, sortStatus])

  const [opened, { open, close }] = useDisclosure()

  const handleRemove = () => {
    removeDuplicates.mutateAsync(undefined, {
      onSuccess: () => notifications.show({ message: "Duplicates are getting removed. Go to the task page to see the result.", }),
      onError: async (error) => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      },
      onSettled: () => close(),
    })
  }

  return (
    <>
      <Group justify="space-between">
        <SectionTitle
          title="Playlist duplicates"
          description={`Playlists with duplicate tracks.\nClick on a row to see the duplicate tracks.`}
        />
        <Button onClick={open} color="secondary.1" disabled={task?.status === TaskStatus.Running}>
          Remove duplicates
        </Button>
      </Group>
      <Table
        columns={[
          {
            accessor: "id",
            title: "",
            width: 52,
            render: playlist => <PlaylistCover playlist={playlist} />,
          },
          { accessor: "name", sortable: true },
          { accessor: "trackAmount", sortable: true },
          { accessor: "owner.name", sortable: true },
          {
            accessor: "duplicates",
            title: "Unique duplicates",
            sortable: true,
            render: ({ duplicates }) => <p>{duplicates.length}</p>
          },
          {
            accessor: "duplicatesTotal",
            title: "Total duplicates",
            sortable: true,
            render: ({ duplicates }) => <p>{duplicates.reduce((acc, curr) => acc + curr.amount, 0)}</p>
          },
        ]}
        rowExpansion={{
          content: ({ record: { duplicates } }) => (
            <Table
              noHeader
              backgroundColor="background.1"
              columns={[
                { accessor: "name" },
                { accessor: "amount" }
              ]}
              records={duplicates}
              height={180}
              className="m-4"
            />
          )
        }}
        records={records}
        noRecordsText="No playlists with duplicate tracks"
        sortStatus={sortStatus}
        onSortStatusChange={setSortStatus}
        fetching={isLoading}
      />

      <Confirm
        opened={opened}
        onClose={close}
        modalTitle="Duplicates"
        title="Remove Duplicates"
        description={`Are you sure you want to remove all duplicates?\nAfter it's finished it'll take a couple of minutes before the changes come through.`}
        onConfirm={handleRemove}
      />
    </>
  )
}
