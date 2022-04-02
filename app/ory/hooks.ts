import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import ory from './sdk';

// Returns a function which will log the user out
export function useLogoutHandler() {
  const [logoutToken, setLogoutToken] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    ory
      .createSelfServiceLogoutFlowUrlForBrowsers()
      .then(({ data }) => {
        setLogoutToken(data.logout_token);
      })
      .catch((err: AxiosError) => {
        switch (err.response?.status) {
          case 401:
            // do nothing, the user is not logged in
            return;
        }

        // Something else happened!
        return Promise.reject(err);
      });
  }, [logoutToken]);

  return () => {
    if (logoutToken) {
      ory.submitSelfServiceLogoutFlow(logoutToken).then(() => router.push('/login'));
    }
  };
}
