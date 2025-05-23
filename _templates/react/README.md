# React + GV Template

A modern frontend application powered by React and built with [GV](https://github.com/struckchure/gv) — a fast Vite-like tool written in Go.

---

## 🛠️ Development

Start the development server with hot-reloading:

```sh
gv sync -c config.yaml
GV_MODE=dev go run .
```

---

## 🏗️ Build

Generate a static production build of your frontend:

```sh
GV_MODE=build go run .
```

> ⚠️ Make sure the `dist/` folder is present in the project root before building the final binary.  
> If you're building outside the original directory, copy the `dist` folder manually:

```sh
cp -r /path/to/project/dist /path/to/build-context/
```

Then build the Go binary:

```sh
go build -o main .
```

---

## 🚀 Production

### Using Docker

Build and run a Docker container:

```sh
docker build -t my-app .
docker run -p 3000:3000 my-app
```

---

## 📁 Folder Structure

```
├── dist/            # Built frontend assets (generated by GV)
├── src/             # React source code
├── config.yaml      # GV configuration
├── main.go          # Main Go entrypoint
├── go.mod           # Go modules
└── Dockerfile       # Optional Docker setup
```

---

## 📦 Dependencies

- React
- GV (Go Vite) - [https://github.com/struckchure/gv](https://github.com/struckchure/gv)
- Echo Web Server
