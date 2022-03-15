import { UiNode, UiNodeInputAttributes } from '@ory/kratos-client';
import { FormikTouched, FormikErrors } from 'formik';
import { TextFieldProps, ButtonProps, CheckboxProps } from '@mui/material';

export type FormDispatcher = () => void;
export type NodeProps = TextFieldProps | ButtonProps | CheckboxProps;
export type Formik = {
  handleBlur: {
    (e: React.FocusEvent<any>): void;
  };
  handleChange: {
    (e: React.ChangeEvent<any>): void;
  };
  touched: FormikTouched<any>;
  setFieldValue: (field: string, value: any, shouldValidate?: boolean) => void;
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
