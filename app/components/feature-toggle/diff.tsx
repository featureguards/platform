import { DateTime } from 'luxon';
import { ReactNode } from 'react';

import { Box, Divider, Typography } from '@mui/material';

import { FeatureToggle } from '../../api';
import { FeatureToggleType } from '../../api/enums';

export type DiffProps = {
  older: FeatureToggle;
  newer: FeatureToggle;
};

export const Diff = ({ older, newer }: DiffProps) => {
  // assertions
  if (older.name !== newer.name || older.toggleType !== newer.toggleType || older.id !== newer.id) {
    // these are unchangeable.
    throw new Error(`Impossible change in feature flag.`);
  }
  const diffs: ReactNode[] = [];
  switch (older.toggleType) {
    case FeatureToggleType.PERCENTAGE:
      const olderPercDef = older?.percentage;
      const olderPercentage = olderPercDef?.on?.weight || 0;
      const newerPercDef = newer?.percentage;
      const newerPercentage = newerPercDef?.on?.weight || 0;
      if (olderPercentage !== newerPercentage) {
        diffs.push(
          <>
            <Typography>Percentage:</Typography>
            <Typography color="green">{newerPercentage}</Typography>
            <Typography color="red">{olderPercentage}</Typography>
          </>
        );
      }
    case FeatureToggleType.ON_OFF:
      const renderOnOff = (weight: number | undefined) => (weight ? 'On' : 'Off');
      if (older.onOff?.on?.weight !== newer.onOff?.on?.weight) {
        diffs.push(
          <>
            <Typography color="green">{renderOnOff(newer.onOff?.on?.weight)}</Typography>
            <Typography color="red">{renderOnOff(older.onOff?.on?.weight)}</Typography>
          </>
        );
      }
  }
  if (!diffs.length) {
    return <></>;
  }
  return (
    <>
      <Typography>
        {DateTime.fromISO(newer.updatedAt!).toLocaleString(DateTime.DATETIME_FULL_WITH_SECONDS)}
      </Typography>
      <Box
        sx={{
          alignItems: 'center',
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'space-between'
        }}
      >
        {diffs.map((diff, i) => (
          <div key={i}>{diff}</div>
        ))}
      </Box>
      <Divider sx={{ my: 2 }}></Divider>
    </>
  );
};
