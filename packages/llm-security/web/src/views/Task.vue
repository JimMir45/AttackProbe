<template>
  <div class="page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>任务管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新建任务
          </el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="任务名称" width="200" />
        <el-table-column label="目标" width="150">
          <template #default="{ row }">
            {{ row.target?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]">
              {{ statusMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="200">
          <template #default="{ row }">
            <el-progress
              :percentage="row.total_count > 0 ? Math.round(row.completed_count / row.total_count * 100) : 0"
              :status="row.status === 2 ? 'success' : (row.status >= 3 ? 'exception' : '')"
            />
            <span class="progress-text">{{ row.completed_count }}/{{ row.total_count }}</span>
          </template>
        </el-table-column>
        <el-table-column label="成功/失败/错误" width="130">
          <template #default="{ row }">
            <span class="success">{{ row.success_count }}</span> /
            <span class="failed">{{ row.failed_count }}</span> /
            <span class="error">{{ row.error_count }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleDetail(row)">详情</el-button>
            <el-button link type="success" @click="handleStart(row)" v-if="row.status === 0">启动</el-button>
            <el-button link type="warning" @click="handleCancel(row)" v-if="row.status === 1">取消</el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="row.status !== 1">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadData"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" title="新建任务" width="500">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="form.name" placeholder="输入任务名称" />
        </el-form-item>
        <el-form-item label="测试目标" prop="target_id">
          <el-select v-model="form.target_id" placeholder="选择目标" style="width: 100%;">
            <el-option
              v-for="t in targets"
              :key="t.id"
              :label="t.name"
              :value="t.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="测试用例">
          <el-select
            v-model="form.testcase_ids"
            placeholder="全部启用的用例"
            multiple
            filterable
            style="width: 100%;"
            clearable
          >
            <el-option
              v-for="tc in testcases"
              :key="tc.id"
              :label="tc.name"
              :value="tc.id"
            />
          </el-select>
          <div class="tips">不选则使用全部启用的用例</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { taskApi, targetApi, testcaseApi } from '../api'

const router = useRouter()

const statusMap = { 0: '待执行', 1: '执行中', 2: '已完成', 3: '已取消', 4: '失败' }
const statusTypeMap = { 0: 'info', 1: 'warning', 2: 'success', 3: 'danger', 4: 'danger' }

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const formRef = ref()
const form = reactive({
  name: '',
  target_id: null,
  testcase_ids: []
})

const rules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  target_id: [{ required: true, message: '请选择测试目标', trigger: 'change' }]
}

const targets = ref([])
const testcases = ref([])

const formatTime = (t) => {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await taskApi.page({
      page: { page: page.value, size: pageSize.value }
    })
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const loadTargets = async () => {
  targets.value = await targetApi.options() || []
}

const loadTestcases = async () => {
  const res = await testcaseApi.page({
    status: 1,
    page: { page: 1, size: 1000 }
  })
  testcases.value = res.list || []
}

const handleAdd = async () => {
  await loadTargets()
  await loadTestcases()
  Object.assign(form, {
    name: `测试任务-${new Date().toLocaleDateString('zh-CN')}`,
    target_id: targets.value[0]?.id || null,
    testcase_ids: []
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  await taskApi.add(form)
  ElMessage.success('任务创建成功')
  dialogVisible.value = false
  loadData()
}

const handleStart = async (row) => {
  await ElMessageBox.confirm('确定启动该任务吗？', '提示')
  await taskApi.start(row.id)
  ElMessage.success('任务已启动')
  loadData()
}

const handleCancel = async (row) => {
  await ElMessageBox.confirm('确定取消该任务吗？', '提示')
  await taskApi.cancel(row.id)
  ElMessage.success('任务已取消')
  loadData()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该任务吗？', '提示')
  await taskApi.delete(row.id)
  ElMessage.success('删除成功')
  loadData()
}

const handleDetail = (row) => {
  router.push(`/task/${row.id}`)
}

let timer = null

onMounted(() => {
  loadData()
  // Auto refresh for running tasks
  timer = setInterval(() => {
    if (list.value.some(t => t.status === 1)) {
      loadData()
    }
  }, 3000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.page {
  height: 100%;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.progress-text {
  font-size: 12px;
  color: #999;
  margin-left: 8px;
}
.success { color: #67c23a; }
.failed { color: #f56c6c; }
.error { color: #e6a23c; }
.tips {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>
