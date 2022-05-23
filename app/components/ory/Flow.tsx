import { FormikTouched, useFormik } from 'formik';
import { ReactNode, useEffect } from 'react';
import { OptionalObjectSchema } from 'yup/lib/object';

import { Grid } from '@mui/material';
import { getNodeId, isUiNodeInputAttributes } from '@ory/integrations/ui';
import {
  SelfServiceLoginFlow,
  SelfServiceRecoveryFlow,
  SelfServiceRegistrationFlow,
  SelfServiceSettingsFlow,
  SelfServiceVerificationFlow,
  SubmitSelfServiceLoginFlowBody,
  SubmitSelfServiceRecoveryFlowBody,
  SubmitSelfServiceRegistrationFlowBody,
  SubmitSelfServiceSettingsFlowBody,
  SubmitSelfServiceVerificationFlowBody,
  UiNode
} from '@ory/kratos-client';

import { nestedValue, setNestedValue } from '../../utils/helpers';
import { useNotifier } from '../hooks';
import { Messages } from './Messages';
import { Node, NodeProps } from './Node';

export type Values = Partial<
  | SubmitSelfServiceLoginFlowBody
  | SubmitSelfServiceRegistrationFlowBody
  | SubmitSelfServiceRecoveryFlowBody
  | SubmitSelfServiceSettingsFlowBody
  | SubmitSelfServiceVerificationFlowBody
>;

export type Methods =
  | 'oidc'
  | 'password'
  | 'profile'
  | 'totp'
  | 'webauthn'
  | 'link'
  | 'lookup_secret';
export type PropsOverrides = { [name: string]: NodeProps | PropsOverrides };
export type AugmentedNodes = { [name: string]: ReactNode };
export type Props<T> = {
  // The flow
  flow?:
    | SelfServiceLoginFlow
    | SelfServiceRegistrationFlow
    | SelfServiceSettingsFlow
    | SelfServiceVerificationFlow
    | SelfServiceRecoveryFlow;
  // Only show certain nodes. We will always render the default nodes for CSRF tokens.
  only?: Methods;
  // Is triggered on submission
  onSubmit: (_values: T) => Promise<void>;
  // Do not show the global messages. Useful when rendering them elsewhere.
  hideGlobalMessages?: boolean;
  notify?: boolean;
  validationSchema?: OptionalObjectSchema<any>;
  nodeProps?: PropsOverrides;
  preNodes?: AugmentedNodes;
  postNodes?: AugmentedNodes;
  childrenNodes?: AugmentedNodes;
  spacing?: number;
  sx?: any;
  initialValues?: T;
};

function emptyState<T>() {
  return {} as T;
}

export function Flow<T extends Values>(props: Props<T>) {
  const { hideGlobalMessages, flow, nodeProps, preNodes, postNodes, childrenNodes } = props;
  const formik = useFormik({
    initialValues: emptyState<T>(),
    validationSchema: props.validationSchema,
    onSubmit: async (values) => {
      return await props.onSubmit(values);
    }
  });
  const notifier = useNotifier();

  const initializeValues = (nodes: Array<UiNode> = []) => {
    // Compute the values
    const values: { [field: string]: any } = {};
    nodes.forEach((node) => {
      // This only makes sense for text nodes
      if (isUiNodeInputAttributes(node.attributes)) {
        if (node.attributes.type === 'button' || node.attributes.type === 'submit') {
          // In order to mimic real HTML forms, we need to skip setting the value
          // for buttons as the button value will (in normal HTML forms) only trigger
          // if the user clicks it.
          return;
        }

        const id = getNodeId(node) as keyof Values;
        setNestedValue(id, values, node.attributes.value);
      }
    });

    // Reset touched after resetting the flow.
    formik.resetForm();
    formik.setValues(values as T, false);
    const touched = Object.fromEntries(Object.entries(values).map(([k]) => [k, false]));
    formik.setTouched(touched as FormikTouched<T>);
  };

  const filterNodes = (): Array<UiNode> => {
    const { flow, only } = props;
    if (!flow) {
      return [];
    }
    return flow.ui.nodes.filter(({ group, meta }) => {
      if (!only) {
        return true;
      }
      // id (InfoNodeLabelID) is also included for oidc. Filter it out.
      if (only === 'oidc' && meta.label?.id === 1070004) {
        return false;
      }
      return group === 'default' || group === only;
    });
  };

  useEffect(() => {
    // Flow has changed, reload the values!
    initializeValues(filterNodes());

    if (props.notify && props.flow?.ui.messages && props.flow?.ui.messages.length) {
      const message = props.flow?.ui.messages[0];
      if (message.type === 'error') {
        notifier.error(message.text);
      } else {
        notifier.info(message.text);
      }
    }
  }, [props.flow, props.notify, notifier]);

  // Filter the nodes - only show the ones we want
  const nodes = filterNodes();

  if (!flow) {
    // No flow was set yet? It's probably still loading...
    //
    // Nodes have only one element? It is probably just the CSRF Token
    // and the filter did not match any elements!
    return null;
  }

  return (
    <form action={flow.ui.action} method={flow.ui.method} onSubmit={formik.handleSubmit}>
      {!hideGlobalMessages ? <Messages messages={flow.ui.messages} /> : null}
      <Grid sx={props.sx} container spacing={props.spacing || 3}>
        {nodes.map((node, k) => {
          const id = getNodeId(node) as keyof Values;
          return (
            <Grid item xs={12} md={12} key={`${id}-${k}`}>
              <Node
                disabled={formik.isSubmitting}
                node={node}
                value={nestedValue(formik.values, id)}
                dispatchSubmit={formik.handleSubmit}
                formik={formik}
                propsOverride={nodeProps ? nodeProps[id] : undefined}
                childrenOverride={childrenNodes ? childrenNodes[id] : undefined}
                preNode={preNodes ? preNodes[id] : undefined}
                postNode={postNodes ? postNodes[id] : undefined}
              />
            </Grid>
          );
        })}
      </Grid>
    </form>
  );
}
