import { LinkButton } from "@/components/atoms/LinkButton";
import { useAuth } from "@/lib/hooks/useAuth";
import { AppShell, Burger, Button, Container, Group, Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { LinkProps } from "@tanstack/react-router";
import { ComponentProps } from "react";
import { FaArrowRightToBracket, FaMusic } from "react-icons/fa6";

type Route = {
  title: string;
  link: LinkProps;
};

const routes: Route[] = [
  {
    title: "Playlists",
    link: { to: "/playlist" },
  },
  {
    title: "Directories",
    link: { to: "/directory" },
  },
];

const NavLink = ({ route, closeNavbar }: { route: Route; closeNavbar?: () => void }) => {
  return (
    <div onClick={closeNavbar} className="w-fit">
      <LinkButton
        to={route.link.to}
        activeProps={{ variant: "filled", c: "white" }}
        variant="subtle"
        size="md"
        c="green.9"
        tt="uppercase"
        radius="md"
        className="font-black tracking-wide"
      >
        {route.title}
      </LinkButton>
    </div>
  );
};

export const NavLayout = ({ children, ...props }: ComponentProps<"div">) => {
  const [opened, { close, toggle }] = useDisclosure();
  const { user, logout } = useAuth()

  return (
    <AppShell
      header={{ height: 80 }}
      navbar={{ width: 300, breakpoint: "md", collapsed: { desktop: true, mobile: !opened } }}
      padding="md"
    >
      <AppShell.Header>
        <Group justify="space-between" h="100%" px="md" wrap="nowrap" align="center">
          <div className="flex items-center w-full">
            <LinkButton to="/" variant="transparent" className="flex-1 flex items-center">
              <div className="flex gap-4 items-center">
                <FaMusic className="size-8" />
                <p className="font-bold text-xl hidden md:block">Organizer</p>
              </div>
            </LinkButton>
            <div className="flex-1 justify-center hidden md:flex">
              <div className="flex items-center gap-4 bg-green-50 rounded-xl p-2">
                {routes.map((route) => (
                  <NavLink key={route.title} route={route} />
                ))}
              </div>
            </div>
            <div className="flex-1 justify-end hidden md:flex">
              <Button variant="subtle" onClick={logout} tt="uppercase" size="md" rightSection={<FaArrowRightToBracket />} className="font-bold tracking-wide">
                <div className="flex gap-2 items-center">
                  {user?.name}
                </div>
              </Button>
            </div>
          </div>
          <Burger opened={opened} onClick={toggle} hiddenFrom="md" size="sm" />
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        <Stack align="flex-start">
          {routes.map((route) => (
            <NavLink key={route.title} route={route} closeNavbar={close} />
          ))}
        </Stack>
      </AppShell.Navbar>

      <AppShell.Main className="bg-green-50 overflow-hidden">
        <Container fluid className="pt-10 container mx-auto" {...props}>
          {children}
        </Container>
      </AppShell.Main>
    </AppShell>
  );
}
