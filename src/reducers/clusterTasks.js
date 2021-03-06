/*
Copyright 2019 The Tekton Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import keyBy from 'lodash.keyby';
import { combineReducers } from 'redux';
import {
  createErrorMessageReducer,
  createIsFetchingReducer
} from './reducerCreators';

const type = 'ClusterTask';

function byName(state = {}, action) {
  switch (action.type) {
    case 'ClusterTaskCreated':
    case 'ClusterTaskUpdated':
      return { [action.payload.metadata.name]: action.payload, ...state };
    case 'ClusterTaskDeleted':
      const newState = { ...state };
      delete newState[action.payload.metadata.name];
      return newState;
    case 'CLUSTER_TASKS_FETCH_SUCCESS':
      return keyBy(action.data, 'metadata.name');
    default:
      return state;
  }
}

const isFetching = createIsFetchingReducer({ type });
const errorMessage = createErrorMessageReducer({ type });

export default combineReducers({
  byName,
  errorMessage,
  isFetching
});

export function getClusterTasks(state) {
  return Object.values(state.byName);
}

export function getClusterTask(state, name) {
  return state.byName[name];
}

export function getClusterTasksErrorMessage(state) {
  return state.errorMessage;
}

export function isFetchingClusterTasks(state) {
  return state.isFetching;
}
