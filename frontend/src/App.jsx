import { useState, useCallback, useEffect, useRef, Component } from 'react';
import { ThemeProvider, useTheme } from './ThemeContext';
import TopBar from './components/TopBar';
import Toolbox from './components/Toolbox';
import Canvas from './components/Canvas';
import ValidationPanel from './components/ValidationPanel';
import QuestPropertiesEditor from './components/QuestPropertiesEditor';
import { useQuests, useReferenceData } from './hooks/useApi';
import * as api from './api/client';

// Error Boundary to catch React rendering errors
class ErrorBoundary extends Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }

  componentDidCatch(error, errorInfo) {
    console.error('React Error Boundary caught an error:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div style={errorBoundaryStyles.container}>
          <div style={errorBoundaryStyles.content}>
            <h1 style={errorBoundaryStyles.title}>Something went wrong</h1>
            <p style={errorBoundaryStyles.message}>
              The quest editor encountered an unexpected error.
            </p>
            <pre style={errorBoundaryStyles.error}>
              {this.state.error?.message || 'Unknown error'}
            </pre>
            <button
              style={errorBoundaryStyles.button}
              onClick={() => window.location.reload()}
            >
              Reload Page
            </button>
          </div>
        </div>
      );
    }
    return this.props.children;
  }
}

const errorBoundaryStyles = {
  container: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    height: '100vh',
    backgroundColor: '#1a1a2e',
    color: '#fff',
  },
  content: {
    textAlign: 'center',
    padding: '40px',
    maxWidth: '500px',
  },
  title: {
    fontSize: '24px',
    marginBottom: '16px',
    color: '#ff6b6b',
  },
  message: {
    fontSize: '16px',
    marginBottom: '16px',
    color: '#a0a0a0',
  },
  error: {
    backgroundColor: '#2a2a3e',
    padding: '12px',
    borderRadius: '4px',
    fontSize: '12px',
    textAlign: 'left',
    overflow: 'auto',
    maxHeight: '150px',
    marginBottom: '20px',
  },
  button: {
    padding: '10px 24px',
    backgroundColor: '#5c6bc0',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '14px',
  },
};

function AppContent() {
  const { theme, toggleTheme } = useTheme();
  const { quests, refresh: refreshQuests } = useQuests();
  const referenceData = useReferenceData();
  const canvasRef = useRef(null);
  
  const [currentQuestId, setCurrentQuestId] = useState(null);
  const [quest, setQuest] = useState(null);
  const [metadata, setMetadata] = useState(null);
  const [validation, setValidation] = useState(null);
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState(null);
  const [highlightedNodeId, setHighlightedNodeId] = useState(null);
  const [showQuestEditor, setShowQuestEditor] = useState(false);

  // Load quest when selected
  const loadQuest = useCallback(async (questId) => {
    if (!questId) {
      setCurrentQuestId(null);
      setQuest(null);
      setMetadata(null);
      setValidation(null);
      setSaveError(null);
      return;
    }

    try {
      const data = await api.fetchQuest(questId);
      setCurrentQuestId(questId);
      setQuest(data.quest);
      setMetadata(data.metadata);
      setSaveError(null);
      
      // Validate on load
      const result = await api.validateQuest(data.quest);
      setValidation(result);
    } catch (e) {
      console.error('Failed to load quest:', e);
      setSaveError(`Failed to load quest: ${e.message}`);
    }
  }, []);

  // Handle quest changes from canvas
  const handleQuestChange = useCallback(async (updatedQuest, updatedMetadata) => {
    setQuest(updatedQuest);
    setMetadata(updatedMetadata);
    setSaveError(null); // Clear save error on edit

    // Validate on change
    try {
      const result = await api.validateQuest(updatedQuest);
      setValidation(result);
      // Clear highlight if the highlighted node no longer has errors or warnings
      setHighlightedNodeId(prev => {
        if (prev === null) return null;
        const stillHasError = result.errors?.some(e => e.nodeId === prev);
        const stillHasWarning = result.warnings?.some(w => w.nodeId === prev);
        return (stillHasError || stillHasWarning) ? prev : null;
      });
    } catch (e) {
      console.error('Validation failed:', e);
    }
  }, []);

  // Save quest
  const handleSave = useCallback(async () => {
    if (!currentQuestId || !quest) return;

    setSaving(true);
    setSaveError(null);
    try {
      const result = await api.saveQuest(currentQuestId, quest, metadata);
      // Only update validation if save succeeded - preserves context
      setValidation(prev => ({
        ...result,
        // Keep track that this is from a save operation
        savedAt: new Date().toISOString(),
      }));
    } catch (e) {
      console.error('Failed to save quest:', e);
      setSaveError(`Failed to save: ${e.message}`);
    } finally {
      setSaving(false);
    }
  }, [currentQuestId, quest, metadata]);

  // Create new quest
  const handleNewQuest = useCallback(() => {
    const questId = prompt('Enter Quest ID (e.g., PAT_New_Quest):');
    if (!questId) return;

    // Validate quest ID format per schema: ^[A-Z][A-Za-z0-9.\-_:]*$
    // Must start with uppercase letter, followed by alphanumeric, dots, hyphens, underscores, or colons
    if (!/^[A-Z][A-Za-z0-9.\-_:]*$/.test(questId)) {
      alert('Invalid Quest ID. Must start with an uppercase letter, followed by letters, numbers, dots, hyphens, underscores, or colons (e.g., PAT_My_Quest).');
      return;
    }
    if (questId.length > 100) {
      alert('Quest ID is too long. Maximum 100 characters allowed.');
      return;
    }

    // Check if quest already exists
    if (quests.includes(questId)) {
      alert(`A quest with ID "${questId}" already exists. Please choose a different ID.`);
      return;
    }

    const newQuest = {
      QuestTypeVersion: 1,
      QuestVersion: 1,
      QuestID: questId,
      QuestType: 'SideQuest',
      DisplayName: { 'en-US': 'New Quest', 'de-DE': 'Neue Aufgabe' },
      Repeatable: 'never',
      QuestNodes: [
        { NodeID: 0, NodeType: 'EntryPoint', NextNodes: [1] },
        { NodeID: 1, NodeType: 'Actions', Actions: ['CompleteQuest'] },
      ],
    };

    setCurrentQuestId(questId);
    setQuest(newQuest);
    setMetadata({ questId, nodePositions: {} });
    setValidation({ valid: true });
  }, [quests]);

  // Drag start handler for toolbox
  const handleDragStart = useCallback((event, nodeType) => {
    event.dataTransfer.setData('application/reactflow', nodeType);
    event.dataTransfer.effectAllowed = 'move';
  }, []);

  // Validation panel handlers
  const handleHoverNode = useCallback((nodeId) => {
    setHighlightedNodeId(nodeId);
  }, []);

  const handleSelectNode = useCallback((nodeId) => {
    if (canvasRef.current) {
      canvasRef.current.centerOnNode(nodeId);
      canvasRef.current.openNodeEditor(nodeId);
    }
  }, []);

  // Quest properties editor handlers
  const handleOpenQuestEditor = useCallback(() => {
    setShowQuestEditor(true);
  }, []);

  const handleSaveQuestProperties = useCallback(async (updatedQuest) => {
    setQuest(updatedQuest);
    // Validate the updated quest
    try {
      const result = await api.validateQuest(updatedQuest);
      setValidation(result);
    } catch (e) {
      console.error('Validation failed:', e);
    }
  }, []);

  return (
    <div style={{ ...styles.app, backgroundColor: theme.canvasBg, color: theme.text }}>
      <TopBar
        questId={currentQuestId}
        quests={quests}
        onSelect={loadQuest}
        onNew={handleNewQuest}
        onSave={handleSave}
        validation={validation}
        saving={saving}
        saveError={saveError}
        onToggleTheme={toggleTheme}
      />
      <div style={styles.main}>
        <Toolbox
          quest={quest}
          onDragStart={handleDragStart}
          onEditQuest={handleOpenQuestEditor}
        />
        <Canvas
          ref={canvasRef}
          quest={quest}
          metadata={metadata}
          referenceData={referenceData}
          onChange={handleQuestChange}
          highlightedNodeId={highlightedNodeId}
        />
        <ValidationPanel
          validation={validation}
          onHoverNode={handleHoverNode}
          onSelectNode={handleSelectNode}
        />
      </div>

      {showQuestEditor && (
        <QuestPropertiesEditor
          quest={quest}
          onSave={handleSaveQuestProperties}
          onClose={() => setShowQuestEditor(false)}
        />
      )}
    </div>
  );
}

function App() {
  return (
    <ErrorBoundary>
      <ThemeProvider>
        <AppContent />
      </ThemeProvider>
    </ErrorBoundary>
  );
}

const styles = {
  app: {
    display: 'flex',
    flexDirection: 'column',
    height: '100vh',
  },
  main: {
    display: 'flex',
    flex: 1,
    overflow: 'hidden',
  },
};

export default App;
