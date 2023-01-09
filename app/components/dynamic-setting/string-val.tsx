import { TextField } from '@mui/material';

import { StringValue } from '../../api';

type Props = {
  val: StringValue;
  setVal: (_n: StringValue) => void;
};

export const StringVal = ({ val, setVal }: Props) => {
  const set = (on: string) => {
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
      onChange={(e) => set(e.target.value)}
      value={val.value || ''}
      variant="outlined"
    />
  );
};
