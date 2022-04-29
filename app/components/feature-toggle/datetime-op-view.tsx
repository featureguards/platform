import { DateTime } from 'luxon';

import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  TextField
} from '@mui/material';
import { DateTimePicker } from '@mui/x-date-pickers';

import { DateTimeOp } from '../../api';
import { DateTimeOperator } from '../../api/enums';
import { dateTimeOperatorName } from '../../utils/display';

export type DateTimeOpProps = {
  dateTimeOp: DateTimeOp;
  setDateTimeOp: (_n: DateTimeOp) => void;
  creating?: boolean;
};

export const DateTimeOpView = ({ dateTimeOp, setDateTimeOp, creating }: DateTimeOpProps) => {
  const setValue = (value: DateTime | null) => {
    setDateTimeOp({
      ...dateTimeOp,
      timestamp: !!value ? value.toISO() : undefined
    });
  };

  const renderValue = () => {
    return dateTimeOp.timestamp
      ? DateTime.fromISO(dateTimeOp.timestamp)
      : DateTime.fromJSDate(new Date());
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
            value={dateTimeOp.op || DateTimeOperator.BEFORE}
            onChange={(e: SelectChangeEvent<any>) => {
              setDateTimeOp({
                ...dateTimeOp,
                op: e.target.value
              });
            }}
          >
            {Object.entries(DateTimeOperator).map((kt) => (
              <MenuItem key={kt[0]} value={kt[1]}>
                {dateTimeOperatorName(kt[1])}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <DateTimePicker
          renderInput={(props) => <TextField size="small" {...props}></TextField>}
          label="Date/time"
          onChange={async (dateTime) => setValue(dateTime)}
          value={renderValue()}
        ></DateTimePicker>
      </Box>
    </>
  );
};
