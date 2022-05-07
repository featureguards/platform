import { useRouter } from 'next/router';
import { useState } from 'react';

import AddIcon from '@mui/icons-material/Add';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  Grid,
  IconButton,
  Typography
} from '@mui/material';
import { SerializedError } from '@reduxjs/toolkit';

import { Environment } from '../../api';
import { Dashboard } from '../../data/api';
import { NewApiKey } from '../api-key/new';
import { ApiKeyView } from '../api-key/view';
import { useApiKeysList, useNotifier } from '../hooks';
import { handleError } from '../hooks/utils';
import SuspenseLoader from '../suspense-loader';
import { CloneEnvironment } from './clone';

export type ViewEnvironmentProps = {
  environment: Environment;
  refetchProject: () => Promise<void>;
};

export const ViewEnvironment = ({
  environment,
  refetchProject,
  ...others
}: ViewEnvironmentProps) => {
  const { apiKeys, loading, refetch } = useApiKeysList({ environmentId: environment.id });
  const [showNewApiKey, setShowNewApiKey] = useState<boolean>(false);
  const [showCloneEnv, setShowCloneEnv] = useState<boolean>(false);
  const [showDeleteEnv, setShowDeleteEnv] = useState<boolean>(false);
  const notifier = useNotifier();
  const router = useRouter();

  if (loading) {
    return <SuspenseLoader />;
  }
  if (!environment.id) {
    return <></>;
  }
  const handleDelete = async () => {
    try {
      await Dashboard.deleteEnvironment(environment.id!);
      await refetchProject();
      setShowDeleteEnv(false);
    } catch (err) {
      if (err) {
        handleError(router, notifier, err as SerializedError);
      }
    }
  };
  return (
    <Card {...others}>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'row',
          alignItems: 'center',
          justifyContent: 'space-between'
        }}
      >
        <CardHeader subheader={environment.description} title={environment.name} />
        <Dialog open={showCloneEnv} onClose={() => setShowCloneEnv(false)}>
          <CloneEnvironment
            id={environment.id}
            onSubmit={async ({ err }: { err?: Error }) => {
              setShowCloneEnv(false);
              if (!err) {
                await refetchProject();
              }
            }}
          ></CloneEnvironment>
        </Dialog>
        <Box>
          <Button
            sx={{ ml: 2, maxHeight: 40 }}
            color="primary"
            onClick={() => setShowCloneEnv(true)}
          >
            Clone
          </Button>
          <Dialog open={showDeleteEnv} onClose={() => setShowDeleteEnv(false)}>
            <DialogTitle>Confirm Deletion</DialogTitle>
            <DialogContent>
              Are you sure you want to delete the environment permanently?
            </DialogContent>
            <DialogActions>
              <Button onClick={() => setShowDeleteEnv(false)}>Cancel</Button>
              <Button color="error" variant="contained" onClick={handleDelete} autoFocus>
                Confirm
              </Button>
            </DialogActions>
          </Dialog>

          <Button
            sx={{ ml: 2, maxHeight: 40 }}
            color="error"
            onClick={() => setShowDeleteEnv(true)}
          >
            Delete
          </Button>
        </Box>
      </Box>

      <CardContent sx={{ my: -5 }}>
        <Grid container>
          <Grid item xs={12}>
            <Dialog open={showNewApiKey} onClose={() => setShowNewApiKey(false)}>
              <DialogTitle>New API Key</DialogTitle>
              <DialogContent>
                <NewApiKey
                  environmentId={environment.id}
                  onSubmit={async ({ err }: { err?: Error }) => {
                    setShowNewApiKey(false);
                    if (!err) {
                      await refetch();
                    }
                  }}
                ></NewApiKey>
              </DialogContent>
            </Dialog>
            <Box
              sx={{
                alignItems: 'center',
                display: 'flex',
                flexDirection: 'row'
              }}
            >
              <Typography variant="button">API Keys</Typography>
              <IconButton
                onClick={async () => {
                  setShowNewApiKey(true);
                }}
              >
                <AddIcon></AddIcon>
              </IconButton>
            </Box>
            <Box sx={{ mb: 1 }}>
              {apiKeys?.map((apiKey, index) => (
                <>
                  {index > 0 && <Divider sx={{ my: 1 }} />}
                  <ApiKeyView
                    key={index}
                    apiKey={apiKey}
                    onDelete={async () => {
                      await refetch();
                    }}
                  ></ApiKeyView>
                </>
              ))}
            </Box>
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
};
