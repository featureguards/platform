import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  TextField
} from '@mui/material';

import { IntOp } from '../../api';
import { IntOperator } from '../../api/enums';
import { intOperatorName } from '../../utils/display';

export type IntOpProps = {
  intOp: IntOp;
  setIntOp: (_n: IntOp) => void;
  creating?: boolean;
};

export const IntOpView = ({ intOp, setIntOp, creating }: IntOpProps) => {
  const setValues = (value: string) => {
    let v: number[] = [];
    switch (intOp.op) {
      case IntOperator.IN:
        v = value.split(',').map((s) => Number(s.trim()));
        break;
      default:
        v = [Number(value)];
    }
    setIntOp({
      ...intOp,
      values: v
    });
  };

  const renderValues = () => {
    switch (intOp.op) {
      case IntOperator.IN:
        return intOp.values?.join(',') || '';
      default:
        return intOp.values?.[0] || '';
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
            value={intOp.op || IntOperator.EQ}
            onChange={(e: SelectChangeEvent<any>) => {
              setIntOp({
                ...intOp,
                op: e.target.value
              });
            }}
          >
            {Object.entries(IntOperator).map((kt) => (
              <MenuItem key={kt[0]} value={kt[1]}>
                {intOperatorName(kt[1])}
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
