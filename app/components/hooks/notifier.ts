import { useSnackbar } from 'notistack';
import { useMemo } from 'react';

import { Notif } from '../../utils/notif';

export function useNotifier() {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const notifier = useMemo(
    () => new Notif({ enqueueSnackbar: enqueueSnackbar, closeSnackbar: closeSnackbar }),
    [enqueueSnackbar, closeSnackbar]
  );
  return notifier;
}
