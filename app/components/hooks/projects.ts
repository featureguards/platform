import { useCallback } from 'react';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetchAll } from '../../features/projects/slice';
import { useNotifier } from '../hooks';

export function useProjects() {
  const notifier = useNotifier();
  const projects = useAppSelector((state) => state.projects.all.items);
  const status = useAppSelector((state) => state.projects.all.status);
  const error = useAppSelector((state) => state.projects.all.error);
  const dispatch = useAppDispatch();

  console.log(`useProjects ${status}`);
  const fetchProjects = useCallback(async () => {
    if (status === 'succeeded' || status === 'failed' || status === 'loading') {
      return;
    }
    try {
      await dispatch(fetchAll()).unwrap();
    } catch (_err) {
      if (error) {
        notifier.error(error);
      }
    }
  }, [dispatch, error, notifier, status]);

  fetchProjects();

  return { projects, loading: status === 'loading' };
}
