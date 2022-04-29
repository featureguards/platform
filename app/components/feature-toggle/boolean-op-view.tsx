import { Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent } from '@mui/material';

import { BoolOp } from '../../api';

export type BoolOpProps = {
  boolOp: BoolOp;
  setBoolOp: (_n: BoolOp) => void;
  creating?: boolean;
};

export const BoolOpView = ({ boolOp, setBoolOp, creating }: BoolOpProps) => {
  const renderValue = (v: number) => (!!v ? 'True' : 'False');
  const setValue = (v: number) =>
    setBoolOp({
      ...boolOp,
      value: !!v
    });

  return (
    <>
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          flexDirection: 'row'
        }}
      >
        <FormControl sx={{ minWidth: 100 }}>
          <InputLabel>Value</InputLabel>
          <Select
            disabled={!creating}
            size="small"
            value={boolOp.value || 0}
            onChange={(e: SelectChangeEvent<any>) => {
              setValue(e.target.value);
            }}
          >
            {[0, 1].map((kt) => (
              <MenuItem key={renderValue(kt)} value={kt}>
                {renderValue(kt)}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>
    </>
  );
};
