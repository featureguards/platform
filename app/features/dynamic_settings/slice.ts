import { AxiosError } from 'axios';

import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';

import { Dashboard } from '../../data/api';
import { SerializeError } from '../utils';

import type { RootState } from '../../data/store';
import type { DynamicSetting, EnvironmentDynamicSetting } from '../../api';

// Define a type for the slice state
interface DynamicSettingsState {
  environment: {
    id: string | null;
    items: DynamicSetting[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
  };
  details: {
    id: string | null;
    items: EnvironmentDynamicSetting[] | null;
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
  };
  history: {
    id: string | null;
    items: DynamicSetting[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
  };
}

// Define the initial state using that type
const initialState: DynamicSettingsState = {
  environment: {
    id: null,
    items: [],
    status: 'idle'
  },
  details: {
    id: null,
    items: [],
    status: 'idle'
  },
  history: {
    id: null,
    items: [],
    status: 'idle'
  }
};

export type SettingID = {
  id: string;
};

export type EnvironmentID = {
  environmentId: string;
};

export type EnvironmentIDs = {
  environmentIds: string[];
};

export type EnvironmentSettingID = SettingID & EnvironmentID;
export type SettingIDEnvironments = SettingID & EnvironmentIDs;

export const details = createAsyncThunk(
  'dynamic_settings/details',
  async (props: SettingIDEnvironments) => {
    try {
      const res = await Dashboard.getDynamicSetting(props.id, props.environmentIds);
      return {
        id: props.id,
        items: res.data.settings
      };
    } catch (err) {
      throw SerializeError(err as AxiosError);
    }
  }
);

export const history = createAsyncThunk(
  'dynamic_settings/history',
  async (props: EnvironmentSettingID) => {
    const res = await Dashboard.getDynamicSettingHistoryForEnvironment(
      props.id,
      props.environmentId
    );
    return { id: props.id, history: res.data.history };
  }
);

export const list = createAsyncThunk(
  'dynamic_settings/environment',
  async (props: EnvironmentID) => {
    try {
      const res = await Dashboard.listDynamicSettings(props.environmentId);
      return { environmentID: props.environmentId, features: res.data.dynamicSettings };
    } catch (err) {
      throw SerializeError(err as AxiosError);
    }
  }
);

export const featureTogglesSlice = createSlice({
  name: 'dynamic_settings',
  // `createSlice` will infer the state type from the `initialState` argument
  initialState,
  reducers: {},
  extraReducers(builder) {
    builder
      .addCase(details.pending, (state, _action) => {
        state.details.status = 'loading';
      })
      .addCase(details.fulfilled, (state, action) => {
        state.details.status = 'succeeded';
        state.details.id = action.payload.id;
        state.details.items = action.payload.items || null;
      })
      .addCase(details.rejected, (state) => {
        state.details.status = 'failed';
        state.details.id = null;
        state.details.items = [];
      });

    builder
      .addCase(list.pending, (state) => {
        state.environment.status = 'loading';
      })
      .addCase(list.fulfilled, (state, action) => {
        state.environment.status = 'succeeded';
        state.environment.id = action.payload.environmentID;
        state.environment.items = action.payload.features || [];
      })
      .addCase(list.rejected, (state) => {
        state.environment.status = 'failed';
        state.environment.id = null;
        state.environment.items = [];
      });

    builder
      .addCase(history.pending, (state, _action) => {
        state.history.status = 'loading';
      })
      .addCase(history.fulfilled, (state, action) => {
        state.history.status = 'succeeded';
        state.history.id = action.payload.id;
        state.history.items = action.payload.history || [];
      })
      .addCase(history.rejected, (state) => {
        state.history.status = 'failed';
        state.history.id = null;
      });
  }
});

// Other code such as selectors can use the imported `RootState` type
export const selectDetails = (state: RootState) => state.featureToggles.details.items;
export const selectEnvironment = (state: RootState) => state.featureToggles.environment.items;
export const selectHistory = (state: RootState) => state.featureToggles.history.items;

export default featureTogglesSlice.reducer;
