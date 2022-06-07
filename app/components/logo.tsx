import PropTypes from 'prop-types';

import { styled } from '@mui/material/styles';
import { Typography, Box } from '@mui/material';

export type LogoProps = {
  variant?: 'light' | 'primary';
};

export const Logo = styled((props: LogoProps) => {
  const { variant } = props;
  // TODO: Have 2 logos for light vs not.
  const src = variant === 'light' ? '/logo.png' : '/logo.png';
  return (
    <Box
      sx={{
        flexDirection: 'row',
        display: 'flex',
        alignItems: 'center'
      }}
    >
      <img width="42" height="42" src={src} />
      <Typography variant="h6">FeatureGuards</Typography>
    </Box>
  );
})``;

Logo.defaultProps = {
  variant: 'primary'
};

Logo.propTypes = {
  variant: PropTypes.oneOf(['light', 'primary'])
};
