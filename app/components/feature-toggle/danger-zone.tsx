import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { useState } from 'react';

import {
  Box,
  Button,
  Card,
  CardContent,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField
} from '@mui/material';

import { Dashboard } from '../../data/api';
import { SerializeError } from '../../features/utils';
import { useFeatureTogglesListLazy, useNotifier } from '../hooks';
import { handleError } from '../hooks/utils';

export type FeatureToggleProps = {
  id: string | undefined;
  environmentId: string | undefined;
};

export const DangerZone = ({ id, environmentId }: FeatureToggleProps) => {
  const [showDelete, setShowDelete] = useState<boolean>(false);
  const [understand, setUnderstand] = useState<string>('');
  const { refetch } = useFeatureTogglesListLazy({ environmentId: environmentId });

  const notifier = useNotifier();
  const router = useRouter();
  const handleDelete = async () => {
    if (!id) return;
    try {
      await Dashboard.deleteFeatureToggle(id);
      await refetch();
      setShowDelete(false);
      router.push('/');
    } catch (err) {
      if (err) {
        handleError(router, notifier, SerializeError(err as AxiosError));
      }
    }
  };

  if (!id) return <></>;

  return (
    <Card>
      <CardContent>
        <Box
          sx={{
            alignItems: 'center',
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'center'
          }}
        >
          <Dialog open={showDelete} onClose={() => setShowDelete(false)}>
            <DialogTitle>Confirm Deletion</DialogTitle>
            <DialogContent>
              <>
                API calls to this feature flag name will start failing. Are you sure you want to do
                this?
              </>
              <TextField
                fullWidth
                size="small"
                helperText="Type 'I understand' to before deletion."
                placeholder="I understand"
                name="understand"
                onChange={(e) => setUnderstand(e.target.value)}
                value={understand}
              />
            </DialogContent>
            <DialogActions>
              <Button onClick={() => setShowDelete(false)}>Cancel</Button>
              <Button
                disabled={understand !== 'I understand'}
                color="error"
                variant="contained"
                onClick={handleDelete}
                autoFocus
              >
                Delete
              </Button>
            </DialogActions>
          </Dialog>

          <Button color="error" fullWidth variant="contained" onClick={() => setShowDelete(true)}>
            Delete
          </Button>
        </Box>
      </CardContent>
    </Card>
  );
};
