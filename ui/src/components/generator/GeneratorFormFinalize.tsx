import { GeneratorSchema } from "@/lib/types/generator";
import { Group, TextInput } from "@mantine/core";
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

        {/* <Group justify="space-between"> */}
        {/*   <Stack gap={0}> */}
        {/*     <p className="font-bold">Create Spotify playlist</p> */}
        {/*     <p className="text-muted text-sm">If set to off, the generator only shows a preview.</p> */}
        {/*   </Stack> */}
        {/*   <Switch {...form.getInputProps("")} /> */}
        {/* </Group> */}
        {/**/}
        {/* <Group justify="space-between"> */}
        {/*   <Stack gap={0}> */}
        {/*     <p className="font-bold">Auto-maintain playlist</p> */}
        {/*     <p className="text-muted text-sm">Keep the playlist in sync with the new listening data.</p> */}
        {/*   </Stack> */}
        {/*   <Switch {...form.getInputProps("")} /> */}
        {/* </Group> */}
      </Section>

      <Group justify="end">
        <Button onClick={prevStep} color="gray">Back</Button>
        <Button onClick={nextStep}>Save</Button>
      </Group>
    </>
  )
}
