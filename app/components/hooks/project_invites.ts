import { useRouter } from 'next/router';
import { useCallback, useEffect } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetchForProject, fetchForUser } from '../../features/project_invites/slice';
import { useNotifier } from '../hooks';
import { handleError } from './utils';

export function useUserInvites() {
  const notifier = useNotifier();
  const invites = useAppSelector((state) => state.projectInvites.forUser.items);
  const status = useAppSelector((state) => state.projectInvites.forUser.status);
  const error = useAppSelector((state) => state.projectInvites.forUser.error);
  const dispatch = useAppDispatch();

  const fetch = useCallback(async () => {
    if (status === 'succeeded' || status === 'failed' || status === 'loading') {
      return;
    }
    try {
      await dispatch(fetchForUser('me')).unwrap();
    } catch (_err) {
      if (error) {
        notifier.error(error);
      }
    }
  }, [dispatch, error, notifier, status]);

  fetch();

  return { invites, loading: status === 'loading' };
}

export function useProjectInvites(projectID: string) {
  const notifier = useNotifier();
  const invites = useAppSelector((state) => state.projectInvites.forProject.items);
  const status = useAppSelector((state) => state.projectInvites.forProject.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const refetch = async () => {
    if (status === 'loading') {
      return;
    }
    try {
      await dispatch(fetchForProject(projectID)).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

  useEffect(() => {
    refetch();
    // This isn't a bug. We only depend on projectID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [projectID]);

  return { invites, loading: status === 'loading', refetch };
}
