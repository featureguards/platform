import { useFormik } from 'formik';
import { useSnackbar } from 'notistack';
import { FC, useState } from 'react';
import * as Yup from 'yup';

import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Chip,
  Divider,
  Grid,
  TextField,
  Typography
} from '@mui/material';

import { Project } from '../../api';
import { Dashboard } from '../../data/api';
import { Notif } from '../../utils/notif';

export type NewProjectProps = {
  onNewProject?: (_props: { project?: Project; err?: Error }) => Promise<void>;
};

export const NewProject: FC<NewProjectProps> = (props) => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const notifier = new Notif({ enqueueSnackbar: enqueueSnackbar, closeSnackbar: closeSnackbar });

  const [state, setState] = useState({
    environments: ['Production', 'QA', 'Development']
  });

  const formik = useFormik({
    initialValues: {
      projectName: '',
      projectDescription: ''
    },
    validationSchema: Yup.object({
      projectName: Yup.string().required('Project name is required')
    }),
    onSubmit: async (values) => {
      const environments = state.environments.map((env) => {
        return { name: env };
      });
      try {
        const res = await Dashboard.createProject({
          name: values.projectName,
          description: values.projectDescription,
          environments: environments
        });
        if (props.onNewProject) {
          await props.onNewProject({ project: res.data });
        }
      } catch (err) {
        notifier.error('Error creating a new project');
        if (props.onNewProject) {
          await props.onNewProject({ err: err as Error });
        }
      }
    }
  });
  const handleEnvDelete = (e: string) => {
    setState({
      ...state,
      environments: state.environments.filter((env) => env !== e)
    });
  };
  return (
    <Card {...props}>
      <CardHeader subheader="Let's create a new project" title="New Project" />
      <Divider />
      <CardContent>
        <Grid container spacing={3}>
          <Grid item md={6} xs={12}>
            <TextField
              fullWidth
              error={Boolean(formik.touched.projectName && formik.errors.projectName)}
              helperText="Please specify the name of the project"
              label="Project name"
              name="projectName"
              onBlur={formik.handleBlur}
              onChange={formik.handleChange}
              required
              value={formik.values.projectName}
              variant="outlined"
            />
          </Grid>
          <Grid item md={6} xs={12}>
            <TextField
              fullWidth
              error={Boolean(formik.touched.projectDescription && formik.errors.projectDescription)}
              helperText="Please the project description"
              label="Description"
              name="projectDescription"
              onBlur={formik.handleBlur}
              onChange={formik.handleChange}
              required
              value={formik.values.projectDescription}
              variant="outlined"
            />
          </Grid>
        </Grid>
        <Divider sx={{ py: 2 }} />
        <Grid item md={6} xs={12}>
          <Typography variant="h6" gutterBottom>
            What environments do you want to create?
          </Typography>

          <Box sx={{ flexDirection: 'row' }}>
            {state.environments.map((e, index) => {
              return (
                <Chip
                  sx={{ mx: 1 }}
                  key={e}
                  label={e}
                  onDelete={() => {
                    handleEnvDelete(e);
                  }}
                  variant="outlined"
                  color={index % 2 == 0 ? 'primary' : 'secondary'}
                />
              );
            })}
          </Box>
        </Grid>
      </CardContent>
    </Card>
  );
};
