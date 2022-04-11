import { useCallback, useEffect } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import {
  details,
  EnvironmentFeatureID,
  EnvironmentID,
  history,
  list
} from '../../features/feature_toggles/slice';
import { useNotifier } from '../hooks';

export type MaybeEnvironmentID = {
  projectID?: string;
  environmentID?: string;
};

export type MaybeEnvironmentFeatureID = {
  projectID?: string;
  environmentID?: string;
  id?: string;
};

export function useFeatureTogglesList(props: MaybeEnvironmentID) {
  const notifier = useNotifier();
  const featureToggles = useAppSelector((state) => state.featureToggles.environment.items);
  //   const storedEnvID = useAppSelector((state) => state.featureToggles.environment.id);
  const status = useAppSelector((state) => state.featureToggles.environment.status);
  const dispatch = useAppDispatch();

  const fetch = async () => {
    if (status === 'loading') {
      return;
    }
    if (!props.environmentID || !props.projectID) {
      return;
    }
    try {
      await dispatch(list(props as EnvironmentID)).unwrap();
    } catch (err) {
      const error = err as SerializedError;
      if (error.message && error.code !== '404') {
        notifier.error(error.message);
      }
    }
  };

  useEffect(() => {
    fetch();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [props.environmentID]);

  return { featureToggles, loading: status === 'loading' };
}

export function useFeatureToggleHistory(props: EnvironmentFeatureID) {
  const notifier = useNotifier();
  const items = useAppSelector((state) => state.featureToggles.history.items);
  const storedID = useAppSelector((state) => state.featureToggles.history.id);
  const status = useAppSelector((state) => state.featureToggles.history.status);
  const error = useAppSelector((state) => state.featureToggles.history.error);
  const dispatch = useAppDispatch();

  const fetch = useCallback(async () => {
    if (
      props.id === storedID &&
      (status === 'succeeded' || status === 'failed' || status === 'loading')
    ) {
      return;
    }
    try {
      await dispatch(history(props)).unwrap();
    } catch (_err) {
      if (error) {
        notifier.error(error);
      }
    }
  }, [dispatch, error, notifier, props, status, storedID]);

  fetch();

  return { items, loading: status === 'loading' };
}

export function useFeatureToggleDetails(props: EnvironmentFeatureID) {
  const notifier = useNotifier();
  const item = useAppSelector((state) => state.featureToggles.details.item);
  const storedID = useAppSelector((state) => state.featureToggles.details.id);
  const status = useAppSelector((state) => state.featureToggles.details.status);
  const error = useAppSelector((state) => state.featureToggles.details.error);
  const dispatch = useAppDispatch();

  const fetch = useCallback(async () => {
    if (
      props.id === storedID &&
      (status === 'succeeded' || status === 'failed' || status === 'loading')
    ) {
      return;
    }
    try {
      await dispatch(details(props)).unwrap();
    } catch (_err) {
      if (error) {
        notifier.error(error);
      }
    }
  }, [dispatch, error, notifier, props, status, storedID]);

  fetch();

  return { item, loading: status === 'loading' };
}
