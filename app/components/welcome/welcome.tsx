import { Fragment, ReactNode, useState } from 'react';

import { Box, Typography } from '@mui/material';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Stepper from '@mui/material/Stepper';

import { ProjectInvite, UserVerifiableAddress } from '../../api';
import { useProjectsLazy } from '../hooks';
import { Invitations } from '../project/invitations';
import { NewProject } from '../project/new-project';
import { Confirmation } from './confirmation';

export type WelcomeProps = {
  addresses: UserVerifiableAddress[];
  pendingInvites: ProjectInvite[];
  showNewProject: boolean;
  refetchInvites: () => Promise<void>;
};
type StepProps = {
  title: string;
  component: ReactNode;
};

export const Welcome = ({
  addresses,
  pendingInvites,
  showNewProject,
  refetchInvites
}: WelcomeProps) => {
  const steps: StepProps[] = [];
  const [activeStep] = useState(0);
  const { refetch } = useProjectsLazy();

  if (addresses.length) {
    steps.push({
      title: 'Email Confirmation',
      component: (
        <>
          {addresses.map((addr) => (
            <Confirmation key={addr.address} address={addr.address || ''} verified={false} />
          ))}
        </>
      )
    });
  }

  if (pendingInvites.length) {
    steps.push({
      title: 'Invitations',
      component: <Invitations invitations={pendingInvites} refetch={refetchInvites} />
    });
  }

  const handleNewProject = async ({ err }: { err?: Error }) => {
    if (!err) {
      await refetch();
    }
  };

  if (showNewProject) {
    steps.push({
      title: 'New Project',
      component: <NewProject onSubmit={handleNewProject}></NewProject>
    });
  }

  if (!steps.length) {
    return <></>;
  }

  return (
    <Box sx={{ width: '100%', maxWidth: 800 }}>
      <Typography gutterBottom variant="h5">
        Let&apos;s Get Started
      </Typography>
      <Stepper activeStep={activeStep}>
        {steps.map(({ title }) => {
          const stepProps: { completed?: boolean } = {};
          const labelProps: {
            optional?: React.ReactNode;
          } = {};
          return (
            <Step key={title} {...stepProps}>
              <StepLabel {...labelProps}>{title}</StepLabel>
            </Step>
          );
        })}
      </Stepper>
      {activeStep === steps.length ? (
        <Fragment>
          <Typography sx={{ mt: 2, mb: 1 }}>All steps completed - you&apos;re finished</Typography>
          <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
            <Box sx={{ flex: '1 1 auto' }} />
          </Box>
        </Fragment>
      ) : (
        <Fragment>{steps[activeStep].component}</Fragment>
      )}
    </Box>
  );
};
