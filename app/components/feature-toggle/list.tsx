import { useState } from 'react';

import AddIcon from '@mui/icons-material/Add';
import { Box, Dialog, DialogContent, DialogTitle, Divider, Fab, Typography } from '@mui/material';
import useMediaQuery from '@mui/material/useMediaQuery';
import { styled, useTheme } from '@mui/system';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { EnvironmentID, list } from '../../features/feature_toggles/slice';
import { useFeatureTogglesList } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { NewFeatureToggle } from './create';

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
    <Box sx={{ flexGrow: 1 }}>
      <Typography variant="h5">Feature Toggles</Typography>
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
      <Divider />
      <Box
        alignItems="center"
        sx={{
          display: 'grid',
          gridAutoColumns: '1fr',
          gap: 1
        }}
      >
        {featureToggles.map((ft, index) => (
          <Typography key={index}>{ft.name}</Typography>
        ))}
      </Box>
    </Box>
  );
};
