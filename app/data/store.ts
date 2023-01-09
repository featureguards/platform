import { configureStore } from '@reduxjs/toolkit';

import dynamicSettingsReducer from '../features/dynamic_settings/slice';
import featureTogglesReducer from '../features/feature_toggles/slice';
import projectInvitesReducer from '../features/project_invites/slice';
import projectMembersReducer from '../features/project_members/slice';
import projectsReducer from '../features/projects/slice';
import usersReducer from '../features/users/slice';

export const store = configureStore({
  reducer: {
    featureToggles: featureTogglesReducer,
    dynamicSettings: dynamicSettingsReducer,
    users: usersReducer,
    projects: projectsReducer,
    projectInvites: projectInvitesReducer,
    projectMembers: projectMembersReducer
  }
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
