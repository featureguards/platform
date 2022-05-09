import { Card, CardContent, Divider, Grid, Typography } from '@mui/material';

import { theme } from '../../theme';
import { useProject, useProjectInvites } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { DangerZone } from './danger-zone';
import { Environments } from './environments';
import { Invitations } from './invitations';
import { Members } from './members';

export type ProjectSettingsProps = {
  projectID: string;
};

export const ProjectSettings = ({ projectID }: ProjectSettingsProps) => {
  const { current, loading, refetch } = useProject({ projectID: projectID });
  const {
    invites,
    loading: invitesLoading,
    refetch: refetchInvites
  } = useProjectInvites(projectID);
  if (loading || invitesLoading) {
    return <SuspenseLoader />;
  }

  if (!current) return <></>;

  return (
    <Card
      sx={{
        backgroundColor: theme.palette.background.default
      }}
    >
      <CardContent>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Grid container spacing={1}>
                  <Grid item xs={2}>
                    <Typography sx={{ mr: 2 }} variant="h6">
                      Project Name
                    </Typography>
                  </Grid>
                  <Grid item xs={10}>
                    <Typography>{current.name}</Typography>
                  </Grid>
                  <Grid item xs={2}>
                    <Typography sx={{ mr: 2 }} variant="h6">
                      Description
                    </Typography>
                  </Grid>
                  <Grid item xs={10}>
                    <Typography variant="body1">{current.description}</Typography>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12}>
            <Environments environments={current?.environments} refetch={refetch} />
          </Grid>

          <Grid item md={6} xs={12}>
            <Members projectID={current?.id}></Members>
          </Grid>

          <Divider sx={{ py: 2 }} />

          <Grid item md={6} xs={12}>
            <Invitations showSend invitations={invites} refetch={refetchInvites} forProject />
          </Grid>
          <Grid item xs={12}>
            <DangerZone projectID={projectID} />
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
};
