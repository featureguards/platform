import { matches } from 'lodash';
import { DateTime } from 'luxon';
import { ReactNode } from 'react';

import { Divider, Typography } from '@mui/material';

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
    throw new Error(`Impossible change in feature toggle.`);
  }
  const diffs: ReactNode[] = [];
  if (older.description !== newer.description) {
    diffs.push(
      <>
        <Typography>Description:</Typography>
        <Typography color="green">{newer.description}</Typography>
        <Typography color="red">{older.description}</Typography>
      </>
    );
  }
  if (older.enabled !== newer.enabled) {
    diffs.push(
      <>
        <Typography>Enabled:</Typography>
        <Typography color="green">{newer.enabled}</Typography>
        <Typography color="red">{older.enabled}</Typography>
      </>
    );
  }
  if (!matches(older.platforms)(newer.platforms)) {
    diffs.push(
      <>
        <Typography>Platforms:</Typography>
        <Typography color="green">{newer.platforms}</Typography>
        <Typography color="red">{older.platforms}</Typography>
      </>
    );
  }
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
      if (older.onOff?.on?.weight !== newer.onOff?.on?.weight) {
        diffs.push(
          <>
            <Typography>On:</Typography>
            <Typography color="green">{newer.onOff?.on?.weight}</Typography>
            <Typography color="red">{!!older.onOff?.on?.weight}</Typography>
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
      {diffs.map((diff, i) => (
        <div key={i}>{diff}</div>
      ))}
      <Divider></Divider>
    </>
  );
};
