import { useState, useCallback, useEffect } from 'react';
import TopBar from './components/TopBar';
import Toolbox from './components/Toolbox';
import Canvas from './components/Canvas';
import ValidationPanel from './components/ValidationPanel';
import { useQuests, useReferenceData } from './hooks/useApi';
import * as api from './api/client';

function App() {
  const { quests, refresh: refreshQuests } = useQuests();
  const referenceData = useReferenceData();
  
  const [currentQuestId, setCurrentQuestId] = useState(null);
  const [quest, setQuest] = useState(null);
  const [metadata, setMetadata] = useState(null);
  const [validation, setValidation] = useState(null);
  const [saving, setSaving] = useState(false);

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

  return (
    <div style={styles.app}>
      <TopBar
        questId={currentQuestId}
        quests={quests}
        onSelect={loadQuest}
        onNew={handleNewQuest}
        onSave={handleSave}
        validation={validation}
        saving={saving}
      />
      <div style={styles.main}>
        <Toolbox onDragStart={handleDragStart} />
        <Canvas
          quest={quest}
          metadata={metadata}
          referenceData={referenceData}
          onChange={handleQuestChange}
        />
        <ValidationPanel validation={validation} />
      </div>
    </div>
  );
}

const styles = {
  app: {
    display: 'flex',
    flexDirection: 'column',
    height: '100vh',
    backgroundColor: '#16162a',
    color: '#fff',
  },
  main: {
    display: 'flex',
    flex: 1,
    overflow: 'hidden',
  },
};

export default App;
