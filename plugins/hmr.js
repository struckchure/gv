async function jsHandler(file) {
  console.log(`[HMR] Reloading JS module: ${file}`);
  try {
    await import(`${file}?t=${Date.now()}`);
    await import(`/dist/main.js?t=${Date.now()}`);
  } catch (err) {
    console.error(`[HMR] Failed to reload JS: ${file}`, err);
  }
}

function cssHandler(file) {
  console.log(`[HMR] Updating CSS: ${file}`);
  const links = [...document.querySelectorAll('link[rel="stylesheet"]')];
  const match = links.find((link) => link.href.includes(file));
  if (match) {
    const newLink = match.cloneNode();
    newLink.href = `${file}?t=${Date.now()}`;
    newLink.onload = () => match.remove();
    match.parentNode.insertBefore(newLink, match.nextSibling);
  }
}

async function jsonHandler(file) {
  console.log(`[HMR] Reloading JSON: ${file}`);
  try {
    const res = await fetch(`${file}?t=${Date.now()}`);
    const data = await res.json();
    console.log(`[HMR] New JSON data:`, data);
    // Optional: Notify components, re-render UI, etc.
  } catch (err) {
    console.error(`[HMR] Failed to reload JSON: ${file}`, err);
  }
}

function assetHandler(file) {
  console.log(`[HMR] Reloading asset: ${file}`);
  const fileName = file.split("/").pop();
  const selectors = [
    `img[src*="${fileName}"]`,
    `source[src*="${fileName}"]`,
    `link[href*="${fileName}"]`, // e.g. fonts, icons
    `video[src*="${fileName}"]`,
    `audio[src*="${fileName}"]`,
  ];
  const elements = document.querySelectorAll(selectors.join(","));
  elements.forEach((el) => {
    const attr = el.tagName === "LINK" ? "href" : "src";
    el[attr] = `${file}?t=${Date.now()}`;
  });
}

async function init() {
  const ws = new WebSocket("/_/ws/");
  const handlers = {
    js: jsHandler,
    css: cssHandler,
    json: jsonHandler,
  };

  ws.onmessage = async (e) => {
    const payload = JSON.parse(e.data);
    if (payload.type === "ping") {
      ws.send(JSON.stringify({ type: "pong" }));
    }

    const [ext, type] = payload.type.split(":", 2);
    if (type !== "update") return;

    try {
      const handler = handlers[ext] || assetHandler;
      handler(payload.file);
    } catch (err) {
      console.warn(err);
    }
  };
}

window.addEventListener("DOMContentLoaded", init);
