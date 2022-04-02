import { Box, Card, CardActions, CardContent, Chip, Typography } from '@mui/material';

import Verification from '../ory/verification';

export type ConfirmationProps = {
  address: string;
  verified: boolean;
};

export const Confirmation = (props: ConfirmationProps) => {
  const { address, verified, ...other } = props;

  return (
    <Card {...other}>
      <CardContent>
        <Box
          sx={{
            alignItems: 'center',
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-between'
          }}
        >
          <Typography color="textPrimary" gutterBottom>
            {address}
          </Typography>
          <Chip
            color={verified ? 'success' : 'warning'}
            label={verified ? 'Confirmed' : 'Unconfirmed'}
          />
        </Box>
      </CardContent>
      <CardActions>
        <Verification color="primary" variant="text">
          Resend
        </Verification>
      </CardActions>
    </Card>
  );
};
