import { AxiosError } from 'axios';
import Head from 'next/head';
import NextLink from 'next/link';
import { NextRouter, useRouter } from 'next/router';
import { useEffect, useState } from 'react';
import * as Yup from 'yup';

import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { Box, Button, Container, Link, Typography } from '@mui/material';
import {
  SelfServiceRegistrationFlow,
  SubmitSelfServiceRegistrationFlowBody
} from '@ory/kratos-client';

import { theme } from '../../app/theme';
import { useNotifier } from '../components/hooks';
import { Flow, PropsOverrides } from '../components/ory/Flow';
import SuspenseLoader from '../components/suspense-loader';
import { handleFlowError } from '../ory/errors';
import ory from '../ory/sdk';

import type { NextPage } from 'next';
// Renders the registration page
const Registration: NextPage = () => {
  const router = useRouter();

  // The "flow" represents a registration process and contains
  // information about the form we need to render (e.g. username + password)
  const [flow, setFlow] = useState<SelfServiceRegistrationFlow>();

  // Get ?flow=... from the URL
  const { flow: flowId, return_to: returnTo } = router.query;

  const notifier = useNotifier();
  const resetFlow = () => {
    setFlow(undefined);
  };
  // In this effect we either initiate a new registration flow, or we fetch an existing registration flow.
  useEffect(() => {
    // If the router is not ready yet, or we already have a flow, do nothing.
    if (!router.isReady || flow) {
      return;
    }

    // If ?flow=.. was in the URL, we fetch it
    if (flowId) {
      ory
        .getSelfServiceRegistrationFlow(String(flowId))
        .then(({ data }) => {
          // We received the flow - let's use its data and render the form!
          setFlow(data);
        })
        .catch(handleFlowError(router, 'register', resetFlow, notifier));
      return;
    }

    // Otherwise we initialize it
    ory
      .initializeSelfServiceRegistrationFlowForBrowsers(returnTo ? String(returnTo) : undefined)
      .then(({ data }) => {
        setFlow(data);
      })
      .catch(handleFlowError(router, 'register', resetFlow, notifier));
  }, [flowId, router, router.isReady, returnTo, flow]);

  const validationSchema = Yup.object({
    password: Yup.string()
      .min(8, 'Password must be at least 8 characters')
      .required('Password is required'),
    traits: Yup.object({
      email: Yup.string().required('Email is required').email('Must be a valid email'),
      first_name: Yup.string().required('First name is required'),
      last_name: Yup.string().required('Last name is required')
    })
  });

  const onSubmit = (values: SubmitSelfServiceRegistrationFlowBody) =>
    router
      // On submission, add the flow ID to the URL but do not navigate. This prevents the user loosing
      // his data when she/he reloads the page.
      .push(`/register?flow=${flow?.id}`, undefined, { shallow: true })
      .then(() =>
        ory
          .submitSelfServiceRegistrationFlow(String(flow?.id), values)
          .then(({ data }) => {
            // If we ended up here, it means we are successfully signed up!
            // For now however we just want to redirect home!
            return router.push(flow?.return_to || '/').then(() => {});
          })
          .catch(handleFlowError(router, 'register', resetFlow, notifier))
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
    method: {
      variant: 'contained',
      sx: {
        mt: -12
      }
    },
    'traits.email': {
      sx: {
        mt: -5
      }
    },
    'traits.email_verified': {
      sx: { display: 'none' }
    },
    'traits.hd': {
      sx: { display: 'none' }
    },
    'traits.profile': {
      sx: { display: 'none' }
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
            />
            <Typography color="textSecondary" align="center" variant="body1" sx={{ pt: 3, pb: 2 }}>
              or sign up with an email address
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
              Already have an account?{' '}
              <NextLink href="/login">
                <Link
                  variant="subtitle2"
                  underline="hover"
                  sx={{
                    cursor: 'pointer'
                  }}
                >
                  Login
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

export default Registration;
