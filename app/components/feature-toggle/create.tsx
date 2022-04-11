import { AxiosError } from 'axios';
import { useFormik } from 'formik';
import { useState } from 'react';
import * as Yup from 'yup';

import {
  Box,
  Button,
  CardContent,
  FormControl,
  FormControlLabel,
  Grid,
  Input,
  InputLabel,
  ListItemIcon,
  ListItemText,
  MenuItem,
  Select,
  Slider,
  Switch,
  TextField
} from '@mui/material';
import { styled } from '@mui/material/styles';

import { FeatureToggle } from '../../api';
import { FeatureToggleType } from '../../api/enums';
import { Dashboard } from '../../data/api';
import { useAppSelector } from '../../data/hooks';
import { SerializeError } from '../../features/utils';
import { FeatureToggleTypeName } from '../../utils/display';
import { useNotifier } from '../hooks';
import { FeatureToggleIcon } from './icon';

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
  // For some reason, the slider is sluggish with formik.
  const [percentage, setPercentage] = useState(0.0);
  const formik = useFormik({
    initialValues: {
      name: '',
      description: '',
      toggleType: FeatureToggleType.ON_OFF,
      enabled: true,
      on: false
    },
    validationSchema: Yup.object({
      name: Yup.string().required('name is required')
    }),
    onSubmit: async (values) => {
      try {
        if (!project?.id) {
          throw new Error('no project id');
        }

        const feature: FeatureToggle = {
          name: values.name,
          description: values.description,
          toggleType: values.toggleType,
          enabled: values.enabled,
          projectId: project?.id
        };
        switch (feature.toggleType) {
          case FeatureToggleType.PERCENTAGE:
            feature.percentage = {
              on: {
                weight: percentage
              },
              off: {
                weight: 100 - percentage
              }
            };
            break;
          case FeatureToggleType.ON_OFF:
            feature.onOff = {
              on: {}
            };
            break;
        }

        await Dashboard.createFeatureToggle(project?.id, {
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
    }
  });

  const renderToggleType = () => {
    switch (formik.values.toggleType) {
      // @ts-ignore
      case FeatureToggleType.PERCENTAGE:
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
        return (
          <FormControlLabel
            control={<Switch name="on" checked={formik.values.on} onChange={formik.handleChange} />}
            label="On"
          />
        );
    }
  };

  const handleSubmit = async () => {
    await formik.submitForm();
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
              error={Boolean(formik.touched.name && formik.errors.name)}
              helperText="Name used to check whether the feature toggle is enabled or not"
              label="Name"
              name="name"
              onChange={(e) =>
                formik.setFieldValue(
                  'name',
                  e.target.value.toUpperCase().replace(/[^a-zA-Z0-9_-]/gi, '')
                )
              }
              onBlur={formik.handleBlur}
              required
              value={formik.values.name}
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12} sm={5}>
            <TextField
              fullWidth
              multiline
              error={Boolean(formik.touched.description && formik.errors.description)}
              helperText="Description for what this feature toggle is used for"
              label="Description"
              name="description"
              onBlur={formik.handleBlur}
              onChange={formik.handleChange}
              value={formik.values.description}
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12} sm={2}>
            <FormControlLabel
              control={
                <Switch
                  name="enabled"
                  checked={formik.values.enabled}
                  onChange={formik.handleChange}
                />
              }
              label="Enabled"
            />
          </Grid>
        </Grid>
        <Grid container spacing={3} alignItems="center" sx={{ mb: 2 }}>
          <Grid item md={5} xs={12}>
            <FormControl>
              <InputLabel>Type</InputLabel>
              <ToggleTypeSelector
                input={
                  <Input
                    startAdornment={<FeatureToggleIcon toggleType={formik.values.toggleType} />}
                  ></Input>
                }
                value={formik.values.toggleType}
                label="Type"
                name="toggleType"
                onBlur={formik.handleBlur}
                onChange={formik.handleChange}
              >
                {Object.entries(FeatureToggleType).map((v) => (
                  <MenuItem key={v[0]} value={v[1]}>
                    <ListItemIcon>
                      <FeatureToggleIcon toggleType={v[1]} />
                    </ListItemIcon>
                    <ListItemText primary={FeatureToggleTypeName(v[1])} />
                  </MenuItem>
                ))}
              </ToggleTypeSelector>
            </FormControl>
          </Grid>
          <Grid item md={5} xs={12} sx={{ my: 2 }}>
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
