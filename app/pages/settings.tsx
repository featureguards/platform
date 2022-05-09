import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { Box, Button, Container } from '@mui/material';
import { SelfServiceSettingsFlow, SubmitSelfServiceSettingsFlowBody } from '@ory/kratos-client';

import { useNotifier } from '../components/hooks';
import { Flow } from '../components/ory/Flow';
import SuspenseLoader from '../components/suspense-loader';
import { handleFlowError, handleGetFlowError } from '../ory/errors';
import ory from '../ory/sdk';
import { theme } from '../theme';

const Settings = () => {
  const [flow, setFlow] = useState<SelfServiceSettingsFlow>();
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
        .getSelfServiceSettingsFlow(String(flowId))
        .then(({ data }) => {
          setFlow(data);
        })
        .catch(handleGetFlowError(router, 'settings', resetFlow, notifier));
      return;
    }

    // Otherwise we initialize it
    ory
      .initializeSelfServiceSettingsFlowForBrowsers(returnTo ? String(returnTo) : undefined)
      .then(({ data }) => {
        setFlow(data);
      })
      .catch(handleFlowError(router, 'settings', resetFlow, notifier));
  }, [flowId, router, router.isReady, returnTo, flow, notifier]);

  const onSubmit = (values: SubmitSelfServiceSettingsFlowBody) =>
    router
      // On submission, add the flow ID to the URL but do not navigate. This prevents the user loosing
      // his data when she/he reloads the page.
      .push(`/settings?flow=${flow?.id}`, undefined, { shallow: true })
      .then(() =>
        ory
          .submitSelfServiceSettingsFlow(String(flow?.id), undefined, values)
          // We logged in successfully! Let's bring the user home.
          .then(() => {
            window.location.href = flow?.return_to || '/';
          })
          .then(() => {})
          .catch(handleFlowError(router, 'settings', resetFlow, notifier))
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

export default Settings;
