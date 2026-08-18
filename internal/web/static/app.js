/* ═══════════════════════════════════════════════
   dotagen — Control Plane app.js
   Technical Precision design system
   ═══════════════════════════════════════════════ */

// ── State ──
let agents = [];
let skills = [];
let config = { targets: [], agents: {}, skills: {} };
let knownTargets = [];
let statusLinks = [];
let currentView = 'overview';
let searchQuery = '';
let categoryFilter = 'all';
let skillCategoryFilter = 'all';
let skillVendorFilter = 'all';
let statusFilter = 'all';
let libraryView = 'matrix';
let drawerMode = null;
let drawerAgent = null;
let selectedAgents = new Set();
let selectedSkills = new Set();
let skillPage = 1;
let agentPage = 1;
const PAGE_SIZE = 50;

const PLATFORM_LABELS = {
    'antigravity': 'AG',
    'claude-code': 'CC',
    'codex': 'CX',
    'cursor': 'CR',
    'gemini-cli': 'GC',
    'github-copilot': 'CP',
    'opencode': 'OC',
    'windsurf': 'WS',
};

const PLATFORM_NAMES = {
    'antigravity': 'Antigravity',
    'claude-code': 'Claude Code',
    'codex': 'Codex',
    'cursor': 'Cursor',
    'gemini-cli': 'Gemini CLI',
    'github-copilot': 'GitHub Copilot',
    'opencode': 'OpenCode',
    'windsurf': 'Windsurf',
};

const PLATFORM_ICONS = {
    'antigravity': 'rocket_launch',
    'claude-code': 'psychology',
    'codex': 'code',
    'cursor': 'ads_click',
    'gemini-cli': 'auto_awesome',
    'github-copilot': 'group',
    'opencode': 'data_object',
    'windsurf': 'air',
};

function catLabel(cat) {
    const labels = {
        'core-development': 'Core Dev',
        'language-specialists': 'Languages',
        'infrastructure': 'Infra',
        'quality-security': 'Quality',
        'data-ai': 'Data & AI',
        'developer-experience': 'DevEx',
        'specialized-domains': 'Specialized',
        'business-product': 'Business',
        'meta-orchestration': 'Orchestration',
        'research-analysis': 'Research',
        'game-development': 'Game Dev',
        'Developer Tools': 'Dev Tools',
        'Product & Strategy': 'Product',
        'AI & Machine Learning': 'AI/ML',
        'Testing & QA': 'Testing',
        'Cloud & Infrastructure': 'Cloud',
        'Backend & APIs': 'Backend',
        'DevOps & Monitoring': 'DevOps',
        'Databases & Data': 'Data',
        'Frontend & UI': 'Frontend',
        'Productivity & Collaboration': 'Productivity',
        'Security': 'Security',
        'Documents & Content': 'Documents',
        'Search & Web': 'Search',
        'Mobile & Desktop': 'Mobile',
        'Game Development': 'Game Dev',
    };
    return labels[cat] || cat || '—';
}

function agentCategories(a) {
    return (a.categories && a.categories.length) ? a.categories : (a.category ? a.category.split(',').map(c => c.trim()).filter(Boolean) : []);
}

function allCategories() {
    const s = new Set();
    agents.forEach(a => agentCategories(a).forEach(c => s.add(c)));
    return [...s].sort();
}

function skillCategories(s) {
    return (s.categories && s.categories.length) ? s.categories : (s.category ? s.category.split(',').map(c => c.trim()).filter(Boolean) : []);
}

function skillVendor(s) {
    return s.frontmatter?.vendor || '';
}

function allSkillVendors() {
    const s = new Set();
    skills.forEach(sk => { const v = skillVendor(sk); if (v) s.add(v); });
    return [...s].sort();
}

function allSkillCategories() {
    const s = new Set();
    skills.forEach(sk => skillCategories(sk).forEach(c => s.add(c)));
    return [...s].sort();
}

function shortName(fullName) {
    const parts = fullName.split('-');
    if (parts[0] === 'dotagent' && parts.length > 2) return parts.slice(2).join('-');
    return fullName;
}

function resolveTargets(entry, platforms) {
    if (!entry || entry.disabled) return [];
    const t = entry.targets || [];
    if (t.length === 1 && t[0] === 'all') return [...platforms];
    return t;
}

// ── API ──
async function api(path, opts = {}) {
    const res = await fetch(path, {
        headers: { 'Content-Type': 'application/json' },
        ...opts,
    });
    if (!res.ok) {
        const err = await res.json().catch(() => ({ error: res.statusText }));
        throw new Error(err.error || res.statusText);
    }
    return res.json();
}

// ── Utils ──
function esc(s) {
    if (!s) return '';
    const d = document.createElement('div');
    d.textContent = s;
    return d.innerHTML;
}

function truncate(s, n = 80) {
    if (!s) return '';
    return s.length <= n ? s : s.slice(0, n - 1).trimEnd() + '…';
}

function showSnackbar(msg, ms = 3000) {
    const el = document.getElementById('snackbar');
    el.textContent = msg;
    el.classList.add('show');
    clearTimeout(el._timer);
    el._timer = setTimeout(() => el.classList.remove('show'), ms);
}

// ═══════════════════════════════════
// VIEW NAVIGATION
// ═══════════════════════════════════
function switchView(view) {
    currentView = view;
    document.querySelectorAll('.nav-item').forEach(item => {
        const active = item.dataset.view === view;
        if (active) {
            item.classList.add('text-primary', 'bg-surface-container-high', 'border-r-2', 'border-primary');
            item.classList.remove('text-on-surface-variant');
        } else {
            item.classList.remove('text-primary', 'bg-surface-container-high', 'border-r-2', 'border-primary');
            item.classList.add('text-on-surface-variant');
        }
    });
    document.querySelectorAll('[id^="view-"]').forEach(v => v.classList.add('hidden'));
    const viewEl = document.getElementById('view-' + view);
    if (viewEl) viewEl.classList.remove('hidden');

    const titles = {
        overview: ['Overview', ''],
        agents: ['Agent Library', ''],
        skills: ['Skill Library', ''],
        status: ['Sync Status', ''],
    };
    document.getElementById('topbar-title').textContent = titles[view]?.[0] || view;
    document.getElementById('topbar-subtitle').textContent = titles[view]?.[1] || '';

    if (view === 'overview') renderOverview();
    else if (view === 'agents') loadAgents();
    else if (view === 'skills') loadSkills();
    else if (view === 'status') loadStatus();
}

function quickCreate() {
    if (currentView === 'skills') showCreateSkill();
    else showCreateAgent();
}

// ═══════════════════════════════════
// DATA LOADING
// ═══════════════════════════════════
async function loadAll() {
    try {
        const [a, c, t, s, sk] = await Promise.all([
            api('/api/agents'),
            api('/api/config'),
            api('/api/targets'),
            api('/api/status'),
            api('/api/skills').catch(() => []),
        ]);
        agents = a || [];
        skills = sk || [];
        config = c || { targets: [], agents: {}, skills: {} };
        if (!config.skills) config.skills = {};
        knownTargets = (t && t.targets) || config.targets || [];
        statusLinks = (s && s.symlinks) || [];
        updateNavBadges();
        updateFooter();
    } catch (e) {
        showSnackbar('Failed to load data: ' + e.message, 5000);
    }
}

async function refreshStatus() {
    try {
        const s = await api('/api/status');
        statusLinks = (s && s.symlinks) || [];
        updateNavBadges();
        updateFooter();
        if (currentView === 'overview') renderOverview();
        if (currentView === 'status') { renderStatusFilters(); renderStatusList(); }
    } catch (e) { /* silent */ }
}

function updateNavBadges() {
    document.getElementById('nav-agent-count').textContent = agents.length || '';
    document.getElementById('nav-skill-count').textContent = skills.length || '';
    const broken = statusLinks.filter(l => l.broken).length;
    const badge = document.getElementById('nav-status-badge');
    if (broken > 0) {
        badge.textContent = broken;
        badge.classList.remove('hidden');
        badge.classList.add('w-2', 'h-2', 'rounded-full', 'bg-red-500');
    } else {
        badge.classList.add('hidden');
    }
}

function updateFooter() {
    const info = document.getElementById('footer-info');
    const synced = statusLinks.filter(l => !l.broken).length;
    const broken = statusLinks.filter(l => l.broken).length;
    info.textContent = `${agents.length} agents | ${skills.length} skills | ${synced} symlinks${broken > 0 ? ` | ${broken} broken` : ''}`;
}

function updateClock() {
    const now = new Date();
    const h = String(now.getUTCHours()).padStart(2, '0');
    const m = String(now.getUTCMinutes()).padStart(2, '0');
    const s = String(now.getUTCSeconds()).padStart(2, '0');
    document.getElementById('footer-clock').textContent = `UTC ${h}:${m}:${s}`;
}

// ═══════════════════════════════════
// OVERVIEW DASHBOARD
// ═══════════════════════════════════
async function renderOverview() {
    if (agents.length === 0 && statusLinks.length === 0) await loadAll();

    // Stats
    document.getElementById('stat-agents').textContent = agents.length;
    document.getElementById('stat-agents-bar').style.width = agents.length > 0 ? '100%' : '0%';

    document.getElementById('stat-skills').textContent = skills.length;
    document.getElementById('stat-skills-bar').style.width = skills.length > 0 ? '100%' : '0%';

    const total = statusLinks.length;
    const healthy = statusLinks.filter(l => !l.broken).length;
    const healthPct = total > 0 ? Math.round((healthy / total) * 100) : (agents.length > 0 ? 100 : 0);
    document.getElementById('stat-health').textContent = healthPct + '%';
    document.getElementById('stat-health-detail').textContent = total > 0 ? `${healthy}/${total}` : 'idle';
    const bar = document.getElementById('stat-health-bar');
    bar.style.width = healthPct + '%';
    bar.className = healthPct === 100 ? 'h-full bg-emerald-500 transition-all duration-500' :
                    healthPct > 50 ? 'h-full bg-amber-500 transition-all duration-500' :
                    'h-full bg-red-500 transition-all duration-500';

    renderOverviewCategories();
    renderOverviewPlatforms();
    renderOverviewRecent();
}

function renderOverviewCategories() {
    const counts = {};
    agents.forEach(a => {
        const cats = agentCategories(a);
        if (cats.length === 0) counts['uncategorized'] = (counts['uncategorized'] || 0) + 1;
        else cats.forEach(c => counts[c] = (counts[c] || 0) + 1);
    });

    const entries = Object.entries(counts).sort((a, b) => b[1] - a[1]).slice(0, 10);
    const total = agents.length || 1;
    const container = document.getElementById('overview-categories');

    if (entries.length === 0) {
        container.innerHTML = '<p class="text-mono-sm text-on-surface-variant col-span-2">No agents loaded.</p>';
        return;
    }

    const half = Math.ceil(entries.length / 2);
    const col1 = entries.slice(0, half);
    const col2 = entries.slice(half);

    const renderCol = (items) => items.map(([cat, n]) => {
        const pct = Math.round((n / total) * 100);
        return `<div class="space-y-1">
            <div class="flex justify-between text-mono-sm">
                <span class="text-on-surface">${esc(catLabel(cat))}</span>
                <span class="text-primary font-bold">${n}</span>
            </div>
            <div class="h-2 bg-surface-container-highest">
                <div class="h-full bg-primary transition-all duration-500" style="width: ${pct}%"></div>
            </div>
        </div>`;
    }).join('');

    container.innerHTML = `<div class="space-y-3">${renderCol(col1)}</div><div class="space-y-3">${renderCol(col2)}</div>`;
}

function renderOverviewPlatforms() {
    const container = document.getElementById('overview-platforms');
    if (knownTargets.length === 0) {
        container.innerHTML = '<p class="text-mono-sm text-on-surface-variant">No platforms configured.</p>';
        return;
    }

    container.innerHTML = knownTargets.map(t => {
        const platformLinks = statusLinks.filter(l => l.platform === t);
        const total = platformLinks.length;
        const healthy = platformLinks.filter(l => !l.broken).length;
        const pct = total > 0 ? Math.round((healthy / total) * 100) : 0;
        const color = pct === 100 ? 'bg-emerald-500' : pct > 50 ? 'bg-amber-500' : total === 0 ? 'bg-surface-container-highest' : 'bg-red-500';
        return `<div class="flex items-center gap-3">
            <div class="w-7 h-7 flex items-center justify-center font-mono-sm text-outline border border-outline-variant text-[10px]">${esc(PLATFORM_LABELS[t] || '?')}</div>
            <div class="flex-1 space-y-1">
                <div class="flex justify-between font-mono-sm">
                    <span class="text-on-surface">${esc(PLATFORM_NAMES[t] || t)}</span>
                    <span class="text-primary">${total > 0 ? pct + '%' : '—'}</span>
                </div>
                <div class="h-1 bg-surface-container-highest">
                    <div class="h-full ${color} transition-all duration-500" style="width: ${total > 0 ? pct + '%' : '0%'}"></div>
                </div>
            </div>
        </div>`;
    }).join('');
}

function renderOverviewRecent() {
    const container = document.getElementById('overview-recent');
    if (statusLinks.length === 0) {
        container.innerHTML = '<tr><td colspan="5" class="px-4 py-8 text-center text-on-surface-variant">No symlinks. Toggle a platform to get started.</td></tr>';
        return;
    }

    const recent = statusLinks.slice(0, 8);
    container.innerHTML = recent.map(l => {
        const statusHtml = l.broken
            ? '<span class="text-red-400 font-bold">BROKEN</span>'
            : '<span class="text-emerald-500 font-bold">SYNCED</span>';
        return `<tr class="hover:bg-surface-container-low transition-colors cursor-pointer" onclick="switchView('status')">
            <td class="px-4 py-2 text-primary">${esc(l.agent)}</td>
            <td class="px-4 py-2"><span class="px-1.5 py-0.5 border border-outline-variant bg-surface-variant text-[10px] uppercase">${esc(l.type || 'agent')}</span></td>
            <td class="px-4 py-2"><span class="px-1.5 py-0.5 border border-outline-variant bg-surface-variant text-[10px]">${esc((PLATFORM_NAMES[l.platform] || l.platform).toUpperCase())}</span></td>
            <td class="px-4 py-2 text-outline truncate max-w-[200px]">${esc(l.path)}</td>
            <td class="px-4 py-2 text-right">${statusHtml}</td>
        </tr>`;
    }).join('');
}

// ═══════════════════════════════════
// AGENT LIBRARY
// ═══════════════════════════════════
async function loadAgents() {
    await loadAll();
    renderAgentFilters();
    renderAgentsTable();
    renderAgentPlatformHeader();
}

function renderAgentPlatformHeader() {
    const th = document.getElementById('agent-platform-header');
    if (knownTargets.length === 0) { th.textContent = 'Platforms'; return; }
    th.innerHTML = knownTargets.map(t =>
        `<span title="${esc(PLATFORM_NAMES[t] || t)}" style="display:inline-block;padding:0 4px;font-size:10px">${esc(PLATFORM_LABELS[t] || t)}</span>`
    ).join('');
}

function renderAgentFilters() {
    const counts = {};
    agents.forEach(a => {
        const cats = agentCategories(a);
        if (cats.length === 0) counts['uncategorized'] = (counts['uncategorized'] || 0) + 1;
        else cats.forEach(c => counts[c] = (counts[c] || 0) + 1);
    });

    const container = document.getElementById('agent-filters');
    let html = `<span class="text-[10px] font-label-caps text-on-surface-variant">FILTER:</span>`;
    html += `<button class="filter-chip ${categoryFilter === 'all' ? 'active' : ''}" onclick="setCategory('all')">All</button>`;
    Object.entries(counts).sort((a, b) => b[1] - a[1]).forEach(([cat, n]) => {
        html += `<button class="filter-chip ${categoryFilter === cat ? 'active' : ''}" onclick="setCategory('${esc(cat)}')">${esc(catLabel(cat))} <span style="opacity:.5">${n}</span></button>`;
    });
    container.innerHTML = html;

    const bulkContainer = document.getElementById('agent-bulk-actions');
    bulkContainer.innerHTML = '';
}

function setCategory(cat) {
    categoryFilter = cat;
    agentPage = 1;
    renderAgentFilters();
    renderAgentsTable();
}

function getFilteredAgents() {
    let list = agents;
    if (categoryFilter !== 'all') {
        list = list.filter(a => {
            const cats = agentCategories(a);
            return cats.includes(categoryFilter) || (cats.length === 0 && categoryFilter === 'uncategorized');
        });
    }
    if (searchQuery) {
        const q = searchQuery.toLowerCase();
        list = list.filter(a =>
            a.name.toLowerCase().includes(q) ||
            (a.description || '').toLowerCase().includes(q) ||
            agentCategories(a).some(c => c.toLowerCase().includes(q))
        );
    }
    return list;
}

function renderAgentsTable() {
    const filtered = getFilteredAgents();
    const empty = document.getElementById('agent-empty');
    const matrixView = document.getElementById('agent-matrix-view');
    const countEl = document.getElementById('agent-count');

    if (agents.length === 0) {
        empty.classList.remove('hidden');
        empty.classList.add('flex');
        matrixView.classList.add('hidden');
        countEl.textContent = '';
        return;
    }
    empty.classList.add('hidden');
    empty.classList.remove('flex');

    if (filtered.length === 0) {
        countEl.textContent = '0 of ' + agents.length;
        return;
    }

    countEl.textContent = `${filtered.length} of ${agents.length} agents`;
    matrixView.classList.remove('hidden');
    renderAgentMatrix(filtered);
    updateAgentBulkBar();
}

function renderAgentMatrix(filtered) {
    const head = document.getElementById('agent-matrix-head');
    const body = document.getElementById('agent-matrix-body');

    let headHtml = '<tr class="bg-surface">';
    headHtml += `<th class="w-10 h-8 border-b border-r border-outline-variant bg-surface sticky left-0 z-40 px-2"></th>`;
    headHtml += `<th class="w-48 h-8 border-b border-r border-outline-variant px-3 text-left font-label-caps text-on-surface-variant bg-surface sticky left-10 z-40 uppercase">Agent Name</th>`;
    headHtml += `<th class="w-24 h-8 border-b border-r border-outline-variant px-2 text-left font-label-caps text-on-surface-variant bg-surface uppercase">Category</th>`;
    knownTargets.forEach(t => {
        headHtml += `<th class="h-8 border-b border-r border-outline-variant px-2 bg-surface group">
            <div class="flex items-center justify-between">
                <span class="font-label-caps text-label-caps text-on-surface-variant uppercase">${esc(PLATFORM_NAMES[t] || t)}</span>
                <div class="flex gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button class="hover:text-primary" title="Toggle all agents for this platform" onclick="event.stopPropagation();agentColumnBulk('${esc(t)}')"><span class="material-symbols-outlined !text-[12px]">checklist</span></button>
                </div>
            </div>
        </th>`;
    });
    headHtml += '</tr>';
    head.innerHTML = headHtml;

    body.innerHTML = filtered.map(a => {
        const entry = config.agents?.[a.name];
        const active = resolveTargets(entry, knownTargets);
        const cats = agentCategories(a);
        const catText = cats.length ? esc(catLabel(cats[0])) : '—';
        const checked = selectedAgents.has(a.name);

        let cells = knownTargets.map(t => {
            const on = active.includes(t);
            const platformLinks = statusLinks.filter(l => l.agent === a.name && l.platform === t);
            const broken = platformLinks.some(l => l.broken);

            let icon, color, fill, label;
            if (!on) { icon = 'remove'; color = 'text-on-surface-variant opacity-20'; fill = false; label = 'Off'; }
            else if (broken) { icon = 'error'; color = 'text-red-400'; fill = true; label = 'Broken'; }
            else { icon = 'check_circle'; color = 'text-emerald-500'; fill = true; label = 'Synced'; }

            return `<td class="border-b border-r border-outline-variant px-2 matrix-cell-focus group" tabindex="0">
                <div class="flex items-center justify-center cursor-pointer" onclick="event.stopPropagation();togglePlatform('${esc(a.name)}','${esc(t)}',${!on})">
                    <span class="material-symbols-outlined matrix-cell-icon ${color}" ${fill ? "style='font-variation-settings: \"FILL\" 1;'" : ''}>${icon}</span>
                </div>
            </td>`;
        }).join('');

        return `<tr class="lib-row h-8 matrix-row ${checked ? 'row-selected' : ''}" onclick="viewAgent('${esc(a.name)}')">
            <td class="border-b border-r border-outline-variant text-center sticky left-0 matrix-sticky" onclick="event.stopPropagation()"><input type="checkbox" class="agent-cb" ${checked ? 'checked' : ''} onchange="toggleAgentSelect('${esc(a.name)}', this.checked)"></td>
            <td class="border-b border-r border-outline-variant px-3 text-primary font-bold sticky left-10 matrix-sticky" style="white-space:nowrap;overflow:hidden;text-overflow:ellipsis" title="${esc(a.name)}">${esc(shortName(a.name))}</td>
            <td class="border-b border-r border-outline-variant px-2 text-on-surface-variant italic text-[11px]" style="opacity:.5;white-space:nowrap;overflow:hidden;text-overflow:ellipsis">${catText}</td>
            ${cells}
        </tr>`;
    }).join('');
}

// ── Platform toggle ──
async function togglePlatform(agentName, platform, enable) {
    try {
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'agent', items: [{ name: agentName, platform, enable }] }) });
        if (res.config) { config.agents = res.config.agents || {}; config.skills = res.config.skills || {}; }
        renderAgentsTable();
        refreshStatus();
        showSnackbar(`${shortName(agentName)}: ${PLATFORM_NAMES[platform] || platform} ${enable ? 'enabled' : 'disabled'}`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

// ── Bulk selection ──
function toggleAgentSelect(name, checked) {
    if (checked) selectedAgents.add(name); else selectedAgents.delete(name);
    renderAgentsTable();
}

function toggleSelectAll(checked) {
    const filtered = getFilteredAgents();
    if (checked) filtered.forEach(a => selectedAgents.add(a.name));
    else filtered.forEach(a => selectedAgents.delete(a.name));
    renderAgentsTable();
}

function clearSelection() { selectedAgents.clear(); renderAgentsTable(); }

function updateAgentBulkBar() {
    const bar = document.getElementById('agent-bulk-bar');
    if (selectedAgents.size === 0) { bar.classList.add('hidden'); return; }
    bar.classList.remove('hidden');
    bar.classList.add('flex');
    document.getElementById('agent-bulk-count').textContent = `${selectedAgents.size} selected`;
    document.getElementById('agent-bulk-enable').innerHTML = knownTargets.map(t =>
        `<button class="bulk-btn" onclick="bulkTogglePlatform('${esc(t)}', true)">${esc(PLATFORM_LABELS[t] || t)}</button>`
    ).join('');
}

async function bulkTogglePlatform(platform, enable) {
    const items = [...selectedAgents].map(name => ({ name, platform, enable }));
    try {
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'agent', items }) });
        if (res.config) { config.agents = res.config.agents || {}; }
        renderAgentsTable();
        refreshStatus();
        showSnackbar(`${selectedAgents.size} agents: ${PLATFORM_NAMES[platform] || platform} ${enable ? 'enabled' : 'disabled'}`);
    } catch (e) { showSnackbar('Bulk failed: ' + e.message, 4000); }
}

async function bulkEnableAll() {
    const items = [];
    for (const name of selectedAgents) for (const t of knownTargets) items.push({ name, platform: t, enable: true });
    try {
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'agent', items }) });
        if (res.config) { config.agents = res.config.agents || {}; }
        renderAgentsTable();
        refreshStatus();
        showSnackbar(`${selectedAgents.size} agents: all platforms enabled`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

async function bulkDisableAll() {
    const items = [];
    for (const name of selectedAgents) for (const t of knownTargets) items.push({ name, platform: t, enable: false });
    try {
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'agent', items }) });
        if (res.config) { config.agents = res.config.agents || {}; }
        renderAgentsTable();
        refreshStatus();
        showSnackbar(`${selectedAgents.size} agents: all platforms disabled`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

// ═══════════════════════════════════
// AGENT DETAIL DRAWER
// ═══════════════════════════════════
function openDrawer() {
    document.getElementById('drawer-overlay').classList.remove('hidden');
    document.getElementById('detail-drawer').classList.add('open');
}

function closeDrawer() {
    document.getElementById('drawer-overlay').classList.add('hidden');
    document.getElementById('detail-drawer').classList.remove('open');
    drawerMode = null;
    drawerAgent = null;
}

async function viewAgent(name) {
    try {
        const a = await api('/api/agents/' + name);
        drawerAgent = a;
        drawerMode = 'view';
        const entry = config.agents?.[a.name];
        const active = resolveTargets(entry, knownTargets);

        document.getElementById('drawer-title').textContent = a.name;
        const cats = agentCategories(a);
        const catEl = document.getElementById('drawer-category');
        if (cats.length) { catEl.textContent = catLabel(cats[0]); catEl.classList.remove('hidden'); }
        else { catEl.classList.add('hidden'); }
        document.getElementById('drawer-path').textContent = `.dotagen/agents/${a.name}.md`;

        renderDrawerTargets(active);

        document.getElementById('drawer-view-mode').classList.remove('hidden');
        document.getElementById('drawer-edit-mode').classList.add('hidden');

        document.getElementById('drawer-meta').innerHTML = `
            ${a.description ? `<div class="detail-row"><span class="detail-label">DESCRIPTION</span><span class="detail-value">${esc(a.description)}</span></div>` : ''}
            <div class="detail-row"><span class="detail-label">CATEGORIES</span><span class="detail-value">${cats.map(c => esc(catLabel(c))).join(', ') || '—'}</span></div>
            <div class="detail-row"><span class="detail-label">PLATFORMS</span><span class="detail-value">${active.length ? active.map(t => esc(PLATFORM_NAMES[t] || t)).join(', ') : 'None'}</span></div>
        `;
        document.getElementById('drawer-content').textContent = a.content || '(empty)';

        document.getElementById('drawer-footer').innerHTML = `
            <button class="h-8 px-3 border border-red-500/50 text-red-400 font-mono-sm hover:bg-red-500/10 transition-colors flex items-center gap-1" onclick="deleteAgent('${esc(a.name)}')">
                <span class="material-symbols-outlined !text-[16px]">delete</span> DELETE
            </button>
            <div class="flex items-center gap-2">
                <button class="h-8 px-3 border border-outline-variant text-on-surface font-mono-sm hover:bg-surface-container-low transition-colors flex items-center gap-1" onclick="duplicateAgent('${esc(a.name)}')">
                    <span class="material-symbols-outlined !text-[16px]">content_copy</span> DUPLICATE
                </button>
                <button class="h-8 px-4 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all flex items-center gap-1" onclick="editAgent('${esc(a.name)}')">
                    <span class="material-symbols-outlined !text-[16px]">edit</span> EDIT
                </button>
            </div>
        `;
        openDrawer();
    } catch (e) { showSnackbar('Failed to load: ' + e.message, 4000); }
}

function renderDrawerTargets(active) {
    const container = document.getElementById('drawer-targets');
    container.innerHTML = knownTargets.map(t => {
        const checked = active.includes(t);
        return `<label class="target-chip ${checked ? 'checked' : ''}" onclick="event.preventDefault();">
            <input type="checkbox" name="drawer-target" value="${esc(t)}" ${checked ? 'checked' : ''} onchange="drawerToggleTarget('${esc(t)}', this.checked); this.closest('label').classList.toggle('checked', this.checked);">
            ${esc(PLATFORM_NAMES[t] || t)}
        </label>`;
    }).join('');
}

async function drawerToggleTarget(platform, enable) {
    if (!drawerAgent) return;
    await togglePlatform(drawerAgent.name, platform, enable);
    const entry = config.agents?.[drawerAgent.name];
    const active = resolveTargets(entry, knownTargets);
    if (drawerMode === 'view') {
        const metaEl = document.getElementById('drawer-meta');
        const rows = metaEl.querySelectorAll('.detail-row');
        if (rows.length >= 3) {
            rows[2].querySelector('.detail-value').textContent = active.length ? active.map(t => PLATFORM_NAMES[t] || t).join(', ') : 'None';
        }
    }
}

function showCreateAgent() {
    drawerMode = 'create';
    drawerAgent = null;
    document.getElementById('drawer-title').textContent = 'New Agent';
    document.getElementById('drawer-category').classList.add('hidden');
    document.getElementById('drawer-path').textContent = '.dotagen/agents/';
    renderDrawerTargets([...knownTargets]);
    renderAgentEditForm({ name: '', description: '', categories: [], content: '', targets: [...knownTargets] });
    document.getElementById('drawer-footer').innerHTML = `
        <span></span>
        <button class="h-8 px-4 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all flex items-center gap-1" onclick="submitCreateAgent()">
            <span class="material-symbols-outlined !text-[16px]">add</span> CREATE
        </button>
    `;
    openDrawer();
    setTimeout(() => { const inp = document.getElementById('form-name'); if (inp) inp.focus(); }, 300);
}

async function editAgent(name) {
    try {
        const a = await api('/api/agents/' + name);
        drawerAgent = a;
        drawerMode = 'edit';
        const entry = config.agents?.[a.name];
        const active = resolveTargets(entry, knownTargets);

        document.getElementById('drawer-title').textContent = a.name;
        document.getElementById('drawer-path').textContent = `.dotagen/agents/${a.name}.md`;
        const cats = agentCategories(a);
        const catEl = document.getElementById('drawer-category');
        if (cats.length) { catEl.textContent = catLabel(cats[0]); catEl.classList.remove('hidden'); }
        else { catEl.classList.add('hidden'); }

        renderDrawerTargets(active);
        renderAgentEditForm({ name: a.name, description: a.description || '', categories: cats, content: a.content || '', targets: active, isEdit: true });
        document.getElementById('drawer-footer').innerHTML = `
            <button class="h-8 px-3 border border-red-500/50 text-red-400 font-mono-sm hover:bg-red-500/10 transition-colors flex items-center gap-1" onclick="deleteAgent('${esc(a.name)}')">
                <span class="material-symbols-outlined !text-[16px]">delete</span> DELETE
            </button>
            <div class="flex items-center gap-2">
                <button class="h-8 px-3 border border-outline-variant text-on-surface font-mono-sm hover:bg-surface-container-low transition-colors" onclick="viewAgent('${esc(a.name)}')">CANCEL</button>
                <button class="h-8 px-4 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all flex items-center gap-1" onclick="submitEditAgent('${esc(a.name)}')">
                    <span class="material-symbols-outlined !text-[16px]" style="font-variation-settings: 'FILL' 1;">save</span> SAVE
                </button>
            </div>
        `;
        openDrawer();
    } catch (e) { showSnackbar('Failed to load: ' + e.message, 4000); }
}

function renderAgentEditForm({ name, description, categories = [], content, targets, isEdit = false }) {
    document.getElementById('drawer-view-mode').classList.add('hidden');
    document.getElementById('drawer-edit-mode').classList.remove('hidden');
    document.getElementById('drawer-edit-mode').classList.add('flex');

    const container = document.getElementById('edit-form-container');
    container.innerHTML = `
        <div class="form-group">
            <label class="form-label" for="form-name">Name</label>
            <input class="form-input" id="form-name" value="${esc(name)}" placeholder="my-agent" ${isEdit ? 'disabled' : ''}>
        </div>
        <div class="form-group">
            <label class="form-label" for="form-desc">Description</label>
            <input class="form-input" id="form-desc" value="${esc(description)}" placeholder="Short description…">
        </div>
        <div class="form-group">
            <label class="form-label">Categories</label>
            <div class="cat-dropdown" id="agent-cat-dropdown"></div>
        </div>
        <div class="form-group">
            <label class="form-label" for="form-content">Prompt (Markdown)</label>
            <textarea class="form-textarea" id="form-content" style="min-height:200px" placeholder="# Agent Name\n\nDescribe what this agent does…">${esc(content)}</textarea>
        </div>
    `;
    renderCategoryDropdown('agent-cat-dropdown', allCategories(), categories, 'form-category');

    // Render preview tabs
    renderPreviewTabs(targets);
}

function renderPreviewTabs(targets) {
    const tabContainer = document.getElementById('preview-tabs');
    if (targets.length === 0) {
        tabContainer.innerHTML = '<span class="px-3 py-2 text-mono-sm text-outline">No targets</span>';
        return;
    }
    tabContainer.innerHTML = targets.map((t, i) =>
        `<button class="preview-tab ${i === 0 ? 'active' : ''}" data-target="${esc(t)}" onclick="loadDrawerPreview('${esc(t)}'); document.querySelectorAll('.preview-tab').forEach(b => b.classList.remove('active')); this.classList.add('active');">${esc(PLATFORM_LABELS[t] || t)}</button>`
    ).join('');
    if (targets.length > 0) loadDrawerPreview(targets[0]);
}

let _previewDebounce = null;
async function loadDrawerPreview(target) {
    if (!drawerAgent && drawerMode !== 'create') return;
    const name = drawerAgent?.name;
    if (!name) {
        document.getElementById('preview-output').textContent = 'Save agent first to preview output.';
        return;
    }
    document.getElementById('preview-output').textContent = 'Loading…';
    try {
        const res = await api('/api/preview/' + name + '/' + target);
        document.getElementById('preview-output').textContent = res.content || '(empty)';
    } catch (e) {
        document.getElementById('preview-output').textContent = 'Error: ' + e.message;
    }
}

function getFormData() {
    const name = document.getElementById('form-name').value.trim();
    const description = document.getElementById('form-desc').value.trim();
    const category = getCategoryDropdownValues().join(',');
    const content = document.getElementById('form-content').value;
    const checks = document.querySelectorAll('input[name="drawer-target"]:checked');
    const targets = Array.from(checks).map(c => c.value);
    return { name, description, category, content, targets };
}

async function submitCreateAgent() {
    const data = getFormData();
    if (!data.name) { showSnackbar('Name is required'); return; }
    if (data.targets.length === 0) { showSnackbar('Select at least one platform'); return; }
    try {
        await api('/api/agents', { method: 'POST', body: JSON.stringify(data) });
        closeDrawer();
        showSnackbar(`Agent "${data.name}" created`);
        loadAgents();
    } catch (e) { showSnackbar('Create failed: ' + e.message, 5000); }
}

async function submitEditAgent(name) {
    const data = getFormData();
    if (data.targets.length === 0) { showSnackbar('Select at least one platform'); return; }
    try {
        await api('/api/agents/' + name, { method: 'PUT', body: JSON.stringify({ content: data.content, description: data.description, category: data.category, targets: data.targets }) });
        closeDrawer();
        showSnackbar(`Agent "${name}" saved`);
        loadAgents();
    } catch (e) { showSnackbar('Save failed: ' + e.message, 5000); }
}

async function deleteAgent(name) {
    const ok = await showConfirm(`Delete "${name}"?`, 'This removes the agent file and its config entry. This cannot be undone.');
    if (!ok) return;
    try {
        await api('/api/agents/' + name, { method: 'DELETE' });
        closeDrawer();
        showSnackbar(`Agent "${name}" deleted`);
        loadAgents();
    } catch (e) { showSnackbar('Delete failed: ' + e.message, 5000); }
}

async function duplicateAgent(name) {
    try {
        const a = await api('/api/agents/' + name);
        const entry = config.agents?.[name];
        const active = resolveTargets(entry, knownTargets);
        drawerMode = 'create';
        drawerAgent = null;
        document.getElementById('drawer-title').textContent = 'Duplicate Agent';
        document.getElementById('drawer-path').textContent = '.dotagen/agents/';
        renderDrawerTargets(active);
        renderAgentEditForm({ name: name + '-copy', description: a.description || '', categories: agentCategories(a), content: a.content || '', targets: active });
        document.getElementById('drawer-footer').innerHTML = `
            <span></span>
            <button class="h-8 px-4 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all flex items-center gap-1" onclick="submitCreateAgent()">
                <span class="material-symbols-outlined !text-[16px]">add</span> CREATE
            </button>
        `;
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

// ═══════════════════════════════════
// SKILL LIBRARY
// ═══════════════════════════════════
async function loadSkills() {
    await loadAll();
    renderSkillFilters();
    renderSkillsTable();
}

function renderSkillFilters() {
    const counts = {};
    skills.forEach(sk => {
        const cats = skillCategories(sk);
        if (cats.length === 0) counts['uncategorized'] = (counts['uncategorized'] || 0) + 1;
        else cats.forEach(c => counts[c] = (counts[c] || 0) + 1);
    });
    const container = document.getElementById('skill-filters');
    let html = `<span class="text-[10px] font-label-caps text-on-surface-variant">FILTER:</span>`;
    html += `<button class="filter-chip ${skillCategoryFilter === 'all' ? 'active' : ''}" onclick="setSkillCategory('all')">All</button>`;
    const sortedCats = Object.entries(counts).sort((a, b) => b[1] - a[1]);
    sortedCats.slice(0, 10).forEach(([cat, n]) => {
        html += `<button class="filter-chip ${skillCategoryFilter === cat ? 'active' : ''}" onclick="setSkillCategory('${esc(cat)}')">${esc(catLabel(cat))} <span style="opacity:.5">${n}</span></button>`;
    });
    if (sortedCats.length > 10) {
        html += `<select class="filter-chip" style="background:transparent;border:1px solid var(--c-outline-variant);color:var(--c-on-surface-variant);font-size:11px;padding:2px 6px;cursor:pointer" onchange="setSkillCategory(this.value);this.value=''"><option value="">More…</option>`;
        sortedCats.slice(10).forEach(([cat, n]) => {
            html += `<option value="${esc(cat)}">${esc(catLabel(cat))} (${n})</option>`;
        });
        html += `</select>`;
    }

    const vendorContainer = document.getElementById('skill-vendor-filter');
    if (vendorContainer) {
        const vendors = allSkillVendors();
        let vh = `<select class="filter-chip" style="background:transparent;border:1px solid var(--c-outline-variant);color:var(--c-on-surface-variant);font-size:11px;padding:2px 6px;cursor:pointer" onchange="setSkillVendor(this.value)">`;
        vh += `<option value="all" ${skillVendorFilter === 'all' ? 'selected' : ''}>All Vendors (${skills.length})</option>`;
        const vendorCounts = {};
        skills.forEach(sk => { const v = skillVendor(sk); if (v) vendorCounts[v] = (vendorCounts[v] || 0) + 1; });
        Object.entries(vendorCounts).sort((a, b) => b[1] - a[1]).forEach(([v, n]) => {
            vh += `<option value="${esc(v)}" ${skillVendorFilter === v ? 'selected' : ''}>${esc(v)} (${n})</option>`;
        });
        vh += `</select>`;
        vendorContainer.innerHTML = vh;
    }

    container.innerHTML = html;
}

function setSkillCategory(cat) { skillCategoryFilter = cat; skillPage = 1; renderSkillFilters(); renderSkillsTable(); }
function setSkillVendor(v) { skillVendorFilter = v; skillPage = 1; renderSkillFilters(); renderSkillsTable(); }

function getFilteredSkills() {
    let list = skills;
    if (skillCategoryFilter !== 'all') {
        list = list.filter(sk => {
            const cats = skillCategories(sk);
            return cats.includes(skillCategoryFilter) || (cats.length === 0 && skillCategoryFilter === 'uncategorized');
        });
    }
    if (skillVendorFilter !== 'all') {
        list = list.filter(sk => skillVendor(sk) === skillVendorFilter);
    }
    if (searchQuery) {
        const q = searchQuery.toLowerCase();
        list = list.filter(sk => sk.name.toLowerCase().includes(q) || (sk.description || '').toLowerCase().includes(q));
    }
    return list;
}

function renderSkillsTable() {
    const filtered = getFilteredSkills();
    const empty = document.getElementById('skill-empty');
    const countEl = document.getElementById('skill-count');
    const body = document.getElementById('skill-matrix-body');
    const head = document.getElementById('skill-matrix-head');

    if (skills.length === 0) {
        empty.classList.remove('hidden'); empty.classList.add('flex');
        head.innerHTML = ''; body.innerHTML = ''; countEl.textContent = '';
        return;
    }
    empty.classList.add('hidden'); empty.classList.remove('flex');

    if (filtered.length === 0) {
        head.innerHTML = '';
        body.innerHTML = `<tr><td colspan="${3 + knownTargets.length}" class="px-4 py-8 text-center text-on-surface-variant">No skills match your filters.</td></tr>`;
        countEl.textContent = '0 of ' + skills.length;
        return;
    }

    const visible = filtered.slice(0, skillPage * PAGE_SIZE);
    countEl.textContent = `${visible.length}/${filtered.length} of ${skills.length} skills`;
    renderSkillMatrix(visible, filtered.length);
}

function renderSkillMatrix(visible, totalFiltered) {
    const head = document.getElementById('skill-matrix-head');
    const body = document.getElementById('skill-matrix-body');

    let headHtml = '<tr class="bg-surface">';
    headHtml += `<th class="w-10 h-8 border-b border-r border-outline-variant bg-surface sticky left-0 z-40 px-2"></th>`;
    headHtml += `<th class="w-48 h-8 border-b border-r border-outline-variant px-3 text-left font-label-caps text-on-surface-variant bg-surface sticky left-10 z-40 uppercase">Skill Name</th>`;
    headHtml += `<th class="w-24 h-8 border-b border-r border-outline-variant px-2 text-left font-label-caps text-on-surface-variant bg-surface uppercase">Category</th>`;
    knownTargets.forEach(t => {
        headHtml += `<th class="h-8 border-b border-r border-outline-variant px-2 bg-surface group">
            <div class="flex items-center justify-between">
                <span class="font-label-caps text-label-caps text-on-surface-variant uppercase">${esc(PLATFORM_NAMES[t] || t)}</span>
                <div class="flex gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button class="hover:text-primary" title="Toggle all skills for this platform" onclick="event.stopPropagation();skillColumnBulk('${esc(t)}')"><span class="material-symbols-outlined !text-[12px]">checklist</span></button>
                </div>
            </div>
        </th>`;
    });
    headHtml += '</tr>';
    head.innerHTML = headHtml;

    body.innerHTML = visible.map(sk => {
        const entry = config.skills?.[sk.name];
        const active = resolveTargets(entry, knownTargets);
        const cats = skillCategories(sk);
        const catText = cats.length ? esc(catLabel(cats[0])) : '—';
        const checked = selectedSkills.has(sk.name);

        let cells = knownTargets.map(t => {
            const on = active.includes(t);
            const platformLinks = statusLinks.filter(l => l.agent === sk.name && l.platform === t);
            const broken = platformLinks.some(l => l.broken);

            let icon, color, fill, label;
            if (!on) { icon = 'remove'; color = 'text-on-surface-variant opacity-20'; fill = false; label = 'Off'; }
            else if (broken) { icon = 'error'; color = 'text-red-400'; fill = true; label = 'Broken'; }
            else { icon = 'check_circle'; color = 'text-emerald-500'; fill = true; label = 'Synced'; }

            return `<td class="border-b border-r border-outline-variant px-2 matrix-cell-focus group" tabindex="0">
                <div class="flex items-center justify-center cursor-pointer" onclick="event.stopPropagation();toggleSkillPlatform('${esc(sk.name)}','${esc(t)}',${!on})">
                    <span class="material-symbols-outlined matrix-cell-icon ${color}" ${fill ? "style='font-variation-settings: \"FILL\" 1;'" : ''}>${icon}</span>
                </div>
            </td>`;
        }).join('');

        return `<tr class="lib-row h-8 matrix-row ${checked ? 'row-selected' : ''}" onclick="viewSkill('${esc(sk.name)}')">
            <td class="border-b border-r border-outline-variant text-center sticky left-0 matrix-sticky" onclick="event.stopPropagation()"><input type="checkbox" class="agent-cb" ${checked ? 'checked' : ''} onchange="toggleSkillSelect('${esc(sk.name)}', this.checked)"></td>
            <td class="border-b border-r border-outline-variant px-3 text-primary font-bold sticky left-10 matrix-sticky" style="white-space:nowrap;overflow:hidden;text-overflow:ellipsis" title="${esc(sk.name)}">${esc(shortName(sk.name))}</td>
            <td class="border-b border-r border-outline-variant px-2 text-on-surface-variant italic text-[11px]" style="opacity:.5;white-space:nowrap;overflow:hidden;text-overflow:ellipsis">${catText}</td>
            ${cells}
        </tr>`;
    }).join('');

    const hasMore = visible.length < totalFiltered;
    if (hasMore) {
        body.innerHTML += `<tr><td colspan="${3 + knownTargets.length}" class="px-4 py-3 text-center"><button class="h-8 px-6 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all" onclick="skillPage++;renderSkillsTable()">Load More (${totalFiltered - visible.length} remaining)</button></td></tr>`;
    }
}

async function toggleSkillPlatform(name, platform, enable) {
    try {
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'skill', items: [{ name, platform, enable }] }) });
        if (res.config) { config.skills = res.config.skills || {}; }
        renderSkillsTable();
        refreshStatus();
        showSnackbar(`${shortName(name)}: ${PLATFORM_NAMES[platform] || platform} ${enable ? 'enabled' : 'disabled'}`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

function toggleSkillSelect(name, checked) { if (checked) selectedSkills.add(name); else selectedSkills.delete(name); renderSkillsTable(); }
function toggleSkillSelectAll(checked) { const f = getFilteredSkills(); if (checked) f.forEach(s => selectedSkills.add(s.name)); else f.forEach(s => selectedSkills.delete(s.name)); renderSkillsTable(); }
function clearSkillSelection() { selectedSkills.clear(); renderSkillsTable(); }

function updateSkillBulkBar() {
    const bar = document.getElementById('skill-bulk-bar');
    if (!bar) return;
    if (selectedSkills.size === 0) { bar.classList.add('hidden'); return; }
    bar.classList.remove('hidden'); bar.classList.add('flex');
    document.getElementById('skill-bulk-count').textContent = `${selectedSkills.size} selected`;
    document.getElementById('skill-bulk-enable').innerHTML = knownTargets.map(t =>
        `<button class="bulk-btn" onclick="skillBulkToggle('${esc(t)}', true)">${esc(PLATFORM_LABELS[t] || t)}</button>`).join('');
}

async function skillBulkToggle(platform, enable) {
    const items = [...selectedSkills].map(name => ({ name, platform, enable }));
    try {
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'skill', items }) });
        if (res.config) { config.skills = res.config.skills || {}; }
        renderSkillsTable();
        refreshStatus();
        showSnackbar(`${selectedSkills.size} skills: ${PLATFORM_NAMES[platform] || platform} ${enable ? 'enabled' : 'disabled'}`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

async function skillBulkEnableAll() {
    const items = [];
    for (const name of selectedSkills) for (const t of knownTargets) items.push({ name, platform: t, enable: true });
    try {
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'skill', items }) });
        if (res.config) { config.skills = res.config.skills || {}; }
        renderSkillsTable();
        refreshStatus();
        showSnackbar(`${selectedSkills.size} skills: all on`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

async function skillBulkDisableAll() {
    const items = [];
    for (const name of selectedSkills) for (const t of knownTargets) items.push({ name, platform: t, enable: false });
    try {
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'skill', items }) });
        if (res.config) { config.skills = res.config.skills || {}; }
        renderSkillsTable();
        refreshStatus();
        showSnackbar(`${selectedSkills.size} skills: all off`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

async function skillColumnBulk(platform) {
    if (selectedSkills.size > 0) {
        await skillBulkToggle(platform, true);
        return;
    }
    const filtered = getFilteredSkills();
    const allOn = filtered.every(s => resolveTargets(config.skills?.[s.name], knownTargets).includes(platform));
    const items = filtered.map(s => ({ name: s.name, platform, enable: !allOn }));
    try {
        showSnackbar(`${filtered.length} skills: ${PLATFORM_NAMES[platform] || platform} ${!allOn ? 'enabling' : 'disabling'}…`);
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'skill', items }) });
        if (res.config) { config.skills = res.config.skills || {}; }
        renderSkillsTable();
        refreshStatus();
        showSnackbar(`${filtered.length} skills: ${PLATFORM_NAMES[platform] || platform} ${!allOn ? 'enabled' : 'disabled'}`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

async function agentColumnBulk(platform) {
    if (selectedAgents.size > 0) {
        await bulkTogglePlatform(platform, true);
        return;
    }
    const filtered = getFilteredAgents();
    const allOn = filtered.every(a => resolveTargets(config.agents?.[a.name], knownTargets).includes(platform));
    const items = filtered.map(a => ({ name: a.name, platform, enable: !allOn }));
    try {
        showSnackbar(`${filtered.length} agents: ${PLATFORM_NAMES[platform] || platform} ${!allOn ? 'enabling' : 'disabling'}…`);
        const res = await api('/api/toggle', { method: 'POST', body: JSON.stringify({ type: 'agent', items }) });
        if (res.config) { config.agents = res.config.agents || {}; }
        renderAgentsTable();
        refreshStatus();
        showSnackbar(`${filtered.length} agents: ${PLATFORM_NAMES[platform] || platform} ${!allOn ? 'enabled' : 'disabled'}`);
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

// Skill drawer
async function viewSkill(name) {
    try {
        const sk = await api('/api/skills/' + name);
        drawerAgent = sk; drawerMode = 'view';
        const entry = config.skills?.[sk.name];
        const active = resolveTargets(entry, knownTargets);
        document.getElementById('drawer-title').textContent = sk.name;
        const cats = skillCategories(sk);
        const catEl = document.getElementById('drawer-category');
        if (cats.length) { catEl.textContent = catLabel(cats[0]); catEl.classList.remove('hidden'); }
        else { catEl.classList.add('hidden'); }
        document.getElementById('drawer-path').textContent = `.dotagen/skills/${sk.name}/SKILL.md`;
        renderDrawerTargets(active);
        document.getElementById('drawer-view-mode').classList.remove('hidden');
        document.getElementById('drawer-edit-mode').classList.add('hidden');
        const refList = (sk.references && sk.references.length) ? sk.references.map(r => esc(r.name)).join(', ') : 'None';
        document.getElementById('drawer-meta').innerHTML = `
            ${sk.description ? `<div class="detail-row"><span class="detail-label">DESCRIPTION</span><span class="detail-value">${esc(sk.description)}</span></div>` : ''}
            <div class="detail-row"><span class="detail-label">CATEGORY</span><span class="detail-value">${cats.map(c => esc(catLabel(c))).join(', ') || '—'}</span></div>
            <div class="detail-row"><span class="detail-label">PLATFORMS</span><span class="detail-value">${active.length ? active.map(t => esc(PLATFORM_NAMES[t] || t)).join(', ') : 'None'}</span></div>
            <div class="detail-row"><span class="detail-label">REFERENCES</span><span class="detail-value">${refList}</span></div>
        `;
        document.getElementById('drawer-content').textContent = sk.content || '(empty)';
        document.getElementById('drawer-footer').innerHTML = `
            <button class="h-8 px-3 border border-red-500/50 text-red-400 font-mono-sm hover:bg-red-500/10 transition-colors flex items-center gap-1" onclick="deleteSkill('${esc(sk.name)}')">
                <span class="material-symbols-outlined !text-[16px]">delete</span> DELETE
            </button>
            <div class="flex items-center gap-2">
                <button class="h-8 px-4 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all flex items-center gap-1" onclick="editSkill('${esc(sk.name)}')">
                    <span class="material-symbols-outlined !text-[16px]">edit</span> EDIT
                </button>
            </div>
        `;
        openDrawer();
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

function showCreateSkill() {
    drawerMode = 'create'; drawerAgent = null;
    document.getElementById('drawer-title').textContent = 'New Skill';
    document.getElementById('drawer-category').classList.add('hidden');
    document.getElementById('drawer-path').textContent = '.dotagen/skills/';
    renderDrawerTargets([...knownTargets]);
    renderSkillEditForm({ name: '', description: '', categories: [], content: '', targets: [...knownTargets] });
    document.getElementById('drawer-footer').innerHTML = `
        <span></span>
        <button class="h-8 px-4 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all flex items-center gap-1" onclick="submitCreateSkill()">
            <span class="material-symbols-outlined !text-[16px]">add</span> CREATE
        </button>
    `;
    openDrawer();
    setTimeout(() => { const inp = document.getElementById('form-name'); if (inp) inp.focus(); }, 300);
}

async function editSkill(name) {
    try {
        const sk = await api('/api/skills/' + name);
        drawerAgent = sk; drawerMode = 'edit';
        const entry = config.skills?.[sk.name];
        const active = resolveTargets(entry, knownTargets);
        document.getElementById('drawer-title').textContent = sk.name;
        document.getElementById('drawer-path').textContent = `.dotagen/skills/${sk.name}/SKILL.md`;
        const cats = skillCategories(sk);
        const catEl = document.getElementById('drawer-category');
        if (cats.length) { catEl.textContent = catLabel(cats[0]); catEl.classList.remove('hidden'); }
        else { catEl.classList.add('hidden'); }
        renderDrawerTargets(active);
        renderSkillEditForm({ name: sk.name, description: sk.frontmatter?.description || '', categories: skillCategories(sk), content: sk.content || '', targets: active, isEdit: true });
        document.getElementById('drawer-footer').innerHTML = `
            <button class="h-8 px-3 border border-red-500/50 text-red-400 font-mono-sm hover:bg-red-500/10 transition-colors flex items-center gap-1" onclick="deleteSkill('${esc(sk.name)}')">
                <span class="material-symbols-outlined !text-[16px]">delete</span> DELETE
            </button>
            <div class="flex items-center gap-2">
                <button class="h-8 px-3 border border-outline-variant text-on-surface font-mono-sm hover:bg-surface-container-low transition-colors" onclick="viewSkill('${esc(sk.name)}')">CANCEL</button>
                <button class="h-8 px-4 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all flex items-center gap-1" onclick="submitEditSkill('${esc(sk.name)}')">
                    <span class="material-symbols-outlined !text-[16px]" style="font-variation-settings: 'FILL' 1;">save</span> SAVE
                </button>
            </div>
        `;
        openDrawer();
    } catch (e) { showSnackbar('Failed: ' + e.message, 4000); }
}

function renderSkillEditForm({ name, description, categories = [], content, targets, isEdit = false }) {
    document.getElementById('drawer-view-mode').classList.add('hidden');
    document.getElementById('drawer-edit-mode').classList.remove('hidden');
    document.getElementById('drawer-edit-mode').classList.add('flex');
    const container = document.getElementById('edit-form-container');
    container.innerHTML = `
        <div class="form-group"><label class="form-label" for="form-name">Name</label><input class="form-input" id="form-name" value="${esc(name)}" placeholder="my-skill" ${isEdit ? 'disabled' : ''}></div>
        <div class="form-group"><label class="form-label" for="form-desc">Description</label><input class="form-input" id="form-desc" value="${esc(description)}" placeholder="When to trigger this skill…"></div>
        <div class="form-group"><label class="form-label">Categories</label><div class="cat-dropdown" id="agent-cat-dropdown"></div></div>
        <div class="form-group"><label class="form-label" for="form-content">SKILL.md Content</label><textarea class="form-textarea" id="form-content" style="min-height:200px" placeholder="# Skill Name\n\n## When to Use\n\n…">${esc(content)}</textarea></div>
    `;
    renderCategoryDropdown('agent-cat-dropdown', allSkillCategories(), categories, 'form-category');
    renderPreviewTabs(targets);
}

function getSkillFormData() {
    const category = getCategoryDropdownValues().join(',');
    return { name: document.getElementById('form-name').value.trim(), description: document.getElementById('form-desc').value.trim(), category, content: document.getElementById('form-content').value, targets: Array.from(document.querySelectorAll('input[name="drawer-target"]:checked')).map(c => c.value) };
}

async function submitCreateSkill() {
    const data = getSkillFormData();
    if (!data.name) { showSnackbar('Name is required'); return; }
    try { await api('/api/skills', { method: 'POST', body: JSON.stringify(data) }); closeDrawer(); showSnackbar(`Skill "${data.name}" created`); loadSkills(); } catch (e) { showSnackbar('Create failed: ' + e.message, 5000); }
}

async function submitEditSkill(name) {
    const data = getSkillFormData();
    try { await api('/api/skills/' + name, { method: 'PUT', body: JSON.stringify({ content: data.content, description: data.description, category: data.category, targets: data.targets }) }); closeDrawer(); showSnackbar(`Skill "${name}" saved`); loadSkills(); } catch (e) { showSnackbar('Save failed: ' + e.message, 5000); }
}

async function deleteSkill(name) {
    const ok = await showConfirm(`Delete "${name}"?`, 'This removes the skill directory and config entry. Cannot be undone.');
    if (!ok) return;
    try { await api('/api/skills/' + name, { method: 'DELETE' }); closeDrawer(); showSnackbar(`Skill "${name}" deleted`); loadSkills(); } catch (e) { showSnackbar('Delete failed: ' + e.message, 5000); }
}

// ═══════════════════════════════════
// STATUS VIEW
// ═══════════════════════════════════
async function loadStatus() {
    try {
        if (statusLinks.length === 0 && agents.length === 0) await loadAll();
        renderStatusFilters();
        renderStatusList();
    } catch (e) { showSnackbar('Failed to load status: ' + e.message, 4000); }
}

function renderStatusFilters() {
    const broken = statusLinks.filter(l => l.broken).length;
    const healthy = statusLinks.length - broken;
    const container = document.getElementById('status-filters');
    container.innerHTML = `
        <span class="text-[10px] font-label-caps text-on-surface-variant">FILTER:</span>
        <button class="filter-chip ${statusFilter === 'all' ? 'active' : ''}" onclick="setStatusFilter('all')">All ${statusLinks.length}</button>
        <button class="filter-chip ${statusFilter === 'healthy' ? 'active' : ''}" onclick="setStatusFilter('healthy')">Synced ${healthy}</button>
        <button class="filter-chip ${statusFilter === 'broken' ? 'active' : ''}" onclick="setStatusFilter('broken')">Broken ${broken}</button>
    `;
}

function setStatusFilter(f) { statusFilter = f; renderStatusFilters(); renderStatusList(); }

function renderStatusList() {
    const tbody = document.getElementById('status-tbody');
    const empty = document.getElementById('status-empty');

    let filtered = statusLinks;
    if (statusFilter === 'healthy') filtered = filtered.filter(l => !l.broken);
    else if (statusFilter === 'broken') filtered = filtered.filter(l => l.broken);

    if (statusLinks.length === 0) {
        tbody.innerHTML = '';
        empty.classList.remove('hidden');
        empty.classList.add('flex');
        return;
    }
    empty.classList.add('hidden');
    empty.classList.remove('flex');

    tbody.innerHTML = filtered.map(l => `
        <tr class="hover:bg-surface-container-low transition-colors">
            <td class="px-4 py-1.5"><span class="status-dot ${l.broken ? 'status-dot-err' : 'status-dot-ok'}"></span></td>
            <td class="px-4 py-1.5 text-primary font-bold">${esc(l.agent)}</td>
            <td class="px-4 py-1.5"><span class="px-1.5 py-0.5 border border-outline-variant bg-surface-variant text-[10px] uppercase">${esc(l.type || 'agent')}</span></td>
            <td class="px-4 py-1.5"><span class="px-1.5 py-0.5 border border-outline-variant bg-surface-variant text-[10px]">${esc((PLATFORM_NAMES[l.platform] || l.platform).toUpperCase())}</span></td>
            <td class="px-4 py-1.5 text-outline truncate max-w-[300px]">${esc(l.path)}</td>
            <td class="px-2 py-1.5"><button class="text-on-surface-variant hover:text-primary transition-colors" title="Preview" onclick="quickPreview('${esc(l.agent)}','${esc(l.platform)}','${esc(l.type || 'agent')}')"><span class="material-symbols-outlined !text-[16px]">visibility</span></button></td>
        </tr>
    `).join('');
}

function quickPreview(agent, platform, type) {
    if (type === 'skill') {
        viewSkill(agent);
    } else {
        viewAgent(agent);
    }
    setTimeout(() => {
        if (drawerMode === 'view') {
            document.getElementById('drawer-view-mode').classList.add('hidden');
            document.getElementById('drawer-edit-mode').classList.remove('hidden');
            document.getElementById('drawer-edit-mode').classList.add('flex');
            renderPreviewTabs([platform]);
        }
    }, 500);
}

// ═══════════════════════════════════
// CLEAN
// ═══════════════════════════════════

async function triggerClean() {
    const ok = await showConfirm('Remove ALL generated files?', 'All symlinks and generated output will be deleted. Source files are kept.');
    if (!ok) return;
    try {
        const res = await api('/api/clean', { method: 'POST' });
        if (res.config) { config.agents = res.config.agents || {}; config.skills = res.config.skills || {}; }
        showSnackbar(`Removed ${res.removed || 0} files`);
        await refreshStatus();
        if (currentView === 'agents') renderAgentsTable();
        if (currentView === 'skills') renderSkillsTable();
        if (currentView === 'overview') renderOverview();
    } catch (e) { showSnackbar('Clean failed: ' + e.message, 5000); }
}

async function triggerCleanBroken() {
    try {
        const res = await api('/api/clean-broken', { method: 'POST' });
        if (res.config) { config.agents = res.config.agents || {}; config.skills = res.config.skills || {}; }
        const parts = [];
        if (res.fixed) parts.push(`fixed ${res.fixed}`);
        if (res.removed) parts.push(`removed ${res.removed}`);
        showSnackbar(parts.length ? `Broken symlinks: ${parts.join(', ')}` : 'No broken symlinks');
        await refreshStatus();
        if (currentView === 'agents') renderAgentsTable();
        if (currentView === 'skills') renderSkillsTable();
        if (currentView === 'overview') renderOverview();
    } catch (e) { showSnackbar('Failed: ' + e.message, 5000); }
}

async function triggerInit() {
    const ok = await showConfirm('Reinitialize ~/.dotagen?', 'This recreates all built-in agents and skills with fresh config. Your existing agents/skills and custom edits will be overwritten.');
    if (!ok) return;
    document.getElementById('footer-status').textContent = 'INITIALIZING…';
    try {
        const res = await api('/api/init', { method: 'POST' });
        showSnackbar(`Reinitialized: ${res.agents || 0} agents, ${res.skills || 0} skills`);
        await loadAll();
        document.getElementById('footer-status').textContent = 'READY';
        if (currentView === 'overview') renderOverview();
        if (currentView === 'agents') renderAgentsTable();
        if (currentView === 'skills') renderSkillsTable();
        if (currentView === 'status') { renderStatusFilters(); renderStatusList(); }
    } catch (e) {
        showSnackbar('Init failed: ' + e.message, 5000);
        document.getElementById('footer-status').textContent = 'ERROR';
    }
}

// ═══════════════════════════════════
// CONFIRM MODAL
// ═══════════════════════════════════
function showConfirm(title, body) {
    return new Promise(resolve => {
        const modal = document.getElementById('modal');
        document.getElementById('modal-title').textContent = title;
        document.getElementById('modal-body').textContent = body;
        document.getElementById('modal-actions').innerHTML = `
            <button class="h-8 px-3 border border-outline-variant text-on-surface font-mono-sm hover:bg-surface-container-low transition-colors" id="confirm-no">CANCEL</button>
            <button class="h-8 px-4 bg-primary text-on-primary font-mono-sm font-bold hover:opacity-90 active:scale-95 transition-all" id="confirm-yes">CONFIRM</button>
        `;
        modal.classList.remove('hidden');
        const cleanup = (result) => { modal.classList.add('hidden'); resolve(result); };
        document.getElementById('confirm-no').onclick = () => cleanup(false);
        document.getElementById('confirm-yes').onclick = () => cleanup(true);
        modal.onclick = (e) => { if (e.target === modal) cleanup(false); };
    });
}

// ═══════════════════════════════════
// CATEGORY DROPDOWN COMPONENT
// ═══════════════════════════════════
let _catDropdownState = { selected: new Set(), allCats: [], containerId: '', inputName: '' };

function renderCategoryDropdown(containerId, allCats, selectedCats, inputName) {
    _catDropdownState = { selected: new Set(selectedCats), allCats: [...new Set([...allCats, ...selectedCats])].sort(), containerId, inputName };
    _rebuildDropdown();
}

function _rebuildDropdown() {
    const { selected, allCats, containerId } = _catDropdownState;
    const container = document.getElementById(containerId);
    if (!container) return;
    const tags = [...selected].map(c => `<span class="cat-tag">${esc(catLabel(c))}<span class="cat-tag-x" data-cat="${esc(c)}">&times;</span></span>`).join('');
    const placeholder = selected.size === 0 ? '<span class="cat-dropdown-placeholder">Select categories…</span>' : '';
    const items = allCats.map(c => { const sel = selected.has(c) ? 'selected' : ''; return `<div class="cat-dropdown-item ${sel}" data-cat="${esc(c)}"><span class="cat-check">${sel ? '✓' : ''}</span>${esc(catLabel(c))}</div>`; }).join('');
    container.innerHTML = `<div class="cat-dropdown-trigger" id="${containerId}-trigger">${tags}${placeholder}<span class="cat-dropdown-arrow">▾</span></div><div class="cat-dropdown-menu" id="${containerId}-menu">${items || '<div style="padding:4px 10px;font-size:11px;color:#444748">No categories yet</div>'}<div class="cat-dropdown-add"><input id="${containerId}-new" placeholder="New category…" onclick="event.stopPropagation()"><button type="button" onclick="event.stopPropagation();_catDropdownAddNew('${containerId}')">Add</button></div></div>`;
    document.getElementById(`${containerId}-trigger`).onclick = (e) => {
        if (e.target.classList.contains('cat-tag-x')) { e.stopPropagation(); _catDropdownToggle(e.target.dataset.cat); return; }
        const menu = document.getElementById(`${containerId}-menu`);
        const trigger = document.getElementById(`${containerId}-trigger`);
        const isOpen = menu.classList.contains('open');
        menu.classList.toggle('open', !isOpen);
        trigger.classList.toggle('open', !isOpen);
    };
    container.querySelectorAll('.cat-dropdown-item').forEach(item => { item.onclick = (e) => { e.stopPropagation(); _catDropdownToggle(item.dataset.cat); }; });
    const newInput = document.getElementById(`${containerId}-new`);
    if (newInput) { newInput.onkeydown = (e) => { if (e.key === 'Enter') { e.preventDefault(); _catDropdownAddNew(containerId); } }; }
}

function _catDropdownToggle(cat) { if (_catDropdownState.selected.has(cat)) _catDropdownState.selected.delete(cat); else _catDropdownState.selected.add(cat); _rebuildDropdown(); }
function _catDropdownAddNew(containerId) { const input = document.getElementById(`${containerId}-new`); if (!input) return; const val = input.value.trim().toLowerCase().replace(/[^a-z0-9-]/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, ''); if (!val) return; if (!_catDropdownState.allCats.includes(val)) { _catDropdownState.allCats.push(val); _catDropdownState.allCats.sort(); } _catDropdownState.selected.add(val); _rebuildDropdown(); }
function getCategoryDropdownValues() { return [..._catDropdownState.selected]; }

document.addEventListener('click', (e) => {
    const dropdown = e.target.closest('.cat-dropdown');
    document.querySelectorAll('.cat-dropdown-menu.open').forEach(menu => { if (!dropdown || !dropdown.contains(menu)) { menu.classList.remove('open'); menu.previousElementSibling?.classList.remove('open'); } });
});

// ═══════════════════════════════════
// SEARCH
// ═══════════════════════════════════
const searchInput = document.getElementById('global-search');
let searchDebounce = null;
searchInput.addEventListener('input', (e) => {
    searchQuery = e.target.value.trim().toLowerCase();
    clearTimeout(searchDebounce);
    searchDebounce = setTimeout(() => {
        skillPage = 1;
        agentPage = 1;
        if (currentView === 'agents') renderAgentsTable();
        else if (currentView === 'skills') renderSkillsTable();
        else if (searchQuery) switchView('agents');
    }, 200);
});

// ═══════════════════════════════════
// KEYBOARD SHORTCUTS
// ═══════════════════════════════════
document.addEventListener('keydown', (e) => {
    const active = document.activeElement;
    const typing = active && (active.tagName === 'INPUT' || active.tagName === 'TEXTAREA' || active.tagName === 'SELECT' || active.isContentEditable);

    if (e.key === 'Escape') {
        if (!document.getElementById('modal').classList.contains('hidden')) { document.getElementById('modal').classList.add('hidden'); }
        else if (document.getElementById('detail-drawer').classList.contains('open')) { closeDrawer(); }
        return;
    }

    if ((e.metaKey || e.ctrlKey) && e.key === 'k') { e.preventDefault(); searchInput.focus(); return; }
    if (typing || e.metaKey || e.ctrlKey) return;

    if (e.key === '/') { e.preventDefault(); searchInput.focus(); return; }
    if (e.key === 'g') { e.preventDefault(); switchView('overview'); return; }
    if (e.key === 'a') { e.preventDefault(); switchView('agents'); return; }
    if (e.key === 's') { e.preventDefault(); switchView('skills'); return; }
    if (e.key === 'n') { e.preventDefault(); quickCreate(); return; }
});

// ═══════════════════════════════════
// INIT
// ═══════════════════════════════════
updateClock();
setInterval(updateClock, 1000);
loadAll().then(() => renderOverview());
