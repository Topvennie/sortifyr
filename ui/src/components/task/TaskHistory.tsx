import { useTaskGetHistory } from "@/lib/api/task";
import { TaskHistoryFilter } from "@/lib/types/task";
import { formatDate } from "@/lib/utils";
import { useState } from "react";
import { LuTriangle } from "react-icons/lu";
import { Table } from "../molecules/Table";
import { TaskResult } from "./TaskResult";

type Props = {
  filter?: TaskHistoryFilter;
}

const formatDuration = (nanos: number) => {
  const ms = Math.floor(nanos / 1_000_000) % 1000
  const s = Math.floor(nanos / 1_000_000_000)

  const msString = `${ms.toString().padStart(3, '0')}ms`
  const sString = `${s.toString().padStart(2, '0')}s `

  return <span>{sString}<span className="text-muted-foreground">{msString}</span></span>
}

export const TaskHistory = ({ filter }: Props) => {
  const { history, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useTaskGetHistory(filter)

  const [expandedRecordIds, setExpandedRecordIds] = useState<number[]>([]);

  const handleBottom = () => {
    if (!hasNextPage) return
    if (isFetchingNextPage) return

    fetchNextPage()
  }

  return (
    <Table
      columns={[
        { accessor: "name", title: "Task" },
        { accessor: "runAt", render: ({ runAt }) => <p className="text-muted">{formatDate(runAt)}</p> },
        { accessor: "duration", render: ({ duration }) => <p className="text-muted">{formatDuration(duration)}</p> },
        {
          accessor: "result",
          title: "Result",
          render: ({ result }) => <TaskResult result={result} />
        },
        {
          accessor: "expanded",
          title: "",
          width: 24,
          render: ({ id, error }) => error && <LuTriangle className={`ml-auto fill-black w-2 transform duration-300 ${expandedRecordIds.includes(id) ? "" : "rotate-180"}`} />
        },
      ]}
      rowExpansion={{
        content: ({ record: { error } }) => (
          <div className="m-4">
            <p className="font-bold text-red-500">Error</p>
            <p className="text-red-500">{error}</p>
          </div>
        ),
        expandable: ({ record: { error } }) => Boolean(error),
        expanded: {
          recordIds: expandedRecordIds,
          onRecordIdsChange: setExpandedRecordIds,
        },
      }}
      records={history}
      noRecordsText="No tasks run yet"
      fetching={isLoading}
      onScrollToBottom={handleBottom}
    />
  )
}

