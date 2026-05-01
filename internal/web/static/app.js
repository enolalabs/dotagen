const API = '';
let allAgents = [];
let selectedCategory = 'all';
let searchQuery = '';
let selectedAgent = null;
let currentPage = 'dashboard';
let lastFocusedElement = null;

const CATEGORIES = {
    'core-dev': 'Core Development',
    'languages': 'Language Specialists',
    'infrastructure': 'Infrastructure',
    'quality-security': 'Quality & Security',
    'data-ai': 'Data & AI',
    'devex': 'Developer Experience',
    'specialized': 'Specialized Domains',
    'business': 'Business & Product',
    'orchestration': 'Meta-Orchestration',
    'research': 'Research & Analysis',
};

const PAGE_TITLES = {
    dashboard: 'Dashboard',
    agents: 'Agents',
    targets: 'Targets',
    preview: 'Preview',
    status: 'Status',
};

async function api(path, opts = {}) {
    const res = await fetch(API + path, {
        headers: { 'Content-Type': 'application/json' },
        ...opts,
    });
    if (!res.ok) {
        const err = await res.json().catch(() => ({ error: res.statusText }));
        throw new Error(err.error || res.statusText);
    }
    return res.json();
}

function escapeHtml(str) {
    if (!str) return '';
    return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

function showLoading(containerId) {
    const el = document.getElementById(containerId);
    if (el) el.innerHTML = '<div class="loading-state"><div class="loading-spinner"></div><p>Loading...</p></div>';
}

function showError(containerId, message) {
    const el = document.getElementById(containerId);
    if (el) el.innerHTML = `<div class="empty-state"><div class="empty-icon">&#9888;</div><h3>Error</h3><p>${escapeHtml(message)}</p></div>`;
}

function renderMarkdown(text) {
    if (!text) return '';
    let html = escapeHtml(text);

    html = html.replace(/```(\w*)\n([\s\S]*?)```/g, (_, lang, code) => {
        return `<pre><code class="lang-${lang}">${code.trim()}</code></pre>`;
    });

    html = html.replace(/`([^`]+)`/g, '<code>$1</code>');

    html = html.replace(/^#### (.+)$/gm, '<h4>$1</h4>');
    html = html.replace(/^### (.+)$/gm, '<h3>$1</h3>');
    html = html.replace(/^## (.+)$/gm, '<h2>$1</h2>');
    html = html.replace(/^# (.+)$/gm, '<h1>$1</h1>');

    html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
    html = html.replace(/\*(.+?)\*/g, '<em>$1</em>');

    html = html.replace(/^- (.+)$/gm, '<li>$1</li>');

    html = html.replace(/^(\d+)\. (.+)$/gm, '<li>$2</li>');

    html = html.replace(/((?:<li>.*<\/li>\n?)+)/g, (match) => {
        return '<ul>' + match + '</ul>';
    });

    html = html.replace(/\n{2,}/g, '</p><p>');
    html = '<p>' + html + '</p>';

    html = html.replace(/<p>\s*<(h[1-4]|pre|ul|ol)/g, '<$1');
    html = html.replace(/<\/(h[1-4]|pre|ul|ol)>\s*<\/p>/g, '</$1>');
    html = html.replace(/<p>\s*<\/p>/g, '');

    return html;
}

function showSnackbar(message, duration = 4000) {
    const el = document.getElementById('snackbar');
    el.textContent = message;
    el.classList.add('show');
    setTimeout(() => el.classList.remove('show'), duration);
}

function toggleNav() {
    document.getElementById('app').classList.toggle('nav-open');
}

function navigateTo(page) {
    currentPage = page;
    document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
    const navItem = document.querySelector(`.nav-item[data-page="${page}"]`);
    if (navItem) navItem.classList.add('active');

    document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
    document.getElementById('page-' + page).classList.add('active');

    document.getElementById('topbar-title').textContent = PAGE_TITLES[page] || page;
    document.getElementById('app').classList.remove('nav-open');

    document.getElementById('main-content').scrollTop = 0;

    loadPage(page);
}

document.querySelectorAll('.nav-item').forEach(link => {
    link.addEventListener('click', e => {
        e.preventDefault();
        navigateTo(link.dataset.page);
    });
});

function loadPage(page) {
    switch (page) {
        case 'dashboard': loadDashboard(); break;
        case 'agents': loadAgents(); break;
        case 'targets': loadTargets(); break;
        case 'preview': loadPreviewOptions(); break;
        case 'status': loadStatus(); break;
    }
}

async function loadDashboard() {
    showLoading('dashboard-stats');
    try {
        const [agents, config, status] = await Promise.all([
            api('/api/agents'),
            api('/api/config'),
            api('/api/status'),
        ]);
        allAgents = agents || [];

        if (allAgents.length === 0) {
            document.getElementById('dashboard-stats').innerHTML = '';
            document.getElementById('dashboard-actions').innerHTML = `
                <div class="empty-state" style="width:100%">
                    <div class="empty-icon">&#128640;</div>
                    <h3>No agents yet</h3>
                    <p>Create your first agent or run <code>dotagen init</code> to get started with built-in agents.</p>
                    <button class="m3-btn m3-btn-filled" onclick="navigateTo('agents')" style="margin-top:16px">Create Agent</button>
                </div>`;
            return;
        }

        document.getElementById('dashboard-stats').innerHTML = `
            <div class="stat-card">
                <div class="stat-value" id="stat-agents">-</div>
                <div class="stat-label">Total Agents</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="stat-targets">-</div>
                <div class="stat-label">Targets</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="stat-synced">-</div>
                <div class="stat-label">Synced</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="stat-categories">-</div>
                <div class="stat-label">Categories</div>
            </div>`;
        document.getElementById('dashboard-actions').innerHTML = `
            <button class="m3-btn m3-btn-filled" onclick="triggerSync()">Sync All</button>
            <button class="m3-btn m3-btn-outlined" onclick="triggerClean()">Clean All</button>
            <button class="m3-btn m3-btn-tonal" onclick="navigateTo('agents')">Manage Agents</button>`;

        document.getElementById('stat-agents').textContent = allAgents.length;
        document.getElementById('stat-targets').textContent = (config.targets || []).length;
        document.getElementById('stat-synced').textContent = (status.symlinks || []).filter(s => !s.broken).length;

        const cats = new Set(allAgents.map(a => a.category).filter(Boolean));
        document.getElementById('stat-categories').textContent = cats.size;
    } catch (e) {
        showError('dashboard-stats', e.message);
    }
}

async function loadAgents() {
    showLoading('agents-grid');
    try {
        const agents = await api('/api/agents');
        allAgents = agents || [];
        renderCategoryChips();
        renderAgentsGrid();
    } catch (e) {
        showError('agents-grid', e.message);
    }
}

function renderCategoryChips() {
    const container = document.getElementById('category-chips');
    const catCounts = {};
    allAgents.forEach(a => {
        const c = a.category || 'uncategorized';
        catCounts[c] = (catCounts[c] || 0) + 1;
    });

    let html = `<button class="chip ${selectedCategory === 'all' ? 'selected' : ''}" onclick="filterCategory('all')">
        <span class="chip-check">&#10003;</span> All (${allAgents.length})
    </button>`;

    const sortedCats = Object.entries(catCounts).sort((a, b) => b[1] - a[1]);
    for (const [cat, count] of sortedCats) {
        const label = CATEGORIES[cat] || cat;
        html += `<button class="chip ${selectedCategory === cat ? 'selected' : ''}" onclick="filterCategory('${escapeHtml(cat)}')">
            <span class="chip-check">&#10003;</span> ${escapeHtml(label)} (${count})
        </button>`;
    }

    container.innerHTML = html;
}

function filterCategory(cat) {
    selectedCategory = cat;
    renderCategoryChips();
    renderAgentsGrid();
}

function handleSearch(query) {
    searchQuery = query.toLowerCase().trim();
    if (currentPage === 'agents') {
        renderAgentsGrid();
    } else if (searchQuery) {
        navigateTo('agents');
    }
}

function getFilteredAgents() {
    let filtered = allAgents;
    if (selectedCategory !== 'all') {
        filtered = filtered.filter(a => a.category === selectedCategory);
    }
    if (searchQuery) {
        filtered = filtered.filter(a =>
            a.name.toLowerCase().includes(searchQuery) ||
            (a.description || '').toLowerCase().includes(searchQuery) ||
            (a.category || '').toLowerCase().includes(searchQuery)
        );
    }
    return filtered;
}

function renderAgentsGrid() {
    const grid = document.getElementById('agents-grid');
    const filtered = getFilteredAgents();

    document.getElementById('agent-count').textContent = `${filtered.length} of ${allAgents.length}`;
    document.getElementById('agents-subtitle').textContent = `${allAgents.length} agents across ${new Set(allAgents.map(a => a.category).filter(Boolean)).size} categories`;

    if (filtered.length === 0) {
        grid.innerHTML = `<div class="empty-state">
            <div class="empty-icon">&#128269;</div>
            <h3>No agents found</h3>
            <p>Try adjusting your search or category filter.</p>
        </div>`;
        return;
    }

    grid.innerHTML = filtered.map(a => {
        const catLabel = CATEGORIES[a.category] || a.category || '';
        const isSelected = selectedAgent && selectedAgent.name === a.name;
        return `<div class="agent-card ${isSelected ? 'selected' : ''}" tabindex="0" onclick="selectAgent('${escapeHtml(a.name)}')" onkeydown="if(event.key==='Enter')selectAgent('${escapeHtml(a.name)}')">
            <div class="agent-card-header">
                <h3>${escapeHtml(a.name)}</h3>
                ${catLabel ? `<span class="agent-card-category">${escapeHtml(catLabel)}</span>` : ''}
            </div>
            <div class="agent-card-desc">${escapeHtml(a.description || 'No description')}</div>
        </div>`;
    }).join('');
}

async function selectAgent(name) {
    try {
        const a = await api('/api/agents/' + name);
        selectedAgent = a;
        openDetailPanel(a);
        renderAgentsGrid();
    } catch (e) {
        showSnackbar('Failed to load agent: ' + e.message, 6000);
    }
}

function openDetailPanel(agent) {
    const app = document.getElementById('app');
    app.classList.add('detail-open');

    document.getElementById('detail-title').textContent = agent.name;

    const catLabel = CATEGORIES[agent.category] || agent.category || 'Uncategorized';
    const body = document.getElementById('detail-body');

    body.innerHTML = `
        <div class="detail-meta">
            ${agent.description ? `<div class="detail-meta-row">
                <span class="meta-label">Description</span>
                <span class="meta-value">${escapeHtml(agent.description)}</span>
            </div>` : ''}
            <div class="detail-meta-row">
                <span class="meta-label">Category</span>
                <span class="meta-value">${escapeHtml(catLabel)}</span>
            </div>
            ${agent.frontmatter && agent.frontmatter.mode ? `<div class="detail-meta-row">
                <span class="meta-label">Mode</span>
                <span class="meta-value">${escapeHtml(agent.frontmatter.mode)}</span>
            </div>` : ''}
            ${agent.frontmatter && agent.frontmatter.targets ? `<div class="detail-meta-row">
                <span class="meta-label">Targets</span>
                <span class="meta-value">${escapeHtml(agent.frontmatter.targets)}</span>
            </div>` : ''}
        </div>
        <div class="md-content">${renderMarkdown(agent.content || '')}</div>
    `;

    document.getElementById('detail-edit-btn').onclick = () => editAgent(agent.name);
    document.getElementById('detail-delete-btn').onclick = () => deleteAgent(agent.name);
}

function closeDetail() {
    document.getElementById('app').classList.remove('detail-open');
    selectedAgent = null;
    renderAgentsGrid();
}

async function showCreateAgent() {
    let targets = [];
    try {
        const res = await api('/api/targets');
        targets = res.targets || [];
    } catch (e) { console.error(e); }

    let targetCheckboxes = targets.map(t =>
        `<label class="checkbox-label"><input type="checkbox" name="agent-target" value="${escapeHtml(t)}" checked> ${escapeHtml(t)}</label>`
    ).join('');

    const categories = Object.entries(CATEGORIES);
    let categoryOptions = '<option value="">None</option>' +
        categories.map(([id, label]) => `<option value="${id}">${escapeHtml(label)}</option>`).join('');

    document.getElementById('modal-title').textContent = 'Create Agent';
    document.getElementById('modal-body').innerHTML = `
        <label for="agent-name">Agent Name</label>
        <input type="text" class="m3-input" id="agent-name" placeholder="my-agent" aria-required="true">
        <div class="m3-input-error-msg" id="agent-name-error" style="display:none"></div>
        <label for="agent-description">Description</label>
        <input type="text" class="m3-input" id="agent-description" placeholder="Short description of the agent">
        <label for="agent-category">Category</label>
        <select class="m3-select" id="agent-category" style="width:100%">${categoryOptions}</select>
        <label>Targets</label>
        <div class="target-checkboxes">${targetCheckboxes}</div>
        <label for="agent-content">Content (Markdown)</label>
        <textarea class="m3-textarea" id="agent-content" placeholder="# My Agent\n\n## Role\n\nDescribe the agent's role.\n\n## Guidelines\n\n- Guideline 1\n- Guideline 2"></textarea>
    `;
    document.getElementById('modal-actions').innerHTML = `
        <button class="m3-btn m3-btn-text" onclick="closeModal()">Cancel</button>
        <button class="m3-btn m3-btn-filled" onclick="createAgent()">Create</button>
    `;
    openModal();
}

function clearValidation() {
    document.querySelectorAll('.m3-input-error').forEach(el => el.classList.remove('m3-input-error'));
    document.querySelectorAll('.m3-input-error-msg').forEach(el => { el.style.display = 'none'; el.textContent = ''; });
}

function showFieldError(fieldId, message) {
    const field = document.getElementById(fieldId);
    const errorEl = document.getElementById(fieldId + '-error');
    if (field) field.classList.add('m3-input-error');
    if (errorEl) {
        errorEl.textContent = message;
        errorEl.style.display = 'block';
    }
}

async function createAgent() {
    clearValidation();
    const name = document.getElementById('agent-name').value.trim();
    const description = document.getElementById('agent-description').value.trim();
    const category = document.getElementById('agent-category').value;
    const content = document.getElementById('agent-content').value;

    if (!name) {
        showFieldError('agent-name', 'Name is required');
        return;
    }
    const cleanName = name.replace(/^da-/, '');
    if (!/^[a-zA-Z0-9][a-zA-Z0-9_-]*$/.test(cleanName)) {
        showFieldError('agent-name', 'Name must contain only alphanumeric characters, hyphens, and underscores');
        return;
    }

    const checked = document.querySelectorAll('input[name="agent-target"]:checked');
    const targets = Array.from(checked).map(cb => cb.value);
    if (targets.length === 0) {
        showSnackbar('Select at least one target platform');
        return;
    }

    try {
        await api('/api/agents', {
            method: 'POST',
            body: JSON.stringify({
                name,
                content: content || '',
                description: description || undefined,
                category: category || undefined,
                targets,
            }),
        });
        closeModal();
        showSnackbar(`Agent "${name}" created`);
        loadAgents();
    } catch (e) {
        showSnackbar('Create failed: ' + e.message, 6000);
    }
}

async function editAgent(name) {
    try {
        const [a, targetsRes] = await Promise.all([
            api('/api/agents/' + name),
            api('/api/targets'),
        ]);
        const allTargets = targetsRes.targets || [];

        const agentConfig = await api('/api/config').catch(() => null);
        const agentEntry = agentConfig?.agents?.[name];
        const currentTargets = agentEntry
            ? resolveAgentTargets(agentEntry, allTargets)
            : allTargets;

        const categories = Object.entries(CATEGORIES);
        let categoryOptions = '<option value="">None</option>' +
            categories.map(([id, label]) =>
                `<option value="${id}" ${a.category === id ? 'selected' : ''}>${escapeHtml(label)}</option>`
            ).join('');

        let targetCheckboxes = allTargets.map(t =>
            `<label class="checkbox-label">
                <input type="checkbox" name="edit-target" value="${escapeHtml(t)}" ${currentTargets.includes(t) ? 'checked' : ''}>
                ${escapeHtml(t)}
            </label>`
        ).join('');

        document.getElementById('modal-title').textContent = 'Edit: ' + name;
        document.getElementById('modal-body').innerHTML = `
            <label for="edit-description">Description</label>
            <input type="text" class="m3-input" id="edit-description" value="${escapeHtml(a.description || '')}">
            <label for="edit-category">Category</label>
            <select class="m3-select" id="edit-category" style="width:100%">${categoryOptions}</select>
            <label>Targets</label>
            <div class="target-checkboxes">${targetCheckboxes}</div>
            <label for="agent-content">Content (Markdown)</label>
            <textarea class="m3-textarea" id="agent-content">${escapeHtml(a.content || '')}</textarea>
        `;
        document.getElementById('modal-actions').innerHTML = `
            <button class="m3-btn m3-btn-text" onclick="closeModal()">Cancel</button>
            <button class="m3-btn m3-btn-filled" onclick="saveAgent('${escapeHtml(name)}')">Save</button>
        `;
        openModal();
    } catch (e) {
        showSnackbar('Failed to load agent: ' + e.message, 6000);
    }
}

async function saveAgent(name) {
    const content = document.getElementById('agent-content').value;
    const description = document.getElementById('edit-description').value.trim();
    const category = document.getElementById('edit-category').value;
    const checked = document.querySelectorAll('input[name="edit-target"]:checked');
    const targets = Array.from(checked).map(cb => cb.value);

    try {
        await api('/api/agents/' + name, {
            method: 'PUT',
            body: JSON.stringify({
                content,
                description,
                category,
                targets,
            }),
        });
        closeModal();
        showSnackbar(`Agent "${name}" updated`);
        if (selectedAgent && selectedAgent.name === name) {
            selectAgent(name);
        }
        loadAgents();
    } catch (e) {
        showSnackbar('Save failed: ' + e.message, 6000);
    }
}

async function deleteAgent(name) {
    if (!confirm('Delete agent "' + name + '"?')) return;
    try {
        await api('/api/agents/' + name, { method: 'DELETE' });
        if (selectedAgent && selectedAgent.name === name) {
            closeDetail();
        }
        showSnackbar(`Agent "${name}" deleted`);
        loadAgents();
    } catch (e) {
        showSnackbar('Delete failed: ' + e.message, 6000);
    }
}

let targetFilter = 'all';
let targetConfig = null;
let allPlatforms = [];

const PLATFORM_ICONS = {
    'claude-code': '\u{1F4AC}',
    'cursor': '\u{1F4BB}',
    'gemini-cli': '\u{2728}',
    'opencode': '\u{26A1}',
};

function resolveAgentTargets(agentEntry, platforms) {
    if (!agentEntry || agentEntry.disabled) return [];
    const t = agentEntry.targets || [];
    if (t.length === 1 && t[0] === 'all') return [...platforms];
    return t;
}

async function loadTargets() {
    showLoading('target-matrix');
    try {
        const [config, agents, validResp] = await Promise.all([
            api('/api/config'),
            api('/api/agents'),
            api('/api/targets'),
        ]);
        targetConfig = config;
        allPlatforms = validResp.targets || [];
        const agentList = agents || [];
        const agentMap = config.agents || {};

        const enabled = agentList.filter(a => {
            const entry = agentMap[a.name];
            return entry && resolveAgentTargets(entry, allPlatforms).length > 0;
        }).length;

        document.getElementById('matrix-stat-total').textContent = agentList.length;
        document.getElementById('matrix-stat-enabled').textContent = enabled;
        document.getElementById('matrix-stat-disabled').textContent = agentList.length - enabled;
        document.getElementById('matrix-stat-platforms').textContent = allPlatforms.length;

        renderTargetFilterChips(enabled, agentList.length - enabled);
        renderTargetMatrix(agentList, allPlatforms, agentMap);
    } catch (e) {
        showError('target-matrix', e.message);
    }
}

function renderTargetFilterChips(enabledCount, disabledCount) {
    const container = document.getElementById('matrix-filter-chips');
    container.innerHTML = `
        <button class="chip ${targetFilter === 'all' ? 'selected' : ''}" onclick="filterTargetView('all')">
            <span class="chip-check">✓</span> All (${enabledCount + disabledCount})
        </button>
        <button class="chip ${targetFilter === 'enabled' ? 'selected' : ''}" onclick="filterTargetView('enabled')">
            <span class="chip-check">✓</span> Enabled (${enabledCount})
        </button>
        <button class="chip ${targetFilter === 'disabled' ? 'selected' : ''}" onclick="filterTargetView('disabled')">
            <span class="chip-check">✓</span> Disabled (${disabledCount})
        </button>
    `;
}

function filterTargetView(filter) {
    targetFilter = filter;
    loadTargets();
}

function renderTargetMatrix(agentList, targets, agentMap) {
    const matrix = document.getElementById('target-matrix');

    const colCounts = {};
    targets.forEach(t => colCounts[t] = 0);
    agentList.forEach(a => {
        const resolved = resolveAgentTargets(agentMap[a.name], targets);
        resolved.forEach(t => colCounts[t] = (colCounts[t] || 0) + 1);
    });

    let filtered = agentList;
    if (targetFilter === 'enabled') {
        filtered = agentList.filter(a => {
            const resolved = resolveAgentTargets(agentMap[a.name], targets);
            return resolved.length > 0;
        });
    } else if (targetFilter === 'disabled') {
        filtered = agentList.filter(a => {
            const resolved = resolveAgentTargets(agentMap[a.name], targets);
            return resolved.length === 0;
        });
    }

    let html = '<table><thead><tr><th class="matrix-agent-header">Agent</th>';
    targets.forEach(t => {
        const icon = PLATFORM_ICONS[t] || '';
        html += `<th class="matrix-platform-header">
            <span class="platform-icon">${icon}</span>
            <span class="platform-name">${escapeHtml(t)}</span>
            <span class="platform-count">${colCounts[t]}</span>
        </th>`;
    });
    html += '</tr></thead><tbody>';

    filtered.forEach(a => {
        const entry = agentMap[a.name] || { targets: [], disabled: false };
        const resolved = resolveAgentTargets(entry, targets);
        const isEnabled = resolved.length > 0;
        const rowClass = isEnabled ? 'matrix-row-enabled' : 'matrix-row-disabled';

        html += `<tr class="${rowClass}"><td class="matrix-agent-name">`;
        if (entry.disabled) {
            html += `<span class="agent-disabled-badge" title="Explicitly disabled">\u{1F6AB}</span> `;
        }
        html += `${escapeHtml(a.name)}</td>`;

        targets.forEach(t => {
            const isTargeted = resolved.includes(t);
            const cellClass = isTargeted ? 'cell-enabled' : (entry.disabled ? 'cell-explicit-disabled' : 'cell-disabled');
            const icon = isTargeted ? '✓' : (entry.disabled ? '\u{1F6AB}' : '—');
            const tooltip = isTargeted ? `Remove ${a.name} from ${t}` : `Add ${a.name} to ${t}`;
            html += `<td class="${cellClass}" title="${tooltip}" onclick="toggleAgentTarget('${escapeHtml(a.name)}', '${escapeHtml(t)}')">${icon}</td>`;
        });
        html += '</tr>';
    });

    html += '</tbody></table>';
    matrix.innerHTML = html;
}

async function toggleAgentTarget(agentName, targetName) {
    if (!targetConfig) return;
    const targets = targetConfig.targets || [];
    const agents = JSON.parse(JSON.stringify(targetConfig.agents || {}));

    if (!agents[agentName]) {
        agents[agentName] = { targets: [targetName], disabled: false };
    } else {
        const entry = agents[agentName];
        let current = entry.targets || [];

        if (current.length === 1 && current[0] === 'all') {
            current = allPlatforms.filter(t => t !== targetName);
            if (current.length === 0) current = [];
        } else {
            const idx = current.indexOf(targetName);
            if (idx >= 0) {
                current.splice(idx, 1);
            } else {
                current.push(targetName);
            }
        }

        entry.targets = current;
    }

    try {
        await api('/api/config', {
            method: 'PUT',
            body: JSON.stringify({ targets, agents }),
        });
        showSnackbar(`${agentName}: ${targetName} updated`);
        loadTargets();
    } catch (e) {
        showSnackbar('Failed to update: ' + e.message, 6000);
    }
}

async function loadPreviewOptions() {
    const output = document.getElementById('preview-output');
    output.textContent = 'Select an agent and target to preview the rendered output.';

    const copyBtn = document.getElementById('preview-copy-btn');
    if (copyBtn) copyBtn.style.display = 'none';

    try {
        const [agents, targets] = await Promise.all([
            api('/api/agents'),
            api('/api/targets'),
        ]);
        const agentSel = document.getElementById('preview-agent');
        const targetSel = document.getElementById('preview-target');
        agentSel.innerHTML = '<option value="">Select agent</option>';
        targetSel.innerHTML = '<option value="">Select target</option>';
        (agents || []).forEach(a => {
            agentSel.innerHTML += `<option value="${escapeHtml(a.name)}">${escapeHtml(a.name)}</option>`;
        });
        (targets.targets || []).forEach(t => {
            targetSel.innerHTML += `<option value="${escapeHtml(t)}">${escapeHtml(t)}</option>`;
        });
    } catch (e) {
        console.error(e);
    }
}

async function loadPreview() {
    const agent = document.getElementById('preview-agent').value;
    const target = document.getElementById('preview-target').value;
    if (!agent || !target) {
        showSnackbar('Select both agent and target');
        return;
    }
    const output = document.getElementById('preview-output');
    output.textContent = 'Loading preview...';
    try {
        const res = await api('/api/preview/' + agent + '/' + target);
        output.textContent = res.content || '';
        const copyBtn = document.getElementById('preview-copy-btn');
        if (copyBtn) copyBtn.style.display = 'inline-flex';
    } catch (e) {
        output.textContent = 'Error: ' + e.message;
    }
}

function copyPreview() {
    const text = document.getElementById('preview-output').textContent;
    navigator.clipboard.writeText(text).then(() => showSnackbar('Copied to clipboard'));
}

let statusFilter = 'all';
let statusSearchQuery = '';
let allStatusLinks = [];

async function loadStatus() {
    showLoading('status-content');
    try {
        const status = await api('/api/status');
        allStatusLinks = status.symlinks || [];

        const controls = document.getElementById('status-controls');
        if (controls) {
            controls.style.display = allStatusLinks.length > 0 ? 'flex' : 'none';
        }

        renderStatusFilterChips();
        renderStatusList();
    } catch (e) {
        showError('status-content', e.message);
    }
}

function renderStatusFilterChips() {
    const container = document.getElementById('status-filter-chips');
    if (!container) return;
    const broken = allStatusLinks.filter(l => l.broken).length;
    const healthy = allStatusLinks.length - broken;
    container.innerHTML = `
        <button class="chip ${statusFilter === 'all' ? 'selected' : ''}" onclick="filterStatusView('all')">
            <span class="chip-check">✓</span> All (${allStatusLinks.length})
        </button>
        <button class="chip ${statusFilter === 'healthy' ? 'selected' : ''}" onclick="filterStatusView('healthy')">
            <span class="chip-check">✓</span> Healthy (${healthy})
        </button>
        <button class="chip ${statusFilter === 'broken' ? 'selected' : ''}" onclick="filterStatusView('broken')">
            <span class="chip-check">✓</span> Broken (${broken})
        </button>
    `;
}

function filterStatusView(filter) {
    statusFilter = filter;
    renderStatusFilterChips();
    renderStatusList();
}

function filterStatus() {
    statusSearchQuery = document.getElementById('status-search').value.toLowerCase().trim();
    renderStatusList();
}

function renderStatusList() {
    const el = document.getElementById('status-content');
    let links = allStatusLinks;

    if (statusFilter === 'healthy') {
        links = links.filter(l => !l.broken);
    } else if (statusFilter === 'broken') {
        links = links.filter(l => l.broken);
    }

    if (statusSearchQuery) {
        links = links.filter(l =>
            l.agent.toLowerCase().includes(statusSearchQuery) ||
            l.platform.toLowerCase().includes(statusSearchQuery) ||
            l.path.toLowerCase().includes(statusSearchQuery)
        );
    }

    if (links.length === 0) {
        el.innerHTML = `<div class="empty-state">
            <div class="empty-icon">&#128196;</div>
            <h3>${allStatusLinks.length === 0 ? 'No symlinks found' : 'No matching symlinks'}</h3>
            <p>${allStatusLinks.length === 0 ? 'Run sync first to generate platform-specific agent files.' : 'Try adjusting your filter or search.'}</p>
        </div>`;
        return;
    }

    el.innerHTML = links.map(l => {
        const cls = l.broken ? 'status-err' : 'status-ok';
        const icon = l.broken ? '&#10005;' : '&#10003;';
        return `<div class="status-item ${cls}">
            <div class="status-icon">${icon}</div>
            <span class="status-text">${escapeHtml(l.agent)} &rarr; ${escapeHtml(l.platform)}</span>
            <span class="status-path">${escapeHtml(l.path)}</span>
        </div>`;
    }).join('');
}

async function triggerSync() {
    if (!confirm('Sync all agents to detected platforms?')) return;
    const btns = document.querySelectorAll('[onclick*="triggerSync"]');
    btns.forEach(b => { b.classList.add('m3-btn-loading'); b.disabled = true; });
    try {
        const res = await api('/api/sync', { method: 'POST' });
        showSnackbar(`Synced ${res.synced || 0} agents`);
        if (currentPage === 'status') loadStatus();
        if (currentPage === 'dashboard') loadDashboard();
    } catch (e) {
        showSnackbar('Sync failed: ' + e.message, 6000);
    } finally {
        btns.forEach(b => { b.classList.remove('m3-btn-loading'); b.disabled = false; });
    }
}

async function triggerClean() {
    if (!confirm('Remove all generated files and symlinks?')) return;
    const btns = document.querySelectorAll('[onclick*="triggerClean"]');
    btns.forEach(b => { b.classList.add('m3-btn-loading'); b.disabled = true; });
    try {
        const res = await api('/api/clean', { method: 'POST' });
        showSnackbar(`Cleaned ${res.removed || 0} files`);
        if (currentPage === 'status') loadStatus();
        if (currentPage === 'dashboard') loadDashboard();
    } catch (e) {
        showSnackbar('Clean failed: ' + e.message, 6000);
    } finally {
        btns.forEach(b => { b.classList.remove('m3-btn-loading'); b.disabled = false; });
    }
}

function openModal() {
    lastFocusedElement = document.activeElement;
    document.getElementById('modal').classList.remove('hidden');
    const firstInput = document.querySelector('.modal-surface input, .modal-surface select, .modal-surface textarea');
    if (firstInput) firstInput.focus();
}

function closeModal() {
    document.getElementById('modal').classList.add('hidden');
    clearValidation();
    if (lastFocusedElement) {
        lastFocusedElement.focus();
        lastFocusedElement = null;
    }
}

document.getElementById('modal').addEventListener('click', e => {
    if (e.target.id === 'modal') closeModal();
});

document.addEventListener('keydown', e => {
    if (e.key === 'Escape') {
        if (!document.getElementById('modal').classList.contains('hidden')) {
            closeModal();
        } else if (document.getElementById('app').classList.contains('detail-open')) {
            closeDetail();
        }
    }
    if (e.key === 'Tab' && !document.getElementById('modal').classList.contains('hidden')) {
        const modal = document.querySelector('.modal-surface');
        const focusable = modal.querySelectorAll('input, select, textarea, button:not([disabled])');
        if (focusable.length === 0) return;
        const first = focusable[0];
        const last = focusable[focusable.length - 1];
        if (e.shiftKey && document.activeElement === first) {
            e.preventDefault();
            last.focus();
        } else if (!e.shiftKey && document.activeElement === last) {
            e.preventDefault();
            first.focus();
        }
    }
});

loadDashboard();
