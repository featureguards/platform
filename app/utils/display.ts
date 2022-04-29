import {
  DateTimeOperator,
  FeatureToggleType,
  FloatOperator,
  KeyType,
  PlatformTypeType,
  StickinessType,
  StringOperator
} from '../api/enums';

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
    case KeyType.DATE_TIME:
      return 'Date/Time';
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

export const stringOperatorName = (v: StringOperator) => {
  switch (v) {
    case StringOperator.EQ:
      return 'Equals';
    case StringOperator.CONTAINS:
      return 'Contains';
    case StringOperator.IN:
      return 'In';
  }
};

export const dateTimeOperatorName = (v: DateTimeOperator) => {
  switch (v) {
    case DateTimeOperator.BEFORE:
      return 'Before';
    case DateTimeOperator.AFTER:
      return 'After';
  }
};

export const floatOperatorName = (v: FloatOperator) => {
  switch (v) {
    case FloatOperator.EQ:
      return '=';
    case FloatOperator.GT:
      return '>';
    case FloatOperator.GTE:
      return '>=';
    case FloatOperator.LT:
      return '<';
    case FloatOperator.LTE:
      return '<=';
    case FloatOperator.NEQ:
      return '!=';
    case FloatOperator.IN:
      return 'In';
  }
};
