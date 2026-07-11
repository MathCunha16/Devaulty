import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createRouter, RouterProvider, createMemoryHistory } from "@tanstack/react-router";
import "./index.css";

import { routeTree } from "./routeTree.gen";

// Setup memory history for WebView / Desktop app context
const memoryHistory = createMemoryHistory({
  initialEntries: ["/"],
});

// Create a new router instance
const router = createRouter({
  routeTree,
  history: memoryHistory,
  defaultPreload: "intent",
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

// Create React Query client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5000,
    },
  },
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>
);
