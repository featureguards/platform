import AddIcon from '@mui/icons-material/Add';
import { Box, FormHelperText, IconButton, Typography } from '@mui/material';

import { Match } from '../../api';
import { KeyType, StringOperator } from '../../api/enums';
import { MatchView } from './match-view';

type MatchCreate = Match & { creating?: boolean };

export type MatchesProps = {
  matches: MatchCreate[];
  setMatches: (_v: MatchCreate[]) => void;
};

export const Matches = ({ matches, setMatches }: MatchesProps) => {
  return (
    <>
      <Box
        sx={{
          alignItems: 'center',
          display: 'flex',
          flexDirection: 'row'
        }}
      >
        <Typography variant="body1">Keys</Typography>
        <IconButton
          onClick={() => {
            setMatches([
              ...(matches || []),
              {
                key: { key: '', keyType: KeyType.STRING },
                stringOp: {
                  op: StringOperator.EQ,
                  values: []
                },
                creating: true
              }
            ]);
          }}
        >
          <AddIcon></AddIcon>
        </IconButton>
      </Box>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column'
        }}
      >
        <FormHelperText sx={{ mb: 2 }}>
          Keys passed in context for affinity (i.e., userId, requestID)
        </FormHelperText>
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column'
          }}
        >
          {matches.map((m, index) => (
            <MatchView
              key={index}
              match={m}
              creating={m.creating}
              setMatch={(m) => {
                matches[index] = m;
                setMatches(matches);
              }}
              onDelete={() => {
                setMatches(matches.filter((_, i) => i !== index));
              }}
            ></MatchView>
          ))}
        </Box>
      </Box>
    </>
  );
};
