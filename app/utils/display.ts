import { FeatureToggleType } from '../api/enums';

export const FeatureToggleTypeName = (v: FeatureToggleType) => {
  switch (v) {
    case FeatureToggleType.EXPERIMENT:
      return 'Experiment';
    case FeatureToggleType.ON_OFF:
      return 'On/Off';
    case FeatureToggleType.PERCENTAGE:
      return 'Percentage';
    case FeatureToggleType.PERMISSION:
      return 'Permission';
    default:
      throw new Error('Unknown feature toggle type');
  }
};
