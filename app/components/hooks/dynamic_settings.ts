import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { SerializedError } from '@reduxjs/toolkit';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import {
  details,
  EnvironmentID,
  EnvironmentSettingID,
  history,
  list,
  SettingIDEnvironments
} from '../../features/dynamic_settings/slice';
import { useNotifier } from '../hooks';
import { handleError, MaybeEnvironmentID } from './utils';

export function useDynamicSettingsList(props: MaybeEnvironmentID) {
  const dynamicSettings = useAppSelector((state) => state.dynamicSettings.environment.items);
  const { refetch, loading } = useDynamicSettingsListLazy(props);

  useEffect(() => {
    refetch();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [props.environmentId]);

  return { dynamicSettings, loading, refetch };
}

export function useDynamicSettingsListLazy(props: MaybeEnvironmentID) {
  const notifier = useNotifier();
  const status = useAppSelector((state) => state.dynamicSettings.environment.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const refetch = async () => {
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

  return { loading: status === 'loading', refetch };
}

export function useDynamicSettingHistory(props: EnvironmentSettingID) {
  const notifier = useNotifier();
  const items = useAppSelector((state) => state.dynamicSettings.history.items);
  const status = useAppSelector((state) => state.dynamicSettings.history.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const refetch = async () => {
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
    refetch();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [props.environmentId, props.id]);

  return { dynamicSettings: items, loading: status === 'loading', refetch };
}

export function useDynamicSettingDetails(props: SettingIDEnvironments) {
  const notifier = useNotifier();
  const items = useAppSelector((state) => state.dynamicSettings.details.items);
  const status = useAppSelector((state) => state.dynamicSettings.details.status);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const refetch = async () => {
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
    refetch();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [props.id, ...props.environmentIds]);

  return { items, loading: status === 'loading', refetch };
}
