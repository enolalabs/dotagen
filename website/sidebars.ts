import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docs: [
    'intro',
    {
      type: 'category',
      label: 'Getting Started',
      items: ['getting-started/installation', 'getting-started/initialize', 'getting-started/configuration'],
    },
    {
      type: 'category',
      label: 'Concepts',
      items: ['concepts/agents', 'concepts/skills', 'concepts/rendering-and-symlinks', 'concepts/built-in-catalog'],
    },
    {
      type: 'category',
      label: 'CLI Reference',
      items: ['cli/overview', 'cli/init-sync-status', 'cli/create-agents', 'cli/manage-skills', 'cli/serve-update-uninstall'],
    },
    {
      type: 'category',
      label: 'Platforms',
      items: ['platforms/overview', 'platforms/antigravity', 'platforms/claude-code', 'platforms/codex', 'platforms/cursor', 'platforms/gemini-cli', 'platforms/github-copilot', 'platforms/opencode', 'platforms/windsurf'],
    },
    {type: 'category', label: 'Dashboard', items: ['dashboard/overview', 'dashboard/rest-api']},
    {type: 'category', label: 'Reference', items: ['reference/config-schema', 'reference/troubleshooting', 'reference/releases-and-verification']},
  ],
};

export default sidebars;
