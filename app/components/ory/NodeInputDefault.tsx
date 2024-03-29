import { TextField, TextFieldProps } from '@mui/material';
import { getNodeId, getNodeLabel } from '@ory/integrations/ui';

import { nestedValue } from '../../utils/helpers';
import { NodeInputProps } from './helpers';

export function NodeInputDefault(props: NodeInputProps) {
  const { node, attributes, value = '', disabled, formik, propsOverride } = props;

  // Some attributes have dynamic JavaScript - this is for example required for WebAuthn.
  const onClick = () => {
    // This section is only used for WebAuthn. The script is loaded via a <script> node
    // and the functions are available on the global window level. Unfortunately, there
    // is currently no better way than executing eval / function here at this moment.
    if (attributes.onclick) {
      const run = new Function(attributes.onclick);
      run();
    }
  };

  const label = (propsOverride as TextFieldProps)?.label?.toString() || getNodeLabel(node);
  const id = getNodeId(node);
  const serverError = !!node.messages.find(({ type }) => type === 'error');
  const touched = nestedValue(formik.touched, id);
  const formikError = nestedValue(formik.errors, id);
  const helperText = () => {
    // Prefer server errors over local errors to not mask them.
    if (!touched && node.messages.length > 0) {
      // Make text user friendly.
      return node.messages.map(({ text }) => text.replaceAll(id, label)).join('\n');
    }
    if (touched) {
      return formikError;
    }
  };

  // Render a generic text input field.
  return (
    <TextField
      label={label}
      onClick={onClick}
      onChange={formik.handleChange}
      onBlur={formik.handleBlur}
      fullWidth
      type={attributes.type}
      name={attributes.name}
      value={value}
      disabled={attributes.disabled || disabled}
      error={(!touched && serverError) || (touched && !!formikError)}
      helperText={helperText()}
      {...(propsOverride as TextFieldProps)}
    />
  );
}
