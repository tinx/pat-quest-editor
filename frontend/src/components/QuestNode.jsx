import { memo } from 'react';
import { Handle, Position } from '@xyflow/react';
import { useTheme } from '../ThemeContext';

const nodeColors = {
  EntryPoint: '#4caf50',
  ConditionWatcher: '#2196f3',
  Dialog: '#9c27b0',
  PlayerDecisionDialog: '#e91e63',
  Actions: '#f44336',
};

const optionColors = ['#e91e63', '#9c27b0', '#673ab7', '#3f51b5', '#2196f3', '#00bcd4', '#009688', '#4caf50'];

function QuestNode({ data, selected }) {
  const { theme } = useTheme();
  const color = nodeColors[data.nodeType] || '#666';
  const isDecisionDialog = data.nodeType === 'PlayerDecisionDialog';
  const isEntryPoint = data.nodeType === 'EntryPoint';
  const options = data.options || [];
  const styles = getStyles(theme);
  const isHighlighted = data.highlighted;
  
  return (
    <div style={{
      ...styles.node,
      borderColor: isHighlighted ? '#ffeb3b' : (selected ? (theme.name === 'dark' ? '#fff' : '#333') : color),
      boxShadow: isHighlighted ? '0 0 15px #ffeb3b, 0 0 30px #ffeb3b' : (selected ? `0 0 10px ${color}` : 'none'),
      maxWidth: isDecisionDialog ? '280px' : '200px',
    }}>
      {!isEntryPoint && <Handle type="target" position={Position.Top} style={styles.handle} />}
      
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
        
        {isDecisionDialog && (
          <>
            {data.speaker && <div style={styles.label}>{data.speaker}</div>}
            {data.text?.['en-US'] && (
              <div style={styles.dialogText}>"{data.text['en-US'].substring(0, 50)}..."</div>
            )}
            <div style={styles.optionsList}>
              {options.map((opt, i) => (
                <div key={i} style={styles.optionRow}>
                  <span 
                    style={{
                      ...styles.optionDot,
                      backgroundColor: optionColors[i % optionColors.length],
                    }}
                  />
                  <span style={styles.optionText}>
                    {opt.Text?.['en-US']?.substring(0, 30) || `Option ${i + 1}`}
                    {opt.Text?.['en-US']?.length > 30 ? '...' : ''}
                  </span>
                  <Handle
                    type="source"
                    position={Position.Right}
                    id={`option-${i}`}
                    style={{
                      ...styles.optionHandle,
                      backgroundColor: optionColors[i % optionColors.length],
                      top: 'auto',
                      right: '-6px',
                    }}
                  />
                </div>
              ))}
              {options.length === 0 && (
                <div style={styles.noOptions}>No options defined</div>
              )}
            </div>
          </>
        )}
        
        {data.nodeType === 'ConditionWatcher' && data.conditions?.length > 0 && (
          <div style={styles.conditions}>
            {data.conditions.map((c, i) => (
              <div key={i} style={styles.condition}>
                {Object.keys(c)[0]}
              </div>
            ))}
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
      
      {/* Only show bottom handle for non-decision dialogs */}
      {!isDecisionDialog && (
        <Handle type="source" position={Position.Bottom} style={styles.handle} />
      )}
    </div>
  );
}

const getStyles = (theme) => ({
  node: {
    backgroundColor: theme.bgSecondary,
    borderRadius: '6px',
    border: '2px solid',
    minWidth: '150px',
    maxWidth: '200px',
    fontSize: '12px',
    color: theme.text,
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
    color: theme.textSecondary,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  dialogText: {
    color: theme.textMuted,
    fontSize: '10px',
    fontStyle: 'italic',
    marginTop: '4px',
    marginBottom: '8px',
  },
  optionsList: {
    display: 'flex',
    flexDirection: 'column',
    gap: '6px',
    marginTop: '4px',
  },
  optionRow: {
    display: 'flex',
    alignItems: 'center',
    gap: '6px',
    position: 'relative',
    backgroundColor: theme.bg,
    padding: '4px 8px',
    borderRadius: '4px',
    paddingRight: '16px',
  },
  optionDot: {
    width: '8px',
    height: '8px',
    borderRadius: '50%',
    flexShrink: 0,
  },
  optionText: {
    fontSize: '10px',
    color: theme.textSecondary,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    flex: 1,
  },
  optionHandle: {
    width: '10px',
    height: '10px',
    position: 'absolute',
  },
  noOptions: {
    color: theme.textDim,
    fontSize: '10px',
    fontStyle: 'italic',
  },
  conditions: {
    display: 'flex',
    flexDirection: 'column',
    gap: '2px',
  },
  condition: {
    backgroundColor: theme.bg,
    padding: '2px 6px',
    borderRadius: '3px',
    fontSize: '10px',
  },
  actions: {
    display: 'flex',
    flexDirection: 'column',
    gap: '2px',
  },
  action: {
    backgroundColor: theme.bg,
    padding: '2px 6px',
    borderRadius: '3px',
    fontSize: '10px',
  },
  handle: {
    backgroundColor: theme.textDim,
    width: '8px',
    height: '8px',
  },
});

export default memo(QuestNode);
