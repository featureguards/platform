import { Card, CardContent, CardHeader, FormControlLabel, Grid, Switch } from '@mui/material';

import { OnOffFeature } from '../../api';
import { Matches } from './matches';

export type OnOffProps = {
  onOff: OnOffFeature;
  setOnOff: (_n: OnOffFeature) => void;
};

export const OnOff = ({ onOff, setOnOff }: OnOffProps) => {
  const setOn = (on: boolean) => {
    const weight = on ? 100 : 0;
    setOnOff({
      ...onOff,
      on: {
        ...onOff.on,
        weight: weight
      },
      off: {
        ...onOff.off,
        weight: 100 - weight
      }
    });
  };

  return (
    <>
      <FormControlLabel
        control={
          <Switch
            name="on"
            checked={!!onOff?.on?.weight}
            onChange={(e) => setOn(e.target.checked)}
          />
        }
        label="On/Off"
      />
      <Grid container pt={5} spacing={2} alignItems="center">
        <Grid item>
          <Card variant="outlined">
            <CardHeader title="Allow list" />
            <CardContent>
              <Matches
                matches={onOff.on?.matches || []}
                setMatches={(matches) =>
                  setOnOff({
                    ...onOff,
                    on: {
                      ...onOff.on,
                      matches: matches
                    }
                  })
                }
              ></Matches>
            </CardContent>
          </Card>
        </Grid>
        <Grid item>
          <Card variant="outlined">
            <CardHeader title="Disallow list" />
            <CardContent>
              <Matches
                matches={onOff.off?.matches || []}
                setMatches={(matches) =>
                  setOnOff({
                    ...onOff,
                    off: {
                      ...onOff.off,
                      matches: matches
                    }
                  })
                }
              ></Matches>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </>
  );
};
