import { FC, Fragment, ReactNode, useState } from 'react';

import { Box, Button, Typography } from '@mui/material';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Stepper from '@mui/material/Stepper';

import { ProjectInviteStatus } from '../../api/enums';
import { useAppSelector } from '../../data/hooks';
import { useProjects, useUserInvites } from '../hooks';
import { NewProject } from '../project/new-project';
import SuspenseLoader from '../suspense-loader';
import { Confirmation } from './confirmation';
import { ProjectInvitations } from './project-invitations';

export type WelcomeProps = {};
type StepProps = {
  title: string;
  component: ReactNode;
};

export const Welcome: FC<WelcomeProps> = () => {
  const me = useAppSelector((state) => state.users.me);
  const steps: StepProps[] = [];
  const { invites, loading: invitesLoading } = useUserInvites();
  const { projects, loading: projectsLoading } = useProjects();
  const [activeStep, setActiveStep] = useState(0);
  if (!me?.addresses?.length) {
    // Not logged in or something wrong fetching the user.
    return <></>;
  }

  if (invitesLoading || projectsLoading) {
    return <SuspenseLoader></SuspenseLoader>;
  }
  if (!me?.addresses[0].verified) {
    steps.push({
      title: 'Email Confirmation',
      component: (
        <Confirmation
          address={me?.addresses[0].address || ''}
          verified={Boolean(me?.addresses[0].verified)}
        ></Confirmation>
      )
    });
  }

  const pendingInvites = invites.filter((el) => el.status === ProjectInviteStatus.PENDING);
  if (pendingInvites.length) {
    steps.push({
      title: 'Invitations',
      component: <ProjectInvitations invitations={pendingInvites}></ProjectInvitations>
    });
  }

  if (!projects.length) {
    steps.push({
      title: 'New Project',
      component: <NewProject></NewProject>
    });
  }

  if (!steps.length) {
    return <></>;
  }

  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleReset = () => {
    setActiveStep(0);
  };

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
            <Button onClick={handleReset}>Reset</Button>
          </Box>
        </Fragment>
      ) : (
        <Fragment>
          {steps[activeStep].component}
          {steps.length > 1 && (
            <>
              <Typography sx={{ mt: 2, mb: 1 }}>Step {activeStep + 1}</Typography>
              <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
                <Button
                  color="inherit"
                  disabled={activeStep === 0}
                  onClick={handleBack}
                  sx={{ mr: 1 }}
                >
                  Back
                </Button>
                <Box sx={{ flex: '1 1 auto' }} />
                <Button onClick={handleNext}>
                  {activeStep === steps.length - 1 ? 'Finish' : 'Next'}
                </Button>
              </Box>
            </>
          )}
        </Fragment>
      )}
    </Box>
  );
};
