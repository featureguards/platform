import { ReactNode, useEffect, useState } from 'react';

import { Box } from '@mui/material';
import { styled } from '@mui/material/styles';

import { useAppSelector } from '../data/hooks';
import { DashboardNavbar } from './dashboard-navbar';
import { DashboardSidebar } from './dashboard-sidebar';
import { useProjectsLazy } from './hooks';
import SuspenseLoader from './suspense-loader';

const DashboardLayoutRoot = styled('div')(({ theme }) => ({
  display: 'flex',
  flex: '1 1 auto',
  maxWidth: '100%',
  paddingTop: 64,
  [theme.breakpoints.up('lg')]: {
    paddingLeft: 280
  }
}));

export type DashboardLayoutProps = {
  children?: ReactNode;
};

export const DashboardLayout = (props: DashboardLayoutProps) => {
  const { children } = props;
  const [isSidebarOpen, setSidebarOpen] = useState(true);
  const { refetch, loading } = useProjectsLazy();
  const projects = useAppSelector((state) => state.projects.all.items);

  useEffect(() => {
    refetch();
  }, []);

  if (loading) {
    return <SuspenseLoader />;
  }
  return (
    <>
      <DashboardLayoutRoot>
        <Box
          sx={{
            display: 'flex',
            flex: '1 1 auto',
            flexDirection: 'column',
            width: '100%'
          }}
        >
          {children}
        </Box>
      </DashboardLayoutRoot>
      <DashboardNavbar projects={projects} onSidebarOpen={() => setSidebarOpen(true)} />
      <DashboardSidebar onClose={() => setSidebarOpen(false)} open={isSidebarOpen} />
    </>
  );
};
