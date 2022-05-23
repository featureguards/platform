import { ChangeEvent } from 'react';

import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import Percent from '@mui/icons-material/Percent';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  FormControl,
  FormControlLabel,
  FormHelperText,
  FormLabel,
  Grid,
  IconButton,
  Input,
  InputLabel,
  MenuItem,
  Radio,
  RadioGroup,
  Select,
  SelectChangeEvent,
  Slider,
  TextField,
  Tooltip,
  Typography
} from '@mui/material';

import { Key, PercentageFeature } from '../../api';
import { KeyType, StickinessType } from '../../api/enums';
import { keyTypeName } from '../../utils/display';
import { Matches } from './matches';

export type PercentageProps = {
  percentage: PercentageFeature;
  setPercentage: (_n: PercentageFeature) => void;
};

export const Percentage = ({ percentage, setPercentage }: PercentageProps) => {
  const setWeight = (weight: number) => {
    setPercentage({
      ...percentage,
      on: {
        ...percentage.on,
        weight: weight
      },
      off: {
        ...percentage.off,
        weight: 100 - weight
      }
    });
  };
  const setSalt = (s: string) => {
    setPercentage({
      ...percentage,
      salt: s
    });
  };
  const setKeys = (keys?: Key[]) => {
    setPercentage({
      ...percentage,
      stickiness: {
        ...percentage.stickiness,
        keys: keys
      }
    });
  };
  const setAffinity = (v: StickinessType) => {
    const keys = v === StickinessType.KEYS ? percentage.stickiness?.keys : [];
    setPercentage({
      ...percentage,
      stickiness: {
        ...percentage.stickiness,
        keys: keys,
        stickinessType: v
      }
    });
  };
  const handleBlur = () => {
    let weight = percentage.on?.weight || 0;
    if (weight < 0) {
      weight = 0;
    } else if (weight > 100) {
      weight = 100;
    }
    setWeight(weight);
  };
  return (
    <>
      <Grid container spacing={2} alignItems="center">
        <Grid item xs={9} sm={6}>
          <Slider
            onChange={(_, val) => setWeight(val as number)}
            value={percentage.on?.weight}
            step={1}
            name="percentage"
            valueLabelDisplay="auto"
          />
        </Grid>
        <Grid item xs={3}>
          <Input
            startAdornment={<Percent />}
            value={percentage.on?.weight}
            size="small"
            onChange={(e) => setWeight(e.target.value === '' ? 0 : Number(e.target.value))}
            onBlur={handleBlur}
            inputProps={{
              step: 0.01,
              min: 0,
              max: 100,
              type: 'number'
            }}
          />
        </Grid>
      </Grid>
      <Grid container pt={5} spacing={2} alignItems="center">
        <Grid item xs={12}>
          <FormControl>
            <Tooltip title="How consistent results are using the same input">
              <FormLabel>Affinity</FormLabel>
            </Tooltip>
            <RadioGroup
              value={percentage.stickiness?.stickinessType || StickinessType.RANDOM}
              onChange={(e: ChangeEvent<HTMLInputElement>) => {
                setAffinity(Number(e.target.value) as StickinessType);
              }}
            >
              <Tooltip title="Checks for the same input often result in different values">
                <FormControlLabel
                  value={StickinessType.RANDOM}
                  control={<Radio />}
                  label="Random"
                />
              </Tooltip>
              <Tooltip title="Checks for the same input get the same output assuming keys match">
                <FormControlLabel value={StickinessType.KEYS} control={<Radio />} label="Sticky" />
              </Tooltip>
            </RadioGroup>
            {percentage.stickiness?.stickinessType === StickinessType.KEYS && (
              <>
                <Box
                  sx={{
                    py: 2,
                    alignItems: 'center',
                    display: 'flex',
                    flexDirection: 'row'
                  }}
                >
                  <Typography variant="body1">Keys</Typography>
                  <IconButton
                    onClick={() => {
                      setKeys([
                        ...(percentage?.stickiness?.keys || []),
                        { key: '', keyType: KeyType.STRING }
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
                  <FormHelperText sx={{ mt: -3, mb: 2 }}>
                    Keys passed in context for affinity (i.e., userId, requestID) will be checked in
                    order.
                  </FormHelperText>
                  {percentage.stickiness?.keys?.map((k, i) => (
                    <Box
                      key={i}
                      sx={{
                        py: 1,
                        display: 'flex',
                        justifyContent: 'space-between',
                        flexDirection: 'row'
                      }}
                    >
                      <TextField
                        required
                        label="Key"
                        size="small"
                        key={i}
                        onChange={async (e) => {
                          k.key = e.target.value;
                          setKeys(percentage.stickiness?.keys);
                        }}
                        value={k.key}
                      ></TextField>
                      <FormControl>
                        <InputLabel>Type</InputLabel>
                        <Select
                          size="small"
                          value={k.keyType}
                          onChange={(e: SelectChangeEvent<any>) => {
                            k.keyType = e.target.value;
                            setKeys(percentage.stickiness?.keys);
                          }}
                        >
                          {[KeyType.STRING, KeyType.FLOAT].map((kt) => (
                            <MenuItem key={kt} value={kt}>
                              {keyTypeName(kt)}
                            </MenuItem>
                          ))}
                        </Select>
                      </FormControl>
                      <IconButton
                        onClick={() => {
                          setKeys(percentage.stickiness?.keys?.filter((_, index) => index !== i));
                        }}
                      >
                        <DeleteIcon></DeleteIcon>
                      </IconButton>
                    </Box>
                  ))}
                  <TextField
                    sx={{ mt: 2 }}
                    size="small"
                    helperText="feature-toggles with the same value target the same population"
                    label="Group"
                    name="salt"
                    onChange={(e) =>
                      setSalt(e.target.value.toUpperCase().replace(/[^a-zA-Z0-9_-]/gi, ''))
                    }
                    value={percentage.salt}
                    variant="outlined"
                  />
                </Box>
              </>
            )}
          </FormControl>
        </Grid>
        <Grid item>
          <Card variant="outlined">
            <CardHeader title="Allow list" />
            <CardContent>
              <Matches
                matches={percentage.on?.matches || []}
                setMatches={(matches) =>
                  setPercentage({
                    ...percentage,
                    on: {
                      ...percentage.on,
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
                matches={percentage.off?.matches || []}
                setMatches={(matches) =>
                  setPercentage({
                    ...percentage,
                    off: {
                      ...percentage.off,
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
