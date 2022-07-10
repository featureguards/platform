import { useFormik } from 'formik';
import { DateTime } from 'luxon';
import * as Yup from 'yup';

import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  FormControl,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  TextField
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers';

import { ApiKey } from '../../api';
import { PlatformTypeType } from '../../api/enums';
import { Dashboard } from '../../data/api';
import { platformTypeName } from '../../utils/display';
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
      platforms: [PlatformTypeType.DEFAULT],
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
            disableMaskedInput
            renderInput={(props) => <TextField size="small" {...props}></TextField>}
            label="Expires at"
            onChange={(dateTime: DateTime | null) => {
              formik.setFieldValue('expiresAt', dateTime?.toISO());
            }}
            value={formik.values.expiresAt ? DateTime.fromISO(formik.values.expiresAt) : null}
          ></DatePicker>
        </Box>
        <Box sx={{ py: 2 }}>
          <FormControl>
            <InputLabel>Platform</InputLabel>
            <Select
              label="Platform"
              name="platforms"
              size="small"
              onBlur={formik.handleBlur}
              onChange={(e) => {
                formik.setFieldValue('platforms', [e.target.value]);
              }}
              value={formik.values.platforms}
              input={<OutlinedInput />}
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
        </Box>
      </CardContent>
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
