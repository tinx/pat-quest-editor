import { useState } from 'react';
import { useTheme } from '../ThemeContext';

const CONDITION_TYPES = [
  { value: 'QuestCompleted', label: 'Quest Completed' },
  { value: 'ResourceAvailability', label: 'Resource Availability' },
  { value: 'FactionStanding', label: 'Faction Standing' },
  { value: 'TimePassed', label: 'Time Passed' },
  { value: 'ItemLost', label: 'Item Lost' },
  { value: 'Inventory', label: 'Inventory Has Items' },
  { value: 'Variable', label: 'Variable Check' },
  { value: 'EventTriggered', label: 'Event Triggered' },
];

const TIME_UNITS = [
  { value: 'h', label: 'Hours' },
  { value: 'd', label: 'Days' },
  { value: 'w', label: 'Weeks' },
  { value: 'M', label: 'Months' },
  { value: 'y', label: 'Years' },
];

const COMPARISONS = [
  { value: 'equal', label: '=' },
  { value: 'not equal', label: '≠' },
  { value: 'greater than', label: '>' },
  { value: 'smaller than', label: '<' },
];

// Get condition type from condition object
function getConditionType(condition) {
  if (!condition) return null;
  const keys = Object.keys(condition);
  return keys.find(k => CONDITION_TYPES.some(ct => ct.value === k)) || null;
}

// Create empty condition of given type
function createEmptyCondition(type) {
  switch (type) {
    case 'QuestCompleted':
      return { QuestCompleted: '' };
    case 'ResourceAvailability':
      return { ResourceAvailability: { Resource: '', Available: true } };
    case 'FactionStanding':
      return { FactionStanding: { Faction: '' } };
    case 'TimePassed':
      return { TimePassed: '1h' };
    case 'ItemLost':
      return { ItemLost: '' };
    case 'Inventory':
      return { Inventory: [{ Type: '' }] };
    case 'Variable':
      return { Variable: { VariableName: '', Comparison: 'equal', Value: 0 } };
    case 'EventTriggered':
      return { EventTriggered: { Event: '', Count: 1 } };
    default:
      return {};
  }
}

// Render summary text for a condition
function getConditionSummary(condition, items, factions, resources) {
  const type = getConditionType(condition);
  if (!type) return 'Unknown condition';

  switch (type) {
    case 'QuestCompleted':
      return `Quest: ${condition.QuestCompleted || '(not set)'}`;
    case 'ResourceAvailability': {
      const ra = condition.ResourceAvailability;
      const res = resources?.find(r => r.ResourceID === ra?.Resource);
      const name = res?.DisplayName?.['en-US'] || ra?.Resource || '(not set)';
      const status = ra?.Available !== false ? 'available' : 'unavailable';
      return `Resource: ${name} (${status})`;
    }
    case 'FactionStanding': {
      const fs = condition.FactionStanding;
      const fac = factions?.find(f => f.FactionID === fs?.Faction);
      const name = fac?.DisplayName?.['en-US'] || fs?.Faction || '(not set)';
      const levels = [];
      if (fs?.MinimumLevel) levels.push(`≥${fs.MinimumLevel}`);
      if (fs?.MaximumLevel) levels.push(`≤${fs.MaximumLevel}`);
      return `Faction: ${name}${levels.length ? ` (${levels.join(', ')})` : ''}`;
    }
    case 'TimePassed':
      return `Time: ${condition.TimePassed || '(not set)'}`;
    case 'ItemLost': {
      const item = items?.find(i => i.ItemID === condition.ItemLost);
      return `Item Lost: ${item?.DisplayName?.['en-US'] || condition.ItemLost || '(not set)'}`;
    }
    case 'Inventory': {
      const inv = condition.Inventory || [];
      if (inv.length === 0) return 'Inventory: (empty)';
      const itemNames = inv.map(i => {
        const item = items?.find(it => it.ItemID === i.Type);
        const name = item?.DisplayName?.['en-US'] || i.Type || '?';
        return i.MinCount ? `${i.MinCount}x ${name}` : name;
      });
      return `Has: ${itemNames.join(', ')}`;
    }
    case 'Variable': {
      const v = condition.Variable;
      const comp = COMPARISONS.find(c => c.value === v?.Comparison)?.label || v?.Comparison;
      return `${v?.VariableName || '?'} ${comp} ${v?.Value ?? '?'}`;
    }
    case 'EventTriggered': {
      const et = condition.EventTriggered;
      return `Event: ${et?.Event || '(not set)'} ×${et?.Count ?? '?'}`;
    }
    default:
      return 'Unknown';
  }
}

// Individual condition type editors
function QuestCompletedEditor({ value, onChange, styles }) {
  return (
    <input
      type="text"
      value={value || ''}
      onChange={e => onChange(e.target.value)}
      placeholder="Quest ID (e.g., PAT_Tutorial:Navigation)"
      style={styles.input}
    />
  );
}

function ResourceAvailabilityEditor({ value, onChange, resources, styles }) {
  const ra = value || { Resource: '', Available: true };
  const update = (field, val) => onChange({ ...ra, [field]: val });

  return (
    <div style={styles.fieldGroup}>
      <select value={ra.Resource || ''} onChange={e => update('Resource', e.target.value)} style={styles.select}>
        <option value="">Select resource...</option>
        {resources?.map(r => (
          <option key={r.ResourceID} value={r.ResourceID}>
            {r.DisplayName?.['en-US'] || r.ResourceID}
          </option>
        ))}
      </select>
      <label style={styles.checkboxLabel}>
        <input
          type="checkbox"
          checked={ra.Available !== false}
          onChange={e => update('Available', e.target.checked)}
        />
        Available
      </label>
    </div>
  );
}

function FactionStandingEditor({ value, onChange, factions, styles }) {
  const fs = value || { Faction: '' };
  const update = (field, val) => onChange({ ...fs, [field]: val });

  return (
    <div style={styles.fieldGroup}>
      <select
        value={fs.Faction || ''}
        onChange={e => update('Faction', e.target.value)}
        style={styles.select}
      >
        <option value="">Select faction...</option>
        {factions?.map(f => (
          <option key={f.FactionID} value={f.FactionID}>
            {f.DisplayName?.['en-US'] || f.FactionID}
          </option>
        ))}
      </select>
      <div style={styles.levelRow}>
        <label style={styles.smallLabel}>Min Level:</label>
        <input
          type="number"
          min="1"
          value={fs.MinimumLevel || ''}
          onChange={e => update('MinimumLevel', e.target.value ? parseInt(e.target.value) : undefined)}
          style={styles.numberInput}
          placeholder="-"
        />
        <label style={styles.smallLabel}>Max Level:</label>
        <input
          type="number"
          min="2"
          value={fs.MaximumLevel || ''}
          onChange={e => update('MaximumLevel', e.target.value ? parseInt(e.target.value) : undefined)}
          style={styles.numberInput}
          placeholder="-"
        />
      </div>
    </div>
  );
}

function TimePassedEditor({ value, onChange, styles }) {
  // Parse value like "36h" into number and unit
  const match = (value || '1h').match(/^(\d+)([hdwMy])$/);
  const num = match ? parseInt(match[1]) : 1;
  const unit = match ? match[2] : 'h';

  const update = (newNum, newUnit) => {
    onChange(`${newNum || 1}${newUnit || 'h'}`);
  };

  return (
    <div style={styles.inlineGroup}>
      <input
        type="number"
        min="1"
        value={num}
        onChange={e => update(parseInt(e.target.value) || 1, unit)}
        style={styles.numberInput}
      />
      <select value={unit} onChange={e => update(num, e.target.value)} style={styles.select}>
        {TIME_UNITS.map(u => (
          <option key={u.value} value={u.value}>{u.label}</option>
        ))}
      </select>
    </div>
  );
}

function ItemLostEditor({ value, onChange, items, styles }) {
  return (
    <select value={value || ''} onChange={e => onChange(e.target.value)} style={styles.select}>
      <option value="">Select item...</option>
      {items?.map(i => (
        <option key={i.ItemID} value={i.ItemID}>
          {i.DisplayName?.['en-US'] || i.ItemID}
        </option>
      ))}
    </select>
  );
}

function InventoryEditor({ value, onChange, items, styles }) {
  const inventory = value || [];

  const updateItem = (index, field, val) => {
    const newInv = [...inventory];
    newInv[index] = { ...newInv[index], [field]: val };
    onChange(newInv);
  };

  const addItem = () => {
    onChange([...inventory, { Type: '' }]);
  };

  const removeItem = (index) => {
    onChange(inventory.filter((_, i) => i !== index));
  };

  return (
    <div style={styles.inventoryList}>
      {inventory.map((inv, i) => (
        <div key={i} style={styles.inventoryRow}>
          <select
            value={inv.Type || ''}
            onChange={e => updateItem(i, 'Type', e.target.value)}
            style={{ ...styles.select, flex: 1 }}
          >
            <option value="">Select item...</option>
            {items?.map(item => (
              <option key={item.ItemID} value={item.ItemID}>
                {item.DisplayName?.['en-US'] || item.ItemID}
              </option>
            ))}
          </select>
          <input
            type="number"
            min="1"
            value={inv.MinCount || ''}
            onChange={e => updateItem(i, 'MinCount', e.target.value ? parseInt(e.target.value) : undefined)}
            style={styles.numberInput}
            placeholder="Count"
          />
          <label style={styles.checkboxLabel}>
            <input
              type="checkbox"
              checked={inv.QuestItem || false}
              onChange={e => updateItem(i, 'QuestItem', e.target.checked || undefined)}
            />
            Quest
          </label>
          <button onClick={() => removeItem(i)} style={styles.removeItemBtn}>×</button>
        </div>
      ))}
      <button onClick={addItem} style={styles.addItemBtn}>+ Add Item</button>
    </div>
  );
}

function VariableEditor({ value, onChange, styles }) {
  const v = value || { VariableName: '', Comparison: 'equal', Value: 0 };
  const update = (field, val) => onChange({ ...v, [field]: val });

  return (
    <div style={styles.fieldGroup}>
      <input
        type="text"
        value={v.VariableName || ''}
        onChange={e => update('VariableName', e.target.value)}
        placeholder="Variable name"
        style={styles.input}
      />
      <div style={styles.inlineGroup}>
        <select
          value={v.Comparison || 'equal'}
          onChange={e => update('Comparison', e.target.value)}
          style={styles.select}
        >
          {COMPARISONS.map(c => (
            <option key={c.value} value={c.value}>{c.label} ({c.value})</option>
          ))}
        </select>
        <input
          type="number"
          value={v.Value ?? 0}
          onChange={e => update('Value', parseInt(e.target.value) || 0)}
          style={styles.numberInput}
        />
      </div>
    </div>
  );
}

function EventTriggeredEditor({ value, onChange, styles }) {
  const et = value || { Event: '', Count: 1 };
  const update = (field, val) => onChange({ ...et, [field]: val });

  return (
    <div style={styles.fieldGroup}>
      <input
        type="text"
        value={et.Event || ''}
        onChange={e => update('Event', e.target.value)}
        placeholder="Event name (e.g., Harvested:Rose)"
        style={styles.input}
      />
      <div style={styles.inlineGroup}>
        <label style={styles.smallLabel}>At least</label>
        <input
          type="number"
          min="1"
          value={et.Count || 1}
          onChange={e => update('Count', parseInt(e.target.value) || 1)}
          style={styles.numberInput}
        />
        <label style={styles.smallLabel}>times</label>
      </div>
    </div>
  );
}

// Single condition editor row
function ConditionRow({ condition, onChange, onRemove, items, factions, resources, expanded, onToggle, styles }) {
  const type = getConditionType(condition);
  const summary = getConditionSummary(condition, items, factions, resources);

  const updateConditionValue = (newValue) => {
    onChange({ [type]: newValue });
  };

  return (
    <div style={styles.conditionRow}>
      <div style={styles.conditionHeader} onClick={onToggle}>
        <span style={styles.conditionType}>
          {CONDITION_TYPES.find(ct => ct.value === type)?.label || type}
        </span>
        <span style={styles.conditionSummary}>{summary}</span>
        <span style={styles.expandIcon}>{expanded ? '▼' : '▶'}</span>
        <button onClick={(e) => { e.stopPropagation(); onRemove(); }} style={styles.removeBtn}>×</button>
      </div>
      
      {expanded && (
        <div style={styles.conditionBody}>
          {type === 'QuestCompleted' && (
            <QuestCompletedEditor value={condition.QuestCompleted} onChange={updateConditionValue} styles={styles} />
          )}
          {type === 'ResourceAvailability' && (
            <ResourceAvailabilityEditor value={condition.ResourceAvailability} onChange={updateConditionValue} resources={resources} styles={styles} />
          )}
          {type === 'FactionStanding' && (
            <FactionStandingEditor value={condition.FactionStanding} onChange={updateConditionValue} factions={factions} styles={styles} />
          )}
          {type === 'TimePassed' && (
            <TimePassedEditor value={condition.TimePassed} onChange={updateConditionValue} styles={styles} />
          )}
          {type === 'ItemLost' && (
            <ItemLostEditor value={condition.ItemLost} onChange={updateConditionValue} items={items} styles={styles} />
          )}
          {type === 'Inventory' && (
            <InventoryEditor value={condition.Inventory} onChange={updateConditionValue} items={items} styles={styles} />
          )}
          {type === 'Variable' && (
            <VariableEditor value={condition.Variable} onChange={updateConditionValue} styles={styles} />
          )}
          {type === 'EventTriggered' && (
            <EventTriggeredEditor value={condition.EventTriggered} onChange={updateConditionValue} styles={styles} />
          )}
        </div>
      )}
    </div>
  );
}

// Main ConditionEditor component
export default function ConditionEditor({ 
  conditions, 
  onChange, 
  items, 
  factions, 
  resources,
  conditionsRequired,
  onConditionsRequiredChange,
  showConditionsRequired = false,
  collapsible = false,
  defaultExpanded = true,
  excludeConditionTypes = [],
}) {
  const { theme } = useTheme();
  const styles = getStyles(theme);
  const [expandedConditions, setExpandedConditions] = useState({});
  const [listExpanded, setListExpanded] = useState(defaultExpanded);
  const [addingType, setAddingType] = useState('');

  // Filter condition types based on excludeConditionTypes prop
  const availableConditionTypes = CONDITION_TYPES.filter(
    ct => !excludeConditionTypes.includes(ct.value)
  );

  const toggleCondition = (index) => {
    setExpandedConditions(prev => ({ ...prev, [index]: !prev[index] }));
  };

  const updateCondition = (index, newCondition) => {
    const newConditions = [...(conditions || [])];
    newConditions[index] = newCondition;
    onChange(newConditions);
  };

  const removeCondition = (index) => {
    onChange((conditions || []).filter((_, i) => i !== index));
  };

  const addCondition = () => {
    if (!addingType) return;
    const newCondition = createEmptyCondition(addingType);
    const newConditions = [...(conditions || []), newCondition];
    onChange(newConditions);
    setExpandedConditions(prev => ({ ...prev, [newConditions.length - 1]: true }));
    setAddingType('');
  };

  const conditionCount = (conditions || []).length;

  // Generate options for ConditionsRequired dropdown
  const requiredOptions = [
    { value: 'all', label: 'All conditions must be met' },
    ...Array.from({ length: conditionCount }, (_, i) => ({
      value: String(i + 1),
      label: `At least ${i + 1} condition${i > 0 ? 's' : ''} must be met`,
    })),
  ];

  const content = (
    <>
      {showConditionsRequired && conditionCount > 1 && (
        <div style={styles.requiredRow}>
          <select
            value={conditionsRequired || 'all'}
            onChange={e => onConditionsRequiredChange(e.target.value)}
            style={styles.select}
          >
            {requiredOptions.map(opt => (
              <option key={opt.value} value={opt.value}>{opt.label}</option>
            ))}
          </select>
        </div>
      )}

      {(conditions || []).map((cond, i) => (
        <ConditionRow
          key={i}
          condition={cond}
          onChange={(newCond) => updateCondition(i, newCond)}
          onRemove={() => removeCondition(i)}
          items={items}
          factions={factions}
          resources={resources}
          expanded={expandedConditions[i] !== false}
          onToggle={() => toggleCondition(i)}
          styles={styles}
        />
      ))}

      {conditionCount === 0 && (
        <div style={styles.noConditions}>No conditions defined</div>
      )}

      <div style={styles.addRow}>
        <select
          value={addingType}
          onChange={e => setAddingType(e.target.value)}
          style={{ ...styles.select, flex: 1 }}
        >
          <option value="">Select condition type...</option>
          {availableConditionTypes.map(ct => (
            <option key={ct.value} value={ct.value}>{ct.label}</option>
          ))}
        </select>
        <button
          onClick={addCondition}
          disabled={!addingType}
          style={{ ...styles.addBtn, opacity: addingType ? 1 : 0.5 }}
        >
          + Add
        </button>
      </div>
    </>
  );

  if (collapsible) {
    return (
      <div style={styles.collapsibleContainer}>
        <div style={styles.collapsibleHeader} onClick={() => setListExpanded(!listExpanded)}>
          <span>{listExpanded ? '▼' : '▶'} Conditions ({conditionCount})</span>
        </div>
        {listExpanded && <div style={styles.collapsibleBody}>{content}</div>}
      </div>
    );
  }

  return <div style={styles.container}>{content}</div>;
}

const getStyles = (theme) => ({
  container: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
  collapsibleContainer: {
    backgroundColor: theme.bg,
    borderRadius: '4px',
    overflow: 'hidden',
  },
  collapsibleHeader: {
    padding: '8px 12px',
    cursor: 'pointer',
    color: theme.textMuted,
    fontSize: '11px',
    fontWeight: 'bold',
    backgroundColor: theme.bgTertiary,
  },
  collapsibleBody: {
    padding: '8px',
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
  requiredRow: {
    marginBottom: '8px',
  },
  conditionRow: {
    backgroundColor: theme.bgTertiary,
    borderRadius: '4px',
    overflow: 'hidden',
  },
  conditionHeader: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    padding: '8px 10px',
    cursor: 'pointer',
  },
  conditionType: {
    color: '#4fc3f7',
    fontSize: '10px',
    fontWeight: 'bold',
    textTransform: 'uppercase',
    flexShrink: 0,
  },
  conditionSummary: {
    color: theme.textSecondary,
    fontSize: '11px',
    flex: 1,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  expandIcon: {
    color: theme.textDim,
    fontSize: '10px',
  },
  removeBtn: {
    background: 'none',
    border: 'none',
    color: '#f44336',
    fontSize: '16px',
    cursor: 'pointer',
    padding: '0 4px',
    flexShrink: 0,
  },
  conditionBody: {
    padding: '8px 10px',
    borderTop: `1px solid ${theme.border}`,
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
  fieldGroup: {
    display: 'flex',
    flexDirection: 'column',
    gap: '6px',
  },
  inlineGroup: {
    display: 'flex',
    gap: '8px',
    alignItems: 'center',
  },
  levelRow: {
    display: 'flex',
    gap: '6px',
    alignItems: 'center',
    flexWrap: 'wrap',
  },
  smallLabel: {
    color: theme.textMuted,
    fontSize: '10px',
  },
  input: {
    padding: '6px 8px',
    backgroundColor: theme.inputBg,
    border: `1px solid ${theme.inputBorder}`,
    borderRadius: '4px',
    color: theme.text,
    fontSize: '12px',
  },
  select: {
    padding: '6px 8px',
    backgroundColor: theme.inputBg,
    border: `1px solid ${theme.inputBorder}`,
    borderRadius: '4px',
    color: theme.text,
    fontSize: '12px',
  },
  numberInput: {
    padding: '6px 8px',
    backgroundColor: theme.inputBg,
    border: `1px solid ${theme.inputBorder}`,
    borderRadius: '4px',
    color: theme.text,
    fontSize: '12px',
    width: '60px',
  },
  inventoryList: {
    display: 'flex',
    flexDirection: 'column',
    gap: '6px',
  },
  inventoryRow: {
    display: 'flex',
    gap: '6px',
    alignItems: 'center',
  },
  checkboxLabel: {
    display: 'flex',
    alignItems: 'center',
    gap: '4px',
    color: theme.textMuted,
    fontSize: '10px',
    flexShrink: 0,
  },
  removeItemBtn: {
    background: 'none',
    border: 'none',
    color: '#f44336',
    fontSize: '14px',
    cursor: 'pointer',
    padding: '0 2px',
  },
  addItemBtn: {
    padding: '4px 8px',
    backgroundColor: theme.bgTertiary,
    color: theme.textMuted,
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '10px',
    alignSelf: 'flex-start',
  },
  noConditions: {
    color: theme.textDim,
    fontSize: '11px',
    fontStyle: 'italic',
    padding: '8px',
    textAlign: 'center',
  },
  addRow: {
    display: 'flex',
    gap: '8px',
    marginTop: '4px',
  },
  addBtn: {
    padding: '6px 12px',
    backgroundColor: '#4caf50',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '11px',
    flexShrink: 0,
  },
});
