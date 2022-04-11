import { useState } from 'react';

import AddIcon from '@mui/icons-material/Add';
import {
  Box,
  Dialog,
  DialogContent,
  DialogTitle,
  Fab,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography
} from '@mui/material';
import useMediaQuery from '@mui/material/useMediaQuery';
import { styled, useTheme } from '@mui/system';

import { FeatureToggle } from '../../api';
import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { EnvironmentID, list } from '../../features/feature_toggles/slice';
import { useFeatureTogglesList } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { NewFeatureToggle } from './create';
import { LiveToggleIcon } from './icon';

export type FeatureTogglesProps = {};

const AddButton = styled(Fab)({
  position: 'fixed',
  bottom: 10,
  right: 10
});

export const FeatureToggles = (_props: FeatureTogglesProps) => {
  const { item: project } = useAppSelector((state) => state.projects.details);
  const environment = useAppSelector((state) => state.projects.environment.item);
  const dispatch = useAppDispatch();
  const listProps = {
    projectID: project?.id,
    environmentID: environment?.id
  };
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down('sm'));
  const { featureToggles, loading } = useFeatureTogglesList(listProps);
  const [openCreate, setOpenCreate] = useState(false);
  const onCreated = async () => {
    setOpenCreate(false);
    // Refetch
    await dispatch(list(listProps as EnvironmentID)).unwrap();
  };

  if (loading) {
    return <SuspenseLoader></SuspenseLoader>;
  }

  const handleAdd = () => {
    setOpenCreate(true);
  };

  return (
    <Box
      sx={{
        pt: 5,
        pb: 5,
        backgroundColor: theme.palette.background.paper
      }}
    >
      <Typography sx={{ pl: 5, pb: 4 }} variant="h5">
        Feature Toggles
      </Typography>
      <Dialog
        maxWidth={'md'}
        fullScreen={fullScreen}
        open={openCreate}
        onClose={() => setOpenCreate(false)}
      >
        <DialogTitle>New Feature Toggle</DialogTitle>
        <DialogContent>
          <NewFeatureToggle onCreate={onCreated}></NewFeatureToggle>
        </DialogContent>
      </Dialog>
      <AddButton color="primary" aria-label="add" onClick={handleAdd}>
        <AddIcon />
      </AddButton>
      <List>
        {featureToggles.map((ft) => (
          <ListItem key={ft.name} disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <LiveToggleIcon featureToggle={ft as FeatureToggle}></LiveToggleIcon>
              </ListItemIcon>
              <ListItemText sx={{ pl: 2 }} primary={ft.name!} />
              <ListItemText secondary={ft.description} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );
};
