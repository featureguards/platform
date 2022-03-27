import PropTypes from 'prop-types';
import { styled } from '@mui/material/styles';
import { ReactNode, FC } from 'react';
import { Theme } from '@mui/material';

type OwnerState = {
  color: 'primary' | 'secondary' | 'error' | 'info' | 'warning' | 'success';
};

const SeverityPillRoot = styled('span')(({ theme, ownerState }: { theme: Theme; ownerState: OwnerState }) => {
  const backgroundColor = theme.palette[ownerState.color].main;
  const color = theme.palette[ownerState.color].contrastText;

  return {
    alignItems: 'center',
    backgroundColor,
    borderRadius: 12,
    color,
    cursor: 'default',
    display: 'inline-flex',
    flexGrow: 0,
    flexShrink: 0,
    fontFamily: theme.typography.fontFamily,
    fontSize: theme.typography.pxToRem(12),
    lineHeight: 2,
    fontWeight: 600,
    justifyContent: 'center',
    letterSpacing: 0.5,
    minWidth: 20,
    paddingLeft: theme.spacing(1),
    paddingRight: theme.spacing(1),
    textTransform: 'uppercase',
    whiteSpace: 'nowrap'
  };
});

export type SeverityPillProps = {
  color?: 'primary' | 'secondary' | 'error' | 'info' | 'warning' | 'success';
  children: ReactNode;
  theme: Theme;
};

export const SeverityPill: FC<SeverityPillProps> = (props) => {
  const { color = 'primary', children, theme, ...other } = props;

  const ownerState = { color };

  return (
    <SeverityPillRoot theme={theme} ownerState={ownerState} {...other}>
      {children}
    </SeverityPillRoot>
  );
};

SeverityPill.propTypes = {
  children: PropTypes.node,
  color: PropTypes.oneOf(['primary', 'secondary', 'error', 'info', 'warning', 'success'])
};
