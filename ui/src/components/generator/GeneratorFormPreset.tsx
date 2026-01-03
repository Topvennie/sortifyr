import { useGeneratorPreview } from "@/lib/api/generator"
import { usePlaylistGetAll } from "@/lib/api/playlist"
import { GeneratorPreset, generatorPresetString, GeneratorSchema, GeneratorWindowSchema } from "@/lib/types/generator"
import { getValueByPath } from "@/lib/utils"
import { Alert, Group, NumberInput, Stack } from "@mantine/core"
import { DatePickerInput, DatesRangeValue } from "@mantine/dates"
import { UseFormReturnType } from "@mantine/form"
import { notifications } from "@mantine/notifications"
import { ReactNode, useEffect, useState } from "react"
import { Button } from "../atoms/Button"
import { Section, SectionTitle } from "../atoms/Page"
import { Table } from "../molecules/Table"
import { GeneratorPlaylistTree } from "./GeneratorPlaylistTree"

type Props = {
  form: UseFormReturnType<GeneratorSchema>
  nextStep: () => void;
  prevStep: () => void;
}

export const GeneratorFormPreset = ({ form, nextStep, prevStep }: Props) => {
  const { mutate: generatorPreview, data: tracks, isPending } = useGeneratorPreview()
  useEffect(() => generatorPreview(form.getValues()), [])

  const handleClickPreset = (p: GeneratorPreset) => {
    form.setFieldValue("params.preset", p)
  }

  const handleRefetchTracks = () => {
    if (form.validateField("params").hasError) {
      notifications.show({ color: "red", message: "Some parameters are invalid" })
      return
    }

    generatorPreview(form.getValues())
  }

  const getPresetArguments = (preset: GeneratorPreset) => {
    switch (preset) {
      case GeneratorPreset.Forgotten:
        return <Forgotten form={form} />
      case GeneratorPreset.Top:
        return <Top form={form} />
      case GeneratorPreset.OldTop:
        return <OldTop form={form} />
      default:
        return <Custom form={form} />
    }
  }

  const [presetArguments, setPresetArguments] = useState<ReactNode>(getPresetArguments(form.getValues().params?.preset ?? GeneratorPreset.Custom))
  form.watch("params.preset", ({ value }) => setPresetArguments(getPresetArguments(value as GeneratorPreset)))

  // TODO: next give arwning if the preview is not the same as the generator

  return (
    <div className="flex-1 flex flex-col md:flex-row gap-4 md:overflow-hidden">
      <Section className="flex-none md:w-[60%]">
        <SectionTitle
          title="Preset & Parameters"
          description="Pick a starting point and then fine tune the filters."
        />

        <Stack gap="xs">
          <p className="text-muted">Preset</p>
          <Group>
            {Object.values(GeneratorPreset).map(p => (
              <Button key={String(p)} onClick={() => handleClickPreset(p)} c={form.getValues().params?.preset === p ? "black" : "gray.6"} color={form.getValues().params?.preset === p ? "secondary.1" : "gray"}>{generatorPresetString[p]}</Button>
            ))}
          </Group>
          {presetArguments}
        </Stack>

        <Stack gap="xs">
          <p className="text-muted">General parameters</p>
          <Group>
            <NumberInput label="Maximum Tracks" allowNegative={false} {...form.getInputProps("params.trackAmount")} />
          </Group>
        </Stack>

        <Stack gap="xs">
          <p className="text-muted">Select playlists</p>
          <Playlists form={form} />
        </Stack>
      </Section>

      <Section>
        <Group justify="space-between">
          <SectionTitle
            title="Preview Tracks"
            description="Adjust parameters and then refresh to update the preview."
          />
          <Button onClick={handleRefetchTracks} color="secondary.1" loading={isPending}>Refresh</Button>
        </Group>

        <Table
          columns={[
            { accessor: "name" },
          ]}
          records={tracks ?? []}
          noRecordsText="No tracks fit the parameters"
          fetching={isPending}
        />
        <Group justify="end">
          <Button onClick={prevStep} color="gray">Cancel</Button>
          <Button onClick={nextStep}>Next: Tracks</Button>
        </Group>
      </Section>
    </div>
  )
}

const Playlists = ({ form }: { form: UseFormReturnType<GeneratorSchema> }) => {
  const { data: playlists, isLoading } = usePlaylistGetAll()
  const [excluded, setExcluded] = useState<number[]>(form.getValues().params?.excludedPlaylistIds ?? [])

  const handleToggle = (playlistIds: number[], isExcluded: boolean) => {
    let newExcluded: number[]

    if (isExcluded) newExcluded = [...excluded, ...playlistIds.filter(p => !excluded.includes(p))]
    else newExcluded = excluded.filter(e => !playlistIds.includes(e))

    form.setFieldValue("params.excludedPlaylistIds", newExcluded)
    setExcluded(newExcluded)
  }

  const handleDeselectAll = () => {
    const newExcluded = playlists?.map(p => p.id) ?? []

    form.setFieldValue("params.excludedPlaylistIds", newExcluded)
    setExcluded(newExcluded)
  }

  const handleSelectAll = () => {
    const newExcluded: number[] = []

    form.setFieldValue("params.excludedPlaylistIds", newExcluded)
    setExcluded(newExcluded)
  }

  return (
    <Stack>
      <Group>
        <Button onClick={handleDeselectAll} color="secondary.1">Deselect all</Button>
        <Button onClick={handleSelectAll} color="secondary.1">Select all</Button>
        <p className="text-muted ml-auto">{`Selected ${(playlists?.length ?? 0) - excluded.length}`}</p>
      </Group>
      <GeneratorPlaylistTree playlists={playlists ?? []} isLoading={isLoading} excluded={excluded} onToggle={handleToggle} />
    </Stack>
  )
}

const Custom = ({ form }: { form: UseFormReturnType<GeneratorSchema> }) => {
  return null
}

const Forgotten = ({ form }: { form: UseFormReturnType<GeneratorSchema> }) => {
  return null
}

const Top = ({ form }: { form: UseFormReturnType<GeneratorSchema> }) => {
  return (
    <Stack>
      <Alert radius="lg" className="whitespace-pre-wrap">
        {`Top finds the tracks you're listening to the most right now.
It looks at your listening history within the selected time range and includes tracks that were played at least the minimum number of times within the given interval.`}
      </Alert>
      <Stack gap={0}>
        <p className="text-muted">Evaluation window</p>
        <Window form={form} path="params.paramsTop.window" />
      </Stack>
    </Stack>
  )
}

const OldTop = ({ form }: { form: UseFormReturnType<GeneratorSchema> }) => {
  return (
    <Stack>
      <Alert radius="lg" className="whitespace-pre-wrap">
        {`Old Top finds tracks you used to listen to on repeat, but don’t play much anymore.
It works in two steps:

1. Historic window:
   Finds tracks that were played at least the minimum number of times within the burst interval somewhere in the historic range.

2. Recent window:
   Filters out tracks that were still played frequently in the recent range.

Example:
If the historic range is 6 months, the minimum plays is 5, and the burst interval is 14 days, it will find tracks you played 5 or more times within any 14-day period during those 6 months. But only if you haven’t played them much recently.`}
      </Alert>
      <Stack gap={0}>
        <p className="text-muted">Historic listening window</p>
        <Window form={form} path="params.paramsOldTop.peakWindow" />
      </Stack>
      <Stack gap={0}>
        <p className="text-muted">Recent listening window</p>
        <Window form={form} path="params.paramsOldTop.recentWindow" />
      </Stack>
    </Stack>
  )
}

const Window = ({ form, path }: { form: UseFormReturnType<GeneratorSchema>, path: string }) => {
  const window = getValueByPath<GeneratorWindowSchema>(form.getValues(), path)

  const [range, setRange] = useState<[Date | null, Date | null]>([window?.start ?? null, window?.end ?? null])
  const [plays, setPlays] = useState<number | string>(window?.minPlays ?? 0)
  const [burst, setBurst] = useState<number | string>((window?.burstIntervalS ?? 0) / (24 * 60 * 60))

  const handleRangeChange = (r: DatesRangeValue) => {
    form.setFieldValue(`${path}.start`, r[0] ?? undefined)
    form.setFieldValue(`${path}.end`, r[1] ?? undefined)

    setRange(r)
  }

  const handlePlayChange = (plays: number | string) => {
    let newValue: number | undefined = undefined
    if (typeof plays === 'number') {
      newValue = plays
    }

    form.setFieldValue(`${path}.minPlays`, newValue)
    setPlays(plays)
  }

  const handleIntervalChange = (days: number | string) => {
    let newValue: number | undefined = undefined
    if (typeof days === 'number') {
      newValue = days * 24 * 60 * 60
    }

    form.setFieldValue(`${path}.burstIntervalS`, newValue)
    setBurst(days)
  }

  return (
    <Group>
      <NumberInput label="Minimum Plays" description="Minimum number of times a track must be played" value={plays} onChange={handlePlayChange} allowNegative={false} />
      <NumberInput label="Burst Interval (days)" description="Time window in which the minimum plays must occur" prefix="Days: " value={burst} onChange={handleIntervalChange} allowNegative={false} />
      <DatePickerInput label="Time range" description="Listening history to analyze" type="range" numberOfColumns={2} value={range} onChange={handleRangeChange} />
    </Group>
  )
}
