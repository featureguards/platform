import { ProviderContext, SnackbarKey } from 'notistack';

export class Notif {
  provider: ProviderContext;
  constructor(provider: ProviderContext) {
    this.provider = provider;
  }

  error(msg: string): SnackbarKey {
    return this.provider.enqueueSnackbar(msg, {
      variant: 'error',
      autoHideDuration: 2000,
      preventDuplicate: true
    });
  }

  info(msg: string): SnackbarKey {
    return this.provider.enqueueSnackbar(msg, {
      variant: 'default',
      autoHideDuration: 1000,
      preventDuplicate: true
    });
  }

  dismiss(key: SnackbarKey) {
    this.provider.closeSnackbar(key);
  }
}
