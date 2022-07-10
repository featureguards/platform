import { AxiosError } from 'axios';
import { useState } from 'react';

import {
  Box,
  Button,
  CardContent,
  Chip,
  FormControl,
  FormControlLabel,
  Grid,
  Input,
  InputLabel,
  ListItemIcon,
  ListItemText,
  MenuItem,
  Select,
  Switch,
  TextField
} from '@mui/material';
import { styled } from '@mui/material/styles';

import { FeatureToggle, OnOffFeature, PercentageFeature } from '../../api';
import { FeatureToggleType, PlatformTypeType, StickinessType } from '../../api/enums';
import { Dashboard } from '../../data/api';
import { useAppSelector } from '../../data/hooks';
import { SerializeError } from '../../features/utils';
import { featureToggleTypeName, platformTypeName } from '../../utils/display';
import { useNotifier } from '../hooks';
import { FeatureToggleIcon } from './icon';
import { Percentage } from './percentage';

const ToggleTypeSelector = styled(Select)(() => ({
  // Weird bug with Mui where it renders the icon on a separate line
  '.MuiListItemIcon-root': {
    display: 'none'
  }
}));

export type NewFeatureToggleProps = {
  onCreate?: () => Promise<void>;
};

export const NewFeatureToggle = (props: NewFeatureToggleProps) => {
  const { item: project } = useAppSelector((state) => state.projects.details);

  const notifier = useNotifier();
  const [percentage, setPercentage] = useState<PercentageFeature>({
    salt: undefined,
    on: {
      weight: 0
    },
    stickiness: {
      stickinessType: StickinessType.KEYS,
      keys: []
    }
  });
  const [onOff, setOnOff] = useState<OnOffFeature>({
    on: {
      weight: 0
    }
  });
  const [feature, setFeatureToggle] = useState<FeatureToggle>({
    name: '',
    description: '',
    toggleType: FeatureToggleType.ON_OFF,
    enabled: true
  });
  const handleSubmit = async () => {
    try {
      if (!project?.id) {
        throw new Error('no project id');
      }

      switch (feature.toggleType) {
        case FeatureToggleType.PERCENTAGE:
          percentage.off = {
            weight: 100 - (percentage.on?.weight || 0)
          };
          // must clear others because someone can switch toggle types.
          feature.percentage = percentage;
          feature.onOff = undefined;
          break;
        case FeatureToggleType.ON_OFF:
          onOff.off = {
            weight: onOff.on?.weight ? 0 : 1000
          };
          // must clear others because someone can switch toggle types.
          feature.onOff = onOff;
          feature.percentage = undefined;
          break;
      }

      await Dashboard.createFeatureToggle({
        projectId: project?.id,
        feature: feature
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

  const renderToggleType = () => {
    switch (feature.toggleType) {
      // @ts-ignore
      case FeatureToggleType.PERCENTAGE:
        return <Percentage percentage={percentage} setPercentage={setPercentage} />;

      case FeatureToggleType.ON_OFF:
        return (
          <FormControlLabel
            control={
              <Switch
                name="on"
                checked={!!onOff.on?.weight}
                onChange={() =>
                  setOnOff({
                    ...onOff,
                    on: {
                      ...onOff.on,
                      weight: !!onOff.on?.weight ? 0 : 100
                    }
                  })
                }
              />
            }
            label="On"
          />
        );
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
              helperText="Name used to check whether the feature toggle is enabled or not"
              label="Name"
              name="name"
              onChange={(e) =>
                setFeatureToggle({
                  ...feature,
                  name: e.target.value.toUpperCase().replace(/[^a-zA-Z0-9_-]/gi, '')
                })
              }
              required
              value={feature.name || ''}
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12} sm={5}>
            <TextField
              fullWidth
              multiline
              helperText="Description for what this feature toggle is used for"
              label="Description"
              name="description"
              onChange={(e) =>
                setFeatureToggle({
                  ...feature,
                  description: e.target.value
                })
              }
              value={feature.description || ''}
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12} sm={2}>
            <FormControlLabel
              control={
                <Switch
                  name="enabled"
                  checked={!!feature.enabled}
                  onChange={(_, checked) => setFeatureToggle({ ...feature, enabled: checked })}
                />
              }
              label="Enabled"
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <FormControl>
              <InputLabel>Platforms</InputLabel>
              <Select
                multiple
                label="Platforms"
                name="platforms"
                size="small"
                onChange={(e) =>
                  setFeatureToggle({
                    ...feature,
                    platforms: e.target.value as number[]
                  })
                }
                value={feature.platforms || [PlatformTypeType.DEFAULT]}
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
                  .filter((v) => [PlatformTypeType.DEFAULT, PlatformTypeType.WEB].includes(v))
                  .map((v) => (
                    <MenuItem key={v} value={v}>
                      {platformTypeName(v)}
                    </MenuItem>
                  ))}
              </Select>
            </FormControl>
          </Grid>{' '}
          <Grid item md={5} xs={12}>
            <FormControl>
              <InputLabel>Type</InputLabel>
              <ToggleTypeSelector
                input={
                  <Input
                    startAdornment={
                      <FeatureToggleIcon toggleType={feature.toggleType as FeatureToggleType} />
                    }
                  ></Input>
                }
                value={feature.toggleType}
                label="Type"
                name="toggleType"
                onChange={(e) =>
                  setFeatureToggle({ ...feature, toggleType: e.target.value as number })
                }
              >
                {Object.entries(FeatureToggleType).map((v) => (
                  <MenuItem key={v[0]} value={v[1]}>
                    <ListItemIcon>
                      <FeatureToggleIcon toggleType={v[1]} />
                    </ListItemIcon>
                    <ListItemText primary={featureToggleTypeName(v[1])} />
                  </MenuItem>
                ))}
              </ToggleTypeSelector>
            </FormControl>
          </Grid>
        </Grid>
        <Grid container spacing={3} alignItems="center" sx={{ mb: 2 }}>
          <Grid item xs={12} sx={{ my: 2 }}>
            {renderToggleType()}
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
        <Button sx={{ ml: 2 }} color="primary" onClick={handleSubmit} variant="contained">
          Create
        </Button>
      </Box>
    </Box>
  );
};
