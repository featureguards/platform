import { useFormik } from 'formik';
import * as Yup from 'yup';

import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Divider,
  Grid,
  TextField
} from '@mui/material';

import { Dashboard } from '../../data/api';
import { useNotifier } from '../hooks';

export type CloneEnvironmentProps = {
  id: string;
  onSubmit?: (_props: { err?: Error }) => Promise<void>;
};

export const CloneEnvironment = (props: CloneEnvironmentProps) => {
  const { onSubmit, ...otherProps } = props;
  const notifier = useNotifier();
  const formik = useFormik({
    initialValues: {
      name: '',
      description: ''
    },
    validationSchema: Yup.object({
      name: Yup.string().required('Environment name is required')
    }),
    onSubmit: async (values) => {
      try {
        await Dashboard.cloneEnvironment(props.id, {
          name: values.name,
          description: values.description
        });
        if (onSubmit) {
          await onSubmit({});
        }
      } catch (err) {
        notifier.error('Error cloning environment');
      }
    }
  });

  const handleSubmit = async () => {
    await formik.submitForm();
  };

  return (
    <Card {...otherProps}>
      <CardHeader title="Clone Environment" />
      <Divider />
      <CardContent>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <TextField
              fullWidth
              error={Boolean(formik.touched.name && formik.errors.name)}
              helperText="Please specify the name of the new environment"
              label="Name"
              name="name"
              onBlur={formik.handleBlur}
              onChange={formik.handleChange}
              required
              value={formik.values.name}
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              multiline
              error={Boolean(formik.touched.description && formik.errors.description)}
              helperText="Please specify new environment description"
              label="Description"
              name="description"
              onBlur={formik.handleBlur}
              onChange={formik.handleChange}
              value={formik.values.description}
              variant="outlined"
            />
          </Grid>
        </Grid>
      </CardContent>
      <Divider />
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'flex-end',
          p: 2
        }}
      >
        <Button
          sx={{ ml: 2 }}
          onClick={() => {
            if (onSubmit) onSubmit({ err: new Error('cancelled') });
          }}
        >
          Cancel
        </Button>
        <Button sx={{ ml: 2 }} color="primary" onClick={handleSubmit} variant="contained">
          Clone
        </Button>
      </Box>
    </Card>
  );
};
