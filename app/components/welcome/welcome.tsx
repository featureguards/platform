import { Fragment, ReactNode, useState } from 'react';

import { Box, Button, Typography } from '@mui/material';
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
};
type StepProps = {
  title: string;
  component: ReactNode;
};

export const Welcome = ({ addresses, pendingInvites, showNewProject }: WelcomeProps) => {
  const steps: StepProps[] = [];
  const [activeStep, setActiveStep] = useState(0);
  const { refetch, loading } = useProjectsLazy();

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
      component: <Invitations invitations={pendingInvites} />
    });
  }

  const handleNewProject = async ({ err }: { err?: Error }) => {
    if (!err) {
      handleNext();
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
