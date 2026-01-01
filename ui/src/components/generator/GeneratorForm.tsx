import { useGeneratorGenerate } from "@/lib/api/generator";
import { GeneratorPreset, generatorPresetString, generatorSchema, GeneratorSchema } from "@/lib/types/generator";
import { ActionIcon, Divider, Group, NumberInput, Stack, Stepper } from "@mantine/core";
import { useForm, UseFormReturnType } from "@mantine/form";
import { useQueryClient } from "@tanstack/react-query";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { LuRotateCcw } from "react-icons/lu";
import { Button } from "../atoms/Button";
import { SectionTitle } from "../atoms/Page";
import { Table } from "../molecules/Table";
import { notifications } from "@mantine/notifications";

const maxSteps = 3

export const GeneratorForm = () => {
  const [active, setActive] = useState(0)

  const form = useForm<GeneratorSchema>({
    initialValues: {
      preset: GeneratorPreset.Top,
      params: {
        trackAmount: 50,
        minPlayCount: 5,
      },
    },
    validate: zod4Resolver(generatorSchema),
  })

  return (
    <div className="flex flex-col gap-4 h-full">
      <Stepper active={active} onStepClick={setActive} allowNextStepsSelect={false} className="flex-1">
        <Stepper.Step label="Preset & Parameters">
          <Preset form={form} />
        </Stepper.Step>
        <Stepper.Step label="Select tracks">
        </Stepper.Step>
        <Stepper.Step label="Finalize & Save">
        </Stepper.Step>
      </Stepper>
      <Divider />
      <Group justify="space-between">
        <p>{`Step ${active + 1} of ${maxSteps}`}</p>
        <Group>
          <Button color="gray">Cancel</Button>
          <Button>Next</Button>
        </Group>
      </Group>
    </div>
  )
}

const Preset = ({ form }: { form: UseFormReturnType<GeneratorSchema> }) => {
  const { data: tracks, isLoading, isRefetching } = useGeneratorGenerate(form.values)
  const queryClient = useQueryClient()

  const handleClickPreset = (p: GeneratorPreset) => {
    form.setFieldValue("preset", p)
  }

  const handleRefetchTracks = () => {
    if (form.validate().hasErrors) {
      notifications.show({ color: "red", message: "Some parameters are invalid" })
      return
    }

    queryClient.invalidateQueries({ queryKey: ["generator", "generate"] })
  }

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
            <Button key={String(p)} onClick={() => handleClickPreset(p)} color={form.values.preset === p ? "primary.3" : "gray"}>{generatorPresetString[p]}</Button>
          ))}
        </Group>
      </Stack>

      <Stack gap="xs">
        <p className="text-muted">Configure parameters</p>
        <Group>
          <NumberInput label="Amount of Tracks" {...form.getInputProps("params.trackAmount")} />
          <NumberInput label="Min Play Count" {...form.getInputProps("params.minPlayCount")} />
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
