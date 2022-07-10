import PropTypes from 'prop-types';
import { useState } from 'react';

import LogoutIcon from '@mui/icons-material/Logout';
import MenuIcon from '@mui/icons-material/Menu';
import {
  Alert,
  AppBar,
  Box,
  IconButton,
  MenuItem,
  Select,
  Theme,
  Toolbar,
  Tooltip,
  Typography
} from '@mui/material';
import { styled } from '@mui/material/styles';

import { Project } from '../api';
import { Selector as SelectorIcon } from '../icons/selector';
import { useLogout } from '../ory/hooks';
import { useProject } from './hooks';
import SuspenseLoader from './suspense-loader';

const DashboardNavbarRoot = styled(AppBar)(({ theme }: { theme: Theme }) => ({
  backgroundColor: theme.palette.background.paper,
  boxShadow: theme.shadows[3]
}));

type DashboardNavbarProps = {
  onSidebarOpen?: () => void;
  projects: Project[];
};

const ProjectSelector = styled(Select)(() => ({
  '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
    border: '0px solid'
  },
  '.MuiOutlinedInput-notchedOutline': {
    border: '0px solid'
  }
}));

export const DashboardNavbar = (props: DashboardNavbarProps) => {
  const { onSidebarOpen, projects, ...other } = props;
  const { logout } = useLogout();
  const [currentIndex, setCurrentIndex] = useState<number>(0);
  const { loading: currentLoading } = useProject({
    projectID: projects?.[currentIndex]?.id
  });
  if (currentLoading) {
    return <SuspenseLoader />;
  }

  return (
    <>
      <DashboardNavbarRoot
        sx={{
          left: {
            lg: 280
          },
          width: {
            lg: 'calc(100% - 280px)'
          }
        }}
        {...other}
      >
        {!!process.env.NEXT_PUBLIC_APP_ENV && (
          <Box sx={{ flexGrow: 1, width: '100%' }}>
            <Alert sx={{ justifyContent: 'center' }} severity="info">
              {process.env.NEXT_PUBLIC_APP_ENV} Environment
            </Alert>
          </Box>
        )}

        <Toolbar
          disableGutters
          sx={{
            minHeight: 64,
            left: 0,
            px: 2
          }}
        >
          <IconButton
            onClick={onSidebarOpen}
            sx={{
              display: {
                xs: 'inline-flex',
                lg: 'none'
              }
            }}
          >
            <MenuIcon fontSize="small" />
          </IconButton>
          {/* <Tooltip title="Search">
            <IconButton sx={{ ml: 1 }}>
              <SearchIcon fontSize="small" />
            </IconButton>
          </Tooltip> */}
          <Box sx={{ flexGrow: 1 }} />
          <Tooltip title="Project">
            {projects.length > 1 ? (
              <ProjectSelector
                sx={{
                  background: 'grey.500',
                  minWidth: 80
                }}
                value={currentIndex}
                onChange={(e) => {
                  setCurrentIndex(Number(e.target.value || 0));
                }}
                IconComponent={() => (
                  <SelectorIcon
                    sx={{
                      color: 'grey.500',
                      width: 14,
                      height: 14
                    }}
                  />
                )}
              >
                {projects.map((p, index) => (
                  <MenuItem key={p.id} value={index}>
                    <Typography color="grey.500" variant="subtitle1">
                      {p.name}
                    </Typography>
                  </MenuItem>
                ))}
              </ProjectSelector>
            ) : (
              <Typography color="grey.500" variant="subtitle1">
                {projects?.[0]?.name}
              </Typography>
            )}
          </Tooltip>
          <Tooltip title="Logout">
            <IconButton sx={{ ml: 1 }} onClick={logout}>
              <LogoutIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Toolbar>
      </DashboardNavbarRoot>
    </>
  );
};

DashboardNavbar.propTypes = {
  onSidebarOpen: PropTypes.func
};
