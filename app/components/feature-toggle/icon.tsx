import Percent from '@mui/icons-material/Percent';
import ToggleOn from '@mui/icons-material/ToggleOn';

import { FeatureToggleType } from '../../api/enums';

export type Props = {
  toggleType: FeatureToggleType;
};

export const FeatureToggleIcon = (props: Props) => {
  switch (props.toggleType) {
    case FeatureToggleType.ON_OFF:
      return <ToggleOn />;
    case FeatureToggleType.PERCENTAGE:
      return <Percent />;
  }
};
