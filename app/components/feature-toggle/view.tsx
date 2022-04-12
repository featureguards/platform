import { Box, Typography } from '@mui/material';
import { useTheme } from '@mui/system';

import { useFeatureToggleDetails } from '../hooks';
import SuspenseLoader from '../suspense-loader';
import { FeatureToggleItem } from './feature-toggle-item';

export type FeatureToggleViewProps = {
  id: string;
  environmentId: string;
};

export const FeatureToggleView = (props: FeatureToggleViewProps) => {
  const theme = useTheme();

  const { items, loading } = useFeatureToggleDetails({
    id: props.id,
    environmentIds: []
  });

  if (loading) {
    return <SuspenseLoader></SuspenseLoader>;
  }

  if (!items?.length) {
    return <></>;
  }

  const ft = items[0].featureToggle;

  return (
    <Box
      sx={{
        pt: 5,
        pb: 5,
        backgroundColor: theme.palette.background.paper
      }}
    >
      <Typography sx={{ pl: 5, pb: 4 }} variant="h5">
        Feature Toggle
      </Typography>
      {ft && <FeatureToggleItem featureToggle={ft}></FeatureToggleItem>}
    </Box>
  );
};
