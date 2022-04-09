import { Box, Card, CardContent, Chip, Divider, Typography } from '@mui/material';

import { ProjectInvite } from '../../api';
import { ProjectInviteStatus } from '../../api/enums';

export type ProjectInvitationsProps = {
  invitations: ProjectInvite[];
};

const Invitation = (props: ProjectInvite & { index: number }) => {
  return (
    <>
      {props.index > 0 && <Divider key={props.index} sx={{ gridColumn: '1/8' }}></Divider>}
      <Typography variant="h6" gutterBottom sx={{ gridColumn: '1/3' }}>
        {props.projectName}
      </Typography>
      <Chip
        sx={{ gridColumn: '5 / 6' }}
        color={props.status === ProjectInviteStatus.ACCEPTED ? 'success' : 'warning'}
        label={props.status === ProjectInviteStatus.ACCEPTED ? 'Accepted' : 'Pending'}
      />
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
          alignItems="center"
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
