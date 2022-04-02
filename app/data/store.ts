import { configureStore } from '@reduxjs/toolkit';

import projectsReducer from '../features/projects/slice';
import usersReducer from '../features/users/slice';

export const store = configureStore({
  reducer: {
    users: usersReducer,
    projects: projectsReducer
  }
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
