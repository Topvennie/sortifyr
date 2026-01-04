import { useGeneratorCreate, useGeneratorEdit } from "@/lib/api/generator";
import { convertGeneratorSchema, Generator, GeneratorPreset, generatorSchema, GeneratorSchema } from "@/lib/types/generator";
import { daysAgo, getErrorMessage } from "@/lib/utils";
import { Stack } from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { useNavigate } from "@tanstack/react-router";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { Confirm } from "../molecules/Confirm";
import { Step, Stepper } from "../molecules/Stepper";
import { GeneratorFormFinalize } from "./GeneratorFormFinalize";
import { GeneratorFormPreset } from "./GeneratorFormPreset";
import { GeneratorFormTrack } from "./GeneratorFormTrack";

type Props = {
  generator?: Generator;
}

// TODO: Before saving an update, warn about deleting

const steps: Step[] = [
  { title: "Preset & Parameters" },
  { title: "Select Tracks" },
  { title: "Finalize & Save" }
]

export const GeneratorForm = ({ generator: initialGenerator }: Props) => {
  const [active, setActive] = useState(0)
  const [opened, { open, close }] = useDisclosure()
  const [submitting, setSubmitting] = useState(false)

  const generatorCreate = useGeneratorCreate()
  const generatorEdit = useGeneratorEdit()

  const navigate = useNavigate()

  const form = useForm<GeneratorSchema>({
    initialValues: initialGenerator ? convertGeneratorSchema(initialGenerator) : {
      name: "",
      description: "",
      createPlaylist: false,
      maintained: false,
      intervalS: 0,
      params: {
        trackAmount: 50,
        excludedPlaylistIds: [],
        excludedTrackIds: [],
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

  const handleSubmit = () => {
    if (form.validate().hasErrors) {
      notifications.show({ color: "red", message: "Some parameters are invalid" })
      console.error(form.errors)
      return
    }

    const update = !!initialGenerator

    let action
    if (update) action = generatorEdit
    else action = generatorCreate

    setSubmitting(true)

    action.mutateAsync(form.getValues(), {
      onSuccess: () => {
        notifications.show({ title: form.getValues().name, message: `Generator ${update ? "updated" : "created"}` })
        navigate({ to: "/generator" })
      },
      onError: async (error) => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      },
      onSettled: () => {
        close()
        setSubmitting(false)
      },
    })
  }

  const handleNextStep = () => {
    if (active < steps.length - 1) {
      setActive(prev => prev + 1)
      return
    }

    if (form.validate().hasErrors) {
      notifications.show({ color: "red", message: "Some parameters are invalid" })
      console.error(form.errors)
      return
    }

    open()
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
    <>
      <Stack className="flex-1 rounded-xl overflow-auto">
        <Stepper steps={steps} activeStep={active} />
        {stepComponent()}
      </Stack>
      <Confirm
        opened={opened}
        onClose={close}
        modalTitle="Generator"
        title="Save Generator"
        description="Are you sure you want to save the generator?"
        onConfirm={handleSubmit}
        loading={submitting}
      />
    </>
  )
}
