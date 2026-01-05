import { NumberInput as MNumberInput, NumberInputProps } from "@mantine/core";

type Props = NumberInputProps

export const NumberInput = (props: Props) => {
  return (
    <MNumberInput
      radius="lg"
      styles={{
        control: {
          borderLeft: "none",
        },
      }}
      {...props}
    />
  )
}
