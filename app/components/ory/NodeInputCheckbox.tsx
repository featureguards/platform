import {
  Checkbox,
  CheckboxProps,
  FormControl,
  FormControlLabel,
  FormHelperText
} from '@mui/material';
import { getNodeId, getNodeLabel } from '@ory/integrations/ui';

import { NodeInputProps } from './helpers';

export function NodeInputCheckbox({
  node,
  attributes,
  disabled,
  value,
  formik,
  propsOverride
}: NodeInputProps) {
  // Render a checkbox.s
  const subtitle = node.messages.map(({ text }) => text).join('\n');
  const id = getNodeId(node);

  return (
    <FormControl
      color={node.messages.find(({ type }) => type === 'error') ? 'error' : undefined}
      disabled={attributes.disabled || disabled}
    >
      <FormControlLabel
        label={getNodeLabel(node)}
        control={
          <Checkbox
            name={attributes.name}
            defaultChecked={value === true}
            onChange={(e) => formik.setFieldValue(id, e.target.checked)}
            {...(propsOverride as CheckboxProps)}
          />
        }
      />
      {subtitle && <FormHelperText>{subtitle}</FormHelperText>}
    </FormControl>
  );
}
