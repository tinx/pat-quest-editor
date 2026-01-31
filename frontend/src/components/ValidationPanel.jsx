import { useTheme } from '../ThemeContext';

export default function ValidationPanel({ validation }) {
  const { theme } = useTheme();
  const styles = getStyles(theme);

  if (!validation || validation.valid) {
    return (
      <div style={styles.container}>
        <div style={styles.valid}>✓ Quest is valid</div>
      </div>
    );
  }

  return (
    <div style={styles.container}>
      <h3 style={styles.title}>Warnings</h3>
      <div style={styles.list}>
        {validation.errors?.map((err, i) => (
          <div key={i} style={styles.error}>
            <span style={styles.icon}>⚠</span>
            <span>
              {err.nodeId !== undefined && <strong>Node {err.nodeId}: </strong>}
              {err.message}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}

const getStyles = (theme) => ({
  container: {
    width: '250px',
    backgroundColor: theme.bg,
    borderLeft: `1px solid ${theme.border}`,
    padding: '12px',
    overflowY: 'auto',
  },
  title: {
    color: '#ff9800',
    fontSize: '12px',
    textTransform: 'uppercase',
    marginBottom: '12px',
  },
  valid: {
    color: '#4caf50',
    padding: '12px',
    textAlign: 'center',
  },
  list: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
  error: {
    display: 'flex',
    alignItems: 'flex-start',
    gap: '8px',
    padding: '8px',
    backgroundColor: theme.bgSecondary,
    borderRadius: '4px',
    fontSize: '12px',
    color: theme.textSecondary,
  },
  icon: {
    color: '#ff9800',
  },
});
