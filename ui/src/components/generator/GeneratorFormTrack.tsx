import { GeneratorSchema } from "@/lib/types/generator"
import { Button, Group, Stack } from "@mantine/core"
import { UseFormReturnType } from "@mantine/form"
import { SectionTitle } from "../atoms/Page"
import { useGeneratorPreview } from "@/lib/api/generator"
import { useQueryClient } from "@tanstack/react-query"
import { Table } from "../molecules/Table"
import { Track } from "@/lib/types/track"
import { useState } from "react"
import { LuRotateCcw } from "react-icons/lu"

type Props = {
  form: UseFormReturnType<GeneratorSchema>
}

export const GeneratorFormTrack = ({ form }: Props) => {
  const { data: tracks, isLoading, isRefetching } = useGeneratorPreview(form.getValues())
  const [selectedRecords, setSelectedRecords] = useState<Track[]>([]);
  const queryClient = useQueryClient()

  const handleRefetch = () => {
    queryClient.invalidateQueries({ queryKey: ["generator", "preview"] })
  }

  const handleExcludeTracks = () => {
    form.setFieldValue("params.excludedTrackIds", [...form.getValues().params?.excludedTrackIds ?? [], ...selectedRecords.map(t => t.id)])
    handleRefetch()
  }

  return (
    <Stack gap="lg">
      <SectionTitle
        title="Select Tracks"
        description="Toggle off any tracks you want to exclude. Only selected tracks will be added to the playlist."
      />

      <Stack gap="xs">
        <Table
          columns={[
            { accessor: "name" },
          ]}
          records={tracks ?? []}
          height={512}
          noRecordsText="No tracks fit the parameters"
          fetching={isLoading || isRefetching}
          selectedRecords={selectedRecords}
          onSelectedRecordsChange={setSelectedRecords}
          selectionTrigger="cell"
        />
        <Group>
          <Button onClick={handleRefetch} variant="outline" leftSection={<LuRotateCcw />}>Reload</Button>
          <Button onClick={handleExcludeTracks} >Exclude tracks</Button>
        </Group>
      </Stack>
    </Stack>
  )
}
