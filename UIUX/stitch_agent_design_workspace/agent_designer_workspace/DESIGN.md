---
name: Agent Designer Workspace
colors:
  surface: '#f7f9fb'
  surface-dim: '#d8dadc'
  surface-bright: '#f7f9fb'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f2f4f6'
  surface-container: '#eceef0'
  surface-container-high: '#e6e8ea'
  surface-container-highest: '#e0e3e5'
  on-surface: '#191c1e'
  on-surface-variant: '#454652'
  inverse-surface: '#2d3133'
  inverse-on-surface: '#eff1f3'
  outline: '#757684'
  outline-variant: '#c6c5d4'
  surface-tint: '#4455b8'
  primary: '#243699'
  on-primary: '#ffffff'
  primary-container: '#3e4fb2'
  on-primary-container: '#c6ccff'
  inverse-primary: '#bbc3ff'
  secondary: '#505f76'
  on-secondary: '#ffffff'
  secondary-container: '#d0e1fb'
  on-secondary-container: '#54647a'
  tertiary: '#693200'
  on-tertiary: '#ffffff'
  tertiary-container: '#8c4500'
  on-tertiary-container: '#ffc39b'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#dee0ff'
  primary-fixed-dim: '#bbc3ff'
  on-primary-fixed: '#000f5d'
  on-primary-fixed-variant: '#2a3c9f'
  secondary-fixed: '#d3e4fe'
  secondary-fixed-dim: '#b7c8e1'
  on-secondary-fixed: '#0b1c30'
  on-secondary-fixed-variant: '#38485d'
  tertiary-fixed: '#ffdcc6'
  tertiary-fixed-dim: '#ffb785'
  on-tertiary-fixed: '#301400'
  on-tertiary-fixed-variant: '#713700'
  background: '#f7f9fb'
  on-background: '#191c1e'
  surface-variant: '#e0e3e5'
typography:
  display:
    fontFamily: Inter
    fontSize: 32px
    fontWeight: '600'
    lineHeight: '1.2'
    letterSpacing: -0.02em
  h1:
    fontFamily: Inter
    fontSize: 24px
    fontWeight: '600'
    lineHeight: '1.3'
  h2:
    fontFamily: Inter
    fontSize: 20px
    fontWeight: '500'
    lineHeight: '1.4'
  body-base:
    fontFamily: Inter
    fontSize: 16px
    fontWeight: '400'
    lineHeight: '1.6'
  body-sm:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '400'
    lineHeight: '1.5'
  label-caps:
    fontFamily: Inter
    fontSize: 12px
    fontWeight: '600'
    lineHeight: '1.2'
    letterSpacing: 0.05em
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  unit: 4px
  container-padding: 24px
  stack-gap: 16px
  inline-gap: 12px
  section-margin: 40px
---

## Brand & Style

The design system is engineered for deep focus and executive-level control within the AI orchestration space. It centers on a **Modern Corporate** aesthetic, heavily influenced by Microsoft’s Fluent Design principles. The goal is to reduce cognitive load by utilizing "Mica" and "Acrylic" layering techniques, creating a sense of physical depth without visual clutter.

The personality is authoritative yet approachable—moving away from the "hacker" aesthetic of early AI tools toward a refined, reliable management environment. The user should feel in total control of complex systems, supported by a UI that feels stable, quiet, and premium.

## Colors

The palette is anchored by **Deep Indigo (#3E4FB2)**, providing a scholarly and professional weight to primary actions. The neutral scale relies on a sequence of off-whites and cool grays to facilitate layering.

- **Primary:** Used for high-priority actions, active states, and focus indicators.
- **Surface Layering:** Transitions between `Slate-50` and `Slate-100` define boundaries, rather than hard strokes.
- **Status Tones:** Success, Warning, and Error states utilize high-luminance backgrounds with high-contrast text to ensure legibility while maintaining the calm, subdued atmosphere of the workspace.

## Typography

This design system utilizes **Inter** for its exceptional readability in data-dense environments. The hierarchy is strictly enforced through weight and scale rather than color. 

- **Headlines:** Reserved for workspace titles and agent names, using a tighter letter-spacing for a modern look.
- **Body:** Optimized for long-form prompt editing and log reading.
- **Labels:** Small-caps are used sparingly for metadata and table headers to provide clear distinction from interactive content.

## Layout & Spacing

The layout follows a **Fluid Grid** model with high-margin "Safe Areas." The workspace is divided into functional zones: Navigation (Left), Canvas (Center), and Inspector (Right).

A 4px base unit governs all dimensions. Elements are spaced generously to prevent the "dashboard fatigue" common in technical tools. Use 24px margins for the main container and 16px gaps for component stacks. Grid columns should be flexible, allowing the Canvas to expand during complex agent configuration tasks.

## Elevation & Depth

This design system uses **Tonal Layering** supplemented by **Mica** effects. Depth is communicated through surface color shifts rather than heavy drop shadows.

- **Level 0 (Background):** Solid off-white (#F8FAFC).
- **Level 1 (Cards/Panels):** Pure white background with a 1px `Slate-200` border at 50% opacity.
- **Level 2 (Floating/Modals):** Pure white with a 12% opacity ambient shadow (0px 8px 24px) to simulate elevation.
- **Acrylic Effect:** Sidebars and utility panels use a 20px backdrop blur with a 70% opacity white tint to maintain context of the underlying canvas.

## Shapes

The shape language is consistently **Rounded**, using an 8px radius for standard components and 12px for larger container panels. This soften's the professional "Slate" palette, making the workspace feel approachable. 

Interactive elements like buttons and input fields must never have sharp corners. Search bars and tags may utilize the **Pill-shaped** (rounded-full) variation to distinguish them from structural content containers.

## Components

- **Buttons:** Primary buttons use the Indigo fill with white text. Secondary actions use **Ghost Buttons** (no fill, slate-600 text) that reveal a subtle gray background only on hover.
- **Cards:** Cards should not have heavy borders. Use a slight background color change (white on a light gray page) and a soft 1px stroke at low opacity to define the edge.
- **Input Fields:** Use a subtle inset shadow or a slightly darker gray background (#F1F5F9) to indicate "emptiness" or "readiness" for input.
- **Agent Chips:** Small, rounded indicators for agent status. Use the subdued success/warning tones for the background and dark text for the label.
- **Nodes/Workflow Elements:** For the agent designer canvas, use rounded-xl containers with a high-blur backdrop to separate the logic flow from the grid background.
- **Checkboxes:** Square with a 4px radius, using the primary indigo for the checked state.