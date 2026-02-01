const API_BASE = '/api';

export async function fetchQuests() {
  const res = await fetch(`${API_BASE}/quests`);
  if (!res.ok) throw new Error('Failed to fetch quests');
  return res.json();
}

export async function fetchQuest(questId) {
  const res = await fetch(`${API_BASE}/quests/${encodeURIComponent(questId)}`);
  if (!res.ok) throw new Error('Failed to fetch quest');
  return res.json();
}

export async function saveQuest(questId, quest, metadata) {
  const res = await fetch(`${API_BASE}/quests/${encodeURIComponent(questId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ quest, metadata }),
  });
  if (!res.ok) throw new Error('Failed to save quest');
  return res.json();
}

export async function deleteQuest(questId) {
  const res = await fetch(`${API_BASE}/quests/${encodeURIComponent(questId)}`, {
    method: 'DELETE',
  });
  if (!res.ok) throw new Error('Failed to delete quest');
}

export async function validateQuest(quest) {
  const res = await fetch(`${API_BASE}/validate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(quest),
  });
  if (!res.ok) throw new Error('Failed to validate quest');
  return res.json();
}

export async function fetchItems() {
  const res = await fetch(`${API_BASE}/items`);
  if (!res.ok) throw new Error('Failed to fetch items');
  return res.json();
}

export async function fetchFactions() {
  const res = await fetch(`${API_BASE}/factions`);
  if (!res.ok) throw new Error('Failed to fetch factions');
  return res.json();
}

export async function fetchResources() {
  const res = await fetch(`${API_BASE}/resources`);
  if (!res.ok) throw new Error('Failed to fetch resources');
  return res.json();
}

export async function fetchNPCs() {
  const res = await fetch(`${API_BASE}/npcs`);
  if (!res.ok) throw new Error('Failed to fetch NPCs');
  return res.json();
}

export async function fetchObjects() {
  const res = await fetch(`${API_BASE}/objects`);
  if (!res.ok) throw new Error('Failed to fetch objects');
  return res.json();
}
