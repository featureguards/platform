import Head from 'next/head';
import { ReactElement } from 'react';

import { Box, Container, Typography } from '@mui/material';

import { DashboardLayout } from '../../components/dashboard-layout';
import { ProjectSettings } from '../../components/project/settings';
import { useAppSelector } from '../../data/hooks';

const Settings = () => {
  const projectDetails = useAppSelector((state) => state.projects.details);
  const currentProject = projectDetails?.item;
  if (!currentProject?.id) return <></>;
  return (
    <>
      <Head>
        <title>Project Settings</title>
      </Head>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          py: 8
        }}
      >
        <Container maxWidth="lg">
          <Typography sx={{ mb: 3 }} variant="h4">
            Settings
          </Typography>
          <ProjectSettings projectID={currentProject?.id} />
        </Container>
      </Box>
    </>
  );
};

Settings.getLayout = (page: ReactElement) => <DashboardLayout>{page}</DashboardLayout>;

export default Settings;
