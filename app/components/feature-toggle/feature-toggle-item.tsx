import { useRouter } from 'next/router';

import { ListItem, ListItemButton, ListItemIcon, ListItemText } from '@mui/material';

import { FeatureToggle } from '../../api';
import { FEATURE_TOGGGLES } from '../../utils/constants';
import { LiveToggleIcon } from './icon';

export type FeatureTogglesItemProps = {
  featureToggle: FeatureToggle;
};

export const FeatureToggleItem = (props: FeatureTogglesItemProps) => {
  const router = useRouter();
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
        <ListItemText sx={{ pl: 2 }} primary={props.featureToggle.name!} />
        <ListItemText secondary={props.featureToggle.description} />
      </ListItemButton>
    </ListItem>
  );
};
