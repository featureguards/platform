import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

import { Dashboard } from '../../data/api';

import type { RootState } from '../../data/store';
import type { Project, Environment } from '../../api';
// Define a type for the slice state
interface ProjectsState {
  details: {
    item: Project | null;
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
  };
  all: {
    items: Project[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
  };
  environment: {
    item: Environment | null;
  };
}

// Define the initial state using that type
const initialState: ProjectsState = {
  details: {
    item: null,
    status: 'idle'
  },
  all: {
    items: [],
    status: 'idle'
  },
  environment: {
    item: null
  }
};

export const fetch = createAsyncThunk('projects/fetch', async (projectID: string) => {
  const res = await Dashboard.getProject(projectID);
  return res.data;
});

export const fetchAll = createAsyncThunk('projects/fetchAll', async () => {
  const res = await Dashboard.listProjects();
  return res.data;
});

export const projectsSlice = createSlice({
  name: 'projects',
  // `createSlice` will infer the state type from the `initialState` argument
  initialState,
  reducers: {
    setDetails: (state, action: PayloadAction<Project, string>) => {
      state.details.item = action.payload;
    },
    setAll: (state, action: PayloadAction<Project[], string>) => {
      state.all.items = action.payload;
    },
    setEnvironment: (state, action: PayloadAction<Environment | null, string>) => {
      state.environment.item = action.payload;
    }
  },
  extraReducers(builder) {
    builder
      .addCase(fetch.pending, (state, _action) => {
        state.details.status = 'loading';
      })
      .addCase(fetch.fulfilled, (state, action) => {
        state.details.status = 'succeeded';
        state.details.item = action.payload;
      })
      .addCase(fetch.rejected, (state) => {
        state.details.status = 'failed';
        state.details.item = null;
      });
    builder
      .addCase(fetchAll.pending, (state, _action) => {
        state.all.status = 'loading';
      })
      .addCase(fetchAll.fulfilled, (state, action) => {
        state.all.status = 'succeeded';
        state.all.items = action.payload.projects || [];
      })
      .addCase(fetchAll.rejected, (state) => {
        state.all.status = 'failed';
        state.all.items = [];
      });
  }
});

// Other code such as selectors can use the imported `RootState` type
export const selectCurrent = (state: RootState) => state.projects.details.item;
export const selectProjects = (state: RootState) => state.projects.all.items;

export default projectsSlice.reducer;
