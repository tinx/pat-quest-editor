import { useState, useEffect } from 'react';
import { useTheme } from '../ThemeContext';

const QUEST_TYPES = [
  'SideQuest',
  'MainQuest',
  'CompanionQuest',
  'FactionQuest',
  'DistrictQuest',
  'MasteryQuest',
];

const REPEATABLE_OPTIONS = [
  { value: 'never', label: 'Never' },
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'always', label: 'Always' },
];

export default function QuestPropertiesEditor({ quest, onSave, onClose }) {
  const { theme } = useTheme();
  const [data, setData] = useState({});

  useEffect(() => {
    if (quest) {
      setData({
        QuestID: quest.QuestID || '',
        QuestType: quest.QuestType || 'SideQuest',
        DisplayName: quest.DisplayName || { 'en-US': '', 'de-DE': '' },
        Repeatable: quest.Repeatable || 'never',
      });
    }
  }, [quest]);

  const handleChange = (field, value) => {
    setData(prev => ({ ...prev, [field]: value }));
  };

  const handleI18nChange = (field, lang, value) => {
    setData(prev => ({
      ...prev,
      [field]: { ...prev[field], [lang]: value },
    }));
  };

  const handleSave = () => {
    const updatedQuest = {
      ...quest,
      QuestType: data.QuestType,
      DisplayName: data.DisplayName,
      Repeatable: data.Repeatable,
    };
    onSave(updatedQuest);
    onClose();
  };

  const styles = getStyles(theme);

  return (
    <div style={styles.overlay} onClick={onClose}>
      <div style={styles.modal} onClick={e => e.stopPropagation()}>
        <div style={styles.header}>
          <h3>Quest Properties</h3>
          <button onClick={onClose} style={styles.closeBtn}>×</button>
        </div>

        <div style={styles.body}>
          <label style={styles.label}>Quest ID</label>
          <input
            type="text"
            value={data.QuestID || ''}
            readOnly
            style={{ ...styles.input, opacity: 0.6, cursor: 'not-allowed' }}
          />

          <label style={styles.label}>Quest Type</label>
          <select
            value={data.QuestType || 'SideQuest'}
            onChange={e => handleChange('QuestType', e.target.value)}
            style={styles.select}
          >
            {QUEST_TYPES.map(type => (
              <option key={type} value={type}>{type}</option>
            ))}
          </select>

          <label style={styles.label}>Display Name (English)</label>
          <input
            type="text"
            value={data.DisplayName?.['en-US'] || ''}
            onChange={e => handleI18nChange('DisplayName', 'en-US', e.target.value)}
            style={styles.input}
            placeholder="Quest name in English..."
          />

          <label style={styles.label}>Display Name (German)</label>
          <input
            type="text"
            value={data.DisplayName?.['de-DE'] || ''}
            onChange={e => handleI18nChange('DisplayName', 'de-DE', e.target.value)}
            style={styles.input}
            placeholder="Quest name in German..."
          />

          <label style={styles.label}>Repeatable</label>
          <select
            value={data.Repeatable || 'never'}
            onChange={e => handleChange('Repeatable', e.target.value)}
            style={styles.select}
          >
            {REPEATABLE_OPTIONS.map(opt => (
              <option key={opt.value} value={opt.value}>{opt.label}</option>
            ))}
          </select>
        </div>

        <div style={styles.footer}>
          <button onClick={onClose} style={styles.cancelBtn}>Cancel</button>
          <button onClick={handleSave} style={styles.saveBtn}>Save</button>
        </div>
      </div>
    </div>
  );
}

const getStyles = (theme) => ({
  overlay: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: theme.overlay,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modal: {
    backgroundColor: theme.bgSecondary,
    borderRadius: '8px',
    width: '500px',
    maxHeight: '85vh',
    display: 'flex',
    flexDirection: 'column',
    boxShadow: `0 4px 20px ${theme.shadow}`,
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '16px 20px',
    borderBottom: `1px solid ${theme.border}`,
    color: theme.text,
  },
  closeBtn: {
    background: 'none',
    border: 'none',
    color: theme.textMuted,
    fontSize: '24px',
    cursor: 'pointer',
  },
  body: {
    padding: '20px',
    overflowY: 'auto',
    display: 'flex',
    flexDirection: 'column',
    gap: '12px',
  },
  label: {
    color: theme.textMuted,
    fontSize: '12px',
    marginBottom: '-8px',
  },
  input: {
    padding: '10px',
    backgroundColor: theme.inputBg,
    border: `1px solid ${theme.inputBorder}`,
    borderRadius: '4px',
    color: theme.text,
  },
  textarea: {
    padding: '10px',
    backgroundColor: theme.inputBg,
    border: `1px solid ${theme.inputBorder}`,
    borderRadius: '4px',
    color: theme.text,
    minHeight: '80px',
    resize: 'vertical',
  },
  select: {
    padding: '10px',
    backgroundColor: theme.inputBg,
    border: `1px solid ${theme.inputBorder}`,
    borderRadius: '4px',
    color: theme.text,
  },
  footer: {
    display: 'flex',
    justifyContent: 'flex-end',
    gap: '12px',
    padding: '16px 20px',
    borderTop: `1px solid ${theme.border}`,
  },
  cancelBtn: {
    padding: '8px 16px',
    backgroundColor: theme.bgTertiary,
    color: theme.text,
    border: `1px solid ${theme.border}`,
    borderRadius: '4px',
    cursor: 'pointer',
  },
  saveBtn: {
    padding: '8px 16px',
    backgroundColor: theme.name === 'dark' ? '#4a4a6a' : '#5c6bc0',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
});
