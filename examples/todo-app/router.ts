import { createBrowserRouter } from "react-router";

import RootLayout from "./pages/layout.js";
import Root, { loader as RootLoad } from "./pages/page.js";

export let router = createBrowserRouter([
  {
    Component: RootLayout,
    children: [{ index: true, Component: Root, loader: RootLoad }],
  },
]);
