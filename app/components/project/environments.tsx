import { Card, CardContent, CardHeader, Divider, Grid } from '@mui/material';

import { useProject } from '../hooks';
import SuspenseLoader from '../suspense-loader';

export type ViewProjectProps = {
  projectID: string;
};

export const ViewProject = (props: ViewProjectProps) => {
  const { current, loading } = useProject({ projectID: props.projectID });
  if (loading) {
    return <SuspenseLoader />;
  }

  return (
    <Card>
      <CardHeader subheader={current?.description} title={current?.name} />
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
