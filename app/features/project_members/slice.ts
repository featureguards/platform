import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

import { Dashboard } from '../../data/api';

import type { RootState } from '../../data/store';
import type { ProjectMember } from '../../api';
// Define a type for the slice state
interface ProjectMembersState {
  id: string;
  members: ProjectMember[];
  status: 'idle' | 'loading' | 'succeeded' | 'failed';
  error: string | null;
}

// Define the initial state using that type
const initialState: ProjectMembersState = {
  id: '',
  members: [],
  status: 'idle',
  error: null
};

type Payload = {
  id: string;
  members: ProjectMember[];
};

export const fetch = createAsyncThunk('project_members/fetch', async (id: string) => {
  const res = await Dashboard.listProjectMembers(id);
  return { id: id, data: res.data };
});

export const projectMembersSlice = createSlice({
  name: 'project_members',
  // `createSlice` will infer the state type from the `initialState` argument
  initialState,
  reducers: {
    current: (state, action: PayloadAction<Payload, string>) => {
      state.members = action.payload.members;
      state.id = action.payload.id;
    }
  },
  extraReducers(builder) {
    builder
      .addCase(fetch.pending, (state, _action) => {
        state.status = 'loading';
      })
      .addCase(fetch.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.id = action.payload.id;
        state.members = action.payload.data.members || [];
        state.error = null;
      })
      .addCase(fetch.rejected, (state, action) => {
        state.status = 'failed';
        state.members = [];
        state.error = action.error.message || null;
      });
  }
});

// Other code such as selectors can use the imported `RootState` type
export const selectProjectMembers = (state: RootState) => {
  return { id: state.projectMembers.id, members: state.projectMembers.members };
};

export default projectMembersSlice.reducer;
