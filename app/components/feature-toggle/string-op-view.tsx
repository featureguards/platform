import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  TextField
} from '@mui/material';

import { StringOp } from '../../api';
import { StringOperator } from '../../api/enums';
import { stringOperatorName } from '../../utils/display';

export type StringOpProps = {
  stringOp: StringOp;
  setStringOp: (_n: StringOp) => void;
  creating?: boolean;
};

export const StringOpView = ({ stringOp, setStringOp, creating }: StringOpProps) => {
  const setValues = (value: string) => {
    let v: string[] = [];
    switch (stringOp.op) {
      case StringOperator.CONTAINS:
      case StringOperator.EQ:
        v = [value];
        break;
      case StringOperator.IN:
        v = value.split(',').map((s) => s.trim());
        break;
    }
    setStringOp({
      ...stringOp,
      values: v
    });
  };

  const renderValues = () => {
    switch (stringOp.op) {
      case StringOperator.CONTAINS:
      case StringOperator.EQ:
        return stringOp.values?.[0] || '';
      case StringOperator.IN:
        return stringOp.values?.join(',') || '';
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
        <FormControl sx={{ px: 1, minWidth: 100 }}>
          <InputLabel>Operator</InputLabel>
          <Select
            disabled={!creating}
            size="small"
            value={stringOp.op}
            onChange={(e: SelectChangeEvent<any>) => {
              setStringOp({
                ...stringOp,
                op: e.target.value
              });
            }}
          >
            {Object.entries(StringOperator).map((kt) => (
              <MenuItem key={kt[0]} value={kt[1]}>
                {stringOperatorName(kt[1])}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <TextField
          sx={{ px: 1 }}
          required
          label="Value(s)"
          size="small"
          onChange={async (e) => setValues(e.target.value)}
          value={renderValues()}
        ></TextField>
      </Box>
    </>
  );
};
