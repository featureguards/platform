import PropTypes from 'prop-types';
import { useState } from 'react';

import LogoutIcon from '@mui/icons-material/Logout';
import MenuIcon from '@mui/icons-material/Menu';
import SearchIcon from '@mui/icons-material/Search';
import {
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

import { Selector as SelectorIcon } from '../icons/selector';
import { useLogoutHandler } from '../ory/hooks';
import { useProject, useProjects } from './hooks';
import SuspenseLoader from './suspense-loader';

const DashboardNavbarRoot = styled(AppBar)(({ theme }: { theme: Theme }) => ({
  backgroundColor: theme.palette.background.paper,
  boxShadow: theme.shadows[3]
}));

type DashboardNavbarProps = {
  onSidebarOpen?: () => void;
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
  const { onSidebarOpen, ...other } = props;
  const onLogout = useLogoutHandler();
  const { projects, loading: projectsLoading } = useProjects({});
  const [currentIndex, setCurrentIndex] = useState<number>(0);
  const { loading: currentLoading } = useProject({
    projectID: projects?.[currentIndex]?.id
  });
  if (projectsLoading || currentLoading) {
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
          <Tooltip title="Search">
            <IconButton sx={{ ml: 1 }}>
              <SearchIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          <Box sx={{ flexGrow: 1 }} />
          <Tooltip title="Project">
            {projects?.length > 1 ? (
              <ProjectSelector
                sx={{
                  background: 'neutral.500',
                  minWidth: 80
                }}
                value={currentIndex}
                onChange={(e) => {
                  setCurrentIndex(Number(e.target.value || 0));
                }}
                IconComponent={() => (
                  <SelectorIcon
                    sx={{
                      color: 'neutral.500',
                      width: 14,
                      height: 14
                    }}
                  />
                )}
              >
                {projects.map((p, index) => (
                  <MenuItem key={p.id} value={index}>
                    <Typography color="neutral.500" variant="subtitle1">
                      {p.name}
                    </Typography>
                  </MenuItem>
                ))}
              </ProjectSelector>
            ) : (
              <Typography color="neutral.500" variant="subtitle1">
                {projects?.[0]?.name}
              </Typography>
            )}
          </Tooltip>
          <Tooltip title="Logout">
            <IconButton sx={{ ml: 1 }} onClick={onLogout}>
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
