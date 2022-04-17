import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Percent from '@mui/icons-material/Percent';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  FormControlLabel,
  Grid,
  Input,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  SelectChangeEvent,
  Slider,
  Switch,
  TextField,
  Typography
} from '@mui/material';
import { useTheme } from '@mui/system';
import { SerializedError } from '@reduxjs/toolkit';

import { EnvironmentFeatureToggle, FeatureToggle } from '../../api';
import { FeatureToggleType } from '../../api/enums';
import { Dashboard } from '../../data/api';
import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { details } from '../../features/feature_toggles/slice';
import { handleError, useNotifier } from '../hooks';
import { EnvFeatureToggleHistoryView } from './view-history';

export type EnvFeatureToggleViewProps = EnvironmentFeatureToggle & { history?: boolean };
export const EnvFeatureToggleView = (props: EnvFeatureToggleViewProps) => {
  const theme = useTheme();
  const [featureToggle, setFeatureToggle] = useState<FeatureToggle | undefined>(
    props.featureToggle
  );
  // useState only picks up the initial props.featureToggle. If the environment is updated, it won't
  // pick up the new props.featureToggle.
  // https://stackoverflow.com/questions/54865764/react-usestate-does-not-reload-state-from-props
  useEffect(() => {
    setFeatureToggle(props.featureToggle);
  }, [props.featureToggle]);
  const [historyExpanded, setHistoryExpanded] = useState(false);
  const [updateDialogOpen, setUpdateDialogOpen] = useState(false);
  const projectDetails = useAppSelector((state) => state.projects.details);
  const currentProject = projectDetails?.item;
  const environments = new Map(currentProject?.environments?.map((env) => [env.id, env]));
  const otherEnvs =
    currentProject?.environments?.filter((env) => env.id !== props.environmentId) || [];
  const [additionalEnvsToUpdate, setAdditionalEnvsToUpdate] = useState<string[]>([]);
  const notifier = useNotifier();
  const dispatch = useAppDispatch();
  const router = useRouter();

  const handleChange = (event: SelectChangeEvent<typeof additionalEnvsToUpdate>) => {
    const {
      target: { value }
    } = event;
    setAdditionalEnvsToUpdate(
      // On autofill we get a stringified value.
      typeof value === 'string' ? value.split(',') : value
    );
  };

  const handleDialogClose = () => {
    setUpdateDialogOpen(false);
  };

  const handleUpdate = async () => {
    try {
      const selectedEnvIDs = [props.environmentId as string, ...additionalEnvsToUpdate];
      await Dashboard.updateFeatureToggle(featureToggle?.id as string, {
        feature: featureToggle,
        environmentIds: selectedEnvIDs
      });
      setUpdateDialogOpen(false);
      // Refetch all envIDs and not just what was updated becauses it's used to show what
      // envIDs are there to render for the feature-toggle.
      const envIDs = currentProject?.environments?.map((env) => env.id) || [];
      await dispatch(
        details({ id: featureToggle?.id as string, environmentIds: envIDs as string[] })
      ).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

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
        backgroundColor: theme.palette.background.paper
      }}
    >
      <Grid container spacing={3} sx={{ pl: 5 }}>
        <Grid xs={12} item>
          <Typography sx={{ pb: 3 }} variant="h5">
            {featureToggle?.name} ({environments.get(props.environmentId)?.name})
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
        <Grid item xs={12} sm={2} sx={{ p: 1 }}>
          <Button variant="contained" onClick={() => setUpdateDialogOpen(true)}>
            Update
          </Button>
        </Grid>
      </Grid>
      <Dialog open={updateDialogOpen} onClose={handleDialogClose}>
        <DialogTitle>Confirm</DialogTitle>
        <DialogContent>
          <Typography>Which environments to update?</Typography>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Chip label={environments.get(props.environmentId)?.name}></Chip>
            </Grid>
            <Grid item xs={12}>
              {otherEnvs.length && (
                <FormControl sx={{ m: 1, width: 300 }}>
                  <InputLabel>More</InputLabel>
                  <Select
                    multiple
                    value={additionalEnvsToUpdate}
                    onChange={handleChange}
                    input={<OutlinedInput label="Chip" />}
                    renderValue={(selected) => (
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {selected.map((id) => (
                          <Chip key={id} label={environments.get(id)?.name} />
                        ))}
                      </Box>
                    )}
                  >
                    {otherEnvs.map((env) => (
                      <MenuItem key={env.id} value={env.id}>
                        {env.name}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDialogClose}>Cancel</Button>
          <Button onClick={handleUpdate} autoFocus>
            Confirm
          </Button>
        </DialogActions>
      </Dialog>
      {props.history && (
        <Accordion
          onChange={(_event: React.SyntheticEvent, expanded: boolean) => {
            setHistoryExpanded(expanded);
          }}
        >
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography>History</Typography>
          </AccordionSummary>
          <AccordionDetails>
            {historyExpanded && (
              <EnvFeatureToggleHistoryView
                environmentId={props.environmentId as string}
                id={props.featureToggle?.id as string}
              ></EnvFeatureToggleHistoryView>
            )}
          </AccordionDetails>
        </Accordion>
      )}
    </Box>
  );
};
