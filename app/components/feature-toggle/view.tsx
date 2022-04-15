import { useEffect, useState } from 'react';

import { FeatureToggle } from '../../api';
import { useFeatureToggleDetails } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { EnvFeatureToggleView } from './view-environment';

export type FeatureToggleViewProps = {
  id: string;
  environmentId: string;
};

export const FeatureToggleView = (props: FeatureToggleViewProps) => {
  const { items, loading } = useFeatureToggleDetails({
    id: props.id,
    environmentIds: []
  });
  const [featureToggle, setFeatureToggle] = useState<FeatureToggle | undefined>();
  //   const [others, setOthers] = useState<EnvironmentFeatureToggle[]>();

  const ft = items?.filter((ft) => ft.environmentId === props.environmentId)?.[0]?.featureToggle;
  useEffect(() => {
    setFeatureToggle(ft);
  }, [ft]);
  if (loading) {
    return <SuspenseLoader></SuspenseLoader>;
  }

  if (!items?.length || !featureToggle) {
    return <></>;
  }

  return (
    <EnvFeatureToggleView
      environmentId={props.environmentId}
      featureToggle={featureToggle}
    ></EnvFeatureToggleView>
  );
};
