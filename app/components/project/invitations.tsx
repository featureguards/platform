import { useFormik } from 'formik';
import { useState } from 'react';
import * as Yup from 'yup';

import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  TextField,
  Typography
} from '@mui/material';

import { ProjectInvite } from '../../api';
import { ProjectInviteStatus } from '../../api/enums';
import { Dashboard } from '../../data/api';
import { useAppSelector } from '../../data/hooks';
import { useNotifier } from '../hooks';

export type ProjectInvitationsProps = {
  invitations: ProjectInvite[];
  showSend?: boolean;
  forProject?: boolean; // per project and not per user.
  refetch?: () => Promise<void>;
};

const Invitation = (props: ProjectInvite & { index: number; forProject: boolean }) => {
  return (
    <>
      {props.index > 0 && <Divider key={props.index} sx={{ gridColumn: '1/8' }}></Divider>}
      <Typography variant="h6" gutterBottom sx={{ gridColumn: '1/3' }}>
        {props.forProject ? props.email : props.projectName}
      </Typography>
      <Chip
        sx={{ gridColumn: '5 / 6' }}
        color={props.status === ProjectInviteStatus.ACCEPTED ? 'success' : 'warning'}
        label={props.status === ProjectInviteStatus.ACCEPTED ? 'Accepted' : 'Pending'}
      />
    </>
  );
};

export const Invitations = ({
  invitations,
  showSend,
  forProject,
  refetch,
  ...others
}: ProjectInvitationsProps) => {
  const [showNewInvite, setShowNewInvite] = useState<boolean>(false);
  const projectDetails = useAppSelector((state) => state.projects.details);
  const currentProject = projectDetails?.item;
  const notifier = useNotifier();

  const formik = useFormik({
    initialValues: {
      email: '',
      firstName: ''
    },
    validationSchema: Yup.object({
      email: Yup.string().email('Must be a valid email').required('Email is required'),
      firstName: Yup.string().required('First Name is required')
    }),
    onSubmit: async (values) => {
      try {
        if (!currentProject?.id) return;
        await Dashboard.createProjectInvite(currentProject?.id, { ...values });
        setShowNewInvite(false);
        if (refetch) {
          await refetch();
        }
      } catch (err) {
        notifier.error('Error sending invite.');
      }
    }
  });
  const handleSubmit = async () => {
    await formik.submitForm();
  };
  return (
    <Card {...others}>
      <CardHeader title="Project Invitations"></CardHeader>
      <CardContent>
        <Box
          alignItems="center"
          sx={{
            display: 'grid',
            gridAutoColumns: '1fr',
            gap: 1
          }}
        >
          {invitations.map((invitation, index) => (
            <Invitation
              index={index}
              key={index}
              forProject={!!forProject}
              {...invitation}
            ></Invitation>
          ))}
        </Box>
      </CardContent>
      <Dialog open={showNewInvite} onClose={() => setShowNewInvite(false)}>
        <DialogTitle>New Invitation</DialogTitle>
        <DialogContent>
          <TextField
            label="Email"
            name="email"
            margin="dense"
            error={Boolean(formik.touched.email && formik.errors.email)}
            sx={{ mr: 2 }}
            onBlur={formik.handleBlur}
            onChange={formik.handleChange}
            helperText={Boolean(formik.touched.email) ? formik.errors.email : ''}
            value={formik.values.email}
            variant="outlined"
          />
          <TextField
            label="First Name"
            name="firstName"
            margin="dense"
            error={Boolean(formik.touched.firstName && formik.errors.firstName)}
            sx={{ mr: 2 }}
            onBlur={formik.handleBlur}
            onChange={formik.handleChange}
            helperText={Boolean(formik.touched.firstName) ? formik.errors.firstName : ''}
            value={formik.values.firstName}
            variant="outlined"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleSubmit}>Send</Button>
        </DialogActions>
      </Dialog>
      {showSend && (
        <CardActions>
          <Button onClick={() => setShowNewInvite(true)}>Invite</Button>
        </CardActions>
      )}
    </Card>
  );
};
