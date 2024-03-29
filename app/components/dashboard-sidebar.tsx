/* eslint-disable unused-imports/no-unused-vars */
import NextLink from 'next/link';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import {
  Avatar,
  Box,
  Divider,
  Drawer,
  Link,
  MenuItem,
  Select,
  SelectChangeEvent,
  Theme,
  Typography,
  useMediaQuery
} from '@mui/material';
import { styled } from '@mui/material/styles';

import { useAppDispatch, useAppSelector } from '../data/hooks';
import { projectsSlice } from '../features/projects/slice';
import { ChartBar as ChartBarIcon } from '../icons/chart-bar';
import { Cog as CogIcon } from '../icons/cog';
import { Selector as SelectorIcon } from '../icons/selector';
import { Sliders as SlidersIcon } from '../icons/sliders';
import { UserCircle as UserCircleIcon } from '../icons/user-circle';
import { Logo } from './logo';
import { NavItem } from './nav-item';

enum ItemType {
  FeatureFlags = 1,
  DynamicSettings,
  ProjectSettings
}

type Item = {
  href: string;
  icon: JSX.Element;
  title: string;
  itemType: ItemType;
};

const items: Item[] = [
  {
    href: '/',
    icon: <SlidersIcon fontSize="small" />,
    title: 'Feature Flags',
    itemType: ItemType.FeatureFlags
  },
  {
    href: '/dynamic-settings',
    icon: <ChartBarIcon fontSize="small" />,
    title: 'Dynamic Settings (Beta)',
    itemType: ItemType.DynamicSettings
  },
  {
    href: '/project/settings',
    icon: <CogIcon fontSize="small" />,
    title: 'Project Settings',
    itemType: ItemType.ProjectSettings
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
  const dispatch = useAppDispatch();
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

  const handleEnvironmentChange = (e: SelectChangeEvent<any>) => {
    setCurrentIndex(Number(e.target.value || 0));
  };

  const [currentIndex, setCurrentIndex] = useState<number>(0);
  const me = useAppSelector((state) => state.users.me);
  const projectDetails = useAppSelector((state) => state.projects.details);
  const currentProject = projectDetails?.item;

  useEffect(() => {
    dispatch(
      projectsSlice.actions.setEnvironment(currentProject?.environments?.[currentIndex] || null)
    );
  }, [currentProject, currentIndex, dispatch]);

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
          {!!currentProject?.environments?.length && (
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
                  onChange={handleEnvironmentChange}
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
                  {currentProject?.environments?.map((p, index) => (
                    <MenuItem key={p.id} value={index}>
                      <Typography color="grey.500" variant="subtitle1">
                        {p.name}
                      </Typography>
                    </MenuItem>
                  ))}
                  {/* <Typography color="grey.400" variant="body2">
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
          {items
            .filter((v) => !!currentProject || v.itemType === ItemType.FeatureFlags)
            .map((item) => (
              <NavItem key={item.title} icon={item.icon} href={item.href} title={item.title} />
            ))}
        </Box>
        <Divider sx={{ borderColor: '#2D3748' }} />
        <Link href="/account" underline="none">
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
            <Typography sx={{ ml: 2 }} color="grey.400">
              {me?.firstName} {me?.lastName}
            </Typography>
          </Box>
        </Link>
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
            backgroundColor: 'grey.900',
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
          backgroundColor: 'grey.900',
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
