import request from '@/utils/request';

export async function queryFakeList(params) {
  return request('/api/appList', {
    params,
  });
}
export async function removeFakeList(params) {
  const { count = 5, ...restParams } = params;
  return request('/api/appList', {
    method: 'POST',
    params: {
      count,
    },
    data: { ...restParams, method: 'delete' },
  });
}
export async function addFakeList(params) {
  const { count = 5, ...restParams } = params;
  return request('/api/appList', {
    method: 'POST',
    params: {
      count,
    },
    data: { ...restParams, method: 'post' },
  });
}
export async function updateFakeList(params) {
  const { count = 5, ...restParams } = params;
  return request('/api/appList', {
    method: 'POST',
    params: {
      count,
    },
    data: { ...restParams, method: 'update' },
  });
}
