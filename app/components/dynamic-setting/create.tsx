import { AxiosError } from 'axios';
import { useState } from 'react';

import {
  Box,
  Button,
  CardContent,
  Chip,
  FormControl,
  Grid,
  Input,
  InputLabel,
  ListItemText,
  MenuItem,
  Select,
  TextField
} from '@mui/material';
import { styled } from '@mui/material/styles';

import { BoolValue, DynamicSetting, FloatValue, IntegerValue, StringValue } from '../../api';
import { DynamicSettingType, PlatformTypeType } from '../../api/enums';
import { Dashboard } from '../../data/api';
import { useAppSelector } from '../../data/hooks';
import { SerializeError } from '../../features/utils';
import { track } from '../../utils/analytics';
import { dynamicSettingTypeName, platformTypeName } from '../../utils/display';
import { useNotifier } from '../hooks';
import { BoolVal } from './bool-val';
import { NumberVal } from './number-val';
import { StringVal } from './string-val';

const ToggleTypeSelector = styled(Select)(() => ({
  // Weird bug with Mui where it renders the icon on a separate line
  '.MuiListItemIcon-root': {
    display: 'none'
  }
}));

export type NewDynamicSettingProps = {
  onCreate?: () => Promise<void>;
  onCancel: () => void;
};

export const NewDynamicSetting = (props: NewDynamicSettingProps) => {
  const { item: project } = useAppSelector((state) => state.projects.details);

  const notifier = useNotifier();
  const [boolValue, setBoolValue] = useState<BoolValue>({
    value: false
  });
  const [intValue, setIntValue] = useState<IntegerValue>({
    value: 0
  });
  const [floatValue, setFloatValue] = useState<FloatValue>({
    value: 0
  });
  const [stringValue, setStringValue] = useState<StringValue>({
    value: ''
  });
  const [setting, setDynamicSetting] = useState<DynamicSetting>({
    name: '',
    description: '',
    settingType: DynamicSettingType.BOOL,
    platforms: [PlatformTypeType.DEFAULT]
  });
  const clear = () => {
    setting.boolValue = undefined;
    setting.integerValue = undefined;
    setting.floatValue = undefined;
    setting.stringValue = undefined;
    setting.listValues = undefined;
    setting.setValues = undefined;
    setting.mapValues = undefined;
    setting.jsonValue = undefined;
  };
  const handleSubmit = async () => {
    try {
      track('newDynamicSetting', {
        type: setting.settingType
      });

      if (!project?.id) {
        throw new Error('no project id');
      }
      clear();
      switch (setting.settingType) {
        case DynamicSettingType.BOOL:
          setting.boolValue = boolValue;
          break;
        case DynamicSettingType.INTEGER:
          setting.integerValue = intValue;
          break;
        case DynamicSettingType.FLOAT:
          setting.floatValue = floatValue;
          break;
        case DynamicSettingType.STRING:
          setting.stringValue = stringValue;
          break;
      }

      await Dashboard.createDynamicSetting({
        projectId: project?.id,
        setting: setting
      });
      notifier.success('Success');
      if (props.onCreate) {
        await props.onCreate();
      }
    } catch (err) {
      const parsed = SerializeError(err as AxiosError);
      if (parsed.message && Number(parsed.code) < 500) {
        notifier.error(parsed.message);
      }
    }
  };

  const renderSettingType = () => {
    switch (setting.settingType) {
      case DynamicSettingType.BOOL:
        return <BoolVal val={boolValue} setVal={setBoolValue} />;
      case DynamicSettingType.INTEGER:
        return <NumberVal val={intValue} setVal={setIntValue} />;
      case DynamicSettingType.FLOAT:
        return <NumberVal val={floatValue} setVal={setFloatValue} />;
      case DynamicSettingType.STRING:
        return <StringVal val={stringValue} setVal={setStringValue} />;
    }
  };

  if (!project) {
    return <></>;
  }

  return (
    <Box>
      <CardContent>
        <Grid container spacing={3} alignItems="top" sx={{ mb: 2 }}>
          <Grid item xs={12} sm={5}>
            <TextField
              fullWidth
              helperText="Name for the dynamic setting"
              label="Name"
              name="name"
              onChange={(e) =>
                setDynamicSetting({
                  ...setting,
                  name: e.target.value.toUpperCase().replace(/[^a-zA-Z0-9_-]/gi, '')
                })
              }
              required
              value={setting.name || ''}
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12} sm={5}>
            <TextField
              fullWidth
              multiline
              helperText="Description for what this dynamic setting is used for"
              label="Description"
              name="description"
              onChange={(e) =>
                setDynamicSetting({
                  ...setting,
                  description: e.target.value
                })
              }
              value={setting.description || ''}
              variant="outlined"
            />
          </Grid>
          <Grid item sm={5} xs={12}>
            <FormControl>
              <InputLabel>Type</InputLabel>
              <ToggleTypeSelector
                value={setting.settingType}
                label="Type"
                name="settingType"
                onChange={(e) =>
                  setDynamicSetting({ ...setting, settingType: e.target.value as number })
                }
              >
                {Object.entries(DynamicSettingType)
                  // TODO: support other types in the UI.
                  .filter((v) => {
                    if (
                      v[1] === DynamicSettingType.BOOL ||
                      v[1] === DynamicSettingType.INTEGER ||
                      v[1] === DynamicSettingType.FLOAT ||
                      v[1] === DynamicSettingType.STRING
                    ) {
                      return true;
                    }
                    return false;
                  })
                  .map((v) => (
                    <MenuItem key={v[0]} value={v[1]}>
                      <ListItemText primary={dynamicSettingTypeName(v[1])} />
                    </MenuItem>
                  ))}
              </ToggleTypeSelector>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={5}>
            <FormControl>
              <InputLabel>Platforms</InputLabel>
              <Select
                multiple
                label="Platforms"
                name="platforms"
                size="small"
                onChange={(e) =>
                  setDynamicSetting({
                    ...setting,
                    platforms: e.target.value as number[]
                  })
                }
                value={setting.platforms}
                input={<Input />}
                renderValue={(selected) => (
                  <Box sx={{ display: 'flex', flexWrap: 'wrap' }}>
                    {selected.map((v) => (
                      <Chip key={v} label={platformTypeName(v as PlatformTypeType)} />
                    ))}
                  </Box>
                )}
              >
                {Object.values(PlatformTypeType)
                  .filter((v) => v === PlatformTypeType.DEFAULT || v === PlatformTypeType.WEB)
                  .map((v) => (
                    <MenuItem key={v} value={v}>
                      {platformTypeName(v)}
                    </MenuItem>
                  ))}
              </Select>
            </FormControl>
          </Grid>
        </Grid>
        <Grid container spacing={3} alignItems="center" sx={{ mb: 2 }}>
          <Grid item xs={12} sx={{ my: 2 }}>
            {renderSettingType()}
          </Grid>
        </Grid>
      </CardContent>
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'flex-end',
          p: 2
        }}
      >
        <Button sx={{ ml: 2 }} onClick={props.onCancel}>
          Cancel
        </Button>
        <Button sx={{ ml: 2 }} color="primary" onClick={handleSubmit} variant="contained">
          Create
        </Button>
      </Box>
    </Box>
  );
};
