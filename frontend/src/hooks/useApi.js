import { useState, useEffect, useCallback } from 'react';
import * as api from '../api/client';

export function useQuests() {
  const [quests, setQuests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.fetchQuests();
      setQuests(data || []);
      setError(null);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { refresh(); }, [refresh]);

  return { quests, loading, error, refresh };
}

export function useReferenceData() {
  const [data, setData] = useState({ items: [], factions: [], resources: [], npcs: [], objects: [] });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.fetchItems(),
      api.fetchFactions(),
      api.fetchResources(),
      api.fetchNPCs(),
      api.fetchObjects(),
    ]).then(([items, factions, resources, npcs, objects]) => {
      setData({ items: items || [], factions: factions || [], resources: resources || [], npcs: npcs || [], objects: objects || [] });
      setLoading(false);
    }).catch(() => setLoading(false));
  }, []);

  return { ...data, loading };
}
