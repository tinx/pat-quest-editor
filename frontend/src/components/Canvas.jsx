import { useState, useCallback, useRef, useEffect } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  addEdge,
  useNodesState,
  useEdgesState,
  MarkerType,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

import QuestNode from './QuestNode';
import NodeEditor from './NodeEditor';

const nodeTypes = { questNode: QuestNode };

const defaultEdgeOptions = {
  style: { strokeWidth: 2 },
  markerEnd: { type: MarkerType.ArrowClosed },
};

// Convert quest data to React Flow nodes/edges
function questToFlow(quest, metadata) {
  const nodes = (quest?.QuestNodes || []).map((qn, idx) => {
    const pos = metadata?.nodePositions?.[qn.NodeID] || { x: 100 + (idx % 5) * 220, y: 100 + Math.floor(idx / 5) * 150 };
    return {
      id: String(qn.NodeID),
      type: 'questNode',
      position: pos,
      data: {
        nodeId: qn.NodeID,
        nodeType: qn.NodeType,
        nextNodes: qn.NextNodes,
        conditions: qn.Conditions,
        conditionsRequired: qn.ConditionsRequired,
        conversationPartner: qn.ConversationPartner,
        speaker: qn.Speaker,
        text: qn.Text,
        options: qn.Options,
        messages: qn.Messages,
        actions: qn.Actions,
        questProgressors: qn.QuestProgressors,
        questStageTitle: qn.QuestStageTitle,
        questStageDescription: qn.QuestStageDescription,
      },
    };
  });

  const edges = [];
  (quest?.QuestNodes || []).forEach(qn => {
    (qn.NextNodes || []).forEach(nextId => {
      edges.push({
        id: `${qn.NodeID}-${nextId}`,
        source: String(qn.NodeID),
        target: String(nextId),
        animated: true,
      });
    });
    // Also handle dialog options
    (qn.Options || []).forEach((opt, i) => {
      (opt.NextNodes || []).forEach(nextId => {
        edges.push({
          id: `${qn.NodeID}-opt${i}-${nextId}`,
          source: String(qn.NodeID),
          target: String(nextId),
          animated: true,
          style: { stroke: '#e91e63' },
        });
      });
    });
  });

  return { nodes, edges };
}

// Convert React Flow nodes/edges back to quest data
function flowToQuest(nodes, edges, originalQuest) {
  const nodeMap = {};
  nodes.forEach(n => { nodeMap[n.id] = []; });
  edges.forEach(e => {
    if (nodeMap[e.source]) {
      nodeMap[e.source].push(parseInt(e.target));
    }
  });

  const questNodes = nodes.map(n => {
    const d = n.data;
    const node = {
      NodeID: d.nodeId,
      NodeType: d.nodeType,
    };

    // PlayerDecisionDialog: NextNodes go in options, not top-level
    if (d.nodeType === 'PlayerDecisionDialog') {
      // Don't set top-level NextNodes for PlayerDecisionDialog
      if (d.options?.length) {
        // Distribute edges to options (simplified: all edges go to first option without NextNodes)
        const edgeTargets = nodeMap[n.id] || [];
        const options = d.options.map((opt, i) => {
          // Keep existing NextNodes if present, otherwise leave empty for user to set
          return { ...opt };
        });
        node.Options = options;
      }
    } else {
      if (nodeMap[n.id]?.length > 0) node.NextNodes = nodeMap[n.id];
    }

    if (d.conditions?.length) node.Conditions = d.conditions;
    if (d.conditionsRequired) node.ConditionsRequired = d.conditionsRequired;
    if (d.conversationPartner) node.ConversationPartner = d.conversationPartner;
    if (d.speaker) node.Speaker = d.speaker;
    if (d.text) node.Text = d.text;
    if (d.nodeType !== 'PlayerDecisionDialog' && d.options?.length) node.Options = d.options;
    if (d.messages?.length) node.Messages = d.messages;
    if (d.actions?.length) node.Actions = d.actions;
    if (d.questProgressors?.length) node.QuestProgressors = d.questProgressors;
    if (d.questStageTitle) node.QuestStageTitle = d.questStageTitle;
    if (d.questStageDescription) node.QuestStageDescription = d.questStageDescription;

    return node;
  });

  return {
    ...originalQuest,
    QuestNodes: questNodes,
  };
}

// Extract metadata (positions) from nodes
function getMetadata(questId, nodes) {
  const nodePositions = {};
  nodes.forEach(n => {
    nodePositions[parseInt(n.id)] = { x: n.position.x, y: n.position.y };
  });
  return { questId, nodePositions };
}

export default function Canvas({ quest, metadata, referenceData, onChange }) {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [editingNode, setEditingNode] = useState(null);
  const reactFlowWrapper = useRef(null);
  const [reactFlowInstance, setReactFlowInstance] = useState(null);
  const isInitialLoad = useRef(false);
  const questIdRef = useRef(null);

  // Load quest into canvas (only when quest ID changes)
  useEffect(() => {
    if (quest && quest.QuestID !== questIdRef.current) {
      questIdRef.current = quest.QuestID;
      isInitialLoad.current = true;
      const { nodes: n, edges: e } = questToFlow(quest, metadata);
      setNodes(n);
      setEdges(e);
      // Reset flag after a tick
      setTimeout(() => { isInitialLoad.current = false; }, 0);
    } else if (!quest) {
      questIdRef.current = null;
      setNodes([]);
      setEdges([]);
    }
  }, [quest, metadata, setNodes, setEdges]);

  // Notify parent of changes (but not on initial load)
  const handleNodesChange = useCallback((changes) => {
    onNodesChange(changes);
    if (!isInitialLoad.current && quest) {
      // Debounce the onChange to avoid too many updates
      setTimeout(() => {
        setNodes(currentNodes => {
          setEdges(currentEdges => {
            const updatedQuest = flowToQuest(currentNodes, currentEdges, quest);
            const updatedMetadata = getMetadata(quest.QuestID, currentNodes);
            onChange(updatedQuest, updatedMetadata);
            return currentEdges;
          });
          return currentNodes;
        });
      }, 100);
    }
  }, [onNodesChange, quest, onChange, setNodes, setEdges]);

  const handleEdgesChange = useCallback((changes) => {
    onEdgesChange(changes);
    if (!isInitialLoad.current && quest) {
      setTimeout(() => {
        setNodes(currentNodes => {
          setEdges(currentEdges => {
            const updatedQuest = flowToQuest(currentNodes, currentEdges, quest);
            const updatedMetadata = getMetadata(quest.QuestID, currentNodes);
            onChange(updatedQuest, updatedMetadata);
            return currentEdges;
          });
          return currentNodes;
        });
      }, 100);
    }
  }, [onEdgesChange, quest, onChange, setNodes, setEdges]);

  const onConnect = useCallback(
    (params) => setEdges((eds) => addEdge({ ...params, animated: true }, eds)),
    [setEdges]
  );

  const onNodeDoubleClick = useCallback((_, node) => {
    setEditingNode(node);
  }, []);

  const onNodeSave = useCallback((updatedNode) => {
    setNodes(nds =>
      nds.map(n => n.id === updatedNode.id ? updatedNode : n)
    );
  }, [setNodes]);

  const onDragOver = useCallback((e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback((e) => {
    e.preventDefault();

    const type = e.dataTransfer.getData('application/reactflow');
    if (!type || !reactFlowInstance) return;

    const position = reactFlowInstance.screenToFlowPosition({
      x: e.clientX,
      y: e.clientY,
    });

    const maxId = nodes.reduce((max, n) => Math.max(max, parseInt(n.id) || 0), -1);
    const newId = maxId + 1;

    const newNode = {
      id: String(newId),
      type: 'questNode',
      position,
      data: {
        nodeId: newId,
        nodeType: type,
      },
    };

    setNodes((nds) => nds.concat(newNode));
  }, [reactFlowInstance, nodes, setNodes]);

  return (
    <div ref={reactFlowWrapper} style={styles.container}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={handleNodesChange}
        onEdgesChange={handleEdgesChange}
        onConnect={onConnect}
        onNodeDoubleClick={onNodeDoubleClick}
        onInit={setReactFlowInstance}
        onDragOver={onDragOver}
        onDrop={onDrop}
        nodeTypes={nodeTypes}
        defaultEdgeOptions={defaultEdgeOptions}
        deleteKeyCode={['Backspace', 'Delete', 'Del']}
        selectionOnDrag
        panOnDrag={[1, 2]}
        selectNodesOnDrag={false}
        fitView
        style={{ backgroundColor: '#16162a' }}
      >
        <Background color="#333" gap={20} />
        <Controls />
      </ReactFlow>

      {editingNode && (
        <NodeEditor
          node={editingNode}
          npcs={referenceData.npcs}
          items={referenceData.items}
          factions={referenceData.factions}
          resources={referenceData.resources}
          onSave={onNodeSave}
          onClose={() => setEditingNode(null)}
        />
      )}
    </div>
  );
}

const styles = {
  container: {
    flex: 1,
    height: '100%',
  },
};
