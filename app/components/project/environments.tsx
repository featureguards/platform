import { Grid, Typography } from '@mui/material';

import { Environment } from '../../api';
import { ViewEnvironment } from '../environment/view';

export type EnvironmentProps = {
  environments: Environment[] | undefined;
  refetch: () => Promise<void>;
};

export const Environments = ({ environments, refetch }: EnvironmentProps) => {
  return (
    <>
      <Typography sx={{ py: 2 }} variant="h5">
        Environments
      </Typography>
      <Grid container spacing={1}>
        {environments?.map((env) => (
          <Grid key={env.id} item xs={12}>
            <ViewEnvironment environment={env} refetchProject={refetch}></ViewEnvironment>
          </Grid>
        ))}
      </Grid>
    </>
  );
};
