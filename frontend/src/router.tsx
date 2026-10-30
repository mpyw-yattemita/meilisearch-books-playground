import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, createRoute, createRouter } from "@tanstack/react-router";
import { RootLayout } from "./routes/root";
import { BookPage, bookQueryOptions } from "./routes/book";
import { SearchPage, validateSearch } from "./routes/search";

export type RouterContext = { queryClient: QueryClient };

const rootRoute = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
});

const searchRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  validateSearch,
  component: SearchPage,
});

const bookRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/books/$id",
  params: {
    parse: ({ id }) => ({ id: Number(id) }),
    stringify: ({ id }) => ({ id: String(id) }),
  },
  loader: ({ context, params }) => context.queryClient.ensureQueryData(bookQueryOptions(params.id)),
  component: BookPage,
});

const routeTree = rootRoute.addChildren([searchRoute, bookRoute]);

export function createAppRouter(queryClient: QueryClient) {
  return createRouter({ routeTree, context: { queryClient }, defaultPreload: "intent" });
}

export { bookRoute, searchRoute };

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof createAppRouter>;
  }
}
