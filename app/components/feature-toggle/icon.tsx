import Percent from '@mui/icons-material/Percent';
import ToggleOff from '@mui/icons-material/ToggleOff';
import ToggleOn from '@mui/icons-material/ToggleOn';
import { Typography } from '@mui/material';

import { FeatureToggle } from '../../api';
import { FeatureToggleType } from '../../api/enums';

export type LiveToggleIconProps = {
  featureToggle: FeatureToggle;
};

export const LiveToggleIcon = (props: LiveToggleIconProps) => {
  switch (props.featureToggle.toggleType) {
    case FeatureToggleType.ON_OFF:
      if (props.featureToggle.onOff?.on?.weight == 100) {
        return <ToggleOn color="secondary" />;
      }
      return <ToggleOff />;

    case FeatureToggleType.PERCENTAGE:
      if (!!props.featureToggle.percentage?.on?.weight) {
        return (
          <>
            <Percent color="secondary" />
            <Typography color="secondary">{props.featureToggle.percentage?.on?.weight}</Typography>
          </>
        );
      }
      return (
        <>
          <Percent />
          <Typography>{props.featureToggle.percentage?.on?.weight}</Typography>
        </>
      );
  }
  return <></>;
};

export type FeatureToggleIconProps = {
  toggleType: FeatureToggleType;
};

export const FeatureToggleIcon = (props: FeatureToggleIconProps) => {
  switch (props.toggleType) {
    case FeatureToggleType.ON_OFF:
      return <ToggleOn />;
    case FeatureToggleType.PERCENTAGE:
      return <Percent />;
  }
  return <></>;
};
