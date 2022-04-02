import { Card, CardContent, CardHeader, Divider, Grid } from '@mui/material';

import { useProject } from '../hooks';
import SuspenseLoader from '../suspense-loader';

export type ProjectSettingsProps = {
  projectID: string;
};

export const ProjectSettings = (props: ProjectSettingsProps) => {
  const { current, loading } = useProject({ projectID: props.projectID });

  if (loading) {
    return <SuspenseLoader />;
  }

  return (
    <Card>
      <CardHeader subheader={current?.description} title="Settings" />
      <Divider />
      <CardContent>
        <Grid container spacing={3}>
          <Grid item md={6} xs={12}>
            {/* <Environments></Environments> */}
          </Grid>

          <Divider sx={{ py: 2 }} />

          <Grid item md={6} xs={12}>
            {/* <ApiKeys></ApiKeys> */}
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
