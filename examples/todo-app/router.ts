import { createBrowserRouter } from "react-router";

import RootLayout from "./pages/layout";
import Root, { loader as RootLoad } from "./pages/page";

export let router = createBrowserRouter([
  {
    Component: RootLayout,
    children: [{ index: true, Component: Root, loader: RootLoad }],
  },
]);
