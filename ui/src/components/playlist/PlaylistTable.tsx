import { Playlist } from "@/lib/types/playlist";
import sortBy from 'lodash/sortBy';
import { type DataTableSortStatus } from 'mantine-datatable';
import { useMemo, useState } from "react";
import { FaCheck, FaX } from "react-icons/fa6";
import { Table } from "../molecules/Table";
import { PlaylistCover } from "./PlaylistCover";

type Props = {
  playlists: Playlist[];
  isLoading: boolean;
}

export const PlaylistTable = ({ playlists, isLoading }: Props) => {
  const [sortStatus, setSortStatus] = useState<DataTableSortStatus<Playlist>>({
    columnAccessor: "name",
    direction: "asc",
  })
  const records = useMemo(() => {
    const data = sortBy(playlists, sortStatus.columnAccessor);
    return sortStatus.direction === "desc" ? data.reverse() : data;
  }, [playlists, sortStatus]);

  return (
    <Table
      idAccessor="id"
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
          accessor: "public",
          textAlign: "right",
          render: ({ public: p }) => {
            if (p === undefined) return null

            return (
              <div className="flex justify-end">
                {p ? <FaCheck /> : <FaX />}
              </div>

            )
          }
        },
      ]}
      records={records}
      noRecordsText="No playlists. Run the playlist task to get started."
      sortStatus={sortStatus}
      onSortStatusChange={setSortStatus}
      fetching={isLoading}
    />
  )
}

