import { useAutoAnimate } from '@formkit/auto-animate/react';
import { DataTable, DataTableProps } from "mantine-datatable";
import { LoadingSpinner } from "./LoadingSpinner";

type Props<T> = {
  animate?: boolean;
} & DataTableProps<T> & Record<string, unknown>;

export const Table = <T,>({ animate = true, rowExpansion, ...props }: Props<T>) => {
  const [bodyRef] = useAutoAnimate<HTMLTableSectionElement>();

  return (
    <DataTable
      customLoader={<LoadingSpinner />}
      minHeight={180}
      backgroundColor="background.0"
      withTableBorder={false}
      textSelectionDisabled={true}
      rowExpansion={rowExpansion}
      styles={{
        root: (theme) => ({
          borderRadius: theme.radius.md,
        }),
        header: (theme) => ({
          background: theme.colors.secondary[1],
        }),
      }}
      bodyRef={(!rowExpansion && animate) ? bodyRef : undefined}
      {...props}
    />
  );
};
