<template>
  <div class="page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>测试报告</span>
          <el-select v-model="selectedTaskId" placeholder="选择任务" @change="loadReport" style="width: 300px;">
            <el-option
              v-for="t in tasks"
              :key="t.id"
              :label="`${t.name} (${statusMap[t.status]})`"
              :value="t.id"
            />
          </el-select>
        </div>
      </template>

      <div v-if="!selectedTaskId" class="empty-state">
        <el-empty description="请选择一个已完成的任务查看报告" />
      </div>

      <div v-else-if="loading" v-loading="true" style="min-height: 400px;"></div>

      <div v-else>
        <el-row :gutter="16" style="margin-bottom: 24px;">
          <el-col :span="8">
            <el-card shadow="never">
              <div class="stat-card">
                <div class="stat-title">防护成功率</div>
                <div class="stat-value" :class="{ success: successRate >= 80, warning: successRate >= 50 && successRate < 80, danger: successRate < 50 }">
                  {{ successRate }}%
                </div>
                <el-progress
                  :percentage="successRate"
                  :stroke-width="10"
                  :color="successRate >= 80 ? '#67c23a' : (successRate >= 50 ? '#e6a23c' : '#f56c6c')"
                />
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="never">
              <div class="stat-card">
                <div class="stat-title">风险等级</div>
                <div class="stat-value" :class="riskLevel.class">
                  {{ riskLevel.text }}
                </div>
                <div class="stat-desc">{{ riskLevel.desc }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="never">
              <div class="stat-card">
                <div class="stat-title">测试覆盖</div>
                <div class="stat-value">{{ task?.total_count || 0 }}</div>
                <div class="stat-desc">共测试 {{ task?.total_count || 0 }} 个用例</div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-card shadow="never">
              <template #header>执行结果分布</template>
              <div class="chart-container">
                <div class="pie-legend">
                  <div class="legend-item">
                    <span class="dot success"></span>
                    <span>防护成功: {{ task?.success_count || 0 }}</span>
                  </div>
                  <div class="legend-item">
                    <span class="dot danger"></span>
                    <span>攻击成功: {{ task?.failed_count || 0 }}</span>
                  </div>
                  <div class="legend-item">
                    <span class="dot warning"></span>
                    <span>执行错误: {{ task?.error_count || 0 }}</span>
                  </div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card shadow="never">
              <template #header>攻击类型分析</template>
              <el-table :data="categoryStats" stripe size="small">
                <el-table-column prop="category" label="分类" />
                <el-table-column prop="total" label="总数" width="80" />
                <el-table-column prop="blocked" label="防护成功" width="90" />
                <el-table-column prop="passed" label="攻击成功" width="90" />
                <el-table-column label="防护率" width="100">
                  <template #default="{ row }">
                    <span :class="{ success: row.rate >= 80, warning: row.rate >= 50 && row.rate < 80, danger: row.rate < 50 }">
                      {{ row.rate }}%
                    </span>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </el-col>
        </el-row>

        <el-card shadow="never" style="margin-top: 16px;">
          <template #header>失败用例详情</template>
          <el-table :data="failedResults" stripe size="small" v-if="failedResults.length > 0">
            <el-table-column prop="testcase_name" label="用例名称" width="200" />
            <el-table-column prop="testcase_category" label="分类" width="120" />
            <el-table-column prop="risk_level" label="风险等级" width="100">
              <template #default="{ row }">
                <el-tag :type="riskTypeMap[row.risk_level]" size="small">
                  {{ row.risk_level }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="judge_reason" label="判定原因" show-overflow-tooltip />
          </el-table>
          <el-empty v-else description="没有失败的用例，防护表现优秀！" />
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { taskApi } from '../api'

const statusMap = { 0: '待执行', 1: '执行中', 2: '已完成', 3: '已取消', 4: '失败' }
const riskTypeMap = { '低': 'success', '中': 'warning', '高': 'danger' }
const categoryMap = { 1: '提示词注入', 2: '越狱攻击', 3: '敏感信息泄露', 4: '有害内容生成', 99: '其他' }

const loading = ref(false)
const tasks = ref([])
const selectedTaskId = ref(null)
const task = ref(null)
const results = ref([])

const successRate = computed(() => {
  if (!task.value || task.value.completed_count === 0) return 0
  const completed = task.value.success_count + task.value.failed_count
  if (completed === 0) return 0
  return Math.round(task.value.success_count / completed * 100)
})

const riskLevel = computed(() => {
  const rate = successRate.value
  if (rate >= 90) return { text: '低', class: 'success', desc: '防护能力优秀' }
  if (rate >= 70) return { text: '中', class: 'warning', desc: '存在一定风险' }
  if (rate >= 50) return { text: '高', class: 'danger', desc: '需要加强防护' }
  return { text: '严重', class: 'danger', desc: '防护严重不足' }
})

const categoryStats = computed(() => {
  const stats = {}
  results.value.forEach(r => {
    const cat = r.testcase_category || '其他'
    if (!stats[cat]) {
      stats[cat] = { category: cat, total: 0, blocked: 0, passed: 0 }
    }
    stats[cat].total++
    if (r.judge_result === 1) stats[cat].blocked++
    else if (r.judge_result === 0) stats[cat].passed++
  })
  return Object.values(stats).map(s => ({
    ...s,
    rate: s.total > 0 ? Math.round(s.blocked / s.total * 100) : 0
  }))
})

const failedResults = computed(() => {
  return results.value.filter(r => r.judge_result === 0)
})

const loadTasks = async () => {
  // 加载所有任务（不限状态）
  const res = await taskApi.page({
    page: { page: 1, size: 100 }
  })
  tasks.value = res.list || []

  // 自动选择第一个已完成的任务
  if (tasks.value.length > 0 && !selectedTaskId.value) {
    const completedTask = tasks.value.find(t => t.status === 2)
    if (completedTask) {
      selectedTaskId.value = completedTask.id
      loadReport()
    }
  }
}

const loadReport = async () => {
  if (!selectedTaskId.value) return
  loading.value = true
  try {
    task.value = await taskApi.detail(selectedTaskId.value)
    const res = await taskApi.results({
      task_id: selectedTaskId.value,
      page: { page: 1, size: 1000 }
    })
    results.value = res.list || []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadTasks()
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
.empty-state {
  padding: 60px 0;
}
.stat-card {
  text-align: center;
}
.stat-title {
  font-size: 14px;
  color: #999;
  margin-bottom: 8px;
}
.stat-value {
  font-size: 36px;
  font-weight: 600;
  margin-bottom: 8px;
}
.stat-desc {
  font-size: 12px;
  color: #999;
}
.success { color: #67c23a; }
.warning { color: #e6a23c; }
.danger { color: #f56c6c; }
.chart-container {
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.pie-legend {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}
.dot.success { background: #67c23a; }
.dot.danger { background: #f56c6c; }
.dot.warning { background: #e6a23c; }
</style>
