import styles from '../styles/Home.module.css';
import { edgeConfig } from '@ory/integrations/next';
import { Configuration, Session, V0alpha2Api } from '@ory/kratos-client';
import { AxiosError, AxiosResponse } from 'axios';
import { GlobalApi, DashboardApi } from '../api/index';
import Head from 'next/head';
import { useEffect, useState } from 'react';
import { DashboardLayout } from '../components/dashboard-layout';
import type { ReactElement } from 'react';
import { NextPageWithLayout } from '../components/common';
import { APP_TITLE } from '../utils/constants';
import { Box, Container, Grid } from '@mui/material';

edgeConfig.basePath = '/identity';
const kratos = new V0alpha2Api(new Configuration(edgeConfig));

const SignedOut = () => (
  <>
    Get started and <a href={'/identity/self-service/registration/browser'}>create an example account</a> or{' '}
    <a href={'/identity/self-service/login/browser'}>sign in</a>,{' '}
    <a href={'/identity/self-service/recovery/browser'}>recover your account</a> or{' '}
    <a href={'/identity/self-service/verification/browser'}>verify your email address</a>! All using open source{' '}
    <a href={'https://github.com/ory/kratos'}>Ory Kratos</a> in minutes with just a{' '}
    <a href={'https://www.ory.sh/login-spa-react-nextjs-authentication-example-api/'}>few lines of code</a>!
  </>
);

const SignedIn = () => {
  return (
    <Box
      component="main"
      sx={{
        flexGrow: 1,
        py: 8
      }}
    >
      <Container maxWidth={false}>
        <Grid container spacing={3}>
          <Grid item lg={3} sm={6} xl={3} xs={12}>
            <h1>Hello 1</h1>
          </Grid>
          <Grid item xl={3} lg={3} sm={6} xs={12}>
            {/* <TotalCustomers /> */}
          </Grid>
          <Grid item xl={3} lg={3} sm={6} xs={12}>
            <h2>Hello 2</h2>
            {/* <TasksProgress /> */}
          </Grid>
          <Grid item xl={3} lg={3} sm={6} xs={12}>
            {/* <TotalProfit sx={{ height: '100%' }} /> */}
          </Grid>
          <Grid item lg={8} md={12} xl={9} xs={12}>
            {/* <Sales /> */}
          </Grid>
          <Grid item lg={4} md={6} xl={3} xs={12}>
            {/* <TrafficByDevice sx={{ height: '100%' }} /> */}
          </Grid>
          <Grid item lg={4} md={6} xl={3} xs={12}>
            {/* <LatestProducts sx={{ height: '100%' }} /> */}
          </Grid>
          <Grid item lg={8} md={12} xl={9} xs={12}>
            {/* <LatestOrders /> */}
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

const Home: NextPageWithLayout = () => {
  // Contains the current session or undefined.
  const [session, setSession] = useState<Session>();

  // The URL we can use to log out.
  const [logoutUrl, setLogoutUrl] = useState<string>();

  // The error message or undefined.
  const [error, setError] = useState<any>();

  async function me(): Promise<AxiosResponse> {
    const api = new DashboardApi(undefined, '');
    const res = await api.dashboardMe();
    console.log(`Me: ${res.data}`);
    return res;
  }

  async function greetGlobal(): Promise<AxiosResponse> {
    const api = new GlobalApi(undefined, '');
    const res = await api.globalSayHello('bar');
    console.log(`Global: ${res.data.message}`);
    return res;
  }

  me().then((res: AxiosResponse) => {
    console.log(res);
  });
  greetGlobal().then((res: AxiosResponse) => {
    console.log(res);
  });

  async function fetchSession() {
    // If the session or error have been loaded, do nothing.
    if (session || error) {
      return;
    }

    // Try to load the session.
    try {
      const res = await kratos.toSession();
      setSession(res.data);

      const logoutUrl = await kratos.createSelfServiceLogoutFlowUrlForBrowsers();
      setLogoutUrl(logoutUrl.data.logout_url);

      await me();
    } catch (err) {
      const axiosErr = err as AxiosError;
      setError({
        error: axiosErr.toString(),
        data: axiosErr.response?.data
      });
    }
  }

  useEffect(() => {
    fetchSession();
  }, [session, error]);

  return (
    <>
      <Head>
        <title>{APP_TITLE}</title>
      </Head>
      <main className={styles.main}>
        <h1 className={styles.title}>
          {session ? (
            <>
              You are signed in using <a href="https://www.ory.sh/">Ory</a>!
            </>
          ) : (
            <>
              Add Auth to <a href={'https://nextjs.org'}>Next.js</a> with <a href="https://www.ory.sh/">Ory</a>!
            </>
          )}
        </h1>

        <p className={styles.description}>{session ? <SignedIn /> : <SignedOut />}</p>

        {session ? (
          <div className={styles.session}>
            <>
              <p>Find your session details below. </p>
              <pre className={styles.pre + ' ' + styles.code}>
                <code data-testid={'session-content'}>{JSON.stringify(session, null, 2)}</code>
              </pre>
            </>
          </div>
        ) : null}
      </main>
    </>
  );
};

Home.getLayout = (page: ReactElement) => {
  return <DashboardLayout>{page}</DashboardLayout>;
};

export default Home;
