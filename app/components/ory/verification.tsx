import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { ReactNode, useEffect, useState } from 'react';

import {
  SelfServiceVerificationFlow,
  SubmitSelfServiceVerificationFlowBody
} from '@ory/kratos-client';

import { useAppSelector } from '../../data/hooks';
import ory, { urlForFlow } from '../../ory/sdk';
import { track } from '../../utils/analytics';
import SuspenseLoader from '../suspense-loader';
import { Flow } from './Flow';

export type VerificationProps = {
  children?: ReactNode;
  color?: 'inherit' | 'primary' | 'secondary' | 'success' | 'error' | 'info' | 'warning';
  variant?: 'text' | 'outlined' | 'contained';
};
const Verification = (props: VerificationProps) => {
  const [flow, setFlow] = useState<SelfServiceVerificationFlow>();
  const me = useAppSelector((state) => state.users.me);
  // Get ?flow=... from the URL
  const router = useRouter();
  const { return_to: returnTo, flow: flowId } = router.query;

  useEffect(() => {
    // If the router is not ready yet, or we already have a flow, do nothing.
    if (!router.isReady || flow) {
      return;
    }

    // If ?flow=.. was in the URL, we fetch it
    if (flowId) {
      ory
        .getSelfServiceVerificationFlow(String(flowId))
        .then(({ data }) => {
          setFlow(data);
        })
        .catch((err: AxiosError) => {
          switch (err.response?.status) {
            case 410:
            // Status code 410 means the request has expired - so let's load a fresh flow!
            case 403:
              // Status code 403 implies some other issue (e.g. CSRF) - let's reload!
              return router.push(urlForFlow('verification'));
          }
          throw err;
        });
      return;
    }

    // Otherwise we initialize it
    ory
      .initializeSelfServiceVerificationFlowForBrowsers(returnTo ? String(returnTo) : undefined)
      .then(({ data }) => {
        setFlow(data);
      })
      .catch((err: AxiosError) => {
        switch (err.response?.status) {
          case 400:
            // Status code 400 implies the user is already signed in
            return router.push('/');
        }

        throw err;
      });
  }, [flowId, router, router.isReady, returnTo, flow]);

  const onSubmit = (values: SubmitSelfServiceVerificationFlowBody) => {
    values.email = me?.addresses?.[0]?.address || '';
    track('verfification', {
      method: values.method
    });

    return ory
      .submitSelfServiceVerificationFlow(String(flow?.id), values)
      .then(({ data }) => {
        setFlow(data);
      })
      .catch((err: AxiosError) => {
        // If the previous handler did not catch the error it's most likely a form validation error
        if (err.response?.status === 400) {
          // Yup, it is!
          setFlow(err.response?.data);
          return;
        }

        throw err;
      });
  };

  if (flow?.state === 'passed_challenge') {
    window.location.href = flow?.return_to || '/';
  }
  return (
    <>
      {flow ? (
        <Flow
          flow={flow}
          spacing={-3}
          notify={true}
          onSubmit={onSubmit}
          hideGlobalMessages={true}
          nodeProps={{
            email: {
              sx: { display: 'none' }
            },
            method: {
              ...props
            }
          }}
          childrenNodes={{
            method: props.children || 'Resend confirmation'
          }}
        />
      ) : (
        <SuspenseLoader></SuspenseLoader>
      )}
    </>
  );
};

export default Verification;
