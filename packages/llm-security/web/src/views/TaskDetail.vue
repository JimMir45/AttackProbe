<template>
  <div class="page">
    <el-page-header @back="goBack" style="margin-bottom: 16px;">
      <template #content>
        <span>任务详情 - {{ task?.name }}</span>
      </template>
    </el-page-header>

    <el-row :gutter="16" style="margin-bottom: 16px;">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="总用例数" :value="task?.total_count || 0" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="防护成功" :value="task?.success_count || 0" value-style="color: #67c23a;" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="攻击成功" :value="task?.failed_count || 0" value-style="color: #f56c6c;" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="执行错误" :value="task?.error_count || 0" value-style="color: #e6a23c;" />
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <span>执行结果</span>
            <el-tag :type="statusTypeMap[task?.status]" style="margin-left: 8px;">
              {{ statusMap[task?.status] }}
            </el-tag>
            <el-progress
              v-if="task?.status === 1"
              :percentage="progress"
              :stroke-width="10"
              style="width: 200px; display: inline-block; margin-left: 16px;"
            />
          </div>
          <div>
            <el-button type="success" @click="handleStart" v-if="task?.status === 0">
              <el-icon><VideoPlay /></el-icon> 启动任务
            </el-button>
            <el-button type="warning" @click="handleCancel" v-if="task?.status === 1">
              <el-icon><VideoPause /></el-icon> 取消任务
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="results" v-loading="loading" stripe>
        <el-table-column prop="testcase_name" label="用例名称" width="200" />
        <el-table-column prop="testcase_category" label="分类" width="120" />
        <el-table-column label="判定结果" width="120">
          <template #default="{ row }">
            <el-tag :type="judgeTypeMap[row.judge_result]" v-if="row.status !== 0">
              {{ judgeMap[row.judge_result] || '待执行' }}
            </el-tag>
            <el-tag type="info" v-else>待执行</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="judge_reason" label="判定原因" min-width="200" show-overflow-tooltip />
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleViewResult(row)" v-if="row.status !== 0">
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadResults"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>

    <el-dialog v-model="resultDialogVisible" title="执行详情" width="800">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用例名称" :span="2">{{ currentResult?.testcase_name }}</el-descriptions-item>
        <el-descriptions-item label="判定结果">
          <el-tag :type="judgeTypeMap[currentResult?.judge_result]">
            {{ judgeMap[currentResult?.judge_result] }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="耗时">{{ currentResult?.duration }} ms</el-descriptions-item>
        <el-descriptions-item label="判定原因" :span="2">{{ currentResult?.judge_reason }}</el-descriptions-item>
      </el-descriptions>

      <el-divider>请求内容</el-divider>
      <pre class="content-pre">{{ formatJson(currentResult?.request_content) }}</pre>

      <el-divider>响应内容</el-divider>
      <pre class="content-pre response">{{ currentResult?.response_content }}</pre>

      <div v-if="currentResult?.error_message">
        <el-divider>错误信息</el-divider>
        <el-alert :title="currentResult?.error_message" type="error" :closable="false" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { taskApi } from '../api'

const route = useRoute()
const router = useRouter()
const taskId = computed(() => parseInt(route.params.id))

const statusMap = { 0: '待执行', 1: '执行中', 2: '已完成', 3: '已取消', 4: '失败' }
const statusTypeMap = { 0: 'info', 1: 'warning', 2: 'success', 3: 'danger', 4: 'danger' }
const judgeMap = { 0: '攻击成功', 1: '防护成功' }
const judgeTypeMap = { 0: 'danger', 1: 'success' }

const loading = ref(false)
const task = ref(null)
const results = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const progress = computed(() => {
  if (!task.value || task.value.total_count === 0) return 0
  return Math.round(task.value.completed_count / task.value.total_count * 100)
})

const resultDialogVisible = ref(false)
const currentResult = ref(null)

const formatJson = (str) => {
  if (!str) return ''
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

const loadTask = async () => {
  task.value = await taskApi.detail(taskId.value)
}

const loadResults = async () => {
  loading.value = true
  try {
    const res = await taskApi.results({
      task_id: taskId.value,
      page: { page: page.value, size: pageSize.value }
    })
    results.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const handleStart = async () => {
  await ElMessageBox.confirm('确定启动该任务吗？', '提示')
  await taskApi.start(taskId.value)
  ElMessage.success('任务已启动')
  loadTask()
}

const handleCancel = async () => {
  await ElMessageBox.confirm('确定取消该任务吗？', '提示')
  await taskApi.cancel(taskId.value)
  ElMessage.success('任务已取消')
  loadTask()
}

const handleViewResult = (row) => {
  currentResult.value = row
  resultDialogVisible.value = true
}

const goBack = () => {
  router.push('/task')
}

let timer = null

onMounted(() => {
  loadTask()
  loadResults()
  // Auto refresh for running task
  timer = setInterval(() => {
    if (task.value?.status === 1) {
      loadTask()
      loadResults()
    }
  }, 2000)
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
.content-pre {
  margin: 0;
  padding: 12px;
  background: #f5f5f5;
  border-radius: 4px;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: monospace;
  font-size: 13px;
  max-height: 200px;
  overflow-y: auto;
}
.content-pre.response {
  max-height: 300px;
}
</style>
