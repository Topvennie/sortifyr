import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { Outlet, Scripts } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";

// <Scripts /> executes any relevant scripts defined in the router

export const App = () => {
  return (
    <>
      <Scripts />
      <Outlet />
      <TanStackRouterDevtools />
      <ReactQueryDevtools />
    </>
  )
}
