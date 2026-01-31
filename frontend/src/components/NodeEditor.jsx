import { useState, useEffect } from 'react';

const optionColors = ['#e91e63', '#9c27b0', '#673ab7', '#3f51b5', '#2196f3', '#00bcd4', '#009688', '#4caf50'];

export default function NodeEditor({ node, npcs, items, factions, resources, onSave, onClose }) {
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
                </div>
              ))}

              {(!data.options || data.options.length === 0) && (
                <div style={styles.noOptions}>
                  No options defined. Click "+ Add Option" to create dialog choices.
                </div>
              )}
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

const styles = {
  overlay: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.7)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modal: {
    backgroundColor: '#2a2a3e',
    borderRadius: '8px',
    width: '600px',
    maxHeight: '85vh',
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '16px 20px',
    borderBottom: '1px solid #444',
    color: '#fff',
  },
  closeBtn: {
    background: 'none',
    border: 'none',
    color: '#888',
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
    color: '#888',
    fontSize: '12px',
    marginBottom: '-8px',
  },
  input: {
    padding: '10px',
    backgroundColor: '#1a1a2e',
    border: '1px solid #444',
    borderRadius: '4px',
    color: '#fff',
  },
  textarea: {
    padding: '10px',
    backgroundColor: '#1a1a2e',
    border: '1px solid #444',
    borderRadius: '4px',
    color: '#fff',
    minHeight: '80px',
    resize: 'vertical',
  },
  select: {
    padding: '10px',
    backgroundColor: '#1a1a2e',
    border: '1px solid #444',
    borderRadius: '4px',
    color: '#fff',
  },
  optionsHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: '16px',
    paddingTop: '16px',
    borderTop: '1px solid #444',
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
    backgroundColor: '#1a1a2e',
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
    color: '#aaa',
    fontSize: '12px',
    marginTop: '4px',
  },
  connectedTo: {
    fontSize: '10px',
    color: '#666',
    marginTop: '4px',
  },
  noOptions: {
    color: '#666',
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
    borderTop: '1px solid #444',
  },
  cancelBtn: {
    padding: '8px 16px',
    backgroundColor: '#444',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  saveBtn: {
    padding: '8px 16px',
    backgroundColor: '#4a4a6a',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
};
