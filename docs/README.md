# GV Documentation Site

This site is powered by **[VitePress](https://vitepress.dev/)** and serves as the official documentation for the **GV** project.

## Prerequisites

- **Node.js** (v18+ recommended) or **[Bun](https://bun.sh/)** (v1+)
- A package manager: `npm`, `pnpm`, `yarn`, or `bun`

Check your environment:

```bash
node -v
npm -v
# or
bun -v
```

## Getting Started

### 1. Clone the Repository

```bash
git clone git@github.com:struckchure/gv.git
cd gv/docs
```

### 2. Install Dependencies

- Using npm:

  ```bash
  npm install
  ```

- Using bun:

  ```bash
  bun install
  ```

### 3. Run the Dev Server

- Using npm:

  ```bash
  npm run docs:dev
  ```

- Using bun:

  ```bash
  bun run docs:dev
  ```

Visit [http://localhost:5173](http://localhost:5173) to view the site.

### 4. Build for Production

- Using npm:

  ```bash
  npm run docs:build
  ```

- Using bun:

  ```bash
  bun run docs:build
  ```

The output will be generated in `.vitepress/dist`.

### 5. Preview the Production Build

- Using npm:

  ```bash
  npm run docs:preview
  ```

- Using bun:

  ```bash
  bun run docs:preview
  ```

## Project Structure

```
gv/
├── docs/
│   ├── .vitepress/
│   │   ├── config.ts
│   │   └── theme/
│   ├── index.md
│   ├── guide/
│   └── reference/
├── server/
├── client/
└── README.md
```

## Scripts

| Script | Description |
|:-------|:------------|
| `docs:dev` | Start local dev server |
| `docs:build` | Build static site for production |
| `docs:preview` | Preview built site locally |

Run using:

```bash
npm run <script-name>
# or
bun run <script-name>
```

## Useful Links

- [GV GitHub Repository](https://github.com/struckchure/gv)
- [VitePress Documentation](https://vitepress.dev/)
- [Bun Documentation](https://bun.sh/docs)

## Contribution

Contributions to the GV documentation are welcome.  
Fork the repository, make your changes, and submit a pull request.
