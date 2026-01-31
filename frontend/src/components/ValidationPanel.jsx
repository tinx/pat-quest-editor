export default function ValidationPanel({ validation }) {
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

const styles = {
  container: {
    width: '250px',
    backgroundColor: '#1a1a2e',
    borderLeft: '1px solid #333',
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
    backgroundColor: '#2a2a3e',
    borderRadius: '4px',
    fontSize: '12px',
    color: '#ddd',
  },
  icon: {
    color: '#ff9800',
  },
};
