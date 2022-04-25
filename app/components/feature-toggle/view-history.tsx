import { ReactNode } from 'react';

import { Box } from '@mui/material';
import { useTheme } from '@mui/system';

import { EnvironmentFeatureID } from '../../features/feature_toggles/slice';
import { useFeatureToggleHistory } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { Diff } from './diff';

export type EnvFeatureToggleHistoryViewProps = EnvironmentFeatureID;
export const EnvFeatureToggleHistoryView = (props: EnvFeatureToggleHistoryViewProps) => {
  const theme = useTheme();
  const { featureToggles, loading } = useFeatureToggleHistory(props);

  if (loading) {
    return <SuspenseLoader></SuspenseLoader>;
  }

  if (featureToggles?.length < 2) {
    return <></>;
  }

  const history: ReactNode[] = [];
  // history is already ordered by created_at desc
  for (let i = 0; i + 1 < featureToggles.length; i++) {
    history.push(
      <Diff
        key={featureToggles[i].updatedAt}
        older={featureToggles[i + 1]}
        newer={featureToggles[i]}
      />
    );
  }

  return (
    <Box
      sx={{
        pt: 5,
        backgroundColor: theme.palette.background.paper
      }}
    >
      {history}
    </Box>
  );
};
