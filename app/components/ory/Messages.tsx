import { UiText } from '@ory/kratos-client';
// import { Alert, AlertContent } from '@ory/themes'
import { Alert } from '@mui/material';

interface MessageProps {
  message: UiText;
}

export const Message = ({ message }: MessageProps) => {
  return (
    <Alert
      severity={message.type === 'error' ? 'error' : 'info'}
      data-testid={`ui/message/${message.id}`}
    >
      {message.text}
    </Alert>
  );
};

interface MessagesProps {
  messages?: Array<UiText>;
}

export const Messages = ({ messages }: MessagesProps) => {
  if (!messages) {
    // No messages? Do nothing.
    return null;
  }

  return (
    <div>
      {messages.map((message) => (
        <Message key={message.id} message={message} />
      ))}
    </div>
  );
};
