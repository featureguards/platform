import { AxiosError } from 'axios';
import NextLink from 'next/link';
import { useRouter } from 'next/router';
import { useSnackbar } from 'notistack';
import { useEffect, useMemo, useState } from 'react';
import * as Yup from 'yup';

import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { Box, Button, Container, Link, Typography } from '@mui/material';
import { SelfServiceLoginFlow, SubmitSelfServiceLoginFlowBody } from '@ory/kratos-client';

import { theme } from '../../app/theme';
import { Flow, PropsOverrides } from '../components/ory/Flow';
import SuspenseLoader from '../components/suspense-loader';
import { handleFlowError, handleGetFlowError } from '../ory/errors';
import ory from '../ory/sdk';
import { Notif } from '../utils/notif';

const Login = () => {
  const [flow, setFlow] = useState<SelfServiceLoginFlow>();
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  // Get ?flow=... from the URL
  const router = useRouter();
  const {
    return_to: returnTo,
    flow: flowId,
    // Refresh means we want to refresh the session. This is needed, for example, when we want to update the password
    // of a user.
    refresh,
    // AAL = Authorization Assurance Level. This implies that we want to upgrade the AAL, meaning that we want
    // to perform two-factor authentication/verification.
    aal
  } = router.query;

  const notifier = useMemo(
    () => new Notif({ enqueueSnackbar: enqueueSnackbar, closeSnackbar: closeSnackbar }),
    [enqueueSnackbar, closeSnackbar]
  );
  useEffect(() => {
    // If the router is not ready yet, or we already have a flow, do nothing.
    if (!router.isReady || flow) {
      return;
    }

    // If ?flow=.. was in the URL, we fetch it
    if (flowId) {
      ory
        .getSelfServiceLoginFlow(String(flowId))
        .then(({ data }) => {
          setFlow(data);
        })
        .catch(handleGetFlowError(router, 'login', setFlow, notifier));
      return;
    }

    // Otherwise we initialize it
    ory
      .initializeSelfServiceLoginFlowForBrowsers(
        Boolean(refresh),
        aal ? String(aal) : undefined,
        returnTo ? String(returnTo) : undefined
      )
      .then(({ data }) => {
        setFlow(data);
      })
      .catch(handleFlowError(router, 'login', setFlow, notifier));
  }, [flowId, router, router.isReady, aal, refresh, returnTo, flow, notifier]);

  const validationSchema = Yup.object({
    password_identifier: Yup.string()
      .email('Must be a valid email')
      .max(255)
      .required('Email is required'),
    password: Yup.string()
      .min(8, 'Password must be at least 8 characters')
      .max(255)
      .required('Password is required')
  });

  const onSubmit = (values: SubmitSelfServiceLoginFlowBody) =>
    router
      // On submission, add the flow ID to the URL but do not navigate. This prevents the user loosing
      // his data when she/he reloads the page.
      .push(`/login?flow=${flow?.id}`, undefined, { shallow: true })
      .then(() =>
        ory
          .submitSelfServiceLoginFlow(String(flow?.id), undefined, values)
          // We logged in successfully! Let's bring the user home.
          .then(() => {
            if (flow?.return_to) {
              window.location.href = flow?.return_to;
              return;
            }
            router.push('/');
          })
          .then(() => {})
          .catch(handleFlowError(router, 'login', setFlow, notifier))
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
  const props: PropsOverrides = {
    method: { variant: 'contained' },
    password_identifier: {
      label: 'Email'
    }
  };
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
            <NextLink href="/" passHref>
              <Button component="a" startIcon={<ArrowBackIcon fontSize="small" />}>
                Dashboard
              </Button>
            </NextLink>
            <Flow
              onSubmit={onSubmit}
              flow={flow}
              nodeProps={{
                provider: {
                  startIcon: <img src="/images/google_logo.svg"></img>,
                  variant: 'outlined'
                }
              }}
              only="oidc"
              childrenNodes={{ provider: 'Sign in with Google' }}
            />
            <Typography color="textSecondary" align="center" variant="body1" sx={{ pt: 3, pb: 2 }}>
              or login with email address
            </Typography>
            <Flow
              onSubmit={onSubmit}
              flow={flow}
              validationSchema={validationSchema}
              nodeProps={props}
              hideGlobalMessages={true}
              only="password"
            />
            <Typography sx={{ pt: 3 }} color="textSecondary" variant="body2">
              Don&apos;t have an account?{' '}
              <NextLink href="/register">
                <Link
                  variant="subtitle2"
                  underline="hover"
                  sx={{
                    cursor: 'pointer'
                  }}
                >
                  Sign Up
                </Link>
              </NextLink>
            </Typography>
          </Container>
        ) : (
          <SuspenseLoader></SuspenseLoader>
        )}
      </Box>
    </>
  );
};

export default Login;
