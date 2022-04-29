import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  TextField
} from '@mui/material';

import { FloatOp } from '../../api';
import { FloatOperator } from '../../api/enums';
import { floatOperatorName } from '../../utils/display';

export type FloatOpProps = {
  floatOp: FloatOp;
  setFloatOp: (_n: FloatOp) => void;
  creating?: boolean;
};

export const FloatOpView = ({ floatOp, setFloatOp, creating }: FloatOpProps) => {
  const setValues = (value: string) => {
    let v: number[] = [];
    switch (floatOp.op) {
      case FloatOperator.IN:
        v = value.split(',').map((s) => Number(s.trim()));
        break;
      default:
        v = [Number(value)];
    }
    setFloatOp({
      ...floatOp,
      values: v
    });
  };

  const renderValues = () => {
    switch (floatOp.op) {
      case FloatOperator.IN:
        return floatOp.values?.join(',') || '';
      default:
        return floatOp.values?.[0] || '';
    }
  };
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
          <InputLabel>Operator</InputLabel>
          <Select
            disabled={!creating}
            size="small"
            value={floatOp.op || FloatOperator.EQ}
            onChange={(e: SelectChangeEvent<any>) => {
              setFloatOp({
                ...floatOp,
                op: e.target.value
              });
            }}
          >
            {Object.entries(FloatOperator).map((kt) => (
              <MenuItem key={kt[0]} value={kt[1]}>
                {floatOperatorName(kt[1])}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <TextField
          sx={{ px: 1 }}
          required
          inputProps={{ inputMode: 'numeric', pattern: '[0-9,.]*' }}
          label="Value(s)"
          size="small"
          helperText="ex: 0.1 or 1,2,3 for In operator"
          onChange={async (e) => setValues(e.target.value)}
          value={renderValues()}
        ></TextField>
      </Box>
    </>
  );
};
