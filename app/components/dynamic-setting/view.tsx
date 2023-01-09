import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { Accordion, AccordionDetails, AccordionSummary, Card, Typography } from '@mui/material';

import { Environment } from '../../api';
import { useAppSelector } from '../../data/hooks';
import { useDynamicSettingDetails } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { EnvDynamicSettingView } from './view-environment';

export type DynamicSettingViewProps = {
  id: string;
  environmentId: string;
};

export const DynamicSettingView = (props: DynamicSettingViewProps) => {
  const projectDetails = useAppSelector((state) => state.projects.details);
  const currentProject = projectDetails?.item;
  const environments = new Map<string, Environment>(
    currentProject?.environments?.map((env) => {
      return [env.id as string, env];
    })
  );
  const { items, loading } = useDynamicSettingDetails({
    id: props.id,
    environmentIds: []
  });
  const dynamicSetting = items?.filter((ft) => ft.environmentId === props.environmentId)?.[0]
    ?.setting;

  if (loading) {
    return <SuspenseLoader></SuspenseLoader>;
  }

  if (!items?.length || !dynamicSetting) {
    return <></>;
  }

  const others = items?.filter((ft) => ft.environmentId !== props.environmentId);

  return (
    <>
      <EnvDynamicSettingView
        environmentId={props.environmentId}
        setting={dynamicSetting}
        history={true}
      ></EnvDynamicSettingView>
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
                <EnvDynamicSettingView
                  environmentId={envFT.environmentId}
                  setting={envFT.setting}
                ></EnvDynamicSettingView>
              </AccordionDetails>
            </Accordion>
          </Card>
        );
      })}
    </>
  );
};
