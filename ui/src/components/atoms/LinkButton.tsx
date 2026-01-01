import { Button, ButtonProps } from "@mantine/core";
import {
  createLink,
  LinkComponent,
} from "@tanstack/react-router";
import React from "react";

interface LinkButtonProps extends Omit<ButtonProps, "href"> {
  permission?: string
}

const ButtonLinkComponent = React.forwardRef<
  HTMLButtonElement,
  LinkButtonProps
>((props, ref) => {
  return <Button c="black" radius="lg" ref={ref} {...props} />
})

const CreatedLinkComponent = createLink(ButtonLinkComponent);

export const LinkButton: LinkComponent<typeof ButtonLinkComponent> = (props) => {
  return <CreatedLinkComponent preload="intent" {...props} />
}
