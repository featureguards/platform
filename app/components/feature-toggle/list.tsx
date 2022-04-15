import { useState } from 'react';

import AddIcon from '@mui/icons-material/Add';
import { Box, Dialog, DialogContent, DialogTitle, Fab, List, Typography } from '@mui/material';
import useMediaQuery from '@mui/material/useMediaQuery';
import { styled, useTheme } from '@mui/system';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { EnvironmentID, list } from '../../features/feature_toggles/slice';
import { useFeatureTogglesList } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { NewFeatureToggle } from './create';
import { FeatureToggleItem } from './feature-toggle-item';

export type FeatureTogglesProps = {};

const AddButton = styled(Fab)({
  position: 'fixed',
  bottom: 10,
  right: 10
});

export const FeatureToggles = (_props: FeatureTogglesProps) => {
  const environment = useAppSelector((state) => state.projects.environment.item);
  const dispatch = useAppDispatch();
  const listProps = {
    environmentId: environment?.id
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
          <FeatureToggleItem key={ft.id} featureToggle={ft}></FeatureToggleItem>
        ))}
      </List>
    </Box>
  );
};