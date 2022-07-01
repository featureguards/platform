import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { Box, Button, Container, Typography } from '@mui/material';
import { SelfServiceRecoveryFlow, SubmitSelfServiceRecoveryFlowBody } from '@ory/kratos-client';

import { useNotifier } from '../../components/hooks';
import { Flow } from '../../components/ory/Flow';
import SuspenseLoader from '../../components/suspense-loader';
import { handleFlowError, handleGetFlowError } from '../../ory/errors';
import ory, { urlForFlow } from '../../ory/sdk';
import { theme } from '../../theme';

const Recovery = () => {
  const [flow, setFlow] = useState<SelfServiceRecoveryFlow>();
  // Get ?flow=... from the URL
  const router = useRouter();
  const { return_to: returnTo, flow: flowId } = router.query;
  const notifier = useNotifier();
  const resetFlow = () => {
    setFlow(undefined);
  };

  useEffect(() => {
    // If the router is not ready yet, or we already have a flow, do nothing.
    if (!router.isReady || flow) {
      return;
    }

    // If ?flow=.. was in the URL, we fetch it
    if (flowId) {
      ory
        .getSelfServiceRecoveryFlow(String(flowId))
        .then(({ data }) => {
          setFlow(data);
        })
        .catch(handleGetFlowError(router, 'recovery', resetFlow, notifier));
      return;
    }

    // Otherwise we initialize it
    ory
      .initializeSelfServiceRecoveryFlowForBrowsers(returnTo ? String(returnTo) : undefined)
      .then(({ data }) => {
        setFlow(data);
      })
      .catch(handleFlowError(router, 'recovery', resetFlow, notifier));
  }, [flowId, router, router.isReady, returnTo, flow, notifier]);

  const onSubmit = (values: SubmitSelfServiceRecoveryFlowBody) =>
    router
      // On submission, add the flow ID to the URL but do not navigate. This prevents the user loosing
      // his data when she/he reloads the page.
      .push(`${urlForFlow('recovery')}?flow=${flow?.id}`, undefined, { shallow: true })
      .then(() =>
        ory
          .submitSelfServiceRecoveryFlow(String(flow?.id), values)
          // We logged in successfully! Let's bring the user home.
          .then(() => {
            notifier.info('Please, check your email for next steps.', 3000);
          })
          .then(() => {})
          .catch(handleFlowError(router, 'recovery', resetFlow, notifier))
          .catch((err: AxiosError) => {
            // If the previous handler did not catch the error it's most likely a form validation error
            if (err.response?.status === 400) {
              // Yup, it is!
              setFlow(err.response?.data);
              return;
            }

            return Promise.reject(err);
          })
      );
  return (
    <>
      <Box
        component="main"
        sx={{
          alignItems: 'center',
          display: 'flex',
          flexGrow: 1,
          minHeight: '100%'
        }}
      >
        {flow ? (
          <Container
            maxWidth="xs"
            sx={{
              backgroundColor: theme.palette.background.paper,
              pt: 5,
              pb: 5,
              borderRadius: 1
            }}
          >
            <Typography color="textSecondary" align="center" variant="h6" sx={{ pt: 3, pb: 2 }}>
              Account Reset
            </Typography>
            <Button
              component="a"
              onClick={() => (window.location.href = '/')}
              startIcon={<ArrowBackIcon fontSize="small" />}
            >
              Dashboard
            </Button>
            <Flow onSubmit={onSubmit} flow={flow} hideGlobalMessages={false} />
          </Container>
        ) : (
          <SuspenseLoader></SuspenseLoader>
        )}
      </Box>
    </>
  );
};

export default Recovery;
