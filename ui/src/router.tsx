import { createRootRouteWithContext, createRoute, createRouter } from "@tanstack/react-router";
import { App } from "./App";
import { Error404 } from "./pages/404";
import { Directories } from "./pages/Directory";
import { Error } from "./pages/Error";
import { Index } from "./pages/Index";
import { Playlists } from "./pages/Playlist";
import { DirectoryEditor } from "./pages/DirectoryEditor";
import { Home } from "./pages/Home";

const root = createRootRouteWithContext()({
  component: App,
})

const index = createRoute({
  getParentRoute: () => root,
  id: "public-layout",
  component: Index,
})

const home = createRoute({
  getParentRoute: () => index,
  path: "/",
  component: Home,
})

const playlist = createRoute({
  getParentRoute: () => index,
  path: "/playlist",
  component: Playlists,
})

const directory = createRoute({
  getParentRoute: () => index,
  path: "/directory",
  component: Directories,
})

const directoryEditor = createRoute({
  getParentRoute: () => index,
  path: "/directory_edit",
  component: DirectoryEditor,
})

const routeTree = root.addChildren([
  index.addChildren([
    home,
    playlist,
    directory,
    directoryEditor,
  ]),
])

export const router = createRouter({
  routeTree,
  defaultPreload: "render",
  defaultPreloadStaleTime: 0, // Data is immediatly marked as stale and will refetch when the user navigates to the page
  scrollRestoration: true,
  defaultErrorComponent: Error,
  defaultNotFoundComponent: Error404,
})

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
