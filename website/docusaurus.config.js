import { themes as prismThemes } from 'prism-react-renderer';

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Terramaid',
  tagline: 'A utility for generating Mermaid diagrams from Terraform configurations',
  favicon: 'img/favicon.ico',

  url: 'https://terramaid.dev',
  baseUrl: '/',

  organizationName: 'RoseSecurity',
  projectName: 'Terramaid',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          path: '../docs', // Path to the docs folder in the root directory
          routeBasePath: '/', // Serve the docs at the site's root URL
          editUrl:
            'https://github.com/RoseSecurity/Terramaid/edit/main/docs/',
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      image: 'img/terramaid.jpg',
      navbar: {
        title: 'Terramaid',
        logo: {
          alt: 'Terramaid Logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            to: '/', // Link to the docs homepage
            label: 'Docs',
            position: 'left',
          },
          {
            href: 'https://github.com/RoseSecurity/Terramaid',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Getting Started',
                to: '/docs/Getting_Started',
              },
              {
                label: 'GitHub Actions Integration',
                to: '/docs/GitHub_Actions_Integration',
              },
              {
                label: 'GitLab Pipeline Integration',
                to: '/docs/GitLab_Pipelines_Integration',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Stack Overflow',
                href: 'https://stackoverflow.com/questions/tagged/terramaid',
              },
              {
                label: 'Discussions',
                href: 'https://github.com/RoseSecurity/Terramaid/discussions',
              },
              {
                label: 'Open an Issue',
                href: 'https://github.com/RoseSecurity/Terramaid/issues/new/choose',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/RoseSecurity/Terramaid',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Terramaid`,
      },
      prism: {
        theme: prismThemes.github,
        darkTheme: prismThemes.dracula,
      },
    }),
};

export default config;

