import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

const config: Config = {
  title: "Awesome CI",
  customFields: {
    openerTitle: "FS DevOps peresents: Awesome CI",
  },
  tagline: "Fullstack applications and DevOps solutions",
  favicon: "https://avatars.githubusercontent.com/u/97617148?s=200&v=4",
  url: "https://fullstack-devops.github.io",
  baseUrl: "/awesome-ci",
  organizationName: "fullstack-devops",
  projectName: "awesome-ci",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          editUrl: "https://github.com/fullstack-devops/awesome-ci/tree/main/",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      },
    ],
  ],

  plugins: [require.resolve("docusaurus-lunr-search")],

  themeConfig: {
    image: "https://fullstack-devops.github.io/img/full-logo.png",
    navbar: {
      title: "Awesome CI",
      logo: {
        alt: "Awesome CI Logo",
        src: "https://fullstack-devops.github.io/img/logo.png",
      },
      items: [
        {
          type: "doc",
          docId: "overview",
          position: "left",
          label: "Overview",
        },
        {
          type: "search",
          position: "right",
        },
        {
          href: "https://github.com/fullstack-devops/awesome-ci",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Docs",
          items: [
            {
              label: "Overview",
              to: "/docs/overview",
            },
          ],
        },
        {
          title: "Community",
          items: [
            {
              label: "Stack Overflow",
              href: "https://stackoverflow.com/questions/tagged/fs-devops",
            },
          ],
        },
        {
          title: "More",
          items: [
            {
              label: "GitHub",
              href: "https://github.com/fullstack-devops/awesome-ci",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Fs DevOps. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
