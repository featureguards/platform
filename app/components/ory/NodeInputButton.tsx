import { getNodeLabel } from '@ory/integrations/ui';
import { Button, ButtonProps } from '@mui/material';

import { NodeInputProps } from './helpers';

export function NodeInputButton<T>({
  node,
  attributes,
  disabled,
  dispatchSubmit,
  formik,
  propsOverride
}: NodeInputProps) {
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

  return (
    <>
      <Button
        name={attributes.name}
        onClick={(e) => {
          // Prevent all native handlers
          e.stopPropagation();
          e.preventDefault();

          onClick();
          formik.setFieldValue(attributes.name, attributes.value);
          dispatchSubmit();
        }}
        fullWidth
        value={attributes.value || ''}
        disabled={attributes.disabled || disabled}
        {...(propsOverride as ButtonProps)}
      >
        {getNodeLabel(node)}
      </Button>
    </>
  );
}
