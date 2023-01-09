import { NextRouter } from 'next/router';

import { SerializedError } from '@reduxjs/toolkit';

import { FeatureToggle } from '../../api';
import { FeatureToggleType, StickinessType } from '../../api/enums';
import { LOGIN } from '../../utils/constants';
import { Notif } from '../../utils/notif';

export type MaybeEnvironmentID = {
  environmentId?: string;
};

class ValidationError {
  message: string;
  code = '400';
  constructor(message: string) {
    this.message = message;
  }
}

export function handleError(router: NextRouter, notifier: Notif, error: SerializedError) {
  const code = Number(error.code);
  if (code === 401) {
    // redirect to login
    router.push(LOGIN);
    return;
  }
  if (error.message && code < 500 && code != 404) {
    notifier.error(error.message);
  }
}

export function validate(featureToggle: FeatureToggle) {
  if (!featureToggle.name?.length) {
    throw new ValidationError('Name was not specified');
  }
  switch (featureToggle.toggleType) {
    case FeatureToggleType.PERCENTAGE:
      const perc = featureToggle.percentage;
      if (!perc || perc.on?.weight === undefined) {
        throw new ValidationError(`Percentage wasn't set`);
      }
      if (perc.off?.weight !== 100 - perc.on?.weight) {
        throw new ValidationError(`off weight isn't set correctly`);
      }
      if (perc.stickiness?.stickinessType === StickinessType.KEYS) {
        if (!perc.stickiness?.keys?.length) {
          throw new ValidationError(`No attributes specified for a sticky feature flag.`);
        }
        if (perc.stickiness?.keys?.some((k) => !k.key?.length)) {
          throw new ValidationError(`Attribute name cannot be empty`);
        }
      }
      break;
    case FeatureToggleType.ON_OFF:
      const onOff = featureToggle.onOff;
      if (!onOff || (onOff.on?.weight !== 100 && onOff.on?.weight !== 0)) {
        throw new ValidationError(`On/off wasn't set correctly`);
      }
      if (onOff.off?.weight !== 100 - onOff.on?.weight) {
        throw new ValidationError(`off weight isn't set correctly`);
      }
      break;
  }
}
