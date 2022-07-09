import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetchForProject, fetchForUser } from '../../features/project_invites/slice';
import { useNotifier } from '../hooks';
import { handleError } from './utils';

export function useUserInvites() {
  const invites = useAppSelector((state) => state.projectInvites.forUser.items);
  const status = useAppSelector((state) => state.projectInvites.forUser.status);
  const dispatch = useAppDispatch();

  const refetch = async () => {
    if (status === 'loading') {
      return;
    }
    try {
      await dispatch(fetchForUser('me')).unwrap();
    } catch (_err) {}
  };

  useEffect(() => {
    refetch();
  }, []);

  return { invites, loading: status === 'loading', refetch };
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
