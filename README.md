## 🚀 GV (Go Vite)

**GV** is a blazing-fast, browser-native development server and build tool written in Go. It’s inspired by [Vite](https://vitejs.dev) but designed to run without Node.js, using the power of native ESM in modern browsers and CDN-based dependency resolution.

> ⚡ Powered by Go. 🔌 Plugin-friendly. 🧠 Node-free.

---

### ✨ Features

- ✅ Native ESM support in modern browsers
- ✅ CDN-based module fetching (e.g., `esm.sh`, `skypack`)
- ✅ Local caching of remote modules
- ✅ Hot Module Replacement (HMR)
- ✅ Zero-config dev server
- ✅ `esbuild`-based fast transpilation
- ✅ Plugin system (inspired by Vite/Rollup)
- ✅ Written in Go with extensibility in mind
- ✅ No Node.js required — ever

---

### 🔧 Getting Started

```bash
git clone https://github.com/yourusername/gv
cd gv
go run ./cmd/gv dev
```

Then open your browser to [http://localhost:3000](http://localhost:3000) and enjoy instant updates.

---

### 📦 How It Works

- 📜 **Transpiling**: Uses `esbuild` internally for `.ts`, `.jsx`, `.tsx`, etc.
- 🌐 **CDN Resolution**: Bare imports (like `react`) are rewritten to point to `https://esm.sh/react` and cached locally.
- 🔥 **HMR**: WebSocket server pushes updates to the browser with minimal reloads.
- 🧩 **Plugins**: Extend GV with hooks like `transform`, `resolveId`, and `load`.

---

### 📁 Example Project Structure

```
my-app/
├── index.html
├── main.ts
└── components/
    └── Hello.tsx
```

Import from CDNs or local files directly:

```ts
import React from "react";
import Hello from "./components/Hello.tsx";
```

---

### 🔌 Plugin API

Coming soon!

---

### 📜 License

MIT © 2025 [Mohammed Al-Ameen](mailto:ameenmohammed2311@gmail.com)
