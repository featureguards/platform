import { useCallback, useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetchAll } from '../../features/projects/slice';
import { useNotifier } from '../hooks';

export function useProjects({ force }: { force?: boolean }) {
  const notifier = useNotifier();
  const projects = useAppSelector((state) => state.projects.all.items);
  const status = useAppSelector((state) => state.projects.all.status);
  const error = useAppSelector((state) => state.projects.all.error);
  const dispatch = useAppDispatch();

  const fetchProjects = useCallback(async () => {
    try {
      await dispatch(fetchAll()).unwrap();
    } catch (_err) {
      if (error) {
        notifier.error(error);
      }
    }
  }, [dispatch, error, notifier]);

  useEffect(() => {
    if (!force && status === 'succeeded') {
      return;
    }
    if (status === 'loading') {
      return;
    }

    fetchProjects();
  }, [force, projects, status, error, fetchProjects]);

  return { projects, loading: status === 'loading' };
}
