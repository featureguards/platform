import Head from 'next/head';
import { useRouter } from 'next/router';

import { Box } from '@mui/material';

import { NextPageWithLayout } from '../../components/common';
import { DashboardLayout } from '../../components/dashboard-layout';
import { NewFeatureToggle } from '../../components/feature-toggle/create';
import { APP_TITLE } from '../../utils/constants';

import type { ReactElement } from 'react';
const Content = () => {
  const router = useRouter();
  return (
    <Box
      component="main"
      sx={{
        flexGrow: 1,
        py: 8
      }}
    >
      <NewFeatureToggle onCancel={() => router.back()}></NewFeatureToggle>
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
