import request from '@/utils/request';

export async function queryServerList(params) {
  return request('/api/serverList', {
    params,
  });
}
export async function deleteServer(params) {
  const { count = 5, ...restParams } = params;
  return request('/api/serverList', {
    method: 'POST',
    params: {
      count,
    },
    data: { ...restParams, method: 'delete' },
  });
}
export async function addServer(params) {
  const { count = 5, ...restParams } = params;
  return request('/api/serverList', {
    method: 'POST',
    params: {
      count,
    },
    data: { ...restParams, method: 'post' },
  });
}
export async function updateServer(params) {
  const { count = 5, ...restParams } = params;
  return request('/api/fake_list', {
    method: 'POST',
    params: {
      count,
    },
    data: { ...restParams, method: 'update' },
  });
}
