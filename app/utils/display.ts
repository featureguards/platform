import {
  DateTimeOperator,
  DynamicSettingType,
  FeatureToggleType,
  FloatOperator,
  IntOperator,
  KeyType,
  PlatformTypeType,
  ProjectMemberRole,
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

export const projectMemberRoleTypeName = (v: ProjectMemberRole) => {
  switch (v) {
    case ProjectMemberRole.ADMIN:
      return 'Admin';
    case ProjectMemberRole.MEMBER:
      return 'Member';
    case ProjectMemberRole.UNKNOWN:
      return 'Unknown';
  }
};

export const keyTypeName = (v: KeyType) => {
  switch (v) {
    case KeyType.BOOLEAN:
      return 'Bool';
    case KeyType.FLOAT:
      return 'Float';
    case KeyType.INT:
      return 'Integer';
    case KeyType.DATE_TIME:
      return 'Date/Time';
    case KeyType.STRING:
      return 'String';
  }
};

export const stickinessTypeName = (v: StickinessType) => {
  switch (v) {
    case StickinessType.KEYS:
      return 'Attributes';
    case StickinessType.RANDOM:
      return 'Random';
  }
};

export const platformTypeName = (v: PlatformTypeType) => {
  switch (v) {
    case PlatformTypeType.DEFAULT:
      return 'Server';
    case PlatformTypeType.ANDROID:
      return 'Android';
    case PlatformTypeType.IOS:
      return 'iOS';
    case PlatformTypeType.WEB:
      return 'Browser';
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

export const intOperatorName = (v: IntOperator) => {
  switch (v) {
    case IntOperator.EQ:
      return '=';
    case IntOperator.GT:
      return '>';
    case IntOperator.GTE:
      return '>=';
    case IntOperator.LT:
      return '<';
    case IntOperator.LTE:
      return '<=';
    case IntOperator.NEQ:
      return '!=';
    case IntOperator.IN:
      return 'In';
  }
};

export const dynamicSettingTypeName = (v: DynamicSettingType) => {
  switch (v) {
    case DynamicSettingType.BOOL:
      return 'Bool';
    case DynamicSettingType.FLOAT:
      return 'Float';
    case DynamicSettingType.INTEGER:
      return 'Integer';
    case DynamicSettingType.JSON:
      return 'JSON';
    case DynamicSettingType.LIST:
      return 'List';
    case DynamicSettingType.MAP:
      return 'Map';
    case DynamicSettingType.SET:
      return 'Set';
    case DynamicSettingType.STRING:
      return 'String';
  }
};
