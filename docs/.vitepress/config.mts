import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "GV",
  description: "Golang Build tool for Modern Web Frameworks",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    sidebar: [
      {
        text: "Introduction",
        items: [
          { text: "What is GV?", link: "/what-is-gv" },
          { text: "Getting Started", link: "/getting-started" },
        ],
      },
      {
        text: "Plugins",
        link: "/plugins/",
        collapsed: true,
        items: [
          { text: "HMR", link: "/plugins/hmr" },
          { text: "HTML", link: "/plugins/html" },
          { text: "CDN Dependecy", link: "/plugins/cdn-dependency" },
          {
            text: "Frameworks",
            collapsed: true,
            items: [
              { text: "React", link: "/plugins/react" },
              { text: "Vue", link: "/plugins/vue" },
              { text: "Svelte", link: "/plugins/svelte" },
            ],
          },
          { text: "Custom Plugin", link: "/plugins/custom-plugin" },
        ],
      },
    ],

    socialLinks: [
      { icon: "github", link: "https://github.com/struckchure/gv" },
    ],
  },
});
