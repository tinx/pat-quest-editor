const nodeTypes = [
  { type: 'EntryPoint', label: 'Entry Point', color: '#4caf50' },
  { type: 'ConditionWatcher', label: 'Condition', color: '#2196f3' },
  { type: 'Dialog', label: 'Dialog', color: '#9c27b0' },
  { type: 'PlayerDecisionDialog', label: 'Decision', color: '#e91e63' },
  { type: 'QuestProgress', label: 'Progress', color: '#ff9800' },
  { type: 'QuestAvailable', label: 'Available', color: '#00bcd4' },
  { type: 'Actions', label: 'Actions', color: '#f44336' },
];

export default function Toolbox({ onDragStart }) {
  return (
    <div style={styles.container}>
      <h3 style={styles.title}>Node Types</h3>
      {nodeTypes.map(node => (
        <div
          key={node.type}
          style={{ ...styles.item, borderLeftColor: node.color }}
          draggable
          onDragStart={e => onDragStart(e, node.type)}
        >
          <span style={{ ...styles.dot, backgroundColor: node.color }} />
          {node.label}
        </div>
      ))}
    </div>
  );
}

const styles = {
  container: {
    width: '180px',
    backgroundColor: '#1a1a2e',
    borderRight: '1px solid #333',
    padding: '12px',
    overflowY: 'auto',
  },
  title: {
    color: '#888',
    fontSize: '12px',
    textTransform: 'uppercase',
    marginBottom: '12px',
  },
  item: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    padding: '10px 12px',
    marginBottom: '6px',
    backgroundColor: '#2a2a3e',
    borderRadius: '4px',
    borderLeft: '3px solid',
    cursor: 'grab',
    color: '#ddd',
    fontSize: '13px',
  },
  dot: {
    width: '8px',
    height: '8px',
    borderRadius: '50%',
  },
};
