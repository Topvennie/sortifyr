import { createRootRouteWithContext, createRoute, createRouter } from "@tanstack/react-router";
import { App } from "./App";
import { Error404 } from "./pages/404";
import { Directories } from "./pages/Directories";
import { Error } from "./pages/Error";
import { Index } from "./pages/Index";
import { Playlists } from "./pages/Playlists";
import { Home } from "./pages/Home";
import { Links } from "./pages/Links";
import { Tasks } from "./pages/Tasks";
import { Tracks } from "./pages/Tracks";
import { Settings } from "./pages/Settings";
import { GeneratorOverview } from "./pages/generator/GeneratorOverview";
import { Generator } from "./pages/generator/Generator";
import { GeneratorCreate } from "./pages/generator/GeneratorCreate";
import { GeneratorEdit } from "./pages/generator/GeneratorEdit";

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

const link = createRoute({
  getParentRoute: () => index,
  path: "/link",
  component: Links,
})

const task = createRoute({
  getParentRoute: () => index,
  path: "/task",
  component: Tasks,
})

const history = createRoute({
  getParentRoute: () => index,
  path: "/history",
  component: Tracks,
})

const setting = createRoute({
  getParentRoute: () => index,
  path: "/setting",
  component: Settings,
})

//
// Generator
//

const generator = createRoute({
  getParentRoute: () => index,
  path: "/generator",
  component: Generator,
})

const generatorOverview = createRoute({
  getParentRoute: () => generator,
  path: "/",
  component: GeneratorOverview,
})

const generatorCreate = createRoute({
  getParentRoute: () => generator,
  path: "/create",
  component: GeneratorCreate,
})

const generatorEdit = createRoute({
  getParentRoute: () => generator,
  path: "/edit/$generatorId",
  component: GeneratorEdit,
})

const routeTree = root.addChildren([
  index.addChildren([
    home,
    playlist,
    directory,
    link,
    task,
    history,
    setting,
    generator.addChildren(([generatorOverview, generatorCreate, generatorEdit])),
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
