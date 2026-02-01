import { useState, useCallback, useRef } from 'react';

const DEFAULT_MAX_HISTORY = 50;

/**
 * Deep clone an object using JSON serialization.
 * Works for plain objects without circular references or special types.
 */
function deepClone(obj) {
  if (obj === null || obj === undefined) return obj;
  return JSON.parse(JSON.stringify(obj));
}

/**
 * Hook for managing undo history of state snapshots.
 * @param {number} maxHistory - Maximum number of history entries to keep
 * @returns {object} - { pushState, undo, clear, canUndo }
 */
export function useUndoHistory(maxHistory = DEFAULT_MAX_HISTORY) {
  const [past, setPast] = useState([]);
  const pastRef = useRef(past);
  pastRef.current = past;

  const pushState = useCallback((state) => {
    if (!state) return;
    // Deep clone to prevent mutations from affecting history
    const clonedState = deepClone(state);
    setPast(prev => {
      const newPast = [...prev, clonedState];
      // Trim to max history size
      if (newPast.length > maxHistory) {
        return newPast.slice(newPast.length - maxHistory);
      }
      return newPast;
    });
  }, [maxHistory]);

  const undo = useCallback(() => {
    const current = pastRef.current;
    if (current.length === 0) return null;
    
    const previous = current[current.length - 1];
    setPast(prev => prev.slice(0, -1));
    return previous;
  }, []);

  const clear = useCallback(() => {
    setPast([]);
  }, []);

  return {
    pushState,
    undo,
    clear,
    canUndo: past.length > 0,
  };
}
