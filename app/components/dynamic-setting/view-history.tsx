import { ReactNode } from 'react';

import { Box } from '@mui/material';

import { EnvironmentSettingID } from '../../features/dynamic_settings/slice';
import { useDynamicSettingHistory } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { Diff } from './diff';

export type EnvDynamicSettingHistoryViewProps = EnvironmentSettingID;
export const EnvDynamicSettingHistoryView = (props: EnvDynamicSettingHistoryViewProps) => {
  const { dynamicSettings, loading } = useDynamicSettingHistory(props);

  if (loading) {
    return <SuspenseLoader />;
  }

  if (dynamicSettings?.length < 2) {
    return <></>;
  }

  const history: ReactNode[] = [];
  // history is already ordered by created_at desc
  for (let i = 0; i + 1 < dynamicSettings.length; i++) {
    history.push(
      <Box key={dynamicSettings[i].updatedAt}>
        <Diff older={dynamicSettings[i + 1]} newer={dynamicSettings[i]} />
      </Box>
    );
  }

  return <>{history}</>;
};
