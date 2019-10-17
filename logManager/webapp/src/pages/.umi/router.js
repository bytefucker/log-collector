import React from 'react';
import { Router as DefaultRouter, Route, Switch } from 'react-router-dom';
import dynamic from 'umi/dynamic';
import renderRoutes from 'umi/lib/renderRoutes';
import history from '@tmp/history';
import RendererWrapper0 from '/Users/yihongzhi/GoProject/logCollect/logManager/webapp/src/pages/.umi/LocaleWrapper.jsx';

const Router = require('dva/router').routerRedux.ConnectedRouter;

const routes = [
  {
    path: '/user',
    component: require('../../layouts/UserLayout').default,
    routes: [
      {
        name: 'login',
        path: '/user/login',
        component: require('../user/login').default,
        exact: true,
      },
      {
        component: () =>
          React.createElement(
            require('/Users/yihongzhi/GoProject/logCollect/logManager/webapp/node_modules/umi-build-dev/lib/plugins/404/NotFound.js')
              .default,
            { pagesPath: 'src/pages', hasRoutesInConfig: true },
          ),
      },
    ],
  },
  {
    path: '/',
    component: require('../../layouts/SecurityLayout').default,
    routes: [
      {
        path: '/',
        component: require('../../layouts/BasicLayout').default,
        authority: ['admin', 'user'],
        routes: [
          {
            path: '/',
            redirect: '/dashboard',
            exact: true,
          },
          {
            name: '控制台',
            icon: 'dashboard',
            path: '/dashboard',
            component: require('../DashboardAnalysis').default,
            exact: true,
          },
          {
            name: '实时日志',
            icon: 'desktop',
            path: '/realtime',
            component: require('../Realtime').default,
            exact: true,
          },
          {
            name: '监控报警',
            icon: 'alert',
            path: '/warnning',
            exact: true,
          },
          {
            name: '设置',
            icon: 'setting',
            path: '/setting',
            routes: [
              {
                name: '服务器管理',
                path: '/setting/server',
                component: require('../Setting/Server').default,
                exact: true,
              },
              {
                name: '应用管理',
                path: '/setting/application',
                component: require('../Setting/Application').default,
                exact: true,
              },
              {
                component: () =>
                  React.createElement(
                    require('/Users/yihongzhi/GoProject/logCollect/logManager/webapp/node_modules/umi-build-dev/lib/plugins/404/NotFound.js')
                      .default,
                    { pagesPath: 'src/pages', hasRoutesInConfig: true },
                  ),
              },
            ],
          },
          {
            component: require('../Exception404').default,
            exact: true,
          },
          {
            component: () =>
              React.createElement(
                require('/Users/yihongzhi/GoProject/logCollect/logManager/webapp/node_modules/umi-build-dev/lib/plugins/404/NotFound.js')
                  .default,
                { pagesPath: 'src/pages', hasRoutesInConfig: true },
              ),
          },
        ],
      },
      {
        component: require('../Exception404').default,
        exact: true,
      },
      {
        component: () =>
          React.createElement(
            require('/Users/yihongzhi/GoProject/logCollect/logManager/webapp/node_modules/umi-build-dev/lib/plugins/404/NotFound.js')
              .default,
            { pagesPath: 'src/pages', hasRoutesInConfig: true },
          ),
      },
    ],
  },
  {
    component: require('../Exception404').default,
    exact: true,
  },
  {
    component: () =>
      React.createElement(
        require('/Users/yihongzhi/GoProject/logCollect/logManager/webapp/node_modules/umi-build-dev/lib/plugins/404/NotFound.js')
          .default,
        { pagesPath: 'src/pages', hasRoutesInConfig: true },
      ),
  },
];
window.g_routes = routes;
const plugins = require('umi/_runtimePlugin');
plugins.applyForEach('patchRoutes', { initialValue: routes });

export { routes };

export default class RouterWrapper extends React.Component {
  unListen() {}

  constructor(props) {
    super(props);

    // route change handler
    function routeChangeHandler(location, action) {
      plugins.applyForEach('onRouteChange', {
        initialValue: {
          routes,
          location,
          action,
        },
      });
    }
    this.unListen = history.listen(routeChangeHandler);
    routeChangeHandler(history.location);
  }

  componentWillUnmount() {
    this.unListen();
  }

  render() {
    const props = this.props || {};
    return (
      <RendererWrapper0>
        <Router history={history}>{renderRoutes(routes, props)}</Router>
      </RendererWrapper0>
    );
  }
}
