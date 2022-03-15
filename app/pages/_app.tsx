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

function MyApp({ Component, pageProps }: AppPropsWithLayout) {
  const getLayout = Component.getLayout ?? ((page) => page);
  return (
    <>
      <Head>
        <title>{APP_TITLE}</title>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <LocalizationProvider dateAdapter={AdapterDateFns}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <SnackbarProvider maxSnack={3}>{getLayout(<Component {...pageProps} />)}</SnackbarProvider>
        </ThemeProvider>
      </LocalizationProvider>
    </>
  );
}

export default MyApp;
