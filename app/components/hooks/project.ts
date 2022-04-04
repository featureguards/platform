import { useCallback, useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetch } from '../../features/projects/slice';
import { useNotifier } from '../hooks';

export function useProject({ projectID }: { projectID?: string }) {
  const notifier = useNotifier();
  const current = useAppSelector((state) => state.projects.details.item);
  const status = useAppSelector((state) => state.projects.details.status);
  const error = useAppSelector((state) => state.projects.details.error);
  const dispatch = useAppDispatch();

  // https://stackoverflow.com/questions/53332321/react-hook-warnings-for-async-function-in-useeffect-useeffect-function-must-ret
  const fetchProject = useCallback(async () => {
    if (!projectID) {
      return;
    }
    try {
      await dispatch(fetch(projectID)).unwrap();
    } catch (_err) {
      if (error) {
        notifier.error(error);
      }
    }
  }, [projectID, dispatch, error, notifier]);

  useEffect(() => {
    if (!projectID) {
      return;
    }
    if (projectID === current?.id) {
      return;
    }
    if (status === 'loading') {
      return;
    }
    if (status === 'succeeded' || status === 'failed') {
      return;
    }
    fetchProject();
  }, [current, status, error, fetchProject, projectID]);

  return { current, loading: status === 'loading' };
}
