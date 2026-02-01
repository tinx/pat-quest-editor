import { useState, useCallback, useRef } from 'react';

const DEFAULT_MAX_HISTORY = 50;

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
    setPast(prev => {
      const newPast = [...prev, state];
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
