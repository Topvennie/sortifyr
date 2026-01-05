import { GeneratorSchema } from "@/lib/types/generator";
import { Group, NumberInput, Stack, Switch, TextInput } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import { Button } from "../atoms/Button";
import { Section, SectionTitle } from "../atoms/Page";

type Props = {
  form: UseFormReturnType<GeneratorSchema>
  nextStep: () => void;
  prevStep: () => void;
}

export const GeneratorFormFinalize = ({ form, nextStep, prevStep }: Props) => {
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
          <Switch color="secondary.1" {...form.getInputProps("createPlaylist")} />
        </Group>

        <NumberInput label="Update interval (days)" description="Frequency to update the playlist if maintained" prefix="Days: " allowNegative={false} {...form.getInputProps("intervalDays")} />
      </Section>

      <Group justify="end">
        <Button onClick={prevStep} color="gray">Back</Button>
        <Button onClick={nextStep}>Save</Button>
      </Group>
    </>
  )
}
