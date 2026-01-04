import { Text, Group, ModalBaseProps, Stack } from "@mantine/core";
import { ModalCenter } from "../atoms/ModalCenter";
import { Button } from "../atoms/Button";

interface Props extends Omit<ModalBaseProps, 'title'> {
  modalTitle: string;
  title: string;
  description?: string;
  onConfirm: () => void;
  loading?: boolean;
}

export const Confirm = ({ modalTitle, title, description, onConfirm, loading = false, ...props }: Props) => {
  return (
    <ModalCenter title={modalTitle} {...props}>
      <Stack>
        <Text fw="bold">{title}</Text>
        {description && <Text className="whitespace-pre-wrap">{description}</Text>}
        <Group justify="end">
          <Button onClick={props.onClose} variant="default">
            Cancel
          </Button>
          <Button onClick={onConfirm} loading={loading}>
            Confirm
          </Button>
        </Group>
      </Stack>
    </ModalCenter>
  )
}
