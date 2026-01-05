import { GeneratorSchema } from "@/lib/types/generator";
import { Group, Slider, Stack, Switch } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import { useState } from "react";
import { Button } from "../atoms/Button";
import { Section, SectionTitle } from "../atoms/Page";
import { TextInput } from "../molecules/TextInput";

type Props = {
  form: UseFormReturnType<GeneratorSchema>
  nextStep: () => void;
  prevStep: () => void;
}

export const GeneratorFormFinalize = ({ form, nextStep, prevStep }: Props) => {
  const [createPlaylist, useCreatePlaylist] = useState(form.getValues().createPlaylist)
  const [interval, setInterval] = useState(form.getValues().intervalDays)

  const handleCreatePlaylistChange = (checked: boolean) => {
    form.setFieldValue("createPlaylist", checked)
    useCreatePlaylist(checked)
  }

  const handleIntervalChange = (value: number) => {
    form.setFieldValue("intervalDays", value)
    setInterval(value)
  }

  return (
    <>
      <Section>
        <SectionTitle
          title="Generator Setup"
        />

        <TextInput label="Name" required {...form.getInputProps("name")} />
        <TextInput label="Description" {...form.getInputProps("description")} />

        <Group justify="space-between">
          <Stack gap={0}>
            <p className="text-sm font-medium">Create Spotify playlist</p>
            <p className="text-muted text-xs">If set to off, the generator only shows a preview.</p>
          </Stack>
          <Switch checked={createPlaylist} onChange={e => handleCreatePlaylistChange(e.target.checked)} color="secondary.1" />
        </Group>


        <Stack gap={0}>
          <p className="text-sm font-medium">Update Interval (days)</p>
          <p className="text-xs text-muted">Frequency to regenerate the tracks</p>
          <Group>
            <Slider
              value={interval}
              onChange={handleIntervalChange}
              color="secondary.1"
              restrictToMarks
              marks={Array.from({ length: 30 }).map((_, i) => ({ value: i + 1 }))}
              min={0}
              max={30}
              className="flex-1"
            />
            <p className="w-[2ch] text-right">{interval.toString().padStart(2, "0")}</p>
          </Group>
        </Stack>
      </Section>

      <Group justify="end">
        <Button onClick={prevStep} color="gray">Back</Button>
        <Button onClick={nextStep}>Save</Button>
      </Group>
    </>
  )
}
