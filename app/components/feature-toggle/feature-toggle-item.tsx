import { useRouter } from 'next/router';

import { ListItem, ListItemButton, ListItemIcon, ListItemText } from '@mui/material';
import useMediaQuery from '@mui/material/useMediaQuery';
import { useTheme } from '@mui/system';

import { FeatureToggle } from '../../api';
import { FEATURE_TOGGGLES } from '../../utils/constants';
import { LiveToggleIcon } from './icon';

export type FeatureTogglesItemProps = {
  featureToggle: FeatureToggle;
};

export const FeatureToggleItem = (props: FeatureTogglesItemProps) => {
  const router = useRouter();
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down('sm'));
  return (
    <ListItem disablePadding>
      <ListItemButton
        onClick={() => {
          router.push(`${FEATURE_TOGGGLES}/${props.featureToggle.id}`);
        }}
      >
        <ListItemIcon>
          <LiveToggleIcon featureToggle={props.featureToggle}></LiveToggleIcon>
        </ListItemIcon>
        <ListItemText
          sx={{ pl: 2, maxWidth: 300 }}
          primary={props.featureToggle.name!}
          primaryTypographyProps={{ noWrap: true }}
        />
        {!fullScreen && <ListItemText sx={{ ml: 2 }} secondary={props.featureToggle.description} />}
      </ListItemButton>
    </ListItem>
  );
};
