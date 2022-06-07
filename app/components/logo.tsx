import PropTypes from 'prop-types';

import { styled } from '@mui/material/styles';
import { Typography, Box } from '@mui/material';

export type LogoProps = {
  variant?: 'light' | 'primary';
};

export const Logo = styled((props: LogoProps) => {
  const { variant, ...other } = props;

  const color = variant === 'light' ? '#C1C4D6' : '#5048E5';

  return (
    <Box
      sx={{
        flexDirection: 'row',
        display: 'flex',
        alignItems: 'center'
      }}
    >
      <img width="42" height="42" src="/logo.png" />
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
