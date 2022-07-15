import Head from 'next/head';
import { useRouter } from 'next/router';
import { ReactElement } from 'react';

import { Box, Container } from '@mui/material';

import { DashboardLayout } from '../../components/dashboard-layout';
import { useProjectsLazy } from '../../components/hooks';
import { NewProject } from '../../components/project/new';
import SuspenseLoader from '../../components/suspense-loader';

const New = () => {
  const router = useRouter();
  const { refetch, loading } = useProjectsLazy();
  const handleNewProject = async ({ err }: { err?: Error }) => {
    if (!err) {
      await refetch();
      router.push('/');
    }
  };
  if (loading) {
    return <SuspenseLoader />;
  }
  return (
    <>
      <Head>
        <title>New Project</title>
      </Head>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          py: 8
        }}
      >
        <Container maxWidth="lg">
          <NewProject onSubmit={handleNewProject} />
        </Container>
      </Box>
    </>
  );
};

New.getLayout = (page: ReactElement) => <DashboardLayout>{page}</DashboardLayout>;

export default New;
