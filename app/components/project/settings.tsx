import { Card, CardContent, CardHeader, Divider, Grid } from '@mui/material';

import { useProject } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { Environments } from './environments';

export type ProjectSettingsProps = {
  projectID: string;
};

export const ProjectSettings = (props: ProjectSettingsProps) => {
  const { current, loading, refetch } = useProject({ projectID: props.projectID });

  if (loading) {
    return <SuspenseLoader />;
  }

  if (!current) return <></>;

  return (
    <Card>
      <CardHeader subheader={current.description} title={current.name} />
      <Divider />
      <CardContent>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Environments environments={current?.environments} refetch={refetch} />
          </Grid>

          <Divider sx={{ py: 2 }} />

          <Grid item md={6} xs={12}>
            {/* <Members></Members> */}
          </Grid>

          <Divider sx={{ py: 2 }} />

          <Grid item md={6} xs={12}>
            {/* <Invitations></Invitations> */}
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
};
