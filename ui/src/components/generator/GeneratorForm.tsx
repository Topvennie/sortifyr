import { GeneratorPreset, generatorSchema, GeneratorSchema } from "@/lib/types/generator";
import { daysAgo } from "@/lib/utils";
import { useForm } from "@mantine/form";
import { useNavigate } from "@tanstack/react-router";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useMemo, useState } from "react";
import { Step, Stepper } from "../molecules/Stepper";
import { GeneratorFormPreset } from "./GeneratorFormPreset";
import { GeneratorFormTrack } from "./GeneratorFormTrack";
import { Stack } from "@mantine/core";
import { GeneratorFormFinalize } from "./GeneratorFormFinalize";

const steps: Step[] = [
  {
    title: "Preset & Parameters",
  },
  {
    title: "Select Tracks",
  },
  {
    title: "Finalize & Save",
  }
]

export const GeneratorForm = () => {
  const [active, setActive] = useState(0)

  const navigate = useNavigate()

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
            minPlays: 2,
            burstIntervalS: 14 * 24 * 60 * 60 // 14 days
          }
        }
      }
    },
    validate: zod4Resolver(generatorSchema),
  })

  const handleNextStep = () => {
    if (active < steps.length - 1) {
      setActive(prev => prev + 1)
      return
    }
  }

  const handlePrevStep = () => {
    if (active === 0) {
      navigate({ to: "/generator" })
      return
    }

    setActive(prev => prev - 1)
  }

  const stepComponent = () => {
    switch (active) {
      case 0:
        return <GeneratorFormPreset form={form} nextStep={handleNextStep} prevStep={handlePrevStep} />
      case 1:
        return <GeneratorFormTrack form={form} nextStep={handleNextStep} prevStep={handlePrevStep} />
      case 2:
        return <GeneratorFormFinalize form={form} nextStep={handleNextStep} prevStep={handlePrevStep} />
      default:
        return null
    }
  }

  return (
    <Stack className="flex-1 rounded-xl overflow-auto">
      <Stepper steps={steps} activeStep={active} />
      {stepComponent()}
    </Stack>
  )
}
