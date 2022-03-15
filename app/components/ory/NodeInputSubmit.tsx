import { getNodeLabel } from '@ory/integrations/ui';
import { Button, ButtonProps } from '@mui/material';

import { NodeInputProps } from './helpers';

export function NodeInputSubmit<T>({
  node,
  attributes,
  disabled,
  dispatchSubmit,
  formik,
  propsOverride,
  childrenOverride
}: NodeInputProps) {
  return (
    <>
      <Button
        name={attributes.name}
        color="primary"
        fullWidth
        onClick={(e) => {
          // Prevent all native handlers
          e.stopPropagation();
          e.preventDefault();

          // On click, we set this value, and once set, dispatch the submission!
          formik.setFieldValue(attributes.name, attributes.value);
          dispatchSubmit();
        }}
        value={attributes.value || ''}
        disabled={attributes.disabled || disabled}
        {...(propsOverride as ButtonProps)}
      >
        {childrenOverride ? childrenOverride : getNodeLabel(node)}
      </Button>
    </>
  );
}
