import Head from 'next/head';

import { Box, Container, Grid } from '@mui/material';

import { NextPageWithLayout } from '../components/common';
import { DashboardLayout } from '../components/dashboard-layout';
import { FeatureToggles } from '../components/feature-toggle/list';
import { Welcome } from '../components/welcome/welcome';
import styles from '../styles/Home.module.css';
import { APP_TITLE } from '../utils/constants';

import type { ReactElement } from 'react';
const SignedIn = () => {
  return (
    <Box
      component="main"
      sx={{
        flexGrow: 1,
        width: '100%',
        py: 8
      }}
    >
      <Container>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Welcome />
          </Grid>
          <Grid item xs={12}>
            <FeatureToggles />
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

const Home: NextPageWithLayout = () => {
  return (
    <>
      <Head>
        <title>{APP_TITLE}</title>
      </Head>
      <main className={styles.main}>
        <SignedIn />
        <div className={styles.session}>
          <>
            <p>Find your session details below. </p>
            <pre className={styles.pre + ' ' + styles.code}>
              {/* <code data-testid={'session-content'}>{JSON.stringify(me, null, 2)}</code> */}
            </pre>
          </>
        </div>
      </main>
    </>
  );
};

Home.getLayout = (page: ReactElement) => {
  return <DashboardLayout>{page}</DashboardLayout>;
};

export default Home;
