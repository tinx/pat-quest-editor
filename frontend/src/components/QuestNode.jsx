import { memo } from 'react';
import { Handle, Position } from '@xyflow/react';

const nodeColors = {
  EntryPoint: '#4caf50',
  ConditionWatcher: '#2196f3',
  Dialog: '#9c27b0',
  PlayerDecisionDialog: '#e91e63',
  QuestProgress: '#ff9800',
  QuestAvailable: '#00bcd4',
  Actions: '#f44336',
};

function QuestNode({ data, selected }) {
  const color = nodeColors[data.nodeType] || '#666';
  
  return (
    <div style={{
      ...styles.node,
      borderColor: selected ? '#fff' : color,
      boxShadow: selected ? `0 0 10px ${color}` : 'none',
    }}>
      <Handle type="target" position={Position.Top} style={styles.handle} />
      
      <div style={{ ...styles.header, backgroundColor: color }}>
        <span style={styles.type}>{data.nodeType}</span>
        <span style={styles.id}>#{data.nodeId}</span>
      </div>
      
      <div style={styles.body}>
        {data.nodeType === 'EntryPoint' && (
          <div style={styles.label}>Quest Start</div>
        )}
        
        {data.nodeType === 'Dialog' && data.conversationPartner && (
          <div style={styles.label}>{data.conversationPartner}</div>
        )}
        
        {data.nodeType === 'PlayerDecisionDialog' && (
          <>
            {data.speaker && <div style={styles.label}>{data.speaker}</div>}
            {data.options?.length > 0 && (
              <div style={styles.options}>
                {data.options.length} option(s)
              </div>
            )}
          </>
        )}
        
        {data.nodeType === 'ConditionWatcher' && data.conditions?.length > 0 && (
          <div style={styles.conditions}>
            {data.conditions.length} condition(s)
          </div>
        )}
        
        {data.nodeType === 'QuestProgress' && data.questStageTitle && (
          <div style={styles.label}>
            {data.questStageTitle['en-US'] || 'Stage'}
          </div>
        )}
        
        {data.nodeType === 'QuestAvailable' && data.questStageTitle && (
          <div style={styles.label}>
            {data.questStageTitle['en-US'] || 'Available'}
          </div>
        )}
        
        {data.nodeType === 'Actions' && data.actions?.length > 0 && (
          <div style={styles.actions}>
            {data.actions.map((a, i) => (
              <div key={i} style={styles.action}>
                {typeof a === 'string' ? a : Object.keys(a)[0]}
              </div>
            ))}
          </div>
        )}
      </div>
      
      <Handle type="source" position={Position.Bottom} style={styles.handle} />
    </div>
  );
}

const styles = {
  node: {
    backgroundColor: '#2a2a3e',
    borderRadius: '6px',
    border: '2px solid',
    minWidth: '150px',
    maxWidth: '200px',
    fontSize: '12px',
    color: '#fff',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '6px 10px',
    borderRadius: '4px 4px 0 0',
  },
  type: {
    fontWeight: 'bold',
    fontSize: '11px',
  },
  id: {
    opacity: 0.7,
    fontSize: '10px',
  },
  body: {
    padding: '8px 10px',
  },
  label: {
    color: '#ccc',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  options: {
    color: '#888',
    fontSize: '10px',
    marginTop: '4px',
  },
  conditions: {
    color: '#888',
    fontSize: '10px',
  },
  actions: {
    display: 'flex',
    flexDirection: 'column',
    gap: '2px',
  },
  action: {
    backgroundColor: '#1a1a2e',
    padding: '2px 6px',
    borderRadius: '3px',
    fontSize: '10px',
  },
  handle: {
    backgroundColor: '#666',
    width: '8px',
    height: '8px',
  },
};

export default memo(QuestNode);
