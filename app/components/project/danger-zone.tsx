import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { useState } from 'react';

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
  TextField
} from '@mui/material';

import { Dashboard } from '../../data/api';
import { SerializeError } from '../../features/utils';
import { useNotifier, useProjectsLazy } from '../hooks';
import { handleError } from '../hooks/utils';

export type EnvironmentProps = {
  projectID: string | undefined;
};

export const DangerZone = ({ projectID }: EnvironmentProps) => {
  const [showDelete, setShowDelete] = useState<boolean>(false);
  const [understand, setUnderstand] = useState<string>('');
  const { refetch } = useProjectsLazy();

  const notifier = useNotifier();
  const router = useRouter();
  const handleDelete = async () => {
    if (!projectID) return;
    try {
      await Dashboard.deleteProject(projectID);
      await refetch();
      setShowDelete(false);
    } catch (err) {
      if (err) {
        handleError(router, notifier, SerializeError(err as AxiosError));
      }
    }
  };

  if (!projectID) return <></>;

  return (
    <Card>
      <CardHeader title="Danger Zone" />
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
              <>This action is permanent. Are you sure you want to do this?</>
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
            Delete Project
          </Button>
        </Box>
      </CardContent>
    </Card>
  );
};
