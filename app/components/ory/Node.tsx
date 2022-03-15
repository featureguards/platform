import {
  isUiNodeAnchorAttributes,
  isUiNodeImageAttributes,
  isUiNodeInputAttributes,
  isUiNodeScriptAttributes,
  isUiNodeTextAttributes
} from '@ory/integrations/ui';
import { UiNode } from '@ory/kratos-client';

import { NodeAnchor } from './NodeAnchor';
import { NodeImage } from './NodeImage';
import { NodeInput } from './NodeInput';
import { NodeScript } from './NodeScript';
import { NodeText } from './NodeText';
import { FormDispatcher, Formik, NodeProps } from './helpers';

export { type NodeProps };

interface Props {
  node: UiNode;
  disabled: boolean;
  value: any;
  dispatchSubmit: FormDispatcher;
  formik: Formik;
  propsOverride?: NodeProps;
  childrenOverride?: React.ReactNode;
  preNode?: React.ReactNode;
  postNode?: React.ReactNode;
}

export const Node = ({
  node,
  value,
  disabled,
  dispatchSubmit,
  formik,
  propsOverride,
  preNode,
  postNode,
  childrenOverride
}: Props) => {
  if (isUiNodeImageAttributes(node.attributes)) {
    return <NodeImage node={node} attributes={node.attributes} />;
  }

  if (isUiNodeScriptAttributes(node.attributes)) {
    return <NodeScript node={node} attributes={node.attributes} />;
  }

  if (isUiNodeTextAttributes(node.attributes)) {
    return <NodeText node={node} attributes={node.attributes} />;
  }

  if (isUiNodeAnchorAttributes(node.attributes)) {
    return <NodeAnchor node={node} attributes={node.attributes} />;
  }

  if (isUiNodeInputAttributes(node.attributes)) {
    return (
      <>
        {!!preNode ? preNode : null}
        <NodeInput
          dispatchSubmit={dispatchSubmit}
          value={value}
          node={node}
          disabled={disabled}
          attributes={node.attributes}
          formik={formik}
          propsOverride={propsOverride}
          childrenOverride={childrenOverride}
        />
        {!!postNode ? postNode : null}
      </>
    );
  }

  return null;
};
