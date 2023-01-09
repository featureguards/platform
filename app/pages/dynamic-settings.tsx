import Head from 'next/head';

import { Box, Container } from '@mui/material';

import { NextPageWithLayout } from '../components/common';
import { DashboardLayout } from '../components/dashboard-layout';
import { DynamicSettings } from '../components/dynamic-setting/list';
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
        <DynamicSettings />
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
      </main>
    </>
  );
};

Home.getLayout = (page: ReactElement) => {
  return <DashboardLayout>{page}</DashboardLayout>;
};

export default Home;
