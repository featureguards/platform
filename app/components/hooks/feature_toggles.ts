import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import {
  details,
  EnvironmentFeatureID,
  EnvironmentID,
  FeatureIDEnvironments,
  history,
  list
} from '../../features/feature_toggles/slice';
import { useNotifier } from '../hooks';
import { handleError } from './utils';

export type MaybeEnvironmentID = {
  environmentId?: string;
};

export function useFeatureTogglesList(props: MaybeEnvironmentID) {
  const notifier = useNotifier();
  const featureToggles = useAppSelector((state) => state.featureToggles.environment.items);
  //   const storedEnvID = useAppSelector((state) => state.featureToggles.environment.id);
  const status = useAppSelector((state) => state.featureToggles.environment.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const fetch = async () => {
    if (status === 'loading') {
      return;
    }
    if (!props.environmentId) {
      return;
    }
    try {
      await dispatch(list(props as EnvironmentID)).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

  useEffect(() => {
    fetch();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [props.environmentId]);

  return { featureToggles, loading: status === 'loading' };
}

export function useFeatureToggleHistory(props: EnvironmentFeatureID) {
  const notifier = useNotifier();
  const items = useAppSelector((state) => state.featureToggles.history.items);
  const status = useAppSelector((state) => state.featureToggles.history.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const fetch = async () => {
    if (status === 'loading') {
      return;
    }
    if (!props.environmentId || !props.id) {
      return;
    }
    try {
      await dispatch(history(props)).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

  useEffect(() => {
    fetch();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [props.environmentId, props.id]);

  return { featureToggles: items, loading: status === 'loading' };
}

export function useFeatureToggleDetails(props: FeatureIDEnvironments) {
  const notifier = useNotifier();
  const items = useAppSelector((state) => state.featureToggles.details.items);
  const status = useAppSelector((state) => state.featureToggles.details.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const fetch = async () => {
    if (status === 'loading') {
      return;
    }

    if (!props.id) {
      return;
    }
    try {
      await dispatch(details(props)).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

  useEffect(() => {
    fetch();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [props.id, ...props.environmentIds]);

  return { items, loading: status === 'loading' };
}
