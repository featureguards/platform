import { AxiosError } from 'axios';

import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';

import { Dashboard } from '../../data/api';
import { SerializeError } from '../utils';

import type { RootState } from '../../data/store';
import type { FeatureToggle, EnvironmentFeatureToggle } from '../../api';
// Define a type for the slice state
interface FeatureTogglesState {
  environment: {
    id: string | null;
    items: FeatureToggle[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
  };
  details: {
    id: string | null;
    items: EnvironmentFeatureToggle[] | null;
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
  };
  history: {
    id: string | null;
    items: FeatureToggle[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
  };
}

// Define the initial state using that type
const initialState: FeatureTogglesState = {
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

export type FeatureID = {
  id: string;
};

export type EnvironmentID = {
  environmentId: string;
};

export type EnvironmentIDs = {
  environmentIds: string[];
};

export type EnvironmentFeatureID = FeatureID & EnvironmentID;
export type FeatureIDEnvironments = FeatureID & EnvironmentIDs;

export const details = createAsyncThunk(
  'feature_toggles/details',
  async (props: FeatureIDEnvironments) => {
    try {
      const res = await Dashboard.getFeatureToggle(props.id, props.environmentIds);
      return {
        id: props.id,
        items: res.data.featureToggles
      };
    } catch (err) {
      throw SerializeError(err as AxiosError);
    }
  }
);

export const history = createAsyncThunk(
  'feature_toggles/history',
  async (props: EnvironmentFeatureID) => {
    const res = await Dashboard.getFeatureToggleHistoryForEnvironment(
      props.id,
      props.environmentId
    );
    return { id: props.id, history: res.data.history };
  }
);

export const list = createAsyncThunk(
  'feature_toggles/environment',
  async (props: EnvironmentID) => {
    try {
      const res = await Dashboard.listFeatureToggles(props.environmentId);
      return { environmentID: props.environmentId, features: res.data.featureToggles };
    } catch (err) {
      throw SerializeError(err as AxiosError);
    }
  }
);

export const featureTogglesSlice = createSlice({
  name: 'feature_toggles',
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
