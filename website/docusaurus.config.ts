import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'dotagen',
  tagline: 'Define sub-agents and skills once, inject everywhere.',
  favicon: 'img/logo.svg',
  url: 'https://dotagen.enolalab.com',
  baseUrl: '/',
  organizationName: 'enolalabs',
  projectName: 'dotagen',
  onBrokenLinks: 'throw',
  presets: [
    [
      'classic',
      {
        docs: {
          routeBasePath: '/',
          editUrl: 'https://github.com/enolalabs/dotagen/edit/main/website/',
          showLastUpdateTime: true,
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],
  themeConfig: {
    image: 'img/logo.svg',
    navbar: {
      title: 'dotagen',
      logo: {
        alt: 'dotagen logo',
        src: 'img/logo.svg',
      },
      items: [
        {label: 'Enolalab', href: 'https://enolalab.com'},
        {label: 'GitHub', href: 'https://github.com/enolalabs/dotagen'},
        {label: 'Releases', href: 'https://github.com/enolalabs/dotagen/releases'},
        {
          label: 'Catalogue',
          href: 'https://github.com/enolalabs/dotagen/blob/main/docs/CATALOG.vi.md',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Project',
          items: [
            {label: 'Enolalab', href: 'https://enolalab.com'},
            {label: 'GitHub', href: 'https://github.com/enolalabs/dotagen'},
            {label: 'Releases', href: 'https://github.com/enolalabs/dotagen/releases'},
            {
              label: 'Catalogue',
              href: 'https://github.com/enolalabs/dotagen/blob/main/docs/CATALOG.vi.md',
            },
          ],
        },
      ],
      copyright: `Copyright ${new Date().getFullYear()} Enolalab.`,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
