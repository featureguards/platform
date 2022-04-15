import { useState } from 'react';

import Percent from '@mui/icons-material/Percent';
import {
  Box,
  FormControlLabel,
  Grid,
  Input,
  Slider,
  Switch,
  TextField,
  Typography
} from '@mui/material';
import { useTheme } from '@mui/system';

import { EnvironmentFeatureToggle, FeatureToggle } from '../../api';
import { FeatureToggleType } from '../../api/enums';

export type EnvFeatureToggleViewProps = EnvironmentFeatureToggle;
export const EnvFeatureToggleView = (props: EnvFeatureToggleViewProps) => {
  const theme = useTheme();
  const [featureToggle, setFeatureToggle] = useState<FeatureToggle | undefined>(
    props.featureToggle
  );
  const renderToggleType = () => {
    switch (featureToggle?.toggleType) {
      case FeatureToggleType.PERCENTAGE:
        const percDef = featureToggle?.percentage;
        const percentage = percDef?.on?.weight || 0;
        const setPercentage = (val: number) => {
          setFeatureToggle({
            ...featureToggle,
            percentage: { ...percDef, on: { ...percDef?.on, weight: val } }
          });
        };
        const handleBlur = () => {
          if (percentage < 0) {
            setPercentage(0);
          } else if (percentage > 100) {
            setPercentage(100);
          }
        };
        return (
          <Grid container spacing={2} alignItems="center">
            <Grid item xs>
              <Slider
                onChange={(_, val) => setPercentage(val as number)}
                value={percentage}
                step={1}
                name="percentage"
                valueLabelDisplay="auto"
              />
            </Grid>
            <Grid item>
              <Input
                startAdornment={<Percent />}
                value={percentage}
                size="small"
                onChange={(e) => setPercentage(e.target.value === '' ? 0 : Number(e.target.value))}
                onBlur={handleBlur}
                inputProps={{
                  step: 0.01,
                  min: 0,
                  max: 100,
                  type: 'number'
                }}
              />
            </Grid>
          </Grid>
        );
      case FeatureToggleType.ON_OFF:
        const onOffDef = featureToggle?.onOff;
        const on = !!onOffDef?.on?.weight;
        const setOnOff = (val: boolean) => {
          const onWeight = val ? 100 : 0;
          const offWeight = val ? 0 : 100;
          setFeatureToggle({
            ...featureToggle,
            onOff: {
              ...onOffDef,
              on: { ...onOffDef?.on, weight: onWeight },
              off: { ...onOffDef?.off, weight: offWeight }
            }
          });
        };

        return (
          <FormControlLabel
            control={<Switch name="on" checked={on} onChange={(e) => setOnOff(e.target.checked)} />}
            label="On"
          />
        );
    }
  };

  return (
    <Box
      sx={{
        pt: 5,
        pb: 5,
        backgroundColor: theme.palette.background.paper
      }}
    >
      <Grid container spacing={3} sx={{ pl: 5 }}>
        <Grid xs={12} item>
          <Typography sx={{ pb: 3 }} variant="h5">
            {featureToggle?.name}
          </Typography>
        </Grid>
        <Grid xs={12} item>
          <Grid container>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                multiline
                label="Description"
                name="description"
                onChange={(v) =>
                  setFeatureToggle({ ...featureToggle, description: v.target.value })
                }
                value={featureToggle?.description}
              />
            </Grid>
            {/* <Grid item sx={{ pl: 3, pt: 1 }} xs={12} sm={6}>
              <FormControlLabel
                control={
                  <Switch
                    name="enabled"
                    checked={featureToggle?.enabled}
                    onChange={(v) =>
                      setFeatureToggle({ ...featureToggle, enabled: v.target.checked })
                    }
                  />
                }
                label="Enabled"
              />
            </Grid> */}
          </Grid>
        </Grid>
        <Grid item md={5} xs={12} sx={{ my: 2 }}>
          {renderToggleType()}
        </Grid>
      </Grid>

      {/* {ft && <FeatureToggleItem featureToggle={ft}></FeatureToggleItem>} */}
    </Box>
  );
};
