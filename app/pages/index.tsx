import Head from 'next/head';

import { Box, Container, Grid } from '@mui/material';

import { NextPageWithLayout } from '../components/common';
import { DashboardLayout } from '../components/dashboard-layout';
import { useProjects } from '../components/hooks';
import { Welcome } from '../components/welcome/welcome';
import styles from '../styles/Home.module.css';
import { APP_TITLE } from '../utils/constants';

import type { ReactElement } from 'react';
const SignedIn = () => {
  const { projects, loading: projectsLoading } = useProjects();
  const { userInvites, loading: userInvitesLoading } = useProjects();

  if (projectsLoading) {
    return <></>;
  }

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
          <Grid item xs={12}>
            <Welcome />
          </Grid>
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
