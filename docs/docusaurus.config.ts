import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'template-go',
  tagline: 'Meigma Go repository template',
  future: {
    v4: true,
  },
  url: 'https://docs.example.com',
  baseUrl: '/',
  organizationName: 'meigma',
  projectName: 'template-go',
  onBrokenLinks: 'throw',
  markdown: {
    hooks: {
      onBrokenMarkdownLinks: 'warn',
    },
  },
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },
  presets: [
    [
      'classic',
      {
        docs: {
          path: 'docs',
          routeBasePath: '/',
          sidebarPath: false,
          breadcrumbs: false,
          editUrl: 'https://github.com/meigma/template-go/edit/master/docs/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],
  themeConfig: {
    colorMode: {
      defaultMode: 'dark',
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'template-go',
      items: [
        {
          href: 'https://github.com/meigma/template-go',
          label: 'GitHub',
          position: 'right',
          className: 'navbar__item--github',
        },
      ],
    },
    footer: {
      style: 'dark',
      copyright: `Copyright © ${new Date().getFullYear()} Meigma. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'json', 'toml', 'yaml'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
