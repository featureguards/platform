import { DateTime } from 'luxon';
import { useState } from 'react';

import CopyIcon from '@mui/icons-material/ContentCopy';
import DeleteIcon from '@mui/icons-material/Delete';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import { sleep } from '../../utils/helpers';
import {
  FormControl,
  Grid,
  IconButton,
  Input,
  InputAdornment,
  Tooltip,
  Typography
} from '@mui/material';

import { ApiKey } from '../../api';
import { Dashboard } from '../../data/api';

export type ApiKeyProps = {
  apiKey: ApiKey;
  onDelete?: () => Promise<void>;
};

export const ApiKeyView = ({ apiKey, onDelete }: ApiKeyProps) => {
  const [showKey, setShowKey] = useState<boolean>(false);
  const [open, setOpen] = useState<boolean>(false);

  const handleTooltipClose = () => {
    setOpen(false);
  };

  const handleTooltipOpen = () => {
    setOpen(true);
  };
  return (
    <Grid container columnSpacing={4} direction="row" justifyContent="center" alignItems="center">
      <Grid item xs={12} sm={3} md={2}>
        <Typography variant="subtitle1">{apiKey.name}</Typography>
      </Grid>
      <Grid item xs={6} sm={4} md={3}>
        <Typography variant="subtitle2">
          Created At: {DateTime.fromISO(apiKey.createdAt || '').toLocaleString(DateTime.DATE_SHORT)}
        </Typography>
      </Grid>
      <Grid item xs={6} sm={4} md={3} sx={{ display: 'flex', flexDirection: 'row' }}>
        <FormControl sx={{ m: 1 }}>
          <Input
            readOnly
            size="small"
            type={showKey ? 'text' : 'password'}
            value={apiKey?.key}
            endAdornment={
              <InputAdornment position="end">
                <IconButton onClick={() => setShowKey(!showKey)} edge="end">
                  {showKey ? <VisibilityOff /> : <Visibility />}
                </IconButton>
              </InputAdornment>
            }
          />
        </FormControl>
        <Tooltip
          PopperProps={{
            disablePortal: true
          }}
          onClose={handleTooltipClose}
          open={open}
          disableFocusListener
          disableHoverListener
          disableTouchListener
          title="Copied"
        >
          <IconButton
            onClick={async () => {
              if (apiKey?.key) {
                navigator.clipboard.writeText(apiKey?.key);
                handleTooltipOpen();
                await sleep(200);
                handleTooltipClose();
              }
            }}
          >
            <CopyIcon />
          </IconButton>
        </Tooltip>
      </Grid>
      <Grid item xs={6} sm={4} md={3}>
        <Typography variant="subtitle2">
          Expires At:{' '}
          {apiKey.expiresAt
            ? DateTime.fromISO(apiKey.expiresAt).toLocaleString(DateTime.DATE_SHORT)
            : 'Never'}
        </Typography>
      </Grid>
      {!!onDelete && (
        <Grid item xs={6} sm={1} md={1}>
          <IconButton
            onClick={async () => {
              if (!apiKey?.id) {
                return;
              }
              await Dashboard.deleteApiKey(apiKey?.id);
              await onDelete();
            }}
          >
            <DeleteIcon></DeleteIcon>
          </IconButton>
        </Grid>
      )}
    </Grid>
  );
};
