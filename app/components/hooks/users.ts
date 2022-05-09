import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { fetchMe, selectMe } from '../../features/users/slice';
import { useNotifier } from '../hooks';
import { handleError } from './utils';

export function useMe() {
  const notifier = useNotifier();
  const me = useAppSelector(selectMe);
  const status = useAppSelector((state) => state.users.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const refetch = async () => {
    if (status === 'loading') {
      return;
    }
    try {
      await dispatch(fetchMe()).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

  useEffect(() => {
    refetch();
    // This isn't a bug. We only depend on projectID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, []);

  return { me, status, refetch };
}
