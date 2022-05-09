import { ReactNode, useEffect, useState } from 'react';

import { useMe } from '../components/hooks';
import SuspenseLoader from '../components/suspense-loader';
import { PUBLIC_PATHS } from '../utils/constants';

import type { Router } from 'next/router';
export type RouteGuardProps = {
  router: Router;
  children: ReactNode | ReactNode[];
};

export const RouteGuard = ({ router, children }: RouteGuardProps) => {
  const [authorized, setAuthorized] = useState(false);
  const { me, status } = useMe();
  useEffect(() => {
    if (status === 'failed') {
      if (PUBLIC_PATHS.includes(router.pathname)) {
        setAuthorized(true);
      } else if (!me) {
        setAuthorized(false);
        const query =
          router.asPath.length && router.asPath !== '/' && router.basePath !== '/login'
            ? { returnUrl: router.asPath }
            : undefined;
        router.push({
          pathname: '/login',
          query: query
        });
      } else {
        setAuthorized(true);
      }
    } else if (status === 'succeeded') {
      setAuthorized(true);
    }
  }, [status, router.pathname, router]);

  if (status === 'loading') {
    return <SuspenseLoader></SuspenseLoader>;
  }
  return <>{authorized && children}</>;
};
