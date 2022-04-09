import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

import { Dashboard } from '../../data/api';

import type { RootState } from '../../data/store';
import type { ProjectInvite } from '../../api';
// Define a type for the slice state
interface ProjectInvitesState {
  forProject: {
    items: ProjectInvite[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
    error: string | null;
  };
  forUser: {
    items: ProjectInvite[];
    status: 'idle' | 'loading' | 'succeeded' | 'failed';
    error: string | null;
  };
}

// Define the initial state using that type
const initialState: ProjectInvitesState = {
  forProject: {
    items: [],
    status: 'idle',
    error: null
  },
  forUser: {
    items: [],
    status: 'idle',
    error: null
  }
};

export const fetch = createAsyncThunk('project_invites/fetch', async (id: string) => {
  const res = await Dashboard.getProjectInvite(id);
  return res.data;
});

export const fetchForProject = createAsyncThunk(
  'project_invites/fetchForProject',
  async (project_id: string) => {
    const res = await Dashboard.listProjectInvites(project_id);
    return res.data;
  }
);

export const fetchForUser = createAsyncThunk(
  'project_invites/fetchForUser',
  async (user_id: string) => {
    const res = await Dashboard.listUserInvites(user_id);
    return res.data;
  }
);

export const projectInvitesSlice = createSlice({
  name: 'project_invites',
  // `createSlice` will infer the state type from the `initialState` argument
  initialState,
  reducers: {
    setForProject: (state, action: PayloadAction<ProjectInvite[], string>) => {
      state.forProject.items = action.payload;
    },
    setForUser: (state, action: PayloadAction<ProjectInvite[], string>) => {
      state.forUser.items = action.payload;
    }
  },
  extraReducers(builder) {
    builder
      .addCase(fetchForProject.pending, (state, _action) => {
        state.forProject.status = 'loading';
      })
      .addCase(fetchForProject.fulfilled, (state, action) => {
        state.forProject.status = 'succeeded';
        state.forProject.items = action.payload.invites || [];
        state.forProject.error = null;
      })
      .addCase(fetchForProject.rejected, (state, action) => {
        state.forProject.status = 'failed';
        state.forProject.error = action.error.message || null;
      });
    builder
      .addCase(fetchForUser.pending, (state, _action) => {
        state.forUser.status = 'loading';
      })
      .addCase(fetchForUser.fulfilled, (state, action) => {
        state.forUser.status = 'succeeded';
        state.forUser.items = action.payload.invites || [];
        state.forUser.error = null;
      })
      .addCase(fetchForUser.rejected, (state, action) => {
        state.forUser.status = 'failed';
        state.forUser.error = action.error.message || null;
      });
  }
});

// Other code such as selectors can use the imported `RootState` type
export const selectCurrent = (state: RootState) => state.projectInvites.details.item;
export const selectProjectInvites = (state: RootState) => state.projectInvites.forProject.items;
export const selectUserInvites = (state: RootState) => state.projectInvites.forUser.items;

export default projectInvitesSlice.reducer;
