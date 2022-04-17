import { format } from 'date-fns';
import { matches } from 'lodash';
import { ReactNode } from 'react';

import { Box, Divider, Typography } from '@mui/material';
import { useTheme } from '@mui/system';

import { FeatureToggle } from '../../api';
import { FeatureToggleType } from '../../api/enums';
import { EnvironmentFeatureID } from '../../features/feature_toggles/slice';
import { useFeatureToggleHistory } from '../hooks';
import SuspenseLoader from '../suspense-loader';

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

  const renderFTDiff = (ft1: FeatureToggle, ft2: FeatureToggle) => {
    // assertions
    if (ft1.name !== ft2.name || ft1.toggleType !== ft2.toggleType || ft1.id !== ft2.id) {
      // these are unchangeable.
      throw new Error(`Impossible change in feature toggle.`);
    }
    const diffs: ReactNode[] = [];
    if (ft1.description !== ft2.description) {
      diffs.push(
        <>
          <Typography>Description:</Typography>
          <Typography color="green">{ft2.description}</Typography>
          <Typography color="red">{ft1.description}</Typography>
        </>
      );
    }
    if (ft1.enabled !== ft2.enabled) {
      diffs.push(
        <>
          <Typography>Enabled:</Typography>
          <Typography color="green">{ft2.enabled}</Typography>
          <Typography color="red">{ft1.enabled}</Typography>
        </>
      );
    }
    if (!matches(ft1.platforms)(ft2.platforms)) {
      diffs.push(
        <>
          <Typography>Platforms:</Typography>
          <Typography color="green">{ft2.platforms}</Typography>
          <Typography color="red">{ft1.platforms}</Typography>
        </>
      );
    }
    switch (ft1.toggleType) {
      case FeatureToggleType.PERCENTAGE:
        const ft1PercDef = ft1?.percentage;
        const ft1Percentage = ft1PercDef?.on?.weight || 0;
        const ft2PercDef = ft2?.percentage;
        const ft2Percentage = ft2PercDef?.on?.weight || 0;
        if (ft1Percentage !== ft2Percentage) {
          diffs.push(
            <>
              <Typography>Percentage:</Typography>
              <Typography color="green">{ft2Percentage}</Typography>
              <Typography color="red">{ft1Percentage}</Typography>
            </>
          );
        }
      case FeatureToggleType.ON_OFF:
        if (ft1.onOff?.on?.weight !== ft2.onOff?.on?.weight) {
          diffs.push(
            <>
              <Typography>On:</Typography>
              <Typography color="green">{ft2.onOff?.on?.weight}</Typography>
              <Typography color="red">{!!ft1.onOff?.on?.weight}</Typography>
            </>
          );
        }
    }
    return diffs;
  };

  const history: ReactNode[] = [];
  // history is already ordered by created_at desc
  for (let i = 0; i + 1 < featureToggles.length; i++) {
    const diff = renderFTDiff(featureToggles[i + 1], featureToggles[i]);
    if (diff.length) {
      history.push(
        <>
          <Typography>{format(Date.parse(featureToggles[i].updatedAt!), 'ff')}</Typography>
          {diff}
          <Divider></Divider>
        </>
      );
    }
  }

  if (!history.length) {
    return <></>;
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
