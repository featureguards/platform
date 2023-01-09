import Head from 'next/head';
import { useState } from 'react';

import AddIcon from '@mui/icons-material/Add';
import { Box, Dialog, DialogContent, DialogTitle, Fab, List, Typography } from '@mui/material';
import useMediaQuery from '@mui/material/useMediaQuery';
import { styled, useTheme } from '@mui/system';

import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { EnvironmentID, list } from '../../features/dynamic_settings/slice';
import { useDynamicSettingsList } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { NewDynamicSetting } from './create';
import { DynamicSettingItem } from './dynamic-setting-item';

export type DynamicSettingsProps = {};

const AddButton = styled(Fab)({
  position: 'fixed',
  bottom: 20,
  right: 15
});

export const DynamicSettings = (_props: DynamicSettingsProps) => {
  const environment = useAppSelector((state) => state.projects.environment.item);
  const dispatch = useAppDispatch();
  const listProps = {
    environmentId: environment?.id
  };
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down('sm'));
  const { dynamicSettings, loading } = useDynamicSettingsList(listProps);
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
    <>
      <Head>
        <title>Dynamic Settings</title>
      </Head>
      <Typography sx={{ pl: 5, pb: 4 }} variant="h5">
        Dynamic Settings
      </Typography>
      <Box
        sx={{
          backgroundColor: theme.palette.background.paper
        }}
      >
        <Dialog
          maxWidth={'lg'}
          fullScreen={fullScreen}
          open={openCreate}
          onClose={() => setOpenCreate(false)}
        >
          <DialogTitle>New Dynamic Setting</DialogTitle>
          <DialogContent>
            <NewDynamicSetting
              onCreate={onCreated}
              onCancel={() => setOpenCreate(false)}
            ></NewDynamicSetting>
          </DialogContent>
        </Dialog>
        <List disablePadding={true}>
          {dynamicSettings.map((ds) => (
            <DynamicSettingItem key={ds.id} dynamicSetting={ds}></DynamicSettingItem>
          ))}
        </List>
        <AddButton color="primary" aria-label="add" onClick={handleAdd}>
          <AddIcon />
        </AddButton>
      </Box>
    </>
  );
};
