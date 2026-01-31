import { useState } from 'react';

export default function TopBar({ questId, quests, onSelect, onNew, onSave, validation, saving }) {
  const [search, setSearch] = useState('');

  const filteredQuests = quests.filter(q =>
    q.toLowerCase().includes(search.toLowerCase())
  );

  const isValid = validation?.valid !== false;

  return (
    <div style={styles.container}>
      <div style={styles.left}>
        <button onClick={onNew} style={styles.button}>New Quest</button>
        <div style={styles.searchContainer}>
          <input
            type="text"
            placeholder="Search quests..."
            value={search}
            onChange={e => setSearch(e.target.value)}
            style={styles.search}
          />
          {search && filteredQuests.length > 0 && (
            <div style={styles.dropdown}>
              {filteredQuests.slice(0, 10).map(q => (
                <div
                  key={q}
                  style={styles.dropdownItem}
                  onClick={() => { onSelect(q); setSearch(''); }}
                >
                  {q}
                </div>
              ))}
            </div>
          )}
        </div>
        <select
          value={questId || ''}
          onChange={e => onSelect(e.target.value)}
          style={styles.select}
        >
          <option value="">Select quest...</option>
          {quests.map(q => <option key={q} value={q}>{q}</option>)}
        </select>
      </div>

      <div style={styles.center}>
        {questId && <span style={styles.questName}>{questId}</span>}
      </div>

      <div style={styles.right}>
        <span style={{ ...styles.status, color: isValid ? '#4caf50' : '#ff9800' }}>
          {isValid ? '✓' : '⚠'}
        </span>
        <button
          onClick={onSave}
          disabled={!questId || saving}
          style={{ ...styles.button, ...styles.saveButton }}
        >
          {saving ? 'Saving...' : 'Save'}
        </button>
      </div>
    </div>
  );
}

const styles = {
  container: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '8px 16px',
    backgroundColor: '#1a1a2e',
    borderBottom: '1px solid #333',
    height: '50px',
    boxSizing: 'border-box',
  },
  left: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
  center: {
    flex: 1,
    textAlign: 'center',
  },
  right: {
    display: 'flex',
    alignItems: 'center',
    gap: '12px',
  },
  button: {
    padding: '6px 12px',
    backgroundColor: '#333',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  saveButton: {
    backgroundColor: '#4a4a6a',
  },
  search: {
    padding: '6px 10px',
    backgroundColor: '#2a2a3e',
    color: '#fff',
    border: '1px solid #444',
    borderRadius: '4px',
    width: '150px',
  },
  searchContainer: {
    position: 'relative',
  },
  dropdown: {
    position: 'absolute',
    top: '100%',
    left: 0,
    right: 0,
    backgroundColor: '#2a2a3e',
    border: '1px solid #444',
    borderRadius: '4px',
    maxHeight: '200px',
    overflowY: 'auto',
    zIndex: 1000,
  },
  dropdownItem: {
    padding: '8px 10px',
    cursor: 'pointer',
    borderBottom: '1px solid #333',
  },
  select: {
    padding: '6px 10px',
    backgroundColor: '#2a2a3e',
    color: '#fff',
    border: '1px solid #444',
    borderRadius: '4px',
  },
  questName: {
    color: '#aaa',
    fontSize: '14px',
  },
  status: {
    fontSize: '20px',
  },
};
