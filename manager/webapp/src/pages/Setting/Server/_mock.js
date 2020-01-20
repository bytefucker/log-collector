const Mock = require('mockjs')

function fakeServer(count) {
  const list = [];
  for (let i = 0; i < count; i += 1) {
    list.push({
      id: Mock.mock('@guid'),
      logo: 'https://gw.alipayobjects.com/zos/rmsportal/kZzEzemZyKLKFsojXItE.png',
      hostname: `HOSTNAME${Mock.mock('@id')}`,
      ip: Mock.mock('@ip'),
      remark: '哈哈哈哈哈啦啦啦啦啦东方红都返回的',
      owner: Mock.mock('@cname'),
      createTime: new Date(),
    });
  }
  return list;
}

let sourceData = [];

function getServerList(req, res) {
  const params = req.query;
  const count = params.count * 1 || 20;
  const result = fakeServer(count);
  sourceData = result;
  return res.json(result);
}

function postServerList(req, res) {
  const {
    /* url = '', */
    body,
  } = req; // const params = getUrlParams(url);

  const { method, id } = body; // const count = (params.count * 1) || 20;

  let result = sourceData || [];

  switch (method) {
    case 'delete':
      result = result.filter(item => item.id !== id);
      break;

    case 'update':
      result.forEach((item, i) => {
        if (item.id === id) {
          result[i] = { ...item, ...body };
        }
      });
      break;

    case 'post':
      result.unshift({
        ...body,
        id: `fake-list-${result.length}`,
        createdAt: new Date().getTime(),
      });
      break;

    default:
      break;
  }

  return res.json(result);
}

export default {
  'POST  /api/serverList': postServerList,
  'GET  /api/serverList': getServerList,
};
