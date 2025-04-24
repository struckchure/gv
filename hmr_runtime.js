const HMR_RELOAD_TYPE = "reload";
const HMR_UPDATE_TYPE = "update";

if (!window.__hmr__) {
  window.__hmr__ = {
    contexts: {},
  };

  const socketURL = new URL(
    "/__hmr__",
    window.location.href.replace(/^http(s)?:/, "ws$1:")
  );
  const socket = (window.__hmr__.socket = new WebSocket(socketURL.href));

  socket.onopen = () => {
    console.log("[HMR] Client connected");
  };
  socket.onmessage = async (event) => {
    const payload = JSON.parse(event.data);

    switch (payload?.type) {
      case HMR_RELOAD_TYPE:
        window.location.reload();
        break;
      case HMR_UPDATE_TYPE:
        if (!payload.updates?.length) return;

        let accepted = false;
        for (const update of payload.updates) {
          console.log("[HMR]", update.id);
          if (window.__hmr__.contexts[update.id]) {
            const modUrl =
              "./" +
              update.url.replaceAll(update.url.split(".").at(-1), "js") +
              "?t=" +
              Date.now();
            accepted = window.__hmr__.contexts[update.id].emit(
              await import(modUrl)
            );
          }

          if (accepted) {
            console.log("[HMR] Updated accepted by", update.id);
          } else {
            console.log("[HMR] Updated rejected, reloading...");
            window.location.reload();
          }
        }

        break;
    }
  };
}

function createHotContext(id) {
  let callback;
  let disposed = false;

  const hot = {
    accept: (cb) => {
      if (disposed) {
        throw new Error("import.meta.hot.accept() called after dispose()");
      }
      if (callback) {
        throw new Error("import.meta.hot.accept() already called");
      }
      callback = cb;
    },
    dispose: () => {
      disposed = true;
      callback = undefined;
    },
    emit(self) {
      if (disposed) {
        throw new Error("import.meta.hot.emit() called after dispose()");
      }

      if (callback) {
        callback(self);
        return true;
      }
      return false;
    },
  };

  window.__hmr__.contexts[id] = hot;

  return hot;
}
