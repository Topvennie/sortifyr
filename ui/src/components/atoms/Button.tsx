import { ButtonProps, Button as MButton } from "@mantine/core"
import { ComponentProps } from "react"

type Props = ButtonProps & ComponentProps<"button">

export const Button = (props: Props) => {
  return <MButton c="black" radius="lg" {...props} />
}
