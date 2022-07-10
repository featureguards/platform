import Head from 'next/head';

import { Box, Container, Grid, Typography } from '@mui/material';

import { NextPageWithLayout } from '../components/common';
import { DashboardLayout } from '../components/dashboard-layout';
import { Settings } from '../components/ory/settings';

const Account: NextPageWithLayout = () => {
  return (
    <>
      <Head>
        <title>Account</title>
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
            Account
          </Typography>
          <Grid container spacing={3}>
            <Grid item lg={8} md={6} xs={12}>
              <Settings />
            </Grid>
          </Grid>
        </Container>
      </Box>
    </>
  );
};

Account.getLayout = (page) => <DashboardLayout>{page}</DashboardLayout>;

export default Account;
