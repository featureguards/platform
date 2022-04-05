import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

import { Dashboard } from '../../data/api';

import type { RootState } from '../../data/store';
import type { User } from '../../api';
// Define a type for the slice state
interface UsersState {
  me: User | null;
  status: 'idle' | 'loading' | 'succeeded' | 'failed';
  error: string | null;
}

// Define the initial state using that type
const initialState: UsersState = {
  me: null,
  status: 'idle',
  error: null
};

export const fetchMe = createAsyncThunk('users/fetchMe', async () => {
  const res = await Dashboard.getUser('me');
  return res.data;
});

export const usersSlice = createSlice({
  name: 'users',
  // `createSlice` will infer the state type from the `initialState` argument
  initialState,
  reducers: {
    // Use the PayloadAction type to declare the contents of `action.payload`
    setMe: (state, action: PayloadAction<User, string>) => {
      state.me = action.payload;
    }
  },
  extraReducers(builder) {
    builder
      .addCase(fetchMe.pending, (state, _action) => {
        state.status = 'loading';
      })
      .addCase(fetchMe.fulfilled, (state, action) => {
        state.status = 'succeeded';
        // Add any fetched posts to the array
        state.me = action.payload;
      })
      .addCase(fetchMe.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.error.message || null;
      });
  }
});

// export const { setMe } = usersSlice.actions;

// Other code such as selectors can use the imported `RootState` type
export const selectMe = (state: RootState) => state.users.me;

export default usersSlice.reducer;
