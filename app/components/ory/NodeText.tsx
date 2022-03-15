import { UiNode, UiNodeTextAttributes } from '@ory/kratos-client';
import { UiText } from '@ory/kratos-client';
import { styled } from '@mui/material/styles';
import { Typography } from '@mui/material';

type ContentProps = {
  node: UiNode;
  attributes: UiNodeTextAttributes;
};
type CodeProps = {
  code: string;
};
const Code = (props: CodeProps) => {
  return <div>{props.code}</div>;
};

export const StyledCode = styled(Code)(({ theme }) => ({
  backgroundColor: theme.palette.grey[50],
  padding: 20,
  borderRadius: 8,
  color: theme.palette.grey[800],
  overflowWrap: 'break-word',
  overflowX: 'auto'
}));

const Content = ({ node, attributes }: ContentProps) => {
  switch (attributes.text.id) {
    case 1050015:
      // This text node contains lookup secrets. Let's make them a bit more beautiful!
      const secrets = (attributes.text.context as any).secrets.map((text: UiText, k: number) => (
        <div key={k} data-testid={`node/text/${attributes.id}/lookup_secret`} className="col-xs-3">
          {/* Used lookup_secret has ID 1050014 */}
          <code>{text.id === 1050014 ? 'Used' : text.text}</code>
        </div>
      ));
      return (
        <div className="container-fluid" data-testid={`node/text/${attributes.id}/text`}>
          <div className="row">{secrets}</div>
        </div>
      );
  }

  return (
    <div data-testid={`node/text/${attributes.id}/text`}>
      <StyledCode code={attributes.text.text} />
    </div>
  );
};

export const NodeText = ({ node, attributes }: ContentProps) => {
  return (
    <>
      <Typography data-testid={`node/text/${attributes.id}/label`}>{node.meta?.label?.text}</Typography>
      <Content node={node} attributes={attributes} />
    </>
  );
};
