import { AxiosError } from 'axios';

import { createAsyncThunk, createSlice, SerializedError } from '@reduxjs/toolkit';

import { Dashboard } from '../../data/api';
import { SerializeError } from '../utils';

import type { RootState } from '../../data/store';
import type { FeatureToggle } from '../../api';
// Define a type for the slice state
interface FeatureTogglesState {
  environment: {
    id: string | null;
    items: FeatureToggle[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
    error: SerializedError | null;
  };
  details: {
    id: string | null;
    item: FeatureToggle | null;
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
    error: SerializedError | null;
  };
  history: {
    id: string | null;
    items: FeatureToggle[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
    error: SerializedError | null;
  };
}

// Define the initial state using that type
const initialState: FeatureTogglesState = {
  environment: {
    id: null,
    items: [],
    status: 'idle',
    error: null
  },
  details: {
    id: null,
    item: null,
    status: 'idle',
    error: null
  },
  history: {
    id: null,
    items: [],
    status: 'idle',
    error: null
  }
};

export type FeatureID = {
  projectID: string;
  id: string;
};

export type EnvironmentID = {
  projectID: string;
  environmentID: string;
};

export type EnvironmentFeatureID = FeatureID & { environmentID: string };

export const details = createAsyncThunk(
  'feature_toggles/details',
  async (props: EnvironmentFeatureID) => {
    try {
      const res = await Dashboard.getFeatureToggle(props.projectID, props.environmentID, props.id);
      return { id: props.id, item: res.data };
    } catch (err) {
      throw SerializeError(err as AxiosError);
    }
  }
);

export const history = createAsyncThunk(
  'feature_toggles/history',
  async (props: EnvironmentFeatureID) => {
    const res = await Dashboard.getFeatureToggleHistory(
      props.projectID,
      props.environmentID,
      props.id
    );
    return { id: props.id, history: res.data.history };
  }
);

export const list = createAsyncThunk(
  'feature_toggles/environment',
  async (props: EnvironmentID) => {
    try {
      const res = await Dashboard.listFeatureToggles(props.projectID, props.environmentID);
      return { environmentID: props.environmentID, features: res.data.features };
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
        state.details.item = action.payload || null;
        state.details.error = null;
      })
      .addCase(details.rejected, (state, action) => {
        state.details.status = 'failed';
        state.details.id = null;
        state.details.item = null;
        state.details.error = action.error || null;
      });

    builder
      .addCase(list.pending, (state, _action) => {
        state.environment.status = 'loading';
      })
      .addCase(list.fulfilled, (state, action) => {
        state.environment.status = 'succeeded';
        state.environment.id = action.payload.environmentID;
        state.environment.items = action.payload.features || [];
        state.environment.error = null;
      })
      .addCase(list.rejected, (state, action) => {
        state.environment.status = 'failed';
        state.environment.id = null;
        state.environment.items = [];
        state.environment.error = action.error || null;
      });

    builder
      .addCase(history.pending, (state, _action) => {
        state.history.status = 'loading';
      })
      .addCase(history.fulfilled, (state, action) => {
        state.history.status = 'succeeded';
        state.history.id = action.payload.id;
        state.history.items = action.payload.history || [];
        state.history.error = null;
      })
      .addCase(history.rejected, (state, action) => {
        state.history.status = 'failed';
        state.history.id = null;
        state.history.error = action.error || null;
      });
  }
});

// Other code such as selectors can use the imported `RootState` type
export const selectDetails = (state: RootState) => state.featureToggles.details.item;
export const selectEnvironment = (state: RootState) => state.featureToggles.environment.items;
export const selectHistory = (state: RootState) => state.featureToggles.history.items;

export default featureTogglesSlice.reducer;
