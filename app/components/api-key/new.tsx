import { useFormik } from 'formik';
import { DateTime } from 'luxon';
import * as Yup from 'yup';

import { Box, Button, Card, CardContent, Divider, TextField } from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers';

import { ApiKey } from '../../api';
import { Dashboard } from '../../data/api';
import { useNotifier } from '../hooks';

export type NewApiKeyProps = {
  environmentId: string;
  onSubmit?: (_props: { err?: Error }) => Promise<void>;
};

export const NewApiKey = (props: NewApiKeyProps) => {
  const { onSubmit, environmentId, ...otherProps } = props;
  const notifier = useNotifier();
  const formik = useFormik<ApiKey>({
    initialValues: {
      name: '',
      environmentId: environmentId
    },
    validationSchema: Yup.object({
      name: Yup.string().required('Name is required')
    }),
    onSubmit: async (values) => {
      try {
        await Dashboard.createApiKey({ ...values, environmentId: environmentId });
        if (onSubmit) {
          onSubmit({});
        }
      } catch (err) {
        notifier.error('Error creating a new api key');
        if (onSubmit) {
          onSubmit({ err: err as Error });
        }
      }
    }
  });

  return (
    <Card {...otherProps}>
      <CardContent>
        <Box
          sx={{
            pt: 1,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            flexDirection: 'row'
          }}
        >
          <TextField
            label="Name"
            name="name"
            size="small"
            error={Boolean(formik.touched.name && formik.errors.name)}
            sx={{ mr: 2 }}
            onBlur={formik.handleBlur}
            onChange={formik.handleChange}
            value={formik.values.name}
            variant="outlined"
          />

          <DatePicker
            renderInput={(props) => <TextField size="small" {...props}></TextField>}
            label="Expires at"
            onChange={(dateTime: DateTime | null) => {
              formik.setFieldValue('expiresAt', dateTime?.toISO());
            }}
            value={formik.values.expiresAt ? DateTime.fromISO(formik.values.expiresAt) : null}
          ></DatePicker>
        </Box>
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
          color="primary"
          onClick={() => formik.handleSubmit()}
          variant="contained"
        >
          Create
        </Button>
      </Box>
    </Card>
  );
};
