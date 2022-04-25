import { FeatureToggleType, KeyType, PlatformTypeType, StickinessType } from '../api/enums';

export const featureToggleTypeName = (v: FeatureToggleType) => {
  switch (v) {
    case FeatureToggleType.ON_OFF:
      return 'On/Off';
    case FeatureToggleType.PERCENTAGE:
      return 'Percentage';
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

export const stickinessTypeName = (v: StickinessType) => {
  switch (v) {
    case StickinessType.KEYS:
      return 'Keys';
    case StickinessType.RANDOM:
      return 'Random';
  }
};

export const platformTypeName = (v: PlatformTypeType) => {
  switch (v) {
    case PlatformTypeType.DEFAULT:
      return 'Server';
    case PlatformTypeType.MOBILE:
      return 'Mobile';
    case PlatformTypeType.WEB:
      return 'Web';
  }
};
