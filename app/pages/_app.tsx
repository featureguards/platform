import '../styles/globals.css';
import Head from 'next/head';
import { CssBaseline } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';
import LocalizationProvider from '@mui/lab/LocalizationProvider';
import AdapterDateFns from '@mui/lab/AdapterDateFns';
import { theme } from '../theme';
import { APP_TITLE } from '../utils/constants';
import { AppPropsWithLayout } from '../components/common';
import { SnackbarProvider } from 'notistack';
import { Provider } from 'react-redux';
import { store } from '../data/store';
import { RouteGuard } from '../providers/protected';

function MyApp({ Component, pageProps, router }: AppPropsWithLayout) {
  const getLayout = Component.getLayout ?? ((page) => page);

  return (
    <>
      <Head>
        <title>{APP_TITLE}</title>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <Provider store={store}>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
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
