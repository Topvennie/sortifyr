import { usePlaylistGetUnplayables } from "@/lib/api/playlist"
import { Playlist } from "@/lib/types/playlist"
import { sortBy } from "lodash"
import { DataTableSortStatus } from "mantine-datatable"
import { useState, useMemo } from "react"
import { SectionTitle } from "../atoms/Page"
import { Table } from "../molecules/Table"
import { PlaylistCover } from "./PlaylistCover"
import { LuTriangle } from "react-icons/lu"

export const PlaylistUnplayables = () => {
  const { data: playlists, isLoading } = usePlaylistGetUnplayables()

  const [expandedRecordIds, setExpandedRecordIds] = useState<number[]>([]);

  const [sortStatus, setSortStatus] = useState<DataTableSortStatus<Playlist>>({
    columnAccessor: "name",
    direction: "asc",
  })
  const records = useMemo(() => {
    const data = sortBy(playlists, sortStatus.columnAccessor);
    return sortStatus.direction === "desc" ? data.reverse() : data;
  }, [playlists, sortStatus])

  return (
    <>
      <SectionTitle
        title="Unplayable tracks"
        description={`Playlists with unplayable tracks.\nClick on a row to see the unplayable tracks.`}
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
            accessor: "id",
            title: "",
            width: 52,
            render: playlist => <PlaylistCover playlist={playlist} />,
          },
          { accessor: "name", sortable: true },
          { accessor: "trackAmount", sortable: true },
          { accessor: "owner.name", sortable: true },
          {
            accessor: "unplayables",
            title: "Total unplayable",
            sortable: true,
            render: ({ unplayables }) => <p>{unplayables.length}</p>
          }
        ]}
        rowExpansion={{
          content: ({ record: { unplayables } }) => (
            <Table
              noHeader
              backgroundColor="background.1"
              columns={[
                { accessor: "name" },
              ]}
              records={unplayables}
              height={180}
              className="m-4"
            />
          ),
          expanded: {
            recordIds: expandedRecordIds,
            onRecordIdsChange: setExpandedRecordIds,
          }
        }}
        records={records}
        noRecordsText="No playlists with unplayable tracks"
        sortStatus={sortStatus}
        onSortStatusChange={setSortStatus}
        fetching={isLoading}
      />
    </>
  )
}
