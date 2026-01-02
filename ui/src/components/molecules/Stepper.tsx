import { LuCheck } from "react-icons/lu";
import { Fragment } from "react/jsx-runtime";

type Props = {
  steps: Step[];
  activeStep: number
}

export type Step = {
  title: string;
}

export const Stepper = ({ steps, activeStep }: Props) => {
  return (
    <div className="flex items-center justify-evenly gap-8 w-full">
      {steps.map((s, idx) => {
        const active = idx === activeStep
        const prev = idx < activeStep

        return (
          <Fragment key={s.title}>
            {idx > 0 && <div className={`flex-1 h-0.5 ${active ? "bg-(--mantine-color-primary-3)" : "bg-gray-300"}`} />}
            <div className="flex items-center gap-2">
              <p className={`flex items-center justify-center w-10 h-10 rounded-full font-semibold transition duration-500 ${prev ? "bg-(--mantine-color-primary-3)" : active ? "bg-(--mantine-color-primary-1)" : "bg-gray-200"} ${active ? "border border-(--mantine-color-primary-3)" : ""}`}>
                {prev
                  ? <LuCheck className="size-6 text-white" />
                  : idx + 1
                }
              </p>
              <p>{s.title}</p>
            </div>
          </Fragment>
        )
      })}
    </div>
  )
}
