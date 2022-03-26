import { useState, useEffect, ReactNode } from 'react';
import { useAppSelector, useAppDispatch } from '../data/hooks';
import { fetchMe } from '../features/users/slice';
import SuspenseLoader from '../components/suspense-loader';
import type { Router } from 'next/router';

export type RouteGuardProps = {
  router: Router;
  children: ReactNode | ReactNode[];
};

export const RouteGuard = ({ router, children }: RouteGuardProps) => {
  const [authorized, setAuthorized] = useState(false);
  const me = useAppSelector((state) => state.users.me);
  const meStatus = useAppSelector((state) => state.users.status);
  const publicPaths = ['/login', '/register'];
  const dispatch = useAppDispatch();
  useEffect(() => {
    if (meStatus === 'idle') {
      dispatch(fetchMe());
    } else if (meStatus === 'failed') {
      if (publicPaths.includes(router.pathname)) {
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
  }, [meStatus, router.pathname]);

  if (meStatus === 'loading') {
    return <SuspenseLoader></SuspenseLoader>;
  }
  return <>{authorized && children}</>;
};
