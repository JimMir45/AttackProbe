import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

api.interceptors.response.use(
  response => {
    const data = response.data
    if (data.code === 0) {
      return data.data
    } else {
      ElMessage.error(data.msg || '请求失败')
      return Promise.reject(new Error(data.msg))
    }
  },
  error => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

// Target API
export const targetApi = {
  page: (params) => api.post('/target/page', params),
  add: (data) => api.post('/target/add', data),
  update: (data) => api.post('/target/update', data),
  delete: (id) => api.post('/target/delete', { id }),
  detail: (id) => api.post('/target/detail', { id }),
  options: () => api.post('/target/options'),
  test: (id) => api.post('/target/test', { id })
}

// TestCase API
export const testcaseApi = {
  page: (params) => api.post('/testcase/page', params),
  add: (data) => api.post('/testcase/add', data),
  update: (data) => api.post('/testcase/update', data),
  delete: (id) => api.post('/testcase/delete', { id }),
  detail: (id) => api.post('/testcase/detail', { id }),
  stats: () => api.post('/testcase/stats'),
  batchStatus: (ids, status) => api.post('/testcase/batch-status', { ids, status })
}

// Task API
export const taskApi = {
  page: (params) => api.post('/task/page', params),
  add: (data) => api.post('/task/add', data),
  delete: (id) => api.post('/task/delete', { id }),
  detail: (id) => api.post('/task/detail', { id }),
  start: (id) => api.post('/task/start', { id }),
  cancel: (id) => api.post('/task/cancel', { id }),
  progress: (id) => api.post('/task/progress', { id }),
  results: (params) => api.post('/task/results', params)
}

// System API
export const systemApi = {
  info: () => api.post('/system/info'),
  getConfig: (key) => api.post('/system/config/get', { key }),
  updateConfig: (key, value) => api.post('/system/config/update', { key, value })
}

export default api
