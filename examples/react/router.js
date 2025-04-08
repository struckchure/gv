import { createBrowserRouter } from "https://esm.sh/react-router@7.5.0";

import RootLayout from "./routes/layout.js";
import Login from "./routes/login/page.js";
import Root from "./routes/page.js";
import Register from "./routes/register/page.js";

export let router = createBrowserRouter([
  {
    Component: RootLayout,
    children: [
      { index: true, Component: Root },
      { path: "/login", Component: Login },
      { path: "/register", Component: Register },
    ],
  },
]);
