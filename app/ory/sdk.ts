import { AxiosError } from 'axios';

import { edgeConfig } from '@ory/integrations/next';
import { Configuration, V0alpha2Api } from '@ory/kratos-client';

// Initialize the Ory Kratos SDK which will connect to the
edgeConfig.basePath = '/identity';

const ory = new V0alpha2Api(new Configuration(edgeConfig));
export default ory;

export const urlForFlow = (
  flowType: 'login' | 'register' | 'settings' | 'recovery' | 'verification'
) => {
  switch (flowType) {
    case 'login':
    case 'register':
      return '/' + flowType;
    case 'settings':
      return '/account/settings';
    case 'recovery':
      return '/account/reset';
    case 'verification':
      return '/account/verify';
  }
};

export type Logout = {
  url: string;
  token: string;
};

export const logoutUrl = (): Promise<Logout> => {
  return ory
    .createSelfServiceLogoutFlowUrlForBrowsers()
    .then(({ data }) => {
      return Promise.resolve({ token: data.logout_token, url: data.logout_url });
    })
    .catch((err: AxiosError) => {
      switch (err.response?.status) {
        case 401:
          // do nothing, the user is not logged in
          return Promise.resolve({ url: '', token: '' });
      }

      // Something else happened!
      return Promise.reject(err);
    });
};

export const logout = () => {
  return logoutUrl().then((url: Logout) => {
    if (url.token.length) {
      ory.submitSelfServiceLogoutFlow(url.token);
    }
    return Promise.resolve();
  });
};
