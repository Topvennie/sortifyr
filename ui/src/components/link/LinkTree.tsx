import { Directory } from "@/lib/types/directory";
import { Side } from "@/lib/types/general";
import { Stack, StackProps } from "@mantine/core";
import { LoadingSpinner } from "../molecules/LoadingSpinner";
import { LinkTreeNode } from "./LinkTreeNode";
import { cn } from "@/lib/utils";
import { Playlist } from "@/lib/types/playlist";
import { LinkTreePlaylist } from "./LinkTreePlaylist";

type Props = {
  roots: Directory[];
  unAssigned: Playlist[]
  side: Side;
  title: string;
  isLoading: boolean;
} & StackProps

export const LinkTree = ({ roots, unAssigned, side, title, isLoading, className, ...props }: Props) => {
  const child = () => {
    if (isLoading) return <LoadingSpinner />

    return (
      <>
        <p className="text-muted">Directories</p>
        {roots.map(r => <LinkTreeNode key={r.id} directory={r} side={side} level={0} />)}
        <p className="text-muted">Unassigned Playlists</p>
        <Stack gap={2} style={{ marginLeft: 6 }}>
          {unAssigned.map(p => <LinkTreePlaylist key={p.id} playlist={p} side={side} />)}
        </Stack>
      </>
    )
  }

  return (
    <Stack p="md" pr={side === "left" ? "xl" : "md"} gap="md" justify="start" bg="background.0" className={cn("overflow-auto rounded-xl relative", className)} {...props}>
      <p className="text-muted">{title}</p>
      {child()}
    </Stack>
  )
}


