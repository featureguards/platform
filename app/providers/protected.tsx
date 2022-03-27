import { ReactNode, useEffect, useState } from 'react';

import SuspenseLoader from '../components/suspense-loader';
import { useAppDispatch, useAppSelector } from '../data/hooks';
import { fetchMe } from '../features/users/slice';
import { PUBLIC_PATHS } from '../utils/constants';

import type { Router } from 'next/router';

export type RouteGuardProps = {
  router: Router;
  children: ReactNode | ReactNode[];
};

export const RouteGuard = ({ router, children }: RouteGuardProps) => {
  const [authorized, setAuthorized] = useState(false);
  const me = useAppSelector((state) => state.users.me);
  const meStatus = useAppSelector((state) => state.users.status);
  const dispatch = useAppDispatch();
  useEffect(() => {
    if (meStatus === 'idle') {
      dispatch(fetchMe());
    } else if (meStatus === 'failed') {
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
    } else if (meStatus === 'succeeded') {
      setAuthorized(true);
    }
  }, [meStatus, router.pathname, dispatch, me, router]);

  if (meStatus === 'loading') {
    return <SuspenseLoader></SuspenseLoader>;
  }
  return <>{authorized && children}</>;
};
