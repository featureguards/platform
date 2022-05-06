import { useState } from 'react';

import AddIcon from '@mui/icons-material/Add';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Dialog,
  DialogContent,
  DialogTitle,
  Divider,
  Grid,
  IconButton,
  Typography
} from '@mui/material';

import { Environment } from '../../api';
import { NewApiKey } from '../api-key/new';
import { ApiKeyView } from '../api-key/view';
import { useApiKeysList } from '../hooks/api_keys';
import SuspenseLoader from '../suspense-loader';

export type ViewEnvironmentProps = {
  environment: Environment;
};

export const ViewEnvironment = ({ environment, ...others }: ViewEnvironmentProps) => {
  const { apiKeys, loading, refetch } = useApiKeysList({ environmentId: environment.id });
  const [showNewApiKey, setShowNewApiKey] = useState<boolean>(false);

  if (loading) {
    return <SuspenseLoader />;
  }
  console.log(apiKeys?.map((a) => a.id));
  if (!environment.id) {
    return <></>;
  }
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
        {/* <Button sx={{ ml: 2, maxHeight: 40 }} color="primary">
          Clone
        </Button> */}
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
