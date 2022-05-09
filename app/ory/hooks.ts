import { useRouter } from 'next/router';

import ory, { logoutUrl } from './sdk';

// Returns a function which will log the user out
export function useLogout() {
  const router = useRouter();
  return {
    logout: async () => {
      const url = await logoutUrl();
      ory.submitSelfServiceLogoutFlow(url.token).then(() => router.push('/login'));
    }
  };
}
