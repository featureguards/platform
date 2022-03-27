import { FormikErrors, FormikTouched } from 'formik';

import { ButtonProps, CheckboxProps, TextFieldProps } from '@mui/material';
import { UiNode, UiNodeInputAttributes } from '@ory/kratos-client';

export type FormDispatcher = () => void;
export type NodeProps = TextFieldProps | ButtonProps | CheckboxProps;
export type Formik = {
  handleBlur: {
    (_e: React.FocusEvent<any>): void;
  };
  handleChange: {
    (_e: React.ChangeEvent<any>): void;
  };
  touched: FormikTouched<any>;
  setFieldValue: (_field: string, _value: any, _shouldValidate?: boolean) => void;
  errors: FormikErrors<any>;
};

export interface NodeInputProps {
  node: UiNode;
  attributes: UiNodeInputAttributes;
  value: any;
  disabled: boolean;
  dispatchSubmit: FormDispatcher;
  formik: Formik;
  propsOverride?: NodeProps;
  childrenOverride?: React.ReactNode;
}
