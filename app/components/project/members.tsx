import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { useState } from 'react';

import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  Typography
} from '@mui/material';

import { ProjectMember } from '../../api';
import { ProjectMemberRole } from '../../api/enums';
import { Dashboard } from '../../data/api';
import { useAppSelector } from '../../data/hooks';
import { SerializeError } from '../../features/utils';
import { projectMemberRoleTypeName } from '../../utils/display';
import { useNotifier } from '../hooks';
import { useProjectMembers } from '../hooks/project_members';
import { handleError } from '../hooks/utils';
import SuspenseLoader from '../suspense-loader';

type MemberProps = {
  member: ProjectMember;
  index: number;
  refetch: () => Promise<void>;
  isAdmin?: boolean;
};

const Member = ({ member, index, isAdmin, refetch }: MemberProps) => {
  const [showDelete, setShowDelete] = useState<boolean>(false);
  const notifier = useNotifier();
  const router = useRouter();
  const handleDelete = async () => {
    try {
      await Dashboard.deleteProjectMember(member.id!);
      await refetch();
      setShowDelete(false);
    } catch (err) {
      if (err) {
        handleError(router, notifier, SerializeError(err as AxiosError));
      }
    }
  };

  return (
    <>
      {index > 0 && <Divider sx={{ gridColumn: '1/8' }} />}
      <Typography variant="h6" gutterBottom sx={{ gridColumn: '1/4' }}>
        {member.user?.firstName + ' ' + member.user?.lastName}
      </Typography>
      <Chip
        sx={{ gridColumn: '4 / 6' }}
        label={projectMemberRoleTypeName(
          (member.role as ProjectMemberRole) ?? ProjectMemberRole.UNKNOWN
        )}
      />
      {isAdmin && (
        <>
          <Dialog open={showDelete} onClose={() => setShowDelete(false)}>
            <DialogTitle>Confirm Removal</DialogTitle>
            <DialogContent>Are you sure you want to remove this member?</DialogContent>
            <DialogActions>
              <Button onClick={() => setShowDelete(false)}>Cancel</Button>
              <Button color="error" variant="contained" onClick={handleDelete} autoFocus>
                Remove
              </Button>
            </DialogActions>
          </Dialog>

          <Button sx={{ gridColumn: '7 / 8' }} color="error" onClick={() => setShowDelete(true)}>
            Remove
          </Button>
        </>
      )}
    </>
  );
};

export type MembersProps = {
  projectID?: string;
};

export const Members = ({ projectID }: MembersProps) => {
  const me = useAppSelector((state) => state.users.me);
  const { members: projectMembers, loading, refetch } = useProjectMembers(projectID);
  if (!projectID) {
    return <></>;
  }
  if (loading) {
    return <SuspenseLoader />;
  }
  const isAdmin = !!projectMembers?.members?.filter(
    (m) => m.user?.id === me?.id && m.role === ProjectMemberRole.ADMIN
  )?.length;
  return (
    <Card>
      <CardHeader title="Members"></CardHeader>
      <CardContent>
        <Box
          alignItems="center"
          sx={{
            display: 'grid',
            gridAutoColumns: '1fr',
            gap: 1
          }}
        >
          {projectMembers?.members.map((member, index) => (
            <Member
              key={member.id}
              index={index}
              isAdmin={isAdmin}
              member={member}
              refetch={refetch}
            ></Member>
          ))}
        </Box>
      </CardContent>
    </Card>
  );
};
