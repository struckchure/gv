import { createBrowserRouter } from "react-router";

import RootLayout from "./routes/layout";
import Login from "./routes/login/page";
import Root from "./routes/page";
import Register from "./routes/register/page";

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
