import { Box, Button, Chip, Card, CardContent, Divider, Typography } from '@mui/material';

export type ProjectInvitation = {
  projectID: string;
  projectName: string;
  url: string;
  accepted: boolean;
};
export type ProjectInvitationsProps = {
  invitations: ProjectInvitation[];
};

const Invitation = (props: ProjectInvitation & { index: number }) => {
  return (
    <>
      {props.index > 0 && <Divider key={props.index} sx={{ gridColumn: '1/8' }}></Divider>}
      <Typography variant="h6" gutterBottom sx={{ gridColumn: '1/3' }}>
        {props.projectName}
      </Typography>
      <Chip
        sx={{ gridColumn: '6 / 7' }}
        color={props.accepted ? 'success' : 'warning'}
        label={props.accepted ? 'Accepted' : 'Pending'}
      />
      {props.accepted ? (
        <Box sx={{ gridColumn: '7 / 8' }}></Box>
      ) : (
        <Button sx={{ gridColumn: '7 / 8' }} color="primary">
          Accept
        </Button>
      )}
    </>
  );
};

export const ProjectInvitations = (props: ProjectInvitationsProps) => {
  return (
    <Card {...props}>
      <CardContent>
        <Typography gutterBottom variant="subtitle2">
          Project Invitations
        </Typography>
        <Box
          alignItems={'center'}
          sx={{
            display: 'grid',
            gridAutoColumns: '1fr',
            gap: 1
          }}
        >
          {props.invitations.map((invitation, index) => (
            <Invitation index={index} key={index} {...invitation}></Invitation>
          ))}
        </Box>
      </CardContent>
    </Card>
  );
};
