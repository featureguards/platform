import { NextRouter } from 'next/router';

import { SerializedError } from '@reduxjs/toolkit';

import { LOGIN } from '../../utils/constants';
import { Notif } from '../../utils/notif';

export function handleError(router: NextRouter, notifier: Notif, error: SerializedError) {
  if (error.code === '401') {
    // redirect to login
    router.push(LOGIN);
  }
  if (error.message && error.code !== '404') {
    notifier.error(error.message);
  }
}
