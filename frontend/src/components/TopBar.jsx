import { useState } from 'react';
import { useTheme } from '../ThemeContext';

export default function TopBar({ questId, quests, onSelect, onNew, onSave, validation, saving, saveError, onToggleTheme }) {
  const { theme, themeName } = useTheme();
  const [search, setSearch] = useState('');

  const filteredQuests = quests.filter(q =>
    q.toLowerCase().includes(search.toLowerCase())
  );

  const isValid = validation?.valid !== false;

  const styles = getStyles(theme);

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
        {saveError && (
          <span style={styles.errorMessage} title={saveError}>
            ⚠ {saveError.length > 40 ? saveError.substring(0, 40) + '...' : saveError}
          </span>
        )}
      </div>

      <div style={styles.right}>
        <button onClick={onToggleTheme} style={styles.themeButton} title="Toggle theme">
          {themeName === 'dark' ? '☀️' : '🌙'}
        </button>
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

const getStyles = (theme) => ({
  container: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '8px 16px',
    backgroundColor: theme.bg,
    borderBottom: `1px solid ${theme.border}`,
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
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '2px',
  },
  right: {
    display: 'flex',
    alignItems: 'center',
    gap: '12px',
  },
  button: {
    padding: '6px 12px',
    backgroundColor: theme.bgTertiary,
    color: theme.text,
    border: `1px solid ${theme.border}`,
    borderRadius: '4px',
    cursor: 'pointer',
  },
  saveButton: {
    backgroundColor: theme.name === 'dark' ? '#4a4a6a' : '#5c6bc0',
    color: '#fff',
  },
  themeButton: {
    padding: '4px 8px',
    backgroundColor: 'transparent',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '18px',
  },
  search: {
    padding: '6px 10px',
    backgroundColor: theme.inputBg,
    color: theme.text,
    border: `1px solid ${theme.inputBorder}`,
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
    backgroundColor: theme.bgSecondary,
    border: `1px solid ${theme.border}`,
    borderRadius: '4px',
    maxHeight: '200px',
    overflowY: 'auto',
    zIndex: 1000,
  },
  dropdownItem: {
    padding: '8px 10px',
    cursor: 'pointer',
    borderBottom: `1px solid ${theme.border}`,
    color: theme.text,
  },
  select: {
    padding: '6px 10px',
    backgroundColor: theme.inputBg,
    color: theme.text,
    border: `1px solid ${theme.inputBorder}`,
    borderRadius: '4px',
  },
  questName: {
    color: theme.textMuted,
    fontSize: '14px',
  },
  errorMessage: {
    color: '#f44336',
    fontSize: '11px',
    maxWidth: '300px',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  status: {
    fontSize: '20px',
  },
});
