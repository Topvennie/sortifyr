import { DatePickerInputProps, DatePickerInput as MDatePickerInput, DatePickerType } from "@mantine/dates";
import { LuCalendar } from "react-icons/lu";

type Props<T extends DatePickerType> = DatePickerInputProps<T>;

export const DatePickerInput = <T extends DatePickerType>(props: Props<T>) => {
  return (
    <MDatePickerInput
      leftSection={<LuCalendar />}
      radius="lg"
      styles={{
        input: {
          background: "var(--mantine-color-background-0)",
          borderColor: "var(--mantine-color-background-0)",
        },
      }}
      className="grow"
      {...props}
    />
  )
}
