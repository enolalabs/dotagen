const API = '';
let allAgents = [];
let selectedCategory = 'all';
let searchQuery = '';
let selectedAgent = null;
let currentPage = 'dashboard';

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

function showSnackbar(message, duration = 3000) {
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
    try {
        const [agents, config, status] = await Promise.all([
            api('/api/agents'),
            api('/api/config'),
            api('/api/status'),
        ]);
        allAgents = agents || [];
        document.getElementById('stat-agents').textContent = allAgents.length;
        document.getElementById('stat-targets').textContent = (config.targets || []).length;
        document.getElementById('stat-synced').textContent = (status.symlinks || []).filter(s => !s.broken).length;

        const cats = new Set(allAgents.map(a => a.category).filter(Boolean));
        document.getElementById('stat-categories').textContent = cats.size;
    } catch (e) {
        console.error(e);
    }
}

async function loadAgents() {
    try {
        const agents = await api('/api/agents');
        allAgents = agents || [];
        renderCategoryChips();
        renderAgentsGrid();
    } catch (e) {
        console.error(e);
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
        return `<div class="agent-card ${isSelected ? 'selected' : ''}" onclick="selectAgent('${escapeHtml(a.name)}')">
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
        showSnackbar('Failed to load agent: ' + e.message);
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
        <label>Agent Name</label>
        <input type="text" class="m3-input" id="agent-name" placeholder="my-agent">
        <label>Description</label>
        <input type="text" class="m3-input" id="agent-description" placeholder="Short description of the agent">
        <label>Category</label>
        <select class="m3-select" id="agent-category" style="width:100%">${categoryOptions}</select>
        <label>Targets</label>
        <div class="target-checkboxes">${targetCheckboxes}</div>
        <label>Content (Markdown)</label>
        <textarea class="m3-textarea" id="agent-content" placeholder="# My Agent\n\n## Role\n\nDescribe the agent's role.\n\n## Guidelines\n\n- Guideline 1\n- Guideline 2"></textarea>
    `;
    document.getElementById('modal-actions').innerHTML = `
        <button class="m3-btn m3-btn-text" onclick="closeModal()">Cancel</button>
        <button class="m3-btn m3-btn-filled" onclick="createAgent()">Create</button>
    `;
    document.getElementById('modal').classList.remove('hidden');
}

async function createAgent() {
    const name = document.getElementById('agent-name').value.trim();
    const description = document.getElementById('agent-description').value.trim();
    const category = document.getElementById('agent-category').value;
    const content = document.getElementById('agent-content').value;

    if (!name) {
        showSnackbar('Name is required');
        return;
    }

    const checked = document.querySelectorAll('input[name="agent-target"]:checked');
    const targets = Array.from(checked).map(cb => cb.value);
    const allTargets = document.querySelectorAll('input[name="agent-target"]');
    const useAll = targets.length === allTargets.length;

    try {
        let fullContent = content;
        if (category || description) {
            let fm = '---\n';
            if (description) fm += `description: ${description}\n`;
            if (category) fm += `category: ${category}\n`;
            fm += '---\n\n';
            fullContent = fm + (content || '');
        }

        await api('/api/agents', {
            method: 'POST',
            body: JSON.stringify({
                name,
                content: fullContent,
                description: description || undefined,
                targets: useAll ? ['all'] : targets,
            }),
        });
        closeModal();
        showSnackbar(`Agent "${name}" created`);
        loadAgents();
    } catch (e) {
        showSnackbar('Create failed: ' + e.message);
    }
}

async function editAgent(name) {
    try {
        const a = await api('/api/agents/' + name);
        document.getElementById('modal-title').textContent = 'Edit: ' + name;
        document.getElementById('modal-body').innerHTML = `
            <label>Content (Markdown)</label>
            <textarea class="m3-textarea" id="agent-content">${escapeHtml(a.content || '')}</textarea>
        `;
        document.getElementById('modal-actions').innerHTML = `
            <button class="m3-btn m3-btn-text" onclick="closeModal()">Cancel</button>
            <button class="m3-btn m3-btn-filled" onclick="saveAgent('${escapeHtml(name)}')">Save</button>
        `;
        document.getElementById('modal').classList.remove('hidden');
    } catch (e) {
        showSnackbar('Failed to load agent: ' + e.message);
    }
}

async function saveAgent(name) {
    const content = document.getElementById('agent-content').value;
    try {
        await api('/api/agents/' + name, {
            method: 'PUT',
            body: JSON.stringify({ content }),
        });
        closeModal();
        showSnackbar(`Agent "${name}" updated`);
        if (selectedAgent && selectedAgent.name === name) {
            selectAgent(name);
        }
        loadAgents();
    } catch (e) {
        showSnackbar('Save failed: ' + e.message);
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
        showSnackbar('Delete failed: ' + e.message);
    }
}

let targetFilter = 'all';
let targetConfig = null;

const PLATFORM_ICONS = {
    'claude-code': '\u{1F4AC}',
    'cursor': '\u{1F4BB}',
    'gemini-cli': '\u{2728}',
    'opencode': '\u{26A1}',
};

function resolveAgentTargets(agentEntry, allTargets) {
    if (!agentEntry || agentEntry.disabled) return [];
    const t = agentEntry.targets || [];
    if (t.length === 1 && t[0] === 'all') return [...allTargets];
    return t;
}

async function loadTargets() {
    try {
        const [config, agents] = await Promise.all([
            api('/api/config'),
            api('/api/agents'),
        ]);
        targetConfig = config;
        const targets = config.targets || [];
        const agentList = agents || [];
        const agentMap = config.agents || {};

        const enabled = agentList.filter(a => {
            const entry = agentMap[a.name];
            return entry && resolveAgentTargets(entry, targets).length > 0;
        }).length;

        document.getElementById('matrix-stat-total').textContent = agentList.length;
        document.getElementById('matrix-stat-enabled').textContent = enabled;
        document.getElementById('matrix-stat-disabled').textContent = agentList.length - enabled;
        document.getElementById('matrix-stat-platforms').textContent = targets.length;

        renderTargetFilterChips(enabled, agentList.length - enabled);
        renderTargetMatrix(agentList, targets, agentMap);
    } catch (e) {
        console.error(e);
    }
}

function renderTargetFilterChips(enabledCount, disabledCount) {
    const container = document.getElementById('matrix-filter-chips');
    container.innerHTML = `
        <button class="chip ${targetFilter === 'all' ? 'selected' : ''}" onclick="filterTargetView('all')">
            <span class="chip-check">\u2713</span> All (${enabledCount + disabledCount})
        </button>
        <button class="chip ${targetFilter === 'enabled' ? 'selected' : ''}" onclick="filterTargetView('enabled')">
            <span class="chip-check">\u2713</span> Enabled (${enabledCount})
        </button>
        <button class="chip ${targetFilter === 'disabled' ? 'selected' : ''}" onclick="filterTargetView('disabled')">
            <span class="chip-check">\u2713</span> Disabled (${disabledCount})
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
            const icon = isTargeted ? '\u2713' : (entry.disabled ? '\u{1F6AB}' : '\u2014');
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
            current = targets.filter(t => t !== targetName);
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
        if (current.length === targets.length) {
            entry.targets = ['all'];
        }
    }

    try {
        await api('/api/config', {
            method: 'PUT',
            body: JSON.stringify({ targets, agents }),
        });
        showSnackbar(`${agentName}: ${targetName} updated`);
        loadTargets();
    } catch (e) {
        showSnackbar('Failed to update: ' + e.message);
    }
}

async function loadPreviewOptions() {
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
    try {
        const res = await api('/api/preview/' + agent + '/' + target);
        document.getElementById('preview-output').textContent = res.content || '';
    } catch (e) {
        document.getElementById('preview-output').textContent = 'Error: ' + e.message;
    }
}

async function loadStatus() {
    try {
        const status = await api('/api/status');
        const el = document.getElementById('status-content');
        const links = status.symlinks || [];
        if (links.length === 0) {
            el.innerHTML = `<div class="empty-state">
                <div class="empty-icon">&#128196;</div>
                <h3>No symlinks found</h3>
                <p>Run sync first to generate platform-specific agent files.</p>
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
    } catch (e) {
        console.error(e);
    }
}

async function triggerSync() {
    try {
        const res = await api('/api/sync', { method: 'POST' });
        showSnackbar(`Synced ${res.synced || 0} agents`);
        if (currentPage === 'status') loadStatus();
        if (currentPage === 'dashboard') loadDashboard();
    } catch (e) {
        showSnackbar('Sync failed: ' + e.message);
    }
}

async function triggerClean() {
    if (!confirm('Remove all generated files and symlinks?')) return;
    try {
        const res = await api('/api/clean', { method: 'POST' });
        showSnackbar(`Cleaned ${res.removed || 0} files`);
        if (currentPage === 'status') loadStatus();
        if (currentPage === 'dashboard') loadDashboard();
    } catch (e) {
        showSnackbar('Clean failed: ' + e.message);
    }
}

function closeModal() {
    document.getElementById('modal').classList.add('hidden');
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
});

loadDashboard();
