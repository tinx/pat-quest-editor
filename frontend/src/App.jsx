import { useState, useCallback, useEffect, useRef } from 'react';
import { ThemeProvider, useTheme } from './ThemeContext';
import TopBar from './components/TopBar';
import Toolbox from './components/Toolbox';
import Canvas from './components/Canvas';
import ValidationPanel from './components/ValidationPanel';
import QuestPropertiesEditor from './components/QuestPropertiesEditor';
import { useQuests, useReferenceData } from './hooks/useApi';
import * as api from './api/client';

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
  const [highlightedNodeId, setHighlightedNodeId] = useState(null);
  const [showQuestEditor, setShowQuestEditor] = useState(false);

  // Load quest when selected
  const loadQuest = useCallback(async (questId) => {
    if (!questId) {
      setCurrentQuestId(null);
      setQuest(null);
      setMetadata(null);
      setValidation(null);
      return;
    }

    try {
      const data = await api.fetchQuest(questId);
      setCurrentQuestId(questId);
      setQuest(data.quest);
      setMetadata(data.metadata);
      
      // Validate on load
      const result = await api.validateQuest(data.quest);
      setValidation(result);
    } catch (e) {
      console.error('Failed to load quest:', e);
    }
  }, []);

  // Handle quest changes from canvas
  const handleQuestChange = useCallback(async (updatedQuest, updatedMetadata) => {
    setQuest(updatedQuest);
    setMetadata(updatedMetadata);

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
    try {
      const result = await api.saveQuest(currentQuestId, quest, metadata);
      setValidation(result);
    } catch (e) {
      console.error('Failed to save quest:', e);
    } finally {
      setSaving(false);
    }
  }, [currentQuestId, quest, metadata]);

  // Create new quest
  const handleNewQuest = useCallback(() => {
    const questId = prompt('Enter Quest ID (e.g., PAT_New_Quest):');
    if (!questId) return;

    const newQuest = {
      QuestTypeVersion: 1,
      QuestVersion: 1,
      QuestID: questId,
      QuestType: 'SideQuest',
      DisplayName: { 'en-US': 'New Quest', 'de-DE': 'Neue Aufgabe' },
      JournalEntry: { 'en-US': '', 'de-DE': '' },
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
  }, []);

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
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
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
