import { Playlist } from "@/lib/types/playlist";
import { usePlaylistGetAll } from "@/lib/api/playlist";
import { DataTable, type DataTableSortStatus } from 'mantine-datatable';
import { useState, useEffect } from "react";
import { FaCheck, FaX } from "react-icons/fa6";
import { LoadingSpinner } from "../molecules/LoadingSpinner";

type SortKey = "name" | "tracks" | "owner.name"

const sortBy = (playlists: Playlist[], key: SortKey): Playlist[] => {
  const getter: Record<SortKey, (p: Playlist) => string | number> = {
    name: p => p.name,
    tracks: p => p.tracks,
    "owner.name": p => p.owner?.name ?? "",
  };

  return [...playlists].sort((a, b) => {
    const av = getter[key](a);
    const bv = getter[key](b);
    return av === bv ? 0 : av > bv ? 1 : -1;
  });
};

export const PlaylistTableView = () => {
  const { data: playlists, isLoading } = usePlaylistGetAll()

  const [sortStatus, setSortStatus] = useState<DataTableSortStatus<Playlist>>({
    columnAccessor: "name",
    direction: "asc",
  })
  const [records, setRecords] = useState(sortBy(playlists ?? [], "name"))

  useEffect(() => {
    const data = sortBy(playlists ?? [], sortStatus.columnAccessor as SortKey);
    setRecords(sortStatus.direction === 'desc' ? data.reverse() : data);
  }, [sortStatus, playlists])

  if (isLoading) return <LoadingSpinner />

  return (
    <div className="max-w-full overflow-x-scroll">
      <DataTable
        striped
        highlightOnHover
        backgroundColor={"none"}
        columns={[
          { accessor: "name", sortable: true },
          { accessor: "tracks", sortable: true },
          { accessor: "owner.name", sortable: true },
          {
            accessor: "public", textAlign: "right", render: ({ public: p }) => (
              <div className="flex justify-end">
                {p ? <FaCheck /> : <FaX />}
              </div>

            )
          },
        ]}
        records={records}
        sortStatus={sortStatus}
        onSortStatusChange={setSortStatus}
      />
    </div>
  )
}

