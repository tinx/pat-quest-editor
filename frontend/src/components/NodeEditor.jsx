import { useState, useEffect } from 'react';

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
    width: '500px',
    maxHeight: '80vh',
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
