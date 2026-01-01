import { SegmentedControl, SegmentedControlProps } from "@mantine/core";

type Props = {
  secondary?: boolean;
} & SegmentedControlProps

export const Segment = ({ secondary, ...props }: Props) => {
  return (
    <SegmentedControl
      radius="lg"
      color={secondary ? "secondary.1" : "primary.3"}
      styles={{
        innerLabel: {
          color: "black",
        },
      }}
      {...props}
    />
  )
}
