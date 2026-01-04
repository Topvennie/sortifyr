import { Text, Group, ModalBaseProps, Stack } from "@mantine/core";
import { ModalCenter } from "../atoms/ModalCenter";
import { Button } from "../atoms/Button";
import { ReactNode } from "react";

interface Props extends Omit<ModalBaseProps, 'title' | 'content'> {
  modalTitle: string;
  title: string;
  description?: string;
  content?: ReactNode;
  onConfirm: () => void;
  loading?: boolean;
}

export const Confirm = ({ modalTitle, title, description, content, onConfirm, loading = false, ...props }: Props) => {
  return (
    <ModalCenter title={modalTitle} {...props}>
      <Stack>
        <Text fw="bold">{title}</Text>
        {description && <Text className="whitespace-pre-wrap">{description}</Text>}
        {content}
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
