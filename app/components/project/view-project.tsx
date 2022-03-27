import { useSnackbar } from 'notistack';
import { useCallback, useEffect, useState } from 'react';

import { Card, CardContent, CardHeader, Divider, Grid } from '@mui/material';

import { Project } from '../../api';
import { Dashboard } from '../../data/api';
import { Notif } from '../../utils/notif';
import SuspenseLoader from '../suspense-loader';

export type ViewProjectProps = {
  projectID: string;
};

export const ViewProject = (props: ViewProjectProps) => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const notifier = new Notif({ enqueueSnackbar: enqueueSnackbar, closeSnackbar: closeSnackbar });
  const [project, setProject] = useState<Project>();
  const [loading, setLoading] = useState<boolean>(false);

  // https://stackoverflow.com/questions/53332321/react-hook-warnings-for-async-function-in-useeffect-useeffect-function-must-ret
  const getProject = useCallback(async () => {
    try {
      setLoading(true);
      const res = await Dashboard.getProject(props.projectID);
      setProject(res.data);
    } catch (err) {
      const msg = (err as Error).message;
      notifier.error(msg);
    } finally {
      setLoading(false);
    }
  }, [props.projectID]);
  useEffect(() => {
    getProject();
  }, [getProject]);

  if (loading) {
    return <SuspenseLoader />;
  }

  return (
    <Card>
      <CardHeader subheader={project?.description} title={project?.name} />
      <Divider />
      <CardContent>
        <Grid container spacing={3}>
          <Grid item md={6} xs={12}>
            {/* <Environments></Environments> */}
          </Grid>

          <Divider sx={{ py: 2 }} />

          <Grid item md={6} xs={12}>
            {/* <FeatureToggles></FeatureToggles> */}
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
};
