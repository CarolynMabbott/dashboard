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

import React, { Component } from 'react';
import { connect } from 'react-redux';
import { injectIntl } from 'react-intl';
import { Table } from '@tektoncd/dashboard-components';
import { fetchDashboardInfo } from '../../actions/about';
import { getSelectedNamespace, isWebSocketConnected } from '../../reducers';

export /* istanbul ignore next */ class About extends Component {
  state = {
    dashboardInfo: []
  };

  componentDidMount() {
    this.fetchDashboardInfo();
  }

  async fetchDashboardInfo() {
    const { namespace } = this.props;
    const dash = await fetchDashboardInfo(namespace);
    this.setState({
      dashboardInfo: dash
    });
  }

  render() {
    const { intl, loading, selectedNamespace } = this.props;

    const initialHeaders = [
      {
        key: 'property',
        header: intl.formatMessage({
          id: 'dashboard.tableHeader.Propety',
          defaultMessage: 'Property'
        })
      },
      {
        key: 'value',
        header: intl.formatMessage({
          id: 'dashboard.tableHeader.Value',
          defaultMessage: 'Value'
        })
      }
    ];

    const initialRows = [];
    this.state.dashboardInfo.forEach(function({
      id,
      dashboardVersion,
      pipelineVersion
    }) {
      let data = '';
      if (
        dashboardVersion !== '' &&
        dashboardVersion !== null &&
        dashboardVersion !== undefined
      ) {
        data = {
          id: `${id}1`,
          property: 'Dashboard Version',
          value: dashboardVersion
        };
      }
      if (
        pipelineVersion !== '' &&
        pipelineVersion !== null &&
        pipelineVersion !== undefined
      ) {
        data = {
          id: `${id}2`,
          property: 'Pipeline Version',
          value: pipelineVersion
        };
      }

      initialRows.push(data);
    });

    return (
      <>
        <h1>About</h1>
        <Table
          headers={initialHeaders}
          rows={initialRows}
          loading={loading}
          selectedNamespace={selectedNamespace}
        />
      </>
    );
  }
}

function mapStateToProps(state, props) {
  const { namespace: namespaceParam } = props.match.params;
  const namespace = namespaceParam || getSelectedNamespace(state);

  return {
    selectedNamespace: namespace,
    webSocketConnected: isWebSocketConnected(state)
  };
}

const mapDispatchToProps = {
  fetchDashboardInfo
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(injectIntl(About));
