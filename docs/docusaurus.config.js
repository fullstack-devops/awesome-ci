// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/vsLight");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').Config} */
const config = {
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
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          editUrl: "https://github.com/fullstack-devops/awesome-ci/tree/main/",
        },
        blog: {
          showReadingTime: true,
          editUrl: "https://github.com/fullstack-devops/awesome-ci/tree/main/",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],

  plugins: [require.resolve("docusaurus-lunr-search")],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
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
            docId: "intro",
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
                to: "/docs/intro",
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
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
