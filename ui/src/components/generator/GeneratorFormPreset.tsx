import { useGeneratorPreview } from "@/lib/api/generator"
import { usePlaylistGetAll } from "@/lib/api/playlist"
import { GeneratorParamsSchema, GeneratorPreset, generatorPresetString, GeneratorSchema, GeneratorWindowSchema } from "@/lib/types/generator"
import { getValueByPath } from "@/lib/utils"
import { Alert, Group, Slider, Stack } from "@mantine/core"
import { DatesRangeValue } from "@mantine/dates"
import { UseFormReturnType } from "@mantine/form"
import { useDisclosure } from "@mantine/hooks"
import { notifications } from "@mantine/notifications"
import { isEqual } from "lodash"
import { ReactNode, useEffect, useState } from "react"
import { Button } from "../atoms/Button"
import { Section, SectionTitle } from "../atoms/Page"
import { Confirm } from "../molecules/Confirm"
import { DatePickerInput } from "../molecules/DatePickerInput"
import { Table } from "../molecules/Table"
import { GeneratorPlaylistTree } from "./GeneratorPlaylistTree"

type Props = {
  form: UseFormReturnType<GeneratorSchema>
  nextStep: () => void;
  prevStep: () => void;
}

export const GeneratorFormPreset = ({ form, nextStep, prevStep }: Props) => {
  const [params, setParams] = useState<Partial<GeneratorParamsSchema> | undefined>(form.getValues().params)

  const { mutate: generatorPreview, data: tracks, isPending } = useGeneratorPreview()
  useEffect(() => {
    const values = form.getValues()
    generatorPreview(values)
    setParams(values.params)
  }, [])

  const [opened, { open, close }] = useDisclosure()
  const handleNextInit = () => {
    if (!isEqual(params, form.getValues().params)) {
      open()
      return
    }

    nextStep()
  }

  const handleNext = () => {
    close()
    nextStep()
  }

  const [maxTracks, setMaxTracks] = useState(form.getValues().params?.trackAmount ?? 50)

  const handleClickPreset = (p: GeneratorPreset) => {
    form.setFieldValue("params.preset", p)
  }

  const handleMaxTracksChange = (amount: number) => {
    form.setFieldValue("params.trackAmount", amount)
    setMaxTracks(amount)
  }

  const handleRefetchTracks = () => {
    if (form.validateField("params").hasError) {
      notifications.show({ color: "red", message: "Some parameters are invalid" })
      console.error(form.errors)
      return
    }

    const values = form.getValues()
    generatorPreview(values)
    setParams(values.params)
  }

  const getPresetArguments = (preset: GeneratorPreset) => {
    switch (preset) {
      case GeneratorPreset.Top:
        return <Top form={form} />
      case GeneratorPreset.OldTop:
        return <OldTop form={form} />
    }
  }

  const [presetArguments, setPresetArguments] = useState<ReactNode>(getPresetArguments(form.getValues().params?.preset ?? GeneratorPreset.Top))
  form.watch("params.preset", ({ value }) => setPresetArguments(getPresetArguments(value as GeneratorPreset)))

  // TODO: next give arwning if the preview is not the same as the generator

  return (
    <>
      <div className="flex-1 flex flex-col md:flex-row gap-4 md:overflow-hidden">
        <Section className="flex-none md:w-[60%]">
          <SectionTitle
            title="Preset & Parameters"
            description="Pick a starting point and then fine tune the filters."
          />

          <p className="text-muted">Preset</p>
          <Group>
            {Object.values(GeneratorPreset).map(p => (
              <Button key={String(p)} onClick={() => handleClickPreset(p)} c={form.getValues().params?.preset === p ? "black" : "gray.6"} color={form.getValues().params?.preset === p ? "secondary.1" : "gray"}>{generatorPresetString[p]}</Button>
            ))}
          </Group>
          {presetArguments}

          <p className="text-muted">General parameters</p>
          <Stack gap={0}>
            <p className="text-sm font-medium">Maximum Tracks</p>
            <p className="text-xs text-muted">The maximum amount of tracks</p>
            <Group>
              <Slider
                value={maxTracks}
                onChange={handleMaxTracksChange}
                color="secondary.1"
                restrictToMarks
                marks={Array.from({ length: 40 }).map((_, i) => ({ value: (i + 1) * 5 }))}
                min={1}
                max={200}
                className="flex-1"
              />
              <p className="w-[2ch] text-right">{maxTracks.toString().padStart(3, "0")}</p>
            </Group>
          </Stack>

          <p className="text-muted">Select playlists</p>
          <Playlists form={form} />
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
        </Section>
      </div>

      <Group justify="end">
        <Button onClick={prevStep} color="gray">Cancel</Button>
        <Button onClick={handleNextInit}>Next: Tracks</Button>
      </Group>

      <Confirm
        opened={opened}
        onClose={close}
        modalTitle="Generator"
        title="Outdated Preview"
        description={`The preview is outdated. The actual playlist will look different.\nAre you sure you want to continue?`}
        onConfirm={handleNext}
      />
    </>
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

const Top = ({ form }: { form: UseFormReturnType<GeneratorSchema> }) => {
  return (
    <Stack>
      <Alert radius="lg" className="whitespace-pre-wrap">
        {`Top finds the tracks you're listening to the most right now.
It looks at your listening history within the selected time range and includes tracks that were played at least the minimum number of times within the given interval.`}
      </Alert>
      <p className="text-muted">Evaluation window</p>
      <Window form={form} path="params.paramsTop.window" />
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
      <p className="text-muted">Historic listening window</p>
      <Window form={form} path="params.paramsOldTop.peakWindow" />
      <p className="text-muted">Recent listening window</p>
      <Window form={form} path="params.paramsOldTop.recentWindow" />
    </Stack>
  )
}

const Window = ({ form, path }: { form: UseFormReturnType<GeneratorSchema>, path: string }) => {
  const window = getValueByPath<GeneratorWindowSchema>(form.getValues(), path)

  const [range, setRange] = useState<[Date | null, Date | null]>([window?.start ?? null, window?.end ?? null])
  const [plays, setPlays] = useState<number>(window?.minPlays ?? 0)
  const [burst, setBurst] = useState<number>(window?.burstIntervalDays ?? 0)

  const handleRangeChange = (r: DatesRangeValue) => {
    form.setFieldValue(`${path}.start`, r[0] ?? undefined)
    form.setFieldValue(`${path}.end`, r[1] ?? undefined)

    setRange(r)
  }

  const handlePlayChange = (plays: number) => {
    form.setFieldValue(`${path}.minPlays`, plays)
    setPlays(plays)
  }

  const handleBurstChange = (burst: number) => {
    form.setFieldValue(`${path}.burstIntervalDays`, burst)
    setBurst(burst)
  }

  return (
    <Stack gap="xs">
      <Stack gap={0}>
        <p className="text-sm font-medium">Plays</p>
        <p className="text-xs text-muted">The minimum amount of plays to trigger the window</p>
        <Group>
          <Slider
            value={plays}
            onChange={handlePlayChange}
            color="secondary.1"
            restrictToMarks
            marks={Array.from({ length: 20 }).map((_, i) => ({ value: i + 1 }))}
            min={1}
            max={20}
            className="flex-1"
          />
          <p className="w-[2ch] text-right">{plays.toString().padStart(2, "0")}</p>
        </Group>
      </Stack>
      <Stack gap={0}>
        <p className="text-sm font-medium">Burst Interval (days)</p>
        <p className="text-xs text-muted">Time window in which the minimum plays must occur</p>
        <Group>
          <Slider
            value={burst}
            onChange={handleBurstChange}
            color="secondary.1"
            restrictToMarks
            marks={Array.from({ length: 30 }).map((_, i) => ({ value: i + 1 }))}
            min={1}
            max={30}
            className="flex-1"
          />
          <p className="w-[2ch] text-right">{burst.toString().padStart(2, "0")}</p>
        </Group>
      </Stack>
      <Stack gap={0}>
        <p className="text-sm font-medium">Time Range</p>
        <p className="text-xs text-muted">The total time range in which to look for the minimum amount of plays in the interval</p>
        <DatePickerInput
          type="range"
          numberOfColumns={2}
          value={range}
          onChange={handleRangeChange}
        />
      </Stack>
    </Stack>
  )
}
