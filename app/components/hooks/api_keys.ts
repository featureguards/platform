import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import { ApiKey } from '../../api';
import { Dashboard } from '../../data/api';
import { SerializeError } from '../../features/utils';
import { useNotifier } from '../hooks';
import { MaybeEnvironmentID } from './feature_toggles';
import { handleError } from './utils';

export function useApiKeysList(props: MaybeEnvironmentID) {
  const notifier = useNotifier();
  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter();
  const [apiKeys, setApiKeys] = useState<ApiKey[] | undefined>();
  const refetch = async () => {
    if (loading) {
      return;
    }
    if (!props.environmentId) {
      return;
    }
    try {
      setLoading(true);
      const res = await Dashboard.listApiKeys(props.environmentId);
      setApiKeys(res.data.apiKeys);
    } catch (err) {
      handleError(router, notifier, SerializeError(err as AxiosError));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refetch();
    // This isn't a bug. We only depend on envrionment ID. Do NOT add other dependencies,
    // it will cause endless loads.
  }, [props.environmentId]);

  return { apiKeys, loading: status === 'loading', refetch };
}
