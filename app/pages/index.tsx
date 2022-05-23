import Head from 'next/head';

import { Box, Container } from '@mui/material';

import { ProjectInviteStatus } from '../api/enums';
import { NextPageWithLayout } from '../components/common';
import { DashboardLayout } from '../components/dashboard-layout';
import { FeatureToggles } from '../components/feature-toggle/list';
import { useUserInvites } from '../components/hooks';
import SuspenseLoader from '../components/suspense-loader';
import { Welcome } from '../components/welcome/welcome';
import { useAppSelector } from '../data/hooks';
import styles from '../styles/Home.module.css';
import { APP_TITLE } from '../utils/constants';

import type { ReactElement } from 'react';
const SignedIn = () => {
  const me = useAppSelector((state) => state.users.me);
  const { invites, loading: invitesLoading, refetch: refetchInvites } = useUserInvites();
  const projects = useAppSelector((state) => state.projects.all.items);

  const unverified = me?.addresses?.filter((a) => !a.verified) || [];
  const pendingInvites = invites.filter((el) => el.status === ProjectInviteStatus.PENDING);
  if (!me?.addresses?.length) {
    // This is impossible
    throw new Error('No email address');
  }

  const showWelcome = !!unverified?.length || !!pendingInvites.length || !projects.length;

  if (invitesLoading) {
    return <SuspenseLoader></SuspenseLoader>;
  }

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
        {showWelcome ? (
          <Welcome
            addresses={unverified}
            pendingInvites={pendingInvites}
            showNewProject={!projects.length}
            refetchInvites={refetchInvites}
          />
        ) : (
          <FeatureToggles />
        )}
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
