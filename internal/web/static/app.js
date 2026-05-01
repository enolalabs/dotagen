/* ═══════════════════════════════════════════════
   dotagen dashboard — app.js
   ═══════════════════════════════════════════════ */

// ── State ──
let agents = [];
let config = { targets: [], agents: {} };
let knownTargets = [];
let statusLinks = [];
let currentTab = 'agents';
let searchQuery = '';
let categoryFilter = 'all';
let statusFilter = 'all';
let panelMode = null; // null | 'view' | 'create' | 'edit'
let panelAgent = null;
let selectedAgents = new Set();

const PLATFORM_LABELS = {
    'claude-code': 'Claude',
    'cursor': 'Cursor',
    'gemini-cli': 'Gemini',
    'opencode': 'OpenCode',
};

const PLATFORM_NAMES = {
    'claude-code': 'Claude Code',
    'cursor': 'Cursor',
    'gemini-cli': 'Gemini CLI',
    'opencode': 'OpenCode',
};

// Category labels — built dynamically from agents
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

function truncate(s, n = 120) {
    if (!s) return '';
    return s.length <= n ? s : s.slice(0, n - 1).trimEnd() + '…';
}

// catLabel defined above with CATEGORIES

function resolveTargets(entry, platforms) {
    if (!entry || entry.disabled) return [];
    const t = entry.targets || [];
    if (t.length === 1 && t[0] === 'all') return [...platforms];
    return t;
}

function showSnackbar(msg, ms = 3000) {
    const el = document.getElementById('snackbar');
    el.textContent = msg;
    el.classList.add('show');
    clearTimeout(el._timer);
    el._timer = setTimeout(() => el.classList.remove('show'), ms);
}

// ── Theme ──
function toggleTheme() {
    const html = document.documentElement;
    const next = html.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', next);
    localStorage.setItem('dotagen-theme', next);
}

function initTheme() {
    const saved = localStorage.getItem('dotagen-theme');
    if (saved) document.documentElement.setAttribute('data-theme', saved);
}

// ── Tab navigation ──
function switchTab(tab) {
    currentTab = tab;
    document.querySelectorAll('.tab').forEach(t => {
        const active = t.dataset.tab === tab;
        t.classList.toggle('active', active);
        t.setAttribute('aria-selected', active);
    });
    document.querySelectorAll('.tab-panel').forEach(p => {
        p.classList.toggle('active', p.id === 'tab-' + tab);
    });
    if (tab === 'agents') loadAgents();
    else if (tab === 'preview') loadPreviewOptions();
    else if (tab === 'status') loadStatus();
}

document.querySelectorAll('.tab').forEach(t => {
    t.addEventListener('click', () => switchTab(t.dataset.tab));
});

// ── Data loading ──
async function loadAll() {
    try {
        const [a, c, t, s] = await Promise.all([
            api('/api/agents'),
            api('/api/config'),
            api('/api/targets'),
            api('/api/status'),
        ]);
        agents = a || [];
        config = c || { targets: [], agents: {} };
        knownTargets = (t && t.targets) || config.targets || [];
        statusLinks = (s && s.symlinks) || [];
        updateBadges();
    } catch (e) {
        showSnackbar('Failed to load data: ' + e.message, 5000);
    }
}

function updateBadges() {
    document.getElementById('badge-agents').textContent = agents.length;
    const broken = statusLinks.filter(l => l.broken).length;
    const badge = document.getElementById('badge-status');
    if (broken > 0) {
        badge.textContent = broken;
        badge.classList.remove('hidden');
    } else {
        badge.classList.add('hidden');
    }
}

// ═══════════════════════════════════
// AGENTS TAB
// ═══════════════════════════════════
async function loadAgents() {
    await loadAll();
    renderCategoryChips();
    renderAgentsTable();
    renderPlatformHeaders();
}

function renderPlatformHeaders() {
    const th = document.getElementById('platforms-header');
    if (knownTargets.length === 0) {
        th.textContent = 'Platforms';
        return;
    }
    th.innerHTML = knownTargets.map(t =>
        `<span title="${esc(PLATFORM_NAMES[t] || t)}" style="display:inline-block;padding:0 4px;font-size:11px">${esc(PLATFORM_LABELS[t] || t)}</span>`
    ).join('');
}

function renderCategoryChips() {
    const counts = {};
    agents.forEach(a => {
        const cats = agentCategories(a);
        if (cats.length === 0) counts['uncategorized'] = (counts['uncategorized'] || 0) + 1;
        else cats.forEach(c => counts[c] = (counts[c] || 0) + 1);
    });

    const container = document.getElementById('category-chips');
    let html = `<button class="chip ${categoryFilter === 'all' ? 'active' : ''}" onclick="setCategory('all')">All</button>`;
    Object.entries(counts).sort((a, b) => b[1] - a[1]).forEach(([cat, n]) => {
        html += `<button class="chip ${categoryFilter === cat ? 'active' : ''}" onclick="setCategory('${esc(cat)}')">${esc(catLabel(cat))} <span style='opacity:.5'>${n}</span></button>`;
    });
    container.innerHTML = html;
}

function setCategory(cat) {
    categoryFilter = cat;
    renderCategoryChips();
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
    const tbody = document.getElementById('agents-tbody');
    const empty = document.getElementById('agents-empty');
    const noResults = document.getElementById('agents-no-results');
    const tableWrap = document.querySelector('.agents-table-wrap');
    const countEl = document.getElementById('agent-count');

    if (agents.length === 0) {
        empty.classList.remove('hidden');
        noResults.classList.add('hidden');
        tableWrap.classList.add('hidden');
        countEl.textContent = '';
        return;
    }

    empty.classList.add('hidden');

    if (filtered.length === 0) {
        noResults.classList.remove('hidden');
        tableWrap.classList.add('hidden');
        countEl.textContent = '';
        return;
    }

    noResults.classList.add('hidden');
    tableWrap.classList.remove('hidden');
    countEl.textContent = `${filtered.length} of ${agents.length}`;

    tbody.innerHTML = filtered.map(a => {
        const entry = config.agents?.[a.name];
        const active = resolveTargets(entry, knownTargets);
        const checked = selectedAgents.has(a.name);

        const dots = knownTargets.map(t => {
            const on = active.includes(t);
            const label = PLATFORM_LABELS[t] || t.slice(0, 2).toUpperCase();
            return `<button class="pdot ${on ? 'pdot-on' : 'pdot-off'}" data-platform="${esc(t)}" title="${on ? 'Disable' : 'Enable'} ${esc(PLATFORM_NAMES[t] || t)}" onclick="event.stopPropagation();togglePlatform('${esc(a.name)}','${esc(t)}',${!on})">${label}</button>`;
        }).join('');

        const cats = agentCategories(a);
        const catBadges = cats.length
            ? cats.map(c => `<span class="agent-cat-badge">${esc(catLabel(c))}</span>`).join(' ')
            : '<span class="agent-cat-badge">—</span>';

        return `<tr class="${checked ? 'row-selected' : ''}" onclick="viewAgent('${esc(a.name)}')">
            <td class="col-check" onclick="event.stopPropagation()"><input type="checkbox" ${checked ? 'checked' : ''} onchange="toggleAgentSelect('${esc(a.name)}', this.checked)"></td>
            <td><span class="agent-name">${esc(a.name)}</span></td>
            <td title="${esc(a.description || '')}"><span class="agent-desc">${esc(truncate(a.description, 100))}</span></td>
            <td>${catBadges}</td>
            <td><div class="platform-dots">${dots}</div></td>
            <td><button class="row-actions-btn" title="Edit" onclick="event.stopPropagation();editAgent('${esc(a.name)}')">✎</button></td>
        </tr>`;
    }).join('');

    updateBulkBar();
    const selectAll = document.getElementById('select-all');
    if (selectAll) selectAll.checked = filtered.length > 0 && filtered.every(a => selectedAgents.has(a.name));
}

// ── Platform toggle ──
async function togglePlatform(agentName, platform, enable) {
    const agentsMap = JSON.parse(JSON.stringify(config.agents || {}));
    const entry = agentsMap[agentName] || { targets: [], disabled: false };
    let current = resolveTargets(entry, knownTargets);

    if (enable) {
        if (!current.includes(platform)) current.push(platform);
    } else {
        current = current.filter(t => t !== platform);
    }
    entry.targets = current;
    entry.disabled = false;
    agentsMap[agentName] = entry;

    try {
        await api('/api/config', {
            method: 'PUT',
            body: JSON.stringify({ targets: config.targets, agents: agentsMap }),
        });
        config.agents = agentsMap;
        renderAgentsTable();
        showSnackbar(`${agentName}: ${PLATFORM_NAMES[platform] || platform} ${enable ? 'enabled' : 'disabled'}`);
    } catch (e) {
        showSnackbar('Failed: ' + e.message, 4000);
    }
}

// ═══════════════════════════════════
// BULK SELECTION
// ═══════════════════════════════════
function toggleAgentSelect(name, checked) {
    if (checked) selectedAgents.add(name);
    else selectedAgents.delete(name);
    renderAgentsTable();
}

function toggleSelectAll(checked) {
    const filtered = getFilteredAgents();
    if (checked) filtered.forEach(a => selectedAgents.add(a.name));
    else filtered.forEach(a => selectedAgents.delete(a.name));
    renderAgentsTable();
}

function clearSelection() {
    selectedAgents.clear();
    renderAgentsTable();
}

function updateBulkBar() {
    const bar = document.getElementById('bulk-bar');
    if (selectedAgents.size === 0) {
        bar.classList.add('hidden');
        return;
    }
    bar.classList.remove('hidden');
    document.getElementById('bulk-count').textContent = `${selectedAgents.size} selected`;

    document.getElementById('bulk-enable-targets').innerHTML = knownTargets.map(t =>
        `<button class="bulk-target-btn" onclick="bulkTogglePlatform('${esc(t)}', true)">${esc(PLATFORM_LABELS[t] || t)}</button>`
    ).join('');

    document.getElementById('bulk-disable-targets').innerHTML = knownTargets.map(t =>
        `<button class="bulk-target-btn bulk-disable" onclick="bulkTogglePlatform('${esc(t)}', false)">${esc(PLATFORM_LABELS[t] || t)}</button>`
    ).join('');
}

async function bulkTogglePlatform(platform, enable) {
    const agentsMap = JSON.parse(JSON.stringify(config.agents || {}));
    for (const name of selectedAgents) {
        const entry = agentsMap[name] || { targets: [], disabled: false };
        let current = resolveTargets(entry, knownTargets);
        if (enable) {
            if (!current.includes(platform)) current.push(platform);
        } else {
            current = current.filter(t => t !== platform);
        }
        entry.targets = current;
        entry.disabled = false;
        agentsMap[name] = entry;
    }
    try {
        await api('/api/config', {
            method: 'PUT',
            body: JSON.stringify({ targets: config.targets, agents: agentsMap }),
        });
        config.agents = agentsMap;
        renderAgentsTable();
        showSnackbar(`${selectedAgents.size} agents: ${PLATFORM_NAMES[platform] || platform} ${enable ? 'enabled' : 'disabled'}`);
    } catch (e) {
        showSnackbar('Bulk update failed: ' + e.message, 4000);
    }
}

async function bulkEnableAll() {
    const agentsMap = JSON.parse(JSON.stringify(config.agents || {}));
    for (const name of selectedAgents) {
        const entry = agentsMap[name] || { targets: [], disabled: false };
        entry.targets = [...knownTargets];
        entry.disabled = false;
        agentsMap[name] = entry;
    }
    try {
        await api('/api/config', {
            method: 'PUT',
            body: JSON.stringify({ targets: config.targets, agents: agentsMap }),
        });
        config.agents = agentsMap;
        renderAgentsTable();
        showSnackbar(`${selectedAgents.size} agents: all platforms enabled`);
    } catch (e) {
        showSnackbar('Bulk update failed: ' + e.message, 4000);
    }
}

async function bulkDisableAll() {
    const agentsMap = JSON.parse(JSON.stringify(config.agents || {}));
    for (const name of selectedAgents) {
        const entry = agentsMap[name] || { targets: [], disabled: false };
        entry.targets = [];
        entry.disabled = false;
        agentsMap[name] = entry;
    }
    try {
        await api('/api/config', {
            method: 'PUT',
            body: JSON.stringify({ targets: config.targets, agents: agentsMap }),
        });
        config.agents = agentsMap;
        renderAgentsTable();
        showSnackbar(`${selectedAgents.size} agents: all platforms disabled`);
    } catch (e) {
        showSnackbar('Bulk update failed: ' + e.message, 4000);
    }
}

// ═══════════════════════════════════
// AGENT PANEL (view / create / edit)
// ═══════════════════════════════════
function openPanel() {
    document.getElementById('panel-overlay').classList.remove('hidden');
    document.getElementById('agent-panel').classList.add('open');
}

function closePanel() {
    document.getElementById('panel-overlay').classList.add('hidden');
    document.getElementById('agent-panel').classList.remove('open');
    panelMode = null;
    panelAgent = null;
}

async function viewAgent(name) {
    try {
        const a = await api('/api/agents/' + name);
        panelAgent = a;
        panelMode = 'view';
        const entry = config.agents?.[a.name];
        const active = resolveTargets(entry, knownTargets);

        document.getElementById('panel-title').textContent = a.name;
        document.getElementById('panel-body').innerHTML = `
            <div class="detail-meta">
                ${a.description ? `<div class="detail-row"><span class="detail-label">Description</span><span class="detail-value">${esc(a.description)}</span></div>` : ''}
                <div class="detail-row"><span class="detail-label">Categories</span><span class="detail-value">${agentCategories(a).map(c => esc(catLabel(c))).join(', ') || '—'}</span></div>
                <div class="detail-row"><span class="detail-label">Platforms</span><span class="detail-value">${active.length ? active.map(t => esc(PLATFORM_NAMES[t] || t)).join(', ') : 'None'}</span></div>
            </div>
            <div class="detail-content">${esc(a.content || '(empty)')}</div>
        `;
        document.getElementById('panel-footer').innerHTML = `
            <button class="btn btn-ghost" onclick="closePanel()">Close</button>
            <button class="btn btn-secondary" onclick="editAgent('${esc(a.name)}')">Edit</button>
            <button class="btn btn-danger-text btn-sm" onclick="deleteAgent('${esc(a.name)}')">Delete</button>
        `;
        openPanel();
    } catch (e) {
        showSnackbar('Failed to load: ' + e.message, 4000);
    }
}

function showCreateAgent() {
    panelMode = 'create';
    panelAgent = null;
    document.getElementById('panel-title').textContent = 'New Agent';
    renderAgentForm({ name: '', description: '', categories: [], content: '', targets: [...knownTargets] });
    document.getElementById('panel-footer').innerHTML = `
        <button class="btn btn-ghost" onclick="closePanel()">Cancel</button>
        <button class="btn btn-primary" onclick="submitCreateAgent()">Create</button>
    `;
    openPanel();
    setTimeout(() => {
        const inp = document.getElementById('form-name');
        if (inp) inp.focus();
    }, 200);
}

async function editAgent(name) {
    try {
        const a = await api('/api/agents/' + name);
        panelMode = 'edit';
        panelAgent = a;
        const entry = config.agents?.[a.name];
        const active = resolveTargets(entry, knownTargets);

        document.getElementById('panel-title').textContent = 'Edit: ' + a.name;
        renderAgentForm({
            name: a.name,
            description: a.description || '',
            categories: agentCategories(a),
            content: a.content || '',
            targets: active,
            isEdit: true,
        });
        document.getElementById('panel-footer').innerHTML = `
            <button class="btn btn-ghost" onclick="closePanel()">Cancel</button>
            <button class="btn btn-primary" onclick="submitEditAgent('${esc(a.name)}')">Save</button>
        `;
        openPanel();
    } catch (e) {
        showSnackbar('Failed to load: ' + e.message, 4000);
    }
}

function renderAgentForm({ name, description, categories = [], content, targets, isEdit = false }) {
    const existingCats = allCategories();
    const catChips = existingCats.map(c => {
        const checked = categories.includes(c);
        return `<label class="target-chip ${checked ? 'checked' : ''}" onclick="this.classList.toggle('checked')">
            <input type="checkbox" name="form-category" value="${esc(c)}" ${checked ? 'checked' : ''}>
            ${esc(catLabel(c))}
        </label>`;
    }).join('');

    const targetChips = knownTargets.map(t => {
        const checked = targets.includes(t);
        return `<label class="target-chip ${checked ? 'checked' : ''}" onclick="this.classList.toggle('checked')">
            <input type="checkbox" name="form-target" value="${esc(t)}" ${checked ? 'checked' : ''}>
            ${esc(PLATFORM_NAMES[t] || t)}
        </label>`;
    }).join('');

    document.getElementById('panel-body').innerHTML = `
        <div class="form-group">
            <label class="form-label" for="form-name">Name</label>
            <input class="form-input" id="form-name" value="${esc(name)}" placeholder="my-agent" ${isEdit ? 'disabled' : ''}>
            <div class="form-error hidden" id="form-name-error"></div>
        </div>
        <div class="form-group">
            <label class="form-label" for="form-desc">Description</label>
            <input class="form-input" id="form-desc" value="${esc(description)}" placeholder="Short description…">
        </div>
        <div class="form-group">
            <label class="form-label">Categories</label>
            <div class="form-targets">${catChips}</div>
            <div style="display:flex;gap:6px;margin-top:8px">
                <input class="form-input" id="form-new-cat" placeholder="New category…" style="flex:1">
                <button class="btn btn-secondary btn-sm" type="button" onclick="addCustomCategory()">Add</button>
            </div>
            <div class="form-hint">Select existing or type a new category name</div>
        </div>
        <div class="form-group">
            <label class="form-label">Platforms</label>
            <div class="form-targets">${targetChips}</div>
        </div>
        <div class="form-group">
            <label class="form-label" for="form-content">Prompt (Markdown)</label>
            <textarea class="form-textarea" id="form-content" placeholder="# Agent Name\n\n## Role\n\nDescribe what this agent does…">${esc(content)}</textarea>
        </div>
    `;
}

function addCustomCategory() {
    const input = document.getElementById('form-new-cat');
    const val = input.value.trim().toLowerCase().replace(/[^a-z0-9-]/g, '-').replace(/-+/g, '-');
    if (!val) return;
    // Check if already exists
    const existing = document.querySelector(`input[name="form-category"][value="${val}"]`);
    if (existing) {
        existing.checked = true;
        existing.closest('.target-chip').classList.add('checked');
        input.value = '';
        return;
    }
    // Add new chip
    const container = input.closest('.form-group').querySelector('.form-targets');
    const chip = document.createElement('label');
    chip.className = 'target-chip checked';
    chip.onclick = function() { this.classList.toggle('checked'); };
    chip.innerHTML = `<input type="checkbox" name="form-category" value="${esc(val)}" checked> ${esc(val)}`;
    container.appendChild(chip);
    input.value = '';
}

function getFormData() {
    const name = document.getElementById('form-name').value.trim();
    const description = document.getElementById('form-desc').value.trim();
    const catChecks = document.querySelectorAll('input[name="form-category"]:checked');
    const category = Array.from(catChecks).map(c => c.value).join(',');
    const content = document.getElementById('form-content').value;
    const checks = document.querySelectorAll('input[name="form-target"]:checked');
    const targets = Array.from(checks).map(c => c.value);
    return { name, description, category, content, targets };
}

async function submitCreateAgent() {
    const data = getFormData();
    const errEl = document.getElementById('form-name-error');
    errEl.classList.add('hidden');
    document.getElementById('form-name').classList.remove('has-error');

    if (!data.name) {
        document.getElementById('form-name').classList.add('has-error');
        errEl.textContent = 'Name is required';
        errEl.classList.remove('hidden');
        return;
    }
    if (data.targets.length === 0) {
        showSnackbar('Select at least one platform');
        return;
    }

    try {
        await api('/api/agents', {
            method: 'POST',
            body: JSON.stringify(data),
        });
        closePanel();
        showSnackbar(`Agent "${data.name}" created`);
        loadAgents();
    } catch (e) {
        showSnackbar('Create failed: ' + e.message, 5000);
    }
}

async function submitEditAgent(name) {
    const data = getFormData();
    if (data.targets.length === 0) {
        showSnackbar('Select at least one platform');
        return;
    }

    try {
        await api('/api/agents/' + name, {
            method: 'PUT',
            body: JSON.stringify({
                content: data.content,
                description: data.description,
                category: data.category,
                targets: data.targets,
            }),
        });
        closePanel();
        showSnackbar(`Agent "${name}" saved`);
        loadAgents();
    } catch (e) {
        showSnackbar('Save failed: ' + e.message, 5000);
    }
}

async function deleteAgent(name) {
    const ok = await showConfirm(
        `Delete "${name}"?`,
        'This removes the agent file and its config entry. This cannot be undone.'
    );
    if (!ok) return;

    try {
        await api('/api/agents/' + name, { method: 'DELETE' });
        closePanel();
        showSnackbar(`Agent "${name}" deleted`);
        loadAgents();
    } catch (e) {
        showSnackbar('Delete failed: ' + e.message, 5000);
    }
}

// ═══════════════════════════════════
// PREVIEW TAB
// ═══════════════════════════════════
async function loadPreviewOptions() {
    try {
        if (agents.length === 0) await loadAll();
        const agentSel = document.getElementById('preview-agent');
        const targetSel = document.getElementById('preview-target');
        const prevAgent = agentSel.value;
        const prevTarget = targetSel.value;

        agentSel.innerHTML = '<option value="">Select agent…</option>' +
            agents.map(a => `<option value="${esc(a.name)}" ${a.name === prevAgent ? 'selected' : ''}>${esc(a.name)}</option>`).join('');
        targetSel.innerHTML = '<option value="">Select platform…</option>' +
            knownTargets.map(t => `<option value="${esc(t)}" ${t === prevTarget ? 'selected' : ''}>${esc(PLATFORM_NAMES[t] || t)}</option>`).join('');
    } catch (e) {
        showSnackbar('Failed to load options: ' + e.message, 4000);
    }
}

let previewDebounce = null;
function autoPreview() {
    clearTimeout(previewDebounce);
    previewDebounce = setTimeout(loadPreview, 150);
}

async function loadPreview() {
    const agent = document.getElementById('preview-agent').value;
    const target = document.getElementById('preview-target').value;
    const output = document.getElementById('preview-output');
    const copyBtn = document.getElementById('preview-copy-btn');

    if (!agent || !target) {
        output.textContent = 'Select an agent and platform to preview rendered output.';
        copyBtn.disabled = true;
        return;
    }

    output.textContent = 'Loading…';
    try {
        const res = await api('/api/preview/' + agent + '/' + target);
        output.textContent = res.content || '(empty output)';
        copyBtn.disabled = false;
    } catch (e) {
        output.textContent = 'Error: ' + e.message;
        copyBtn.disabled = true;
    }
}

function copyPreview() {
    const text = document.getElementById('preview-output').textContent;
    navigator.clipboard.writeText(text).then(() => showSnackbar('Copied'));
}

// ═══════════════════════════════════
// STATUS TAB
// ═══════════════════════════════════
async function loadStatus() {
    try {
        if (statusLinks.length === 0 && agents.length === 0) await loadAll();
        renderStatusChips();
        renderStatusList();
    } catch (e) {
        showSnackbar('Failed to load status: ' + e.message, 4000);
    }
}

function renderStatusChips() {
    const broken = statusLinks.filter(l => l.broken).length;
    const healthy = statusLinks.length - broken;
    const container = document.getElementById('status-chips');
    container.innerHTML = `
        <button class="chip ${statusFilter === 'all' ? 'active' : ''}" onclick="setStatusFilter('all')">All ${statusLinks.length}</button>
        <button class="chip ${statusFilter === 'healthy' ? 'active' : ''}" onclick="setStatusFilter('healthy')">Healthy ${healthy}</button>
        <button class="chip ${statusFilter === 'broken' ? 'active' : ''}" onclick="setStatusFilter('broken')">Broken ${broken}</button>
    `;
}

function setStatusFilter(f) {
    statusFilter = f;
    renderStatusChips();
    renderStatusList();
}

function renderStatusList() {
    const list = document.getElementById('status-list');
    const empty = document.getElementById('status-empty');

    let filtered = statusLinks;
    if (statusFilter === 'healthy') filtered = filtered.filter(l => !l.broken);
    else if (statusFilter === 'broken') filtered = filtered.filter(l => l.broken);

    if (filtered.length === 0) {
        list.innerHTML = '';
        empty.classList.remove('hidden');
        return;
    }
    empty.classList.add('hidden');

    list.innerHTML = filtered.map(l => `
        <div class="status-row">
            <span class="status-dot ${l.broken ? 'status-dot-err' : 'status-dot-ok'}"></span>
            <span class="status-agent">${esc(l.agent)}</span>
            <span class="status-platform">${esc(PLATFORM_NAMES[l.platform] || l.platform)}</span>
            <span class="status-path">${esc(l.path)}</span>
            <div class="status-row-actions">
                <button class="btn btn-ghost btn-sm" onclick="quickPreview('${esc(l.agent)}','${esc(l.platform)}')">Preview</button>
            </div>
        </div>
    `).join('');
}

function quickPreview(agent, platform) {
    switchTab('preview');
    setTimeout(() => {
        document.getElementById('preview-agent').value = agent;
        document.getElementById('preview-target').value = platform;
        loadPreview();
    }, 100);
}

// ═══════════════════════════════════
// SYNC / CLEAN
// ═══════════════════════════════════
async function triggerSync() {
    const ok = await showConfirm(
        'Sync all agents?',
        'This generates platform-specific files for all enabled agents.'
    );
    if (!ok) return;

    const btn = document.getElementById('btn-sync');
    btn.classList.add('btn-loading');
    btn.disabled = true;

    try {
        const res = await api('/api/sync', { method: 'POST' });
        showSnackbar(`Synced ${res.synced || 0} agents`);
        await loadAll();
        if (currentTab === 'status') renderStatusList();
        if (currentTab === 'agents') renderAgentsTable();
    } catch (e) {
        showSnackbar('Sync failed: ' + e.message, 5000);
    } finally {
        btn.classList.remove('btn-loading');
        btn.disabled = false;
    }
}

async function triggerClean() {
    const ok = await showConfirm(
        'Remove generated files?',
        'Source agent files will be kept. Only generated output and symlinks will be deleted.'
    );
    if (!ok) return;

    try {
        const res = await api('/api/clean', { method: 'POST' });
        showSnackbar(`Removed ${res.removed || 0} files`);
        await loadAll();
        if (currentTab === 'status') { renderStatusChips(); renderStatusList(); }
    } catch (e) {
        showSnackbar('Clean failed: ' + e.message, 5000);
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
            <button class="btn btn-ghost" id="confirm-no">Cancel</button>
            <button class="btn btn-primary" id="confirm-yes">Confirm</button>
        `;
        modal.classList.remove('hidden');

        const cleanup = (result) => {
            modal.classList.add('hidden');
            resolve(result);
        };

        document.getElementById('confirm-no').onclick = () => cleanup(false);
        document.getElementById('confirm-yes').onclick = () => cleanup(true);
        modal.onclick = (e) => { if (e.target === modal) cleanup(false); };
    });
}

// ═══════════════════════════════════
// SEARCH
// ═══════════════════════════════════
const searchInput = document.getElementById('global-search');
searchInput.addEventListener('input', (e) => {
    searchQuery = e.target.value.trim().toLowerCase();
    if (currentTab === 'agents') {
        renderAgentsTable();
    } else if (searchQuery) {
        switchTab('agents');
    }
});

// ═══════════════════════════════════
// KEYBOARD SHORTCUTS
// ═══════════════════════════════════
document.addEventListener('keydown', (e) => {
    const active = document.activeElement;
    const typing = active && (active.tagName === 'INPUT' || active.tagName === 'TEXTAREA' || active.tagName === 'SELECT');

    if (e.key === 'Escape') {
        if (!document.getElementById('modal').classList.contains('hidden')) {
            document.getElementById('modal').classList.add('hidden');
        } else if (document.getElementById('agent-panel').classList.contains('open')) {
            closePanel();
        }
        return;
    }

    if (typing || e.metaKey || e.ctrlKey) return;

    if (e.key === '/') { e.preventDefault(); searchInput.focus(); return; }
    if (e.key === 'n') { e.preventDefault(); showCreateAgent(); return; }
    if (e.key === '1') { switchTab('agents'); return; }
    if (e.key === '2') { switchTab('preview'); return; }
    if (e.key === '3') { switchTab('status'); return; }
});

// ═══════════════════════════════════
// INIT
// ═══════════════════════════════════
initTheme();
loadAgents();
