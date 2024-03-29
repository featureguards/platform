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
  InputLabel,
  MenuItem,
  OutlinedInput,
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
          <OutlinedInput
            startAdornment={<Percent />}
            value={percentage.on?.weight}
            size="small"
            sx={{ maxWidth: 110 }}
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
      <Grid container spacing={2} alignItems="center">
        <Grid item xs={12}>
          <FormControl>
            <Tooltip title="Controls whether the results are random every time or consistent based on the attributes provided.">
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
                <TextField
                  sx={{ mt: 2 }}
                  size="small"
                  helperText="By default, we use random values so that different feature flags target different populations for the same %. If multiple feature flags need to be on/off together, set the value here the same across all of them."
                  label="Group"
                  name="salt"
                  onChange={(e) =>
                    setSalt(e.target.value.toUpperCase().replace(/[^a-zA-Z0-9_-]/gi, ''))
                  }
                  value={percentage.salt}
                  variant="outlined"
                />

                <Box
                  sx={{
                    py: 2,
                    alignItems: 'center',
                    display: 'flex',
                    flexDirection: 'row'
                  }}
                >
                  <Typography variant="body1">Attributes</Typography>
                  <Tooltip
                    title={
                      'Attribute to evaluate when determining whether the feature is on/off.' +
                      ' Use a high-cardinality attribute to ensure enough distribution.'
                    }
                  >
                    <IconButton
                      color="primary"
                      onClick={() => {
                        setKeys([
                          ...(percentage?.stickiness?.keys || []),
                          { key: '', keyType: KeyType.STRING }
                        ]);
                      }}
                    >
                      <AddIcon></AddIcon>
                    </IconButton>
                  </Tooltip>
                </Box>
                <Box
                  sx={{
                    display: 'flex',
                    flexDirection: 'column'
                  }}
                >
                  <FormHelperText sx={{ mt: -3, mb: 2 }}>
                    Attributes passed in context for affinity (i.e., userId, requestID) will be
                    checked in order. The first match will be used.
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
                        label="Attribute"
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
                          {Object.values(KeyType).map((kt) => (
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
                </Box>
              </>
            )}
          </FormControl>
        </Grid>
        <Grid item sm={12}>
          <Typography variant="h6">Additional Constraints</Typography>
          <FormHelperText>
            In addition to the controls above, additional conditions can be used to allow/disallow a
            subset of the population matching the conditions below.
          </FormHelperText>
        </Grid>
        <Grid item>
          <Card variant="outlined">
            <CardHeader title="Allow list" />
            <CardContent sx={{ pt: 0 }}>
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
            <CardContent sx={{ pt: 0 }}>
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
