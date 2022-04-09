import { useCallback } from 'react';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetchForProject, fetchForUser } from '../../features/project_invites/slice';
import { useNotifier } from '../hooks';

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
  const invites = useAppSelector((state) => state.projectInvites.forUser.items);
  const status = useAppSelector((state) => state.projectInvites.forUser.status);
  const error = useAppSelector((state) => state.projectInvites.forUser.error);
  const dispatch = useAppDispatch();

  const fetch = useCallback(async () => {
    if (status === 'succeeded' || status === 'failed' || status === 'loading') {
      return;
    }
    try {
      await dispatch(fetchForProject(projectID)).unwrap();
    } catch (_err) {
      if (error) {
        notifier.error(error);
      }
    }
  }, [dispatch, error, notifier, status, projectID]);

  fetch();

  return { invites, loading: status === 'loading' };
}
