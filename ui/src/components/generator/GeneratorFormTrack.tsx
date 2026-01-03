import { useGeneratorPreview } from "@/lib/api/generator"
import { useTrackGetAllById } from "@/lib/api/track"
import { GeneratorSchema } from "@/lib/types/generator"
import { Track } from "@/lib/types/track"
import { ActionIcon, Group } from "@mantine/core"
import { UseFormReturnType } from "@mantine/form"
import { useEffect, useState } from "react"
import { Button } from "../atoms/Button"
import { Section, SectionTitle } from "../atoms/Page"
import { Table } from "../molecules/Table"
import { LuTrash2, LuUndo2 } from "react-icons/lu"

type Props = {
  form: UseFormReturnType<GeneratorSchema>
  nextStep: () => void;
  prevStep: () => void;
}

export const GeneratorFormTrack = ({ form, nextStep, prevStep }: Props) => {
  const { mutate: generatorPreview, data: tracksInitial, isPending } = useGeneratorPreview()
  useEffect(() => generatorPreview(form.getValues()), [])
  const [tracks, setTracks] = useState<Track[]>(tracksInitial ?? [])
  useEffect(() => {
    if (!tracksInitial) return
    setTracks(tracksInitial)
  }, [tracksInitial])

  const { data: excludedTracksInitial, isLoading } = useTrackGetAllById()
  const [excludedTracks, setExcludedTracks] = useState<Track[]>(excludedTracksInitial ?? [])
  useEffect(() => {
    if (!excludedTracksInitial) return
    setExcludedTracks(excludedTracksInitial)
  }, [excludedTracksInitial])

  const handleRefetchTracks = () => {
    generatorPreview(form.getValues())
  }

  const handleExcludeTrack = (track: Track) => {
    form.setFieldValue("params.excludedTrackIds", [...(form.getValues().params?.excludedTrackIds ?? []), track.id])
    setTracks(prev => prev.filter(t => t.id !== track.id))
    setExcludedTracks(prev => [...prev, track].sort((a, b) => a.name.localeCompare(b.name)))
  }

  const handleIncludeTrack = (track: Track) => {
    form.setFieldValue("params.excludedTrackIds", form.getValues().params?.excludedTrackIds?.filter(t => t !== track.id))
    setTracks(prev => [...prev, track].sort((a, b) => a.name.localeCompare(b.name)))
    setExcludedTracks(prev => prev.filter(t => t.id !== track.id))
  }

  return (
    <>
      <div className="flex-1 flex flex-col md:flex-row gap-4 md:overflow-hidden">
        <Section className="flex-none md:w-[60%]">
          <Group justify="space-between">
            <SectionTitle
              title="Generated Tracks"
              description={`Exclude tracks from the playlist.\nRefresh to replace them with new tracks.`}
            />
            <Button onClick={handleRefetchTracks} color="secondary.1">Refresh</Button>
          </Group>

          <Table
            columns={[
              { accessor: "name", width: "100%" },
              {
                accessor: "actions",
                title: "",
                textAlign: "right",
                render: track => <ActionIcon onClick={() => handleExcludeTrack(track)} variant="subtle" color="red"><LuTrash2 /></ActionIcon>
              },
            ]}
            records={tracks ?? []}
            noRecordsText="No tracks fit the parameters"
            fetching={isPending || isLoading}
            animate={false}
          />
        </Section>

        <Section>
          <SectionTitle
            title="Excluded Tracks"
            description={`Readd tracks back to the playlist.\nRefresh to remove any excess tracks.`}
          />

          <Table
            columns={[
              { accessor: "name" },
              {
                accessor: "actions",
                title: "",
                textAlign: "right",
                render: track => <ActionIcon onClick={() => handleIncludeTrack(track)} variant="subtle" color="black"><LuUndo2 /></ActionIcon>
              },
            ]}
            records={excludedTracks ?? []}
            noRecordsText="No excluded tracks"
            fetching={isPending || isLoading}
            animate={false}
          />
        </Section>
      </div>

      <Group justify="end">
        <Button onClick={prevStep} color="gray">Back</Button>
        <Button onClick={nextStep}>Next: Finalize</Button>
      </Group>
    </>
  )
}
