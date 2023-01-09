import { DateTime } from 'luxon';
import { ReactNode } from 'react';

import { Box, Divider, Typography } from '@mui/material';

import { DynamicSetting } from '../../api';
import { DynamicSettingType } from '../../api/enums';

export type DiffProps = {
  older: DynamicSetting;
  newer: DynamicSetting;
};

export const Diff = ({ older, newer }: DiffProps) => {
  // assertions
  if (
    older.name !== newer.name ||
    older.settingType !== newer.settingType ||
    older.id !== newer.id
  ) {
    // these are unchangeable.
    throw new Error(`Impossible change in dynamic settings.`);
  }
  const diffs: ReactNode[] = [];
  switch (older.settingType) {
    case DynamicSettingType.BOOL:
      const oldBool = !!older?.boolValue?.value;
      const newBool = !!newer?.boolValue?.value;
      if (oldBool !== newBool) {
        diffs.push(
          <>
            <Typography>Value:</Typography>
            <Typography color="green">{String(newBool)}</Typography>
            <Typography color="red">{String(oldBool)}</Typography>
          </>
        );
      }
    case DynamicSettingType.INTEGER:
      const oldInt = older?.integerValue?.value;
      const newInt = newer?.integerValue?.value;
      if (oldInt !== newInt) {
        diffs.push(
          <>
            <Typography>Value:</Typography>
            <Typography color="green">{newInt}</Typography>
            <Typography color="red">{oldInt}</Typography>
          </>
        );
      }
    case DynamicSettingType.FLOAT:
      const oldFloat = older?.floatValue?.value;
      const newFloat = newer?.floatValue?.value;
      if (oldFloat !== newFloat) {
        diffs.push(
          <>
            <Typography>Value:</Typography>
            <Typography color="green">{newFloat}</Typography>
            <Typography color="red">{oldFloat}</Typography>
          </>
        );
      }
    case DynamicSettingType.STRING:
      const oldString = older?.stringValue?.value;
      const newString = newer?.stringValue?.value;
      if (oldString !== newString) {
        diffs.push(
          <>
            <Typography>Value:</Typography>
            <Typography color="green">{newString}</Typography>
            <Typography color="red">{oldString}</Typography>
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
