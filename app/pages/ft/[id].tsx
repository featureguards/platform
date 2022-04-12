import Head from 'next/head';
import { useRouter } from 'next/router';

import { Box } from '@mui/material';

import { NextPageWithLayout } from '../../components/common';
import { DashboardLayout } from '../../components/dashboard-layout';
import { FeatureToggleView } from '../../components/feature-toggle/view';
import { useAppSelector } from '../../data/hooks';
import { APP_TITLE } from '../../utils/constants';

import type { ReactElement } from 'react';
const Content = () => {
  const router = useRouter();
  const id = router.query.id;
  const environment = useAppSelector((state) => state.projects.environment.item);

  if (!id || !environment || !environment.id) {
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
      <FeatureToggleView id={id as string} environmentId={environment.id} />
    </Box>
  );
};

const Page: NextPageWithLayout = () => {
  return (
    <>
      <Head>
        <title>{APP_TITLE}</title>
      </Head>
      <Content />
    </>
  );
};

Page.getLayout = (page: ReactElement) => {
  return <DashboardLayout>{page}</DashboardLayout>;
};

export default Page;
