import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { Accordion, AccordionDetails, AccordionSummary, Card, Typography } from '@mui/material';

import { Environment } from '../../api';
import { useAppSelector } from '../../data/hooks';
import { useFeatureToggleDetails } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { EnvFeatureToggleView } from './view-environment';

export type FeatureToggleViewProps = {
  id: string;
  environmentId: string;
};

export const FeatureToggleView = (props: FeatureToggleViewProps) => {
  const projectDetails = useAppSelector((state) => state.projects.details);
  const currentProject = projectDetails?.item;
  const environments = new Map<string, Environment>(
    currentProject?.environments?.map((env) => {
      return [env.id as string, env];
    })
  );
  const { items, loading } = useFeatureToggleDetails({
    id: props.id,
    environmentIds: []
  });
  const featureToggle = items?.filter((ft) => ft.environmentId === props.environmentId)?.[0]
    ?.featureToggle;

  if (loading) {
    return <SuspenseLoader></SuspenseLoader>;
  }

  if (!items?.length || !featureToggle) {
    return <></>;
  }

  const others = items?.filter((ft) => ft.environmentId !== props.environmentId);

  return (
    <>
      <EnvFeatureToggleView
        environmentId={props.environmentId}
        featureToggle={featureToggle}
        history={true}
      ></EnvFeatureToggleView>
      {others.length && (
        <Typography sx={{ pt: 5, pl: 2, pb: 1 }} variant="h5">
          Other Environments
        </Typography>
      )}
      {others.map((envFT) => {
        return (
          <Card sx={{ mx: 2 }} key={envFT.environmentId}>
            <Accordion>
              <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography>{environments.get(envFT.environmentId as string)?.name}</Typography>
              </AccordionSummary>
              <AccordionDetails>
                <EnvFeatureToggleView
                  environmentId={envFT.environmentId}
                  featureToggle={envFT.featureToggle}
                ></EnvFeatureToggleView>
              </AccordionDetails>
            </Accordion>
          </Card>
        );
      })}
    </>
  );
};
