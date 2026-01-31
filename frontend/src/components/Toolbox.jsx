import { useTheme } from '../ThemeContext';

const nodeTypes = [
  { type: 'EntryPoint', label: 'Entry Point', color: '#4caf50' },
  { type: 'ConditionWatcher', label: 'Condition', color: '#2196f3' },
  { type: 'Dialog', label: 'Dialog', color: '#9c27b0' },
  { type: 'PlayerDecisionDialog', label: 'Decision', color: '#e91e63' },
  { type: 'Actions', label: 'Actions', color: '#f44336' },
];

export default function Toolbox({ onDragStart }) {
  const { theme } = useTheme();
  const styles = getStyles(theme);

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
