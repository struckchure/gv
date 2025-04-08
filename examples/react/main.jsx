import { createRoot } from "https://esm.sh/react-dom@19.1.0/client";
import { RouterProvider } from "https://esm.sh/react-router@7.5.0";
import React from "https://esm.sh/react@19.1.0";

import { router } from "./router.js";

createRoot(document.getElementById("root")).render(
  <RouterProvider router={router} />
);
