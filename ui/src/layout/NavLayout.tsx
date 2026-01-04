import { Avatar } from "@/components/atoms/Avatar";
import { Button } from "@/components/atoms/Button";
import { LinkButton } from "@/components/atoms/LinkButton";
import { useAuth } from "@/lib/hooks/useAuth";
import { cn, getBuildTime } from "@/lib/utils";
import { ActionIcon, AppShell, Burger, Divider, Group, ScrollArea, Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { LinkProps } from "@tanstack/react-router";
import { ComponentProps, ReactNode, useState } from "react";
import { LuClock, LuFolderTree, LuLink2, LuListMusic, LuMusic, LuMusic3, LuSettings, LuSlidersHorizontal, LuTriangle } from "react-icons/lu";

type Props = ComponentProps<"div">

type Route = {
  title: string;
  icon: ReactNode;
  link: LinkProps;
};

const routes: Route[] = [
  {
    title: "Tracks",
    icon: <LuMusic3 className="size-5" />,
    link: { to: "/history" },
  },
  {
    title: "Playlists",
    icon: <LuListMusic className="size-5" />,
    link: { to: "/playlist" },
  },
  {
    title: "Directories",
    icon: <LuFolderTree className="size-5" />,
    link: { to: "/directory" },
  },
  {
    title: "Links",
    icon: <LuLink2 className="size-5" />,
    link: { to: "/link" },
  },
  {
    title: "Generator",
    icon: <LuSlidersHorizontal className="size-5" />,
    link: { to: "/generator" },
  },
  {
    title: "Background Tasks",
    icon: <LuClock className="size-5" />,
    link: { to: "/task" },
  },
  {
    title: "Settings",
    icon: <LuSettings className="size-5" />,
    link: { to: "/setting" },
  },
];

const NavLink = ({ route, close }: { route: Route, close?: () => void }) => {
  return (
    <div onClick={close}>
      <LinkButton
        to={route.link.to}
        activeProps={{ variant: "filled", bg: "primary.1" }}
        variant="subtle"
        size="md"
        radius="md"
        c="black"
        fullWidth
        justify="start"
        leftSection={route.icon}
      >
        {route.title}
      </LinkButton>
    </div>
  );
};

export const NavLayout = ({ className, children, ...props }: Props) => {
  const [opened, { close, toggle }] = useDisclosure();
  const { user, logout } = useAuth()

  const [userExpanded, setUserExpanded] = useState(false)

  const buildTime = getBuildTime()

  return (
    <AppShell
      header={{ height: { base: 60, lg: 0 } }}
      navbar={{ width: 300, breakpoint: "lg", collapsed: { mobile: !opened } }}
      padding="md"
    >
      <AppShell.Header p="md" hiddenFrom="lg" className="bg-[#eef6f4]">
        <Group justify="space-between">
          <LuMusic color="black" className="size-8 stroke-[#3b4a49]" />
          <Burger opened={opened} onClick={toggle} />
        </Group>
      </AppShell.Header>
      <AppShell.Navbar p="lg" className="bg-[#eef6f4]">
        <AppShell.Section p="md" h="92px" className="flex items-center">
          <Group gap="md">
            <LuMusic color="black" className="size-8 stroke-[#3b4a49]" />
            <p className="font-bold text-2xl text-[#3b4a49]">Sortifyr</p>
          </Group>
        </AppShell.Section>
        <AppShell.Section grow my="md" component={ScrollArea} px="md">
          <Stack p="sm" gap="xs" className="rounded-xl bg-white">
            {routes.map(r => <NavLink key={r.title} route={r} close={close} />)}
          </Stack>
        </AppShell.Section>
        <AppShell.Section p="md">
          <Stack p="sm" gap="xs" className="rounded-xl bg-white">
            <Group>
              <Avatar user={user} />
              <p className="font-bold">{user?.name}</p>
              <ActionIcon onClick={() => setUserExpanded(prev => !prev)} variant="light" color="black" size="sm" className="ml-auto">
                <LuTriangle className={`fill-black w-2 ${userExpanded && "rotate-180"}`} />
              </ActionIcon>
            </Group>
            {userExpanded && (
              <>
                <Divider mt="sm" />
                <Button onClick={logout} c="" variant="subtle" pl={0} justify="start" className="text-muted">
                  Log out
                </Button>
                <Divider />
                <p className="text-sm text-muted">{`Built: ${buildTime}`}</p>
              </>
            )}
          </Stack>
        </AppShell.Section>
      </AppShell.Navbar>
      <AppShell.Main py={{ lg: 0 }} bg="background.0" className={cn("h-screen overflow-auto border border-red-500", className)} {...props}>
        {children}
      </AppShell.Main>
    </AppShell>
  )
}
