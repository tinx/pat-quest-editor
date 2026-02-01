import { useState, useCallback, useRef, useEffect, forwardRef, useImperativeHandle } from 'react';
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
import { useTheme } from '../ThemeContext';

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
        nextNodesIfTrue: qn.NextNodesIfTrue,
        nextNodesIfFalse: qn.NextNodesIfFalse,
        conditions: qn.Conditions,
        conditionsRequired: qn.ConditionsRequired,
        conversationPartner: qn.ConversationPartner,
        speaker: qn.Speaker,
        text: qn.Text,
        options: qn.Options,
        messages: qn.Messages,
        actions: qn.Actions,
        questStageTitle: qn.QuestStageTitle,
        questStageDescription: qn.QuestStageDescription,
      },
    };
  });

  const optionColors = ['#e91e63', '#9c27b0', '#673ab7', '#3f51b5', '#2196f3', '#00bcd4', '#009688', '#4caf50'];

  const edges = [];
  (quest?.QuestNodes || []).forEach(qn => {
    // For Decision, edges come from options only
    if (qn.NodeType === 'Decision') {
      (qn.Options || []).forEach((opt, i) => {
        (opt.NextNodes || []).forEach(nextId => {
          edges.push({
            id: `${qn.NodeID}-opt${i}-${nextId}`,
            source: String(qn.NodeID),
            sourceHandle: `option-${i}`,
            target: String(nextId),
            animated: true,
            style: { stroke: optionColors[i % optionColors.length], strokeWidth: 2 },
            interactionWidth: 20,
          });
        });
      });
    } else if (qn.NodeType === 'ConditionBranch') {
      // ConditionBranch has two branches: true and false
      (qn.NextNodesIfTrue || []).forEach(nextId => {
        edges.push({
          id: `${qn.NodeID}-true-${nextId}`,
          source: String(qn.NodeID),
          sourceHandle: 'branch-true',
          target: String(nextId),
          animated: true,
          style: { stroke: '#4caf50', strokeWidth: 2 },
          interactionWidth: 20,
        });
      });
      (qn.NextNodesIfFalse || []).forEach(nextId => {
        edges.push({
          id: `${qn.NodeID}-false-${nextId}`,
          source: String(qn.NodeID),
          sourceHandle: 'branch-false',
          target: String(nextId),
          animated: true,
          style: { stroke: '#f44336', strokeWidth: 2 },
          interactionWidth: 20,
        });
      });
    } else {
      // Regular nodes use top-level NextNodes
      (qn.NextNodes || []).forEach(nextId => {
        edges.push({
          id: `${qn.NodeID}-${nextId}`,
          source: String(qn.NodeID),
          target: String(nextId),
          animated: true,
        });
      });
    }
  });

  return { nodes, edges };
}

// Convert React Flow nodes/edges back to quest data
function flowToQuest(nodes, edges, originalQuest) {
  // Build map of regular edges (no sourceHandle)
  const nodeMap = {};
  nodes.forEach(n => { nodeMap[n.id] = []; });
  
  // Build map of option edges (with sourceHandle like "option-0")
  const optionEdgeMap = {}; // nodeId -> { optionIndex -> [targets] }
  
  // Build map of branch edges (with sourceHandle like "branch-true" or "branch-false")
  const branchTrueMap = {}; // nodeId -> [targets]
  const branchFalseMap = {}; // nodeId -> [targets]
  
  edges.forEach(e => {
    if (e.sourceHandle && e.sourceHandle.startsWith('option-')) {
      const optIndex = parseInt(e.sourceHandle.split('-')[1], 10);
      if (!optionEdgeMap[e.source]) optionEdgeMap[e.source] = {};
      if (!optionEdgeMap[e.source][optIndex]) optionEdgeMap[e.source][optIndex] = [];
      optionEdgeMap[e.source][optIndex].push(parseInt(e.target, 10));
    } else if (e.sourceHandle === 'branch-true') {
      if (!branchTrueMap[e.source]) branchTrueMap[e.source] = [];
      branchTrueMap[e.source].push(parseInt(e.target, 10));
    } else if (e.sourceHandle === 'branch-false') {
      if (!branchFalseMap[e.source]) branchFalseMap[e.source] = [];
      branchFalseMap[e.source].push(parseInt(e.target, 10));
    } else {
      if (nodeMap[e.source]) {
        nodeMap[e.source].push(parseInt(e.target, 10));
      }
    }
  });

  const questNodes = nodes.map(n => {
    const d = n.data;
    const node = {
      NodeID: d.nodeId,
      NodeType: d.nodeType,
    };

    // Decision: NextNodes go in options, not top-level
    if (d.nodeType === 'Decision') {
      const options = (d.options || []).map((opt, i) => {
        const optEdges = optionEdgeMap[n.id]?.[i] || [];
        return {
          ...opt,
          NextNodes: optEdges.length > 0 ? optEdges : opt.NextNodes,
        };
      });
      if (options.length > 0) node.Options = options;
    } else if (d.nodeType === 'ConditionBranch') {
      // ConditionBranch uses NextNodesIfTrue and NextNodesIfFalse
      const trueEdges = branchTrueMap[n.id] || [];
      const falseEdges = branchFalseMap[n.id] || [];
      if (trueEdges.length > 0) node.NextNodesIfTrue = trueEdges;
      if (falseEdges.length > 0) node.NextNodesIfFalse = falseEdges;
    } else {
      if (nodeMap[n.id]?.length > 0) node.NextNodes = nodeMap[n.id];
      if (d.options?.length) node.Options = d.options;
    }

    if (d.conditions?.length) node.Conditions = d.conditions;
    if (d.conditionsRequired) node.ConditionsRequired = d.conditionsRequired;
    if (d.conversationPartner) node.ConversationPartner = d.conversationPartner;
    if (d.speaker) node.Speaker = d.speaker;
    if (d.text) node.Text = d.text;
    if (d.messages?.length) node.Messages = d.messages;
    if (d.actions?.length) node.Actions = d.actions;
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
    nodePositions[parseInt(n.id, 10)] = { x: n.position.x, y: n.position.y };
  });
  return { questId, nodePositions };
}

export default forwardRef(function Canvas({ quest, metadata, questVersion, referenceData, onChange, highlightedNodeId }, ref) {
  const { theme } = useTheme();
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [editingNode, setEditingNode] = useState(null);
  const reactFlowWrapper = useRef(null);
  const [reactFlowInstance, setReactFlowInstance] = useState(null);
  const isInitialLoad = useRef(false);
  const questIdRef = useRef(null);
  const questVersionRef = useRef(questVersion);
  const questRef = useRef(quest);
  const changeTimeoutRef = useRef(null);

  // Keep questRef in sync with latest quest prop
  useEffect(() => {
    questRef.current = quest;
  }, [quest]);

  // Expose methods to parent via ref
  useImperativeHandle(ref, () => ({
    centerOnNode: (nodeId) => {
      if (reactFlowInstance) {
        const node = nodes.find(n => n.id === String(nodeId));
        if (node) {
          reactFlowInstance.setCenter(node.position.x + 100, node.position.y + 50, { zoom: 1, duration: 500 });
        }
      }
    },
    openNodeEditor: (nodeId) => {
      const node = nodes.find(n => n.id === String(nodeId));
      if (node) {
        setEditingNode(node);
      }
    },
  }), [reactFlowInstance, nodes]);

  // Load quest into canvas (when quest ID or version changes)
  useEffect(() => {
    const questIdChanged = quest && quest.QuestID !== questIdRef.current;
    const versionChanged = questVersion !== questVersionRef.current;
    
    if (quest && (questIdChanged || versionChanged)) {
      questIdRef.current = quest.QuestID;
      questVersionRef.current = questVersion;
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
  }, [quest, metadata, questVersion, setNodes, setEdges]);

  // Update node highlighting when highlightedNodeId changes
  useEffect(() => {
    setNodes(nds => nds.map(n => ({
      ...n,
      data: { ...n.data, highlighted: n.id === String(highlightedNodeId) },
    })));
  }, [highlightedNodeId, setNodes]);

  // Notify parent of changes (but not on initial load)
  // Use a shared debounce function to avoid race conditions
  const notifyChange = useCallback(() => {
    if (changeTimeoutRef.current) {
      clearTimeout(changeTimeoutRef.current);
    }
    changeTimeoutRef.current = setTimeout(() => {
      setNodes(currentNodes => {
        setEdges(currentEdges => {
          const currentQuest = questRef.current;
          if (currentQuest) {
            const updatedQuest = flowToQuest(currentNodes, currentEdges, currentQuest);
            const updatedMetadata = getMetadata(currentQuest.QuestID, currentNodes);
            onChange(updatedQuest, updatedMetadata);
          }
          return currentEdges;
        });
        return currentNodes;
      });
    }, 100);
  }, [onChange, setNodes, setEdges]);

  const handleNodesChange = useCallback((changes) => {
    onNodesChange(changes);
    if (!isInitialLoad.current && questRef.current) {
      notifyChange();
    }
  }, [onNodesChange, notifyChange]);

  const handleEdgesChange = useCallback((changes) => {
    onEdgesChange(changes);
    if (!isInitialLoad.current && questRef.current) {
      notifyChange();
    }
  }, [onEdgesChange, notifyChange]);

  const onConnect = useCallback(
    (params) => setEdges((eds) => addEdge({ ...params, animated: true }, eds)),
    [setEdges]
  );

  const onNodeDoubleClick = useCallback((_, node) => {
    setEditingNode(node);
  }, []);

  const onNodeSave = useCallback((updatedNode) => {
    setNodes(nds => nds.map(n => n.id === updatedNode.id ? updatedNode : n));
    // Use the same notifyChange mechanism as canvas changes
    notifyChange();
  }, [setNodes, notifyChange]);

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

    // Calculate max ID from existing nodes, handling NaN gracefully
    // Use a reasonable upper bound to prevent overflow issues
    const MAX_SAFE_NODE_ID = 999999;
    const maxId = nodes.reduce((max, n) => {
      const parsed = parseInt(n.id, 10);
      if (Number.isNaN(parsed) || parsed < 0) return max;
      return Math.max(max, parsed);
    }, -1);

    // Check if we can safely create a new ID
    if (maxId >= MAX_SAFE_NODE_ID) {
      console.error('Cannot create more nodes: maximum node ID limit reached');
      return;
    }

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
        style={{ backgroundColor: theme.canvasBg }}
      >
        <Background color={theme.canvasGrid} gap={20} />
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
})

const styles = {
  container: {
    flex: 1,
    height: '100%',
  },
};
