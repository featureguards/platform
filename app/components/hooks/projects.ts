import { useCallback } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetchAll } from '../../features/projects/slice';
import { useNotifier } from '../hooks';

export function useProjects() {
  const notifier = useNotifier();
  const projects = useAppSelector((state) => state.projects.all.items);
  const status = useAppSelector((state) => state.projects.all.status);
  const dispatch = useAppDispatch();

  const fetchProjects = useCallback(async () => {
    if (status === 'succeeded' || status === 'failed' || status === 'loading') {
      return;
    }
    try {
      await dispatch(fetchAll()).unwrap();
    } catch (err) {
      const error = err as SerializedError;
      if (error.message && error.code !== '404') {
        notifier.error(error.message);
      }
    }
  }, [dispatch, notifier, status]);

  fetchProjects();

  return { projects, loading: status === 'loading' };
}
