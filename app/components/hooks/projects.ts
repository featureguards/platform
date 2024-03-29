import { useRouter } from 'next/router';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetchAll } from '../../features/projects/slice';
import { useNotifier } from '../hooks';
import { handleError } from './utils';

export function useProjectsLazy() {
  const notifier = useNotifier();
  const status = useAppSelector((state) => state.projects.all.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const refetch = async () => {
    if (status === 'loading') {
      return;
    }
    try {
      await dispatch(fetchAll()).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

  return { loading: status === 'loading', refetch };
}
