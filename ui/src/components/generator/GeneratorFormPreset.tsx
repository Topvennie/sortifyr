import { useGeneratorPreview } from "@/lib/api/generator"
import { GeneratorPreset, generatorPresetString, GeneratorSchema, GeneratorWindowSchema } from "@/lib/types/generator"
import { getValueByPath } from "@/lib/utils"
import { ActionIcon, Alert, Group, NumberInput, Stack } from "@mantine/core"
import { DatePickerInput, DatesRangeValue } from "@mantine/dates"
import { UseFormReturnType } from "@mantine/form"
import { notifications } from "@mantine/notifications"
import { useQueryClient } from "@tanstack/react-query"
import { ReactNode, useState } from "react"
import { LuRotateCcw } from "react-icons/lu"
import { Button } from "../atoms/Button"
import { SectionTitle } from "../atoms/Page"
import { Table } from "../molecules/Table"

type Props = {
  form: UseFormReturnType<GeneratorSchema>
}

export const GeneratorFormPreset = ({ form }: Props) => {
  const { data: tracks, isLoading, isRefetching } = useGeneratorPreview(form.values)
  const queryClient = useQueryClient()

  const handleClickPreset = (p: GeneratorPreset) => {
    form.setFieldValue("params.preset", p)
  }

  const handleRefetchTracks = () => {
    if (form.validateField("params").hasError) {
      notifications.show({ color: "red", message: "Some parameters are invalid" })
      console.error(form.errors)
      return
    }

    queryClient.invalidateQueries({ queryKey: ["generator", "preview"] })
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

  const [presetArguments, setPresetArguments] = useState<ReactNode>(getPresetArguments(form.values.params?.preset ?? GeneratorPreset.Custom))
  form.watch("params.preset", ({ value }) => setPresetArguments(getPresetArguments(value as GeneratorPreset)))

  return (
    <Stack gap="lg">
      <SectionTitle
        title="Preset & Parameters"
        description="Pick a starting point, then fine tune the filters."
      />

      <Stack gap="xs">
        <p className="text-muted">Choose a preset</p>
        <Group>
          {Object.values(GeneratorPreset).map(p => (
            <Button key={String(p)} onClick={() => handleClickPreset(p)} color={form.values.params?.preset === p ? "primary.3" : "gray"}>{generatorPresetString[p]}</Button>
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
        <Group gap="xs">
          <p className="text-muted">Preview</p>
          <ActionIcon onClick={handleRefetchTracks} variant="subtle" c="black" loading={isLoading || isRefetching}>
            <LuRotateCcw />
          </ActionIcon>
        </Group>
        <Table
          columns={[
            { accessor: "name" },
          ]}
          records={tracks ?? []}
          height={512}
          noRecordsText="No tracks fit the parameters"
          fetching={isLoading || isRefetching}
          noHeader
        />
      </Stack>
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
      <Alert radius="lg" title="Top" className="whitespace-pre-wrap">
        {`Top finds the tracks you're listening to the most right now.

It looks at your listening history within the selected time range and includes tracks that were played at least the minimum number of times within the given interval.`}
      </Alert>
      <Stack gap={0}>
        <p className="font-semibold">Evaluation window</p>
        <Window form={form} path="params.paramsTop.window" />
      </Stack>
    </Stack>
  )
}

const OldTop = ({ form }: { form: UseFormReturnType<GeneratorSchema> }) => {
  return (
    <Stack>
      <Alert radius="lg" title="Old Top" className="whitespace-pre-wrap">
        {`Old Top finds tracks you used to listen to on repeat, but don’t play much anymore.

It works in two steps:

1. Historic window:
   Finds tracks that were played at least the minimum number of times within the burst interval somewhere in the historic range.

2. Recent window:
   Filters out tracks that were still played frequently in the recent range.

Example:
If the historic range is 6 months, the minimum plays is 5, and the burst interval is 14 days, it will find tracks you played 5 or more times within any 14-day period during those 6 months — but only if you haven’t played them much recently.`}
      </Alert>
      <Stack gap={0}>
        <p className="font-semibold">Historic listening window</p>
        <Window form={form} path="params.paramsOldTop.peakWindow" />
      </Stack>
      <Stack gap={0}>
        <p className="font-semibold">Recent listening window</p>
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
