import { FeatureToggleType, KeyType } from '../api/enums';

export const featureToggleTypeName = (v: FeatureToggleType) => {
  switch (v) {
    case FeatureToggleType.ON_OFF:
      return 'On/Off';
    case FeatureToggleType.PERCENTAGE:
      return 'Percentage';
    default:
      throw new Error('Unknown feature toggle type');
  }
};

export const keyTypeName = (v: KeyType) => {
  switch (v) {
    case KeyType.BOOLEAN:
      return 'Bool';
    case KeyType.FLOAT:
      return 'Float';
    case KeyType.INTEGER:
      return 'Integer';
    case KeyType.STRING:
      return 'String';
  }
};
