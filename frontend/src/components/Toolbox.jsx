import { useTheme } from '../ThemeContext';

const nodeTypes = [
  { type: 'EntryPoint', label: 'Entry Point', color: '#4caf50' },
  { type: 'ConditionWatcher', label: 'Condition', color: '#2196f3' },
  { type: 'Dialog', label: 'Dialog', color: '#9c27b0' },
  { type: 'PlayerDecisionDialog', label: 'Decision', color: '#e91e63' },
  { type: 'Actions', label: 'Actions', color: '#f44336' },
];

export default function Toolbox({ quest, onDragStart, onEditQuest }) {
  const { theme } = useTheme();
  const styles = getStyles(theme);

  const questDisplayName = quest?.DisplayName?.['en-US'] || '(Untitled Quest)';

  return (
    <div style={styles.container}>
      {quest && (
        <>
          <h3 style={styles.title}>Quest</h3>
          <div style={styles.questCard}>
            <div style={styles.questName}>{questDisplayName}</div>
            <button onClick={onEditQuest} style={styles.editBtn}>Edit</button>
          </div>
        </>
      )}

      <h3 style={{ ...styles.title, marginTop: quest ? '20px' : 0 }}>Node Types</h3>
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

const getStyles = (theme) => ({
  container: {
    width: '180px',
    backgroundColor: theme.bg,
    borderRight: `1px solid ${theme.border}`,
    padding: '12px',
    overflowY: 'auto',
  },
  title: {
    color: theme.textMuted,
    fontSize: '12px',
    textTransform: 'uppercase',
    marginBottom: '12px',
  },
  questCard: {
    backgroundColor: theme.bgSecondary,
    borderRadius: '4px',
    padding: '10px 12px',
    marginBottom: '6px',
  },
  questName: {
    color: theme.text,
    fontSize: '13px',
    fontWeight: '500',
    marginBottom: '8px',
    wordBreak: 'break-word',
  },
  editBtn: {
    padding: '4px 10px',
    backgroundColor: theme.name === 'dark' ? '#4a4a6a' : '#5c6bc0',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    fontSize: '12px',
    cursor: 'pointer',
    width: '100%',
  },
  item: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    padding: '10px 12px',
    marginBottom: '6px',
    backgroundColor: theme.bgSecondary,
    borderRadius: '4px',
    borderLeft: '3px solid',
    cursor: 'grab',
    color: theme.textSecondary,
    fontSize: '13px',
  },
  dot: {
    width: '8px',
    height: '8px',
    borderRadius: '50%',
  },
});
