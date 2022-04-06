import { useCallback } from 'react';

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
    if (status === 'loading' || status === 'succeeded' || status === 'failed') {
      return;
    }
    try {
      await dispatch(fetch(projectID)).unwrap();
    } catch (_err) {
      if (error) {
        notifier.error(error);
      }
    }
  }, [projectID, dispatch, error, notifier, status]);

  fetchProject();

  return { current, loading: status === 'loading' };
}
