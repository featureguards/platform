import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import DeleteIcon from '@mui/icons-material/Delete';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  Grid,
  IconButton,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  SelectChangeEvent,
  TextField,
  Typography
} from '@mui/material';
import { SerializedError } from '@reduxjs/toolkit';

import {
  BoolValue,
  DynamicSetting,
  EnvironmentDynamicSetting,
  FloatValue,
  IntegerValue,
  StringValue
} from '../../api';
import { DynamicSettingType, PlatformTypeType } from '../../api/enums';
import { Dashboard } from '../../data/api';
import { useAppDispatch, useAppSelector } from '../../data/hooks';
import { details } from '../../features/dynamic_settings/slice';
import { platformTypeName } from '../../utils/display';
import { handleError, useNotifier, validate } from '../hooks';
import { BoolVal } from './bool-val';
import { DangerZone } from './danger-zone';
import { NumberVal } from './number-val';
import { StringVal } from './string-val';
import { EnvDynamicSettingHistoryView } from './view-history';

export type EnvDynamicSettingViewProps = EnvironmentDynamicSetting & { history?: boolean };
export const EnvDynamicSettingView = (props: EnvDynamicSettingViewProps) => {
  const [dynamicSetting, setDynamicSetting] = useState<DynamicSetting | undefined>(props.setting);
  const [envsToUpdate, setEnvsToUpdate] = useState<string[]>([props.environmentId!]);
  const [showDelete, setShowDelete] = useState<boolean>(false);

  // useState only picks up the initial props.dynamicSetting. If the environment is updated, it won't
  // pick up the new props.dynamicSetting.
  // https://stackoverflow.com/questions/54865764/react-usestate-does-not-reload-state-from-props
  useEffect(() => {
    setDynamicSetting(props.setting);
    setEnvsToUpdate([props.environmentId!]);
  }, [props.setting, props.environmentId]);
  const [historyExpanded, setHistoryExpanded] = useState(false);
  const [updateDialogOpen, setUpdateDialogOpen] = useState(false);
  const projectDetails = useAppSelector((state) => state.projects.details);
  const currentProject = projectDetails?.item;
  const environments = new Map(currentProject?.environments?.map((env) => [env.id, env]));
  const otherEnvs =
    currentProject?.environments?.filter((env) => env.id !== props.environmentId) || [];
  const notifier = useNotifier();
  const dispatch = useAppDispatch();
  const router = useRouter();

  const handleChange = (event: SelectChangeEvent<typeof envsToUpdate>) => {
    const {
      target: { value }
    } = event;
    const values = typeof value === 'string' ? value.split(',') : value;
    setEnvsToUpdate([props.environmentId!, ...values.filter((id) => id !== props.environmentId)]);
  };

  const handleDialogClose = () => {
    setUpdateDialogOpen(false);
  };

  const handleUpdate = async () => {
    try {
      if (!dynamicSetting) {
        return;
      }
      validate(dynamicSetting);
      await Dashboard.updateDynamicSetting(dynamicSetting.id as string, {
        setting: dynamicSetting,
        environmentIds: envsToUpdate
      });
      setUpdateDialogOpen(false);
      // Refetch all envIDs and not just what was updated becauses it's used to show what
      // envIDs are there to render for the feature-toggle.
      const envIDs = currentProject?.environments?.map((env) => env.id) || [];
      await dispatch(
        details({ id: dynamicSetting.id as string, environmentIds: envIDs as string[] })
      ).unwrap();
    } catch (err) {
      handleError(router, notifier, err as SerializedError);
    }
  };

  const renderToggleType = () => {
    switch (dynamicSetting?.settingType) {
      case DynamicSettingType.BOOL:
        const setBoolValue = (val: BoolValue) => {
          setDynamicSetting({
            ...dynamicSetting,
            boolValue: val
          });
        };
        if (!dynamicSetting?.boolValue) {
          return <></>;
        }
        return <BoolVal val={dynamicSetting?.boolValue} setVal={setBoolValue} />;
      case DynamicSettingType.INTEGER:
        const setIntegerValue = (val: IntegerValue) => {
          setDynamicSetting({
            ...dynamicSetting,
            integerValue: val
          });
        };
        if (!dynamicSetting?.integerValue) {
          return <></>;
        }
        return <NumberVal val={dynamicSetting?.integerValue} setVal={setIntegerValue} />;
      case DynamicSettingType.FLOAT:
        const setFloatValue = (val: FloatValue) => {
          setDynamicSetting({
            ...dynamicSetting,
            floatValue: val
          });
        };
        if (!dynamicSetting?.floatValue) {
          return <></>;
        }
        return <NumberVal val={dynamicSetting?.floatValue} setVal={setFloatValue} />;
      case DynamicSettingType.STRING:
        const setStringValue = (val: StringValue) => {
          setDynamicSetting({
            ...dynamicSetting,
            stringValue: val
          });
        };
        if (!dynamicSetting?.stringValue) {
          return <></>;
        }
        return <StringVal val={dynamicSetting?.stringValue} setVal={setStringValue} />;
    }
  };

  return (
    <Card
      sx={{
        m: 2
      }}
    >
      <CardHeader
        title={dynamicSetting?.name}
        subheader={environments.get(props.environmentId)?.name}
        action={
          <IconButton aria-label="delete" onClick={() => setShowDelete(true)}>
            <DeleteIcon />
          </IconButton>
        }
      />
      <CardContent>
        <DangerZone
          id={dynamicSetting?.id}
          environmentId={props.environmentId}
          showDelete={showDelete}
          setShowDelete={setShowDelete}
        />
        <Grid container spacing={3}>
          <Grid xs={12} item>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={7}>
                <TextField
                  fullWidth
                  multiline
                  label="Description"
                  name="description"
                  onChange={(v) =>
                    setDynamicSetting({ ...dynamicSetting, description: v.target.value })
                  }
                  value={dynamicSetting?.description}
                />
              </Grid>
              <Grid item xs={12} sm={5}>
                <FormControl>
                  <InputLabel>Platforms</InputLabel>
                  <Select
                    disabled
                    label="Platforms"
                    name="platforms"
                    size="small"
                    value={dynamicSetting?.platforms || []}
                    input={<OutlinedInput />}
                    renderValue={(selected) => (
                      <Box sx={{ display: 'flex', flexWrap: 'wrap' }}>
                        {selected.map((v) => (
                          <Chip key={v} label={platformTypeName(v as PlatformTypeType)} />
                        ))}
                      </Box>
                    )}
                  >
                    {Object.values(PlatformTypeType).map((v) => (
                      <MenuItem key={v} value={v}>
                        {platformTypeName(v)}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              {/* <Grid item sx={{ pl: 3, pt: 1 }} xs={12} sm={6}>
              <FormControlLabel
                control={
                  <Switch
                    name="enabled"
                    checked={dynamicSetting?.enabled}
                    onChange={(v) =>
                      setDynamicSetting({ ...dynamicSetting, enabled: v.target.checked })
                    }
                  />
                }
                label="Enabled"
              />
            </Grid> */}
            </Grid>
          </Grid>
          <Grid item xs sx={{ my: 2 }}>
            {renderToggleType()}
          </Grid>
          <Grid item xs={12} sx={{ p: 1 }}>
            <Button variant="contained" onClick={() => setUpdateDialogOpen(true)}>
              Update
            </Button>
          </Grid>
        </Grid>
        <Dialog open={updateDialogOpen} onClose={handleDialogClose}>
          <DialogTitle>Confirm</DialogTitle>
          <DialogContent>
            <Typography>Which environments to update?</Typography>
            <Grid container>
              <Grid item xs={12}>
                {otherEnvs.length && (
                  <FormControl sx={{ m: 1, width: 300 }}>
                    <InputLabel>More</InputLabel>
                    <Select
                      multiple
                      value={envsToUpdate}
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
                      {currentProject?.environments?.map((env) => (
                        <MenuItem
                          key={env.id}
                          value={env.id}
                          disabled={env.id === props.environmentId}
                        >
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
                <EnvDynamicSettingHistoryView
                  environmentId={props.environmentId as string}
                  id={props.setting?.id as string}
                ></EnvDynamicSettingHistoryView>
              )}
            </AccordionDetails>
          </Accordion>
        )}
      </CardContent>
    </Card>
  );
};
