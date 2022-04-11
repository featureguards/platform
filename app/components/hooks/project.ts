import { useEffect } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetch } from '../../features/projects/slice';
import { useNotifier } from '../hooks';

export function useProject({ projectID }: { projectID?: string }) {
  const notifier = useNotifier();
  const current = useAppSelector((state) => state.projects.details.item);
  const status = useAppSelector((state) => state.projects.details.status);
  const dispatch = useAppDispatch();

  // https://stackoverflow.com/questions/53332321/react-hook-warnings-for-async-function-in-useeffect-useeffect-function-must-ret
  const fetchProject = async () => {
    if (!projectID) {
      return;
    }
    if (status === 'loading') {
      return;
    }
    try {
      await dispatch(fetch(projectID)).unwrap();
    } catch (err) {
      const error = err as SerializedError;
      if (error.message && error.code !== '404') {
        notifier.error(error.message);
      }
    }
  };

  useEffect(() => {
    fetchProject();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [projectID]);

  return { current, loading: status === 'loading' };
}
