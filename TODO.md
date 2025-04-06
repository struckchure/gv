## ✅ TODO: Build Vite-like Tool from Scratch (in Go, No Node)

### 🧠 Core Concepts to Understand First

- [ ] Study **native ESM** in browsers
- [ ] Learn how **esbuild** and **Rollup** work internally
- [ ] Understand how HMR (Hot Module Replacement) works
- [ ] Review dev server patterns (like `koa`, `express`, or `connect` middlewares)
- [ ] Explore how **Deno**, **WMR**, or **Vite CDN** resolve bare imports

---

### ⚙️ Phase 1: Development Server

#### 📦 Module Resolution

- [ ] Implement basic ESM resolver for `.js`, `.ts`, `.jsx`, `.tsx`, `.vue`, etc.
- [ ] Fetch bare imports from a CDN (e.g., `esm.sh`, `skypack.dev`, `unpkg.com`)
- [ ] Cache CDN modules locally on first request
- [ ] Rewrite import paths in served files to point to CDN/local cache
- [ ] Add aliasing support (like `@/components`)

#### 🔄 File Watching & HMR

- [ ] Watch for file changes (e.g., using `fsnotify`)
- [ ] Track dependency graph per module
- [ ] Push updates via WebSocket to the browser
- [ ] Create a lightweight HMR runtime client for browser injection

#### 🧪 Transpiling

- [ ] Integrate `esbuild` for fast on-the-fly transpilation (e.g., `.ts`, `.jsx`)
- [ ] Use `esbuild` or custom plugins to support `.vue`, `.svelte`, etc.
- [ ] Cache transpiled modules to avoid repeated work

#### 🌐 Dev Server (in Go)

- [ ] Create a generic `Server` interface to support different Go HTTP backends
- [ ] Implement adapters for `net/http`, `fasthttp`, etc.
- [ ] Serve static files (HTML, images, styles)
- [ ] Serve ESM modules and inject rewritten imports
- [ ] Inject HMR runtime when serving HTML or JS

---

### 📦 Phase 2: Build Process (Production)

#### 🏗️ Bundler

- [ ] Use **Rollup** or build a custom bundler in Go
- [ ] Setup plugin hooks for transforms and optimization
- [ ] Support:
  - Code splitting
  - Tree-shaking
  - CSS extraction
  - Asset optimization

#### ⚙️ Build Options

- [ ] Implement `vite.config.ts`-like config file (e.g., `build.config.json` or `.toml`)
- [ ] Support custom output directory, base path, etc.
- [ ] Handle different environments (production/dev)

---

### 🔌 Plugin System

- [ ] Define a plugin interface (inspired by Vite/Rollup)
- [ ] Support hooks like:
  - `transform`
  - `resolveId`
  - `load`
  - `handleHotUpdate`
- [ ] Allow build-time plugins and dev-server plugins

---

### 📁 Static Assets & CSS

- [ ] Inline small assets as base64
- [ ] Emit larger assets with hashed filenames
- [ ] Support PostCSS and CSS modules
- [ ] Handle `import './styles.css'` in JS/TS files

---

### 🌍 Environment & Config

- [ ] Support `.env` file parsing
- [ ] Expose only `VITE_`-prefixed variables to the browser
- [ ] Allow `.env.{mode}` variations for environment-specific config

---

### 🧪 Testing & DX

- [ ] Add CLI (`dev`, `build`, `preview`) using Go packages like `urfave/cli`
- [ ] Build clear logging and error overlays
- [ ] Show good stack traces and file/line info in dev
- [ ] Support source maps for better debugging

---

### 💡 Bonus Features (Optional)

- [ ] SSR (Server Side Rendering) support
- [ ] Type checking with `tsc` or custom Go-based parser
- [ ] Static site generation (SSG) hooks
- [ ] File-based routing (like `pages/index.tsx`)
- [ ] Plugin ecosystem/registry
