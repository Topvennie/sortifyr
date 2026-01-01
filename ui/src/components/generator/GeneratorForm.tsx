import { GeneratorPreset, generatorSchema, GeneratorSchema } from "@/lib/types/generator";
import { daysAgo } from "@/lib/utils";
import { Divider, Group, Stepper } from "@mantine/core";
import { useForm } from "@mantine/form";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { Button } from "../atoms/Button";
import { GeneratorFormPreset } from "./GeneratorFormPreset";

const maxSteps = 3

export const GeneratorForm = () => {
  const [active, setActive] = useState(0)

  const form = useForm<GeneratorSchema>({
    initialValues: {
      name: "",
      description: undefined,
      params: {
        trackAmount: 50,
        excludedPlaylistIds: undefined,
        excludedTrackIds: undefined,
        preset: GeneratorPreset.Top,
        paramsCustom: {},
        paramsForgotten: {},
        paramsTop: {
          window: {
            start: daysAgo(14), // 14 days ago
            end: new Date(),
            minPlays: 5,
            burstIntervalS: 14 * 24 * 60 * 60 // 14 days
          },
        },
        paramsOldTop: {
          peakWindow: {
            start: daysAgo(365), // 365 days ago
            end: daysAgo(100), // 100 days ago
            minPlays: 5,
            burstIntervalS: 14 * 24 * 60 * 60 // 14 days
          },
          recentWindow: {
            start: daysAgo(14), // 14 days ago
            end: new Date(),
            minPlays: 5,
            burstIntervalS: 14 * 24 * 60 * 60 // 14 days
          }
        }
      }
    },
    validate: zod4Resolver(generatorSchema),
  })

  return (
    <div className="flex flex-col gap-4 h-full">
      <Stepper active={active} onStepClick={setActive} allowNextStepsSelect={false} className="flex-1 overflow-y-auto">
        <Stepper.Step label="Preset & Parameters">
          <GeneratorFormPreset form={form} />
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

