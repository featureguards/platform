import { useRouter } from 'next/router';

import { ListItem, ListItemButton, ListItemText } from '@mui/material';
import useMediaQuery from '@mui/material/useMediaQuery';
import { useTheme } from '@mui/system';

import { DynamicSetting } from '../../api';
import { DYNAMIC_SETTINGS } from '../../utils/constants';

export type DynamicSettingsItemProps = {
  dynamicSetting: DynamicSetting;
};

export const DynamicSettingItem = (props: DynamicSettingsItemProps) => {
  const router = useRouter();
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down('sm'));
  return (
    <ListItem disablePadding>
      <ListItemButton
        onClick={() => {
          router.push(`${DYNAMIC_SETTINGS}/${props.dynamicSetting.id}`);
        }}
      >
        <ListItemText
          sx={{ pl: 2, maxWidth: 300 }}
          primary={props.dynamicSetting.name!}
          primaryTypographyProps={{ noWrap: true }}
        />
        {!fullScreen && (
          <ListItemText sx={{ ml: 2 }} secondary={props.dynamicSetting.description} />
        )}
      </ListItemButton>
    </ListItem>
  );
};
