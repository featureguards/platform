import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetch, selectProjectMembers } from '../../features/project_members/slice';
import { useNotifier } from '../hooks';
import { handleError } from './utils';

export function useProjectMembers(projectID?: string) {
  const notifier = useNotifier();
  const members = useAppSelector(selectProjectMembers);
  const status = useAppSelector((state) => state.projectMembers.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const refetch = async () => {
    if (!projectID) return;

    if (status === 'loading') {
      return;
    }
    try {
      await dispatch(fetch(projectID)).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

  useEffect(() => {
    refetch();
    // This isn't a bug. We only depend on projectID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [projectID]);

  return { members, loading: status === 'loading', refetch };
}
