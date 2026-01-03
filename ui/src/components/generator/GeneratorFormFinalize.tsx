import { GeneratorSchema } from "@/lib/types/generator";
import { Group, NumberInput, Stack, Switch, TextInput } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import { Button } from "../atoms/Button";
import { Section, SectionTitle } from "../atoms/Page";
import { useState } from "react";

type Props = {
  form: UseFormReturnType<GeneratorSchema>
  nextStep: () => void;
  prevStep: () => void;
}

export const GeneratorFormFinalize = ({ form, nextStep, prevStep }: Props) => {
  const [playlist, setPlaylist] = useState(form.getValues().playlist)
  const [maintained, setMaintained] = useState(form.getValues().maintained)
  const [interval, setInterval] = useState<number | string>(form.getValues().intervalS)

  const handleTogglePlaylist = () => {
    form.setFieldValue("playlist", !playlist)
    setPlaylist(prev => !prev)

    if (!playlist) {
      form.setFieldValue("maintained", false)
      form.setFieldValue("intervalS", 0)
      setMaintained(false)
    }
  }

  const handleToggleMaintained = () => {
    form.setFieldValue("maintained", !maintained)
    setMaintained(prev => !prev)

    if (!maintained) {
      form.setFieldValue("intervalS", 0)
    }
  }

  const handleChangeInterval = (interval: number | string) => {
    let newValue = 0
    if (typeof interval === 'number') {
      newValue = interval * 24 * 60 * 60
    }

    form.setFieldValue("intervalS", newValue)
    setInterval(interval)
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
            <p className="text-sm">Create Spotify playlist</p>
            <p className="text-muted text-xs">If set to off, the generator only shows a preview.</p>
          </Stack>
          <Switch checked={playlist} onChange={handleTogglePlaylist} />
        </Group>

        <Group justify="space-between">
          <Stack gap={0}>
            <p className={`text-sm ${!playlist && "text-gray-400"}`}>Auto-maintain playlist</p>
            <p className="text-muted text-xs">Keep the playlist in sync with the new listening data.</p>
          </Stack>
          <Switch checked={maintained} onChange={handleToggleMaintained} disabled={!playlist} />
        </Group>

        <NumberInput label="Update interval (days)" description="Frequency to update the playlist if maintained" prefix="Days: " value={interval} onChange={handleChangeInterval} allowNegative={false} disabled={!maintained} />
      </Section>

      <Group justify="end">
        <Button onClick={prevStep} color="gray">Back</Button>
        <Button onClick={nextStep}>Save</Button>
      </Group>
    </>
  )
}
