import { useFormik } from 'formik';
import { useState } from 'react';
import * as Yup from 'yup';

import AddIcon from '@mui/icons-material/Add';
import {
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
  Divider,
  Grid,
  IconButton,
  TextField,
  Typography
} from '@mui/material';

import { Project } from '../../api';
import { Dashboard } from '../../data/api';
import { track } from '../../utils/analytics';
import { useNotifier } from '../hooks';

export type NewProjectProps = {
  onSubmit?: (_props: { err?: Error }) => Promise<void>;
};

export const NewProject = (props: NewProjectProps) => {
  const { onSubmit, ...otherProps } = props;
  const notifier = useNotifier();
  const [showAdd, setShowAdd] = useState<boolean>(false);
  const [newEnvName, setNewEnvName] = useState<string>('');
  const [state, setState] = useState<{ environments: string[]; project: Project | null }>({
    environments: ['Production', 'QA', 'Development'],
    project: null
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
      track('newProject', { envs: state.environments?.length });
      try {
        const res = await Dashboard.createProject({
          name: values.projectName,
          description: values.projectDescription,
          environments: environments
        });
        setState({ ...state, project: res.data });
        if (onSubmit) {
          await onSubmit({});
        }
      } catch (err) {
        notifier.error('Error creating a new project');
        if (onSubmit) {
          await onSubmit({ err: err as Error });
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

  const handleSubmit = async () => {
    await formik.submitForm();
  };

  const handleAddEnv = () => {
    setState({
      ...state,
      environments: [...state.environments, newEnvName]
    });
    setShowAdd(false);
  };

  return (
    <Card {...otherProps}>
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
              multiline
              error={Boolean(formik.touched.projectDescription && formik.errors.projectDescription)}
              helperText="Please specify the project description"
              label="Description"
              name="projectDescription"
              onBlur={formik.handleBlur}
              onChange={formik.handleChange}
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
                  key={index}
                  label={e}
                  onDelete={() => {
                    handleEnvDelete(e);
                  }}
                  variant="outlined"
                  color={index % 2 == 0 ? 'primary' : 'secondary'}
                />
              );
            })}
            <IconButton
              color="primary"
              onClick={() => {
                setNewEnvName('');
                setShowAdd(true);
              }}
            >
              <AddIcon></AddIcon>
            </IconButton>
            <Dialog open={showAdd} onClose={() => setShowAdd(false)}>
              <DialogTitle>Environment Name</DialogTitle>
              <DialogContent>
                <TextField
                  size="small"
                  helperText="Name of the new environment."
                  name="environment"
                  onChange={(e) => setNewEnvName(e.target.value)}
                  value={newEnvName}
                />
              </DialogContent>
              <DialogActions>
                <Button onClick={() => setShowAdd(false)}>Cancel</Button>
                <Button
                  disabled={newEnvName === ''}
                  color="primary"
                  variant="contained"
                  onClick={() => handleAddEnv()}
                  autoFocus
                >
                  Add
                </Button>
              </DialogActions>
            </Dialog>
          </Box>
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
        <Button sx={{ ml: 2 }} color="primary" onClick={handleSubmit} variant="contained">
          Create
        </Button>
      </Box>
    </Card>
  );
};
