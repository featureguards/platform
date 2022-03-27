import { Box, Button, Chip, Card, CardActions, CardContent, Typography } from '@mui/material';

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
          <Chip color={verified ? 'success' : 'warning'} label="Verified" />
        </Box>
      </CardContent>
      {verified && (
        <CardActions>
          <Button color="primary" variant="text">
            Resend
          </Button>
        </CardActions>
      )}
    </Card>
  );
};
