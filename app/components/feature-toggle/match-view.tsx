import DeleteIcon from '@mui/icons-material/Delete';
import {
  Box,
  FormControl,
  IconButton,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  TextField
} from '@mui/material';

import { Match } from '../../api';
import { KeyType } from '../../api/enums';
import { keyTypeName } from '../../utils/display';
import { StringOpView } from './string-op-view';

export type MatchProps = {
  match: Match;
  setMatch: (_n: Match) => void;
  creating?: boolean;
  onDelete?: () => void;
};

export const MatchView = ({ match, setMatch, creating, onDelete }: MatchProps) => {
  const renderOperator = () => {
    switch (match.key?.keyType) {
      case KeyType.STRING:
        return (
          <StringOpView
            creating={creating}
            stringOp={match.stringOp || {}}
            setStringOp={(v) => setMatch({ ...match, stringOp: v })}
          ></StringOpView>
        );
    }
  };
  return (
    <>
      <Box
        sx={{
          pt: 1,
          display: 'flex',
          justifyContent: 'space-between',
          flexDirection: 'row'
        }}
      >
        <TextField
          sx={{ mr: 1 }}
          required
          disabled={!creating}
          label="Key"
          size="small"
          onChange={async (e) => {
            setMatch({
              ...match,
              key: {
                ...match.key,
                key: e.target.value
              }
            });
          }}
          value={match.key?.key || ''}
        ></TextField>
        <FormControl sx={{ mx: 1 }}>
          <InputLabel>Type</InputLabel>
          <Select
            disabled={!creating}
            size="small"
            value={match.key?.keyType}
            onChange={(e: SelectChangeEvent<any>) => {
              setMatch({
                ...match,
                key: {
                  ...match.key,
                  keyType: e.target.value
                }
              });
            }}
          >
            {Object.entries(KeyType).map((kt) => (
              <MenuItem key={kt[0]} value={kt[1]}>
                {keyTypeName(kt[1])}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        {renderOperator()}
        {onDelete && (
          <IconButton onClick={() => onDelete()}>
            <DeleteIcon></DeleteIcon>
          </IconButton>
        )}
      </Box>
    </>
  );
};
