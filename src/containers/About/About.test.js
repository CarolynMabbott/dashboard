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

import React from 'react';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';
import configureStore from 'redux-mock-store';
import { Route } from 'react-router-dom';
import { urls } from '@tektoncd/dashboard-utils';
import { renderWithRouter } from '../../utils/test';
import * as ABOUT from '../../actions/about';
import About from '.';

const middleware = [thunk];
const mockStore = configureStore(middleware);
const namespaces = {
  byName: {
    default: {
      metadata: {
        name: 'default'
      }
    }
  },
  errorMessage: null,
  isFetching: false,
  selected: '*'
};

it('About renders correctly', async () => {
  jest
    .spyOn(ABOUT, 'fetchDashboardInfo')
    .mockImplementation(() => [
      { id: 'dashboardVersion', dashboardVersion: 'v0.100.0' },
      { id: 'pipelineVersion', pipelineVersion: 'v0.10.0' }
    ]);

  const store = mockStore({
    namespaces,
    notifications: {}
  });

  const { getByText, queryByText } = renderWithRouter(
    <Provider store={store}>
      <Route path={urls.about()} render={props => <About {...props} />} />
    </Provider>,
    { route: urls.about() }
  );

  await expect(ABOUT.fetchDashboardInfo).toHaveBeenCalledTimes(1);
  expect(queryByText('Property')).toBeTruthy();
  expect(getByText('Value')).toBeTruthy();
  expect(getByText('Dashboard Version')).toBeTruthy();
  expect(getByText('Pipeline Version')).toBeTruthy();
  expect(getByText('v0.100.0')).toBeTruthy();
  expect(getByText('v0.10.0')).toBeTruthy();
});
