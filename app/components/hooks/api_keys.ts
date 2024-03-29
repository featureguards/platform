import axios, { AxiosError, CancelTokenSource } from 'axios';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import { ApiKey } from '../../api';
import { Dashboard } from '../../data/api';
import { SerializeError } from '../../features/utils';
import { useNotifier } from '../hooks';
import { handleError, MaybeEnvironmentID } from './utils';

export function useApiKeysList(props: MaybeEnvironmentID) {
  const notifier = useNotifier();
  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter();
  const [apiKeys, setApiKeys] = useState<ApiKey[] | undefined>();
  const refetch = async (cancelled: boolean, source?: CancelTokenSource) => {
    if (loading) {
      return;
    }
    if (!props.environmentId) {
      return;
    }
    try {
      setLoading(true);
      const res = await Dashboard.listApiKeys(props.environmentId, { cancelToken: source?.token });
      setApiKeys(res.data.apiKeys);
    } catch (err) {
      if (axios.isCancel(err)) {
        return;
      }
      handleError(router, notifier, SerializeError(err as AxiosError));
    } finally {
      if (!cancelled) {
        setLoading(false);
      }
    }
  };

  useEffect(() => {
    const source = axios.CancelToken.source();
    let cancelled = false;

    refetch(cancelled, source);
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
    return () => {
      cancelled = true;
      source.cancel();
    };
  }, [props.environmentId]);

  return { apiKeys, loading: status === 'loading', refetch };
}
