import NextLink from 'next/link';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import {
  Avatar,
  Box,
  Divider,
  Drawer,
  MenuItem,
  Select,
  Theme,
  Typography,
  useMediaQuery
} from '@mui/material';
import { styled } from '@mui/material/styles';

import { useAppSelector } from '../data/hooks';
import { ChartBar as ChartBarIcon } from '../icons/chart-bar';
import { Cog as CogIcon } from '../icons/cog';
import { Lock as LockIcon } from '../icons/lock';
import { Selector as SelectorIcon } from '../icons/selector';
import { User as UserIcon } from '../icons/user';
import { UserAdd as UserAddIcon } from '../icons/user-add';
import { UserCircle as UserCircleIcon } from '../icons/user-circle';
import { XCircle as XCircleIcon } from '../icons/x-circle';
import { useProject, useProjects } from './hooks';
import { Logo } from './logo';
import { NavItem } from './nav-item';
import SuspenseLoader from './suspense-loader';

const items = [
  {
    href: '/',
    icon: <ChartBarIcon fontSize="small" />,
    title: 'Dashboard'
  },
  {
    href: '/account',
    icon: <UserIcon fontSize="small" />,
    title: 'Account'
  },
  {
    href: '/settings',
    icon: <CogIcon fontSize="small" />,
    title: 'Settings'
  },
  {
    href: '/login',
    icon: <LockIcon fontSize="small" />,
    title: 'Login'
  },
  {
    href: '/register',
    icon: <UserAddIcon fontSize="small" />,
    title: 'Register'
  },
  {
    href: '/404',
    icon: <XCircleIcon fontSize="small" />,
    title: 'Error'
  }
];

const EnvironmentSelector = styled(Select)(() => ({
  '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
    border: '0px solid'
  },
  '.MuiOutlinedInput-notchedOutline': {
    border: '0px solid'
  }
}));

type DashboardProps = {
  onClose?: () => void;
  open: boolean;
};

export const DashboardSidebar = (props: DashboardProps) => {
  const { open, onClose } = props;
  const router = useRouter();
  const lgUp = useMediaQuery((theme: Theme) => theme.breakpoints.up('lg'), {
    defaultMatches: true,
    noSsr: false
  });
  useEffect(() => {
    if (!router.isReady) {
      return;
    }

    if (open) {
      onClose?.();
    }
  }, [onClose, open, router.asPath, router.isReady]);

  const { projects, loading: projectsLoading } = useProjects();
  const [currentIndex, setCurrentIndex] = useState<number>(0);
  const { current, loading: currentLoading } = useProject({
    projectID: projects?.[currentIndex]?.id
  });
  const me = useAppSelector((state) => state.users.me);

  if (projectsLoading || currentLoading) {
    return <SuspenseLoader />;
  }

  const content = (
    <>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          height: '100%'
        }}
      >
        <div>
          <Box sx={{ p: 3 }}>
            <NextLink href="/" passHref>
              <a>
                <Logo
                  sx={{
                    height: 42,
                    width: 42
                  }}
                />
              </a>
            </NextLink>
          </Box>
          {!!current?.environments?.length && (
            <Box sx={{ px: 2 }}>
              <Box
                sx={{
                  alignItems: 'center',
                  backgroundColor: 'rgba(255, 255, 255, 0.04)',
                  cursor: 'pointer',
                  display: 'flex',
                  justifyContent: 'space-between',
                  px: 3,
                  py: '11px',
                  borderRadius: 1
                }}
              >
                <EnvironmentSelector
                  fullWidth
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
                  {current?.environments?.map((p, index) => (
                    <MenuItem key={p.id} value={index}>
                      <Typography color="neutral.500" variant="subtitle1">
                        {p.name}
                      </Typography>
                    </MenuItem>
                  ))}
                  {/* <Typography color="neutral.400" variant="body2">
                Your tier : Premium
              </Typography> */}
                </EnvironmentSelector>
              </Box>
            </Box>
          )}
        </div>
        <Divider
          sx={{
            borderColor: '#2D3748',
            my: 3
          }}
        />
        <Box sx={{ flexGrow: 1 }}>
          {items.map((item) => (
            <NavItem key={item.title} icon={item.icon} href={item.href} title={item.title} />
          ))}
        </Box>
        <Divider sx={{ borderColor: '#2D3748' }} />
        <Box
          display="flex"
          alignItems="center"
          sx={{
            px: 2,
            py: 3,
            flexGrow: 1
          }}
        >
          <Avatar
            sx={{
              height: 40,
              width: 40,
              ml: 1
            }}
            src={me?.profile || ''}
          >
            <UserCircleIcon fontSize="small" />
          </Avatar>
          <Typography sx={{ ml: 2 }} color="neutral.400">
            {me?.firstName} {me?.lastName}
          </Typography>
        </Box>
      </Box>
    </>
  );

  if (lgUp) {
    return (
      <Drawer
        anchor="left"
        open
        PaperProps={{
          sx: {
            backgroundColor: 'neutral.900',
            color: '#FFFFFF',
            width: 280
          }
        }}
        variant="permanent"
      >
        {content}
      </Drawer>
    );
  }

  return (
    <Drawer
      anchor="left"
      onClose={onClose}
      open={open}
      PaperProps={{
        sx: {
          backgroundColor: 'neutral.900',
          color: '#FFFFFF',
          width: 280
        }
      }}
      sx={{ zIndex: (theme: Theme) => theme.zIndex.appBar + 100 }}
      variant="temporary"
    >
      {content}
    </Drawer>
  );
};
