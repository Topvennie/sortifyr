import { TextInput as MTextInput, TextInputProps } from "@mantine/core";

type Props = TextInputProps

export const TextInput = (props: Props) => {
  return (
    <MTextInput
      radius="lg"
      {...props}
    />
  )
}
