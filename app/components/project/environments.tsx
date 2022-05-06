import { Grid, Typography } from '@mui/material';

import { useAppSelector } from '../../data/hooks';
import { ViewEnvironment } from '../environment/view';

export const Environments = () => {
  const projectDetails = useAppSelector((state) => state.projects.details);
  const currentProject = projectDetails?.item;

  return (
    <>
      <Typography variant="h5">Environments</Typography>
      <Grid container spacing={3}>
        {currentProject?.environments?.map((env) => (
          <Grid key={env.id} item xs={12}>
            <ViewEnvironment environment={env}></ViewEnvironment>
          </Grid>
        ))}
      </Grid>
    </>
  );
};
