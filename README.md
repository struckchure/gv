[Discord](https://discord.gg/DzPC7D8T)

## 🚀 GV (Go Vite)

**GV** is a blazing-fast, browser-native development server and build tool written in Go. It’s inspired by [Vite](https://vitejs.dev) but designed to run without Node.js, using the power of native ESM in modern browsers and CDN-based dependency resolution.

> ⚡ Powered by Go. 🔌 Plugin-friendly. 🧠 Node-free.

---

### 💎 Philosophy

We should be able to use JavaScript frameworks, without 'node' or any runtime, just your browser.

---

### ✨ Features

- [x] Native ESM support in modern browsers
- [x] CDN-based module fetching (e.g., `esm.sh`, `skypack`)
- [ ] Local caching of remote modules
- [x] Hot Module Replacement (HMR)
- [x] Zero-config dev server
- [x] `esbuild`-based transpilation
- [x] `babel`-based transpilation
- [x] Plugin system (inspired by Vite/Rollup)
- [x] Written in Go with extensibility in mind
- [x] No Node.js required — ever
- [x] TypeScript support
- [ ] Adapter to support any http server

---

# Installation

For Linux, MacOS and Windows

```sh
curl -fsSL https://raw.githubusercontent.com/struckchure/gv/main/scripts/install.sh | bash
```

For Windows

```sh
irm https://raw.githubusercontent.com/struckchure/gv/main/scripts/install.ps1 | iex
```

---

### 🔧 Getting Started

```bash
git clone https://github.com/struckchure/gv
cd gv/examples/react
go run .
```

Then open your browser to [http://localhost:3000](http://localhost:3000).

---

### 📦 How It Works

- 📜 **Transpiling**: Uses `babel` internally for `.ts`, `.jsx`, `.tsx`, etc.
- 🌐 **CDN Resolution**: Bare imports (like `react`) are rewritten to point to `https://esm.sh/react` and cached locally.
- 🔥 **HMR**: WebSocket server pushes updates to the browser with minimal reloads.
- 🧩 **Plugins**: Extend GV with hooks like `transform`, `resolveId`, and `load`.

---

### 📁 Project Structure

Well, your project structure can be anyhow you want, but here's a sample react project

```
.
├── index.html
├── main.go
├── main.jsx
├── router.js
└── routes
    ├── layout.jsx
    ├── login
    │   └── page.jsx
    ├── page.jsx
    └── register
        └── page.jsx
```

Import from CDNs or local files directly:

```jsx
import { createRoot } from "https://esm.sh/react-dom@19.1.0/client";
import { RouterProvider } from "https://esm.sh/react-router@7.5.0";
import React from "https://esm.sh/react@19.1.0";

import { router } from "./router.js";

createRoot(document.getElementById("root")).render(
  <RouterProvider router={router} />
);
```

---

### 🔌 Plugin API

Check [here](./docs/writing-your-own-gv-plugin.md).

---

### 📜 License

MIT © 2025 [Mohammed Al-Ameen](mailto:ameenmohammed2311@gmail.com)
