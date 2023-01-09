import { Switch } from '@mui/material';

import { BoolValue } from '../../api';

type Props = {
  val: BoolValue;
  setVal: (_n: BoolValue) => void;
};

export const BoolVal = ({ val, setVal }: Props) => {
  const set = (on: boolean) => {
    setVal({
      ...val,
      value: on
    });
  };

  return <Switch name="on" checked={!!val?.value} onChange={(e) => set(e.target.checked)} />;
};
