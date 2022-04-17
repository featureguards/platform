import { NextRouter } from 'next/router';

import { SerializedError } from '@reduxjs/toolkit';

import { LOGIN } from '../../utils/constants';
import { Notif } from '../../utils/notif';

export function handleError(router: NextRouter, notifier: Notif, error: SerializedError) {
  const code = Number(error.code);
  if (code === 401) {
    // redirect to login
    router.push(LOGIN);
    return;
  }
  if (error.message && code < 500) {
    notifier.error(error.message);
  }
}
