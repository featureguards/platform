import { TextField } from '@mui/material';

import { FloatValue, IntegerValue } from '../../api';

type Props = {
  val: IntegerValue | FloatValue;
  setVal: (_n: IntegerValue | FloatValue) => void;
};

export const NumberVal = ({ val, setVal }: Props) => {
  const set = (on: number) => {
    setVal({
      ...val,
      value: on
    });
  };

  return (
    <TextField
      sx={{ mt: 2 }}
      size="small"
      label="Value"
      name="value"
      type={'number'}
      onChange={(e) => set(Number(e.target.value))}
      value={val.value || 0}
      variant="outlined"
    />
  );
};
