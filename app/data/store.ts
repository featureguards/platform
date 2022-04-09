import { configureStore } from '@reduxjs/toolkit';

import projectInvitesReducer from '../features/project_invites/slice';
import projectsReducer from '../features/projects/slice';
import usersReducer from '../features/users/slice';

export const store = configureStore({
  reducer: {
    users: usersReducer,
    projects: projectsReducer,
    projectInvites: projectInvitesReducer
  }
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
