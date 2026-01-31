import { useState, useEffect } from 'react';
import ConditionEditor from './ConditionEditor';
import { useTheme } from '../ThemeContext';

const optionColors = ['#e91e63', '#9c27b0', '#673ab7', '#3f51b5', '#2196f3', '#00bcd4', '#009688', '#4caf50'];

export default function NodeEditor({ node, npcs, items, factions, resources, onSave, onClose }) {
  const { theme } = useTheme();
  const [data, setData] = useState(node?.data || {});

  useEffect(() => {
    setData(node?.data || {});
  }, [node]);

  if (!node) return null;

  const handleChange = (field, value) => {
    setData(prev => ({ ...prev, [field]: value }));
  };

  const handleI18nChange = (field, lang, value) => {
    setData(prev => ({
      ...prev,
      [field]: { ...prev[field], [lang]: value },
    }));
  };

  const handleOptionChange = (index, field, lang, value) => {
    setData(prev => {
      const options = [...(prev.options || [])];
      if (lang) {
        options[index] = {
          ...options[index],
          [field]: { ...options[index][field], [lang]: value },
        };
      } else {
        options[index] = { ...options[index], [field]: value };
      }
      return { ...prev, options };
    });
  };

  const handleOptionConditionsChange = (index, conditions) => {
    setData(prev => {
      const options = [...(prev.options || [])];
      options[index] = { ...options[index], Conditions: conditions?.length ? conditions : undefined };
      return { ...prev, options };
    });
  };

  const addOption = () => {
    setData(prev => ({
      ...prev,
      options: [
        ...(prev.options || []),
        {
          Text: { 'en-US': '', 'de-DE': '' },
          NextNodes: [],
        },
      ],
    }));
  };

  const removeOption = (index) => {
    setData(prev => ({
      ...prev,
      options: prev.options.filter((_, i) => i !== index),
    }));
  };

  const handleSave = () => {
    onSave({ ...node, data });
    onClose();
  };

  const styles = getStyles(theme);

  return (
    <div style={styles.overlay} onClick={onClose}>
      <div style={styles.modal} onClick={e => e.stopPropagation()}>
        <div style={styles.header}>
          <h3>Edit {data.nodeType} Node #{data.nodeId}</h3>
          <button onClick={onClose} style={styles.closeBtn}>×</button>
        </div>
        
        <div style={styles.body}>
          {(data.nodeType === 'Dialog' || data.nodeType === 'PlayerDecisionDialog') && (
            <>
              <label style={styles.label}>Conversation Partner</label>
              <select
                value={data.conversationPartner || ''}
                onChange={e => handleChange('conversationPartner', e.target.value)}
                style={styles.select}
              >
                <option value="">Select NPC...</option>
                {npcs?.map(npc => (
                  <option key={npc.NPCID} value={npc.NPCID}>
                    {npc.DisplayName?.['en-US'] || npc.NPCID}
                  </option>
                ))}
              </select>
            </>
          )}

          {data.nodeType === 'PlayerDecisionDialog' && (
            <>
              <label style={styles.label}>Speaker</label>
              <select
                value={data.speaker || ''}
                onChange={e => handleChange('speaker', e.target.value)}
                style={styles.select}
              >
                <option value="">Select Speaker...</option>
                {npcs?.map(npc => (
                  <option key={npc.NPCID} value={npc.NPCID}>
                    {npc.DisplayName?.['en-US'] || npc.NPCID}
                  </option>
                ))}
              </select>

              <label style={styles.label}>Text (English)</label>
              <textarea
                value={data.text?.['en-US'] || ''}
                onChange={e => handleI18nChange('text', 'en-US', e.target.value)}
                style={styles.textarea}
              />

              <label style={styles.label}>Text (German)</label>
              <textarea
                value={data.text?.['de-DE'] || ''}
                onChange={e => handleI18nChange('text', 'de-DE', e.target.value)}
                style={styles.textarea}
              />

              <div style={styles.optionsHeader}>
                <label style={styles.label}>Dialog Options</label>
                <button onClick={addOption} style={styles.addBtn}>+ Add Option</button>
              </div>

              {(data.options || []).map((opt, i) => (
                <div key={i} style={styles.optionCard}>
                  <div style={styles.optionHeader}>
                    <span style={{ ...styles.optionBadge, backgroundColor: optionColors[i % optionColors.length] }}>
                      Option {i + 1}
                    </span>
                    <button onClick={() => removeOption(i)} style={styles.removeBtn}>×</button>
                  </div>
                  
                  <label style={styles.label}>Text (English)</label>
                  <input
                    type="text"
                    value={opt.Text?.['en-US'] || ''}
                    onChange={e => handleOptionChange(i, 'Text', 'en-US', e.target.value)}
                    style={styles.input}
                    placeholder="Option text in English..."
                  />
                  
                  <label style={styles.label}>Text (German)</label>
                  <input
                    type="text"
                    value={opt.Text?.['de-DE'] || ''}
                    onChange={e => handleOptionChange(i, 'Text', 'de-DE', e.target.value)}
                    style={styles.input}
                    placeholder="Option text in German..."
                  />
                  
                  <label style={styles.checkboxLabel}>
                    <input
                      type="checkbox"
                      checked={opt.DefaultOption || false}
                      onChange={e => handleOptionChange(i, 'DefaultOption', null, e.target.checked)}
                    />
                    Default Option (used if player exits dialog early)
                  </label>

                  {opt.NextNodes?.length > 0 && (
                    <div style={styles.connectedTo}>
                      Connected to: Node {opt.NextNodes.join(', ')}
                    </div>
                  )}

                  <ConditionEditor
                    conditions={opt.Conditions}
                    onChange={(conds) => handleOptionConditionsChange(i, conds)}
                    items={items}
                    factions={factions}
                    resources={resources}
                    collapsible={true}
                    defaultExpanded={false}
                  />
                </div>
              ))}

              {(!data.options || data.options.length === 0) && (
                <div style={styles.noOptions}>
                  No options defined. Click "+ Add Option" to create dialog choices.
                </div>
              )}
            </>
          )}

          {data.nodeType === 'ConditionWatcher' && (
            <>
              <label style={styles.label}>Conditions</label>
              <ConditionEditor
                conditions={data.conditions}
                onChange={(conds) => handleChange('conditions', conds?.length ? conds : undefined)}
                items={items}
                factions={factions}
                resources={resources}
                showConditionsRequired={true}
                conditionsRequired={data.conditionsRequired}
                onConditionsRequiredChange={(val) => handleChange('conditionsRequired', val)}
              />
            </>
          )}

          {(data.nodeType === 'QuestProgress' || data.nodeType === 'QuestAvailable') && (
            <>
              <label style={styles.label}>Stage Title (English)</label>
              <input
                type="text"
                value={data.questStageTitle?.['en-US'] || ''}
                onChange={e => handleI18nChange('questStageTitle', 'en-US', e.target.value)}
                style={styles.input}
              />

              <label style={styles.label}>Stage Title (German)</label>
              <input
                type="text"
                value={data.questStageTitle?.['de-DE'] || ''}
                onChange={e => handleI18nChange('questStageTitle', 'de-DE', e.target.value)}
                style={styles.input}
              />
            </>
          )}

          {data.nodeType === 'QuestProgress' && (
            <>
              <label style={styles.label}>Stage Description (English)</label>
              <textarea
                value={data.questStageDescription?.['en-US'] || ''}
                onChange={e => handleI18nChange('questStageDescription', 'en-US', e.target.value)}
                style={styles.textarea}
              />

              <label style={styles.label}>Stage Description (German)</label>
              <textarea
                value={data.questStageDescription?.['de-DE'] || ''}
                onChange={e => handleI18nChange('questStageDescription', 'de-DE', e.target.value)}
                style={styles.textarea}
              />
            </>
          )}
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
    width: '600px',
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
  optionsHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: '16px',
    paddingTop: '16px',
    borderTop: `1px solid ${theme.border}`,
  },
  addBtn: {
    padding: '6px 12px',
    backgroundColor: '#4caf50',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '12px',
  },
  optionCard: {
    backgroundColor: theme.bg,
    borderRadius: '6px',
    padding: '12px',
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
  optionHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  optionBadge: {
    padding: '4px 8px',
    borderRadius: '4px',
    fontSize: '11px',
    fontWeight: 'bold',
    color: '#fff',
  },
  removeBtn: {
    background: 'none',
    border: 'none',
    color: '#f44336',
    fontSize: '20px',
    cursor: 'pointer',
    padding: '0 4px',
  },
  checkboxLabel: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    color: theme.textMuted,
    fontSize: '12px',
    marginTop: '4px',
  },
  connectedTo: {
    fontSize: '10px',
    color: theme.textDim,
    marginTop: '4px',
  },
  noOptions: {
    color: theme.textDim,
    fontSize: '12px',
    fontStyle: 'italic',
    padding: '12px',
    textAlign: 'center',
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
