import '../styles/globals.css';

import Head from 'next/head';
import { SnackbarProvider } from 'notistack';
import { Provider } from 'react-redux';

import { CssBaseline } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterLuxon } from '@mui/x-date-pickers/AdapterLuxon';

import { AppPropsWithLayout } from '../components/common';
import { store } from '../data/store';
import { RouteGuard } from '../providers/protected';
import { theme } from '../theme';
import { APP_TITLE } from '../utils/constants';

function MyApp({ Component, pageProps, router }: AppPropsWithLayout) {
  const getLayout = Component.getLayout ?? ((page) => page);

  return (
    <>
      <Head>
        <title>{APP_TITLE}</title>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <Provider store={store}>
        <LocalizationProvider dateAdapter={AdapterLuxon}>
          <ThemeProvider theme={theme}>
            <CssBaseline />
            <SnackbarProvider maxSnack={3}>
              <RouteGuard router={router}>{getLayout(<Component {...pageProps} />)}</RouteGuard>
            </SnackbarProvider>
          </ThemeProvider>
        </LocalizationProvider>
      </Provider>
    </>
  );
}

export default MyApp;
