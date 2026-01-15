<template>
  <div class="page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>测试用例</span>
          <div>
            <el-button @click="loadStats">
              <el-icon><DataAnalysis /></el-icon> 统计
            </el-button>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon> 新增用例
            </el-button>
          </div>
        </div>
      </template>

      <el-form :inline="true" style="margin-bottom: 16px;">
        <el-form-item label="分类">
          <el-select v-model="query.category" placeholder="全部" clearable @change="loadData" style="width: 140px;">
            <el-option label="提示词注入" :value="1" />
            <el-option label="越狱攻击" :value="2" />
            <el-option label="敏感信息泄露" :value="3" />
            <el-option label="有害内容生成" :value="4" />
            <el-option label="其他" :value="99" />
          </el-select>
        </el-form-item>
        <el-form-item label="风险等级">
          <el-select v-model="query.risk_level" placeholder="全部" clearable @change="loadData" style="width: 100px;">
            <el-option label="低" :value="1" />
            <el-option label="中" :value="2" />
            <el-option label="高" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源">
          <el-select v-model="query.is_builtin" placeholder="全部" clearable @change="loadData" style="width: 100px;">
            <el-option label="内置" :value="1" />
            <el-option label="自定义" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" width="200" />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            <el-tag>{{ categoryMap[row.category] || '其他' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="风险等级" width="100">
          <template #default="{ row }">
            <el-tag :type="riskTypeMap[row.risk_level]">
              {{ riskMap[row.risk_level] || '-' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="来源" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_builtin === 1 ? 'info' : 'warning'" size="small">
              {{ row.is_builtin === 1 ? '内置' : '自定义' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleView(row)">查看</el-button>
            <el-button link type="primary" @click="handleEdit(row)" v-if="row.is_builtin !== 1">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="row.is_builtin !== 1">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" v-if="!viewMode">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="用例名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="form.category" placeholder="选择分类">
            <el-option label="提示词注入" :value="1" />
            <el-option label="越狱攻击" :value="2" />
            <el-option label="敏感信息泄露" :value="3" />
            <el-option label="有害内容生成" :value="4" />
            <el-option label="其他" :value="99" />
          </el-select>
        </el-form-item>
        <el-form-item label="风险等级" prop="risk_level">
          <el-select v-model="form.risk_level" placeholder="选择风险等级">
            <el-option label="低" :value="1" />
            <el-option label="中" :value="2" />
            <el-option label="高" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="攻击类型" prop="attack_type">
          <el-input v-model="form.attack_type" placeholder="如: direct_injection" />
        </el-form-item>
        <el-form-item label="System Prompt">
          <el-input v-model="form.system_prompt" type="textarea" :rows="3" placeholder="可选的系统提示词" />
        </el-form-item>
        <el-form-item label="攻击内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="6" placeholder="攻击内容" />
        </el-form-item>
        <el-form-item label="判定配置" prop="judge_config">
          <el-input v-model="form.judge_config" type="textarea" :rows="3" placeholder='{"keywords":["我无法","抱歉"],"mode":"any"}' />
        </el-form-item>
      </el-form>

      <el-descriptions v-else :column="2" border>
        <el-descriptions-item label="名称" :span="2">{{ form.name }}</el-descriptions-item>
        <el-descriptions-item label="分类">{{ categoryMap[form.category] }}</el-descriptions-item>
        <el-descriptions-item label="风险等级">{{ riskMap[form.risk_level] }}</el-descriptions-item>
        <el-descriptions-item label="攻击类型">{{ form.attack_type }}</el-descriptions-item>
        <el-descriptions-item label="来源">{{ form.is_builtin === 1 ? '内置' : '自定义' }}</el-descriptions-item>
        <el-descriptions-item label="System Prompt" :span="2">
          <pre class="content-pre">{{ form.system_prompt || '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="攻击内容" :span="2">
          <pre class="content-pre">{{ form.content }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="判定配置" :span="2">
          <pre class="content-pre">{{ form.judge_config }}</pre>
        </el-descriptions-item>
      </el-descriptions>

      <template #footer v-if="!viewMode">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="statsVisible" title="用例统计" width="500">
      <el-descriptions :column="2" border v-if="stats">
        <el-descriptions-item label="总数">{{ stats.total }}</el-descriptions-item>
        <el-descriptions-item label="内置">{{ stats.builtin }}</el-descriptions-item>
      </el-descriptions>
      <el-table :data="stats?.category_stats || []" stripe style="margin-top: 16px;">
        <el-table-column prop="category_name" label="分类" />
        <el-table-column prop="count" label="数量" width="100" />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { testcaseApi } from '../api'

const categoryMap = { 1: '提示词注入', 2: '越狱攻击', 3: '敏感信息泄露', 4: '有害内容生成', 99: '其他' }
const riskMap = { 1: '低', 2: '中', 3: '高' }
const riskTypeMap = { 1: 'success', 2: 'warning', 3: 'danger' }

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const query = reactive({ category: null, risk_level: null, is_builtin: null })

const dialogVisible = ref(false)
const viewMode = ref(false)
const formRef = ref()
const form = reactive({
  id: null,
  name: '',
  category: 1,
  risk_level: 2,
  attack_type: '',
  system_prompt: '',
  content: '',
  judge_config: '{"keywords":["我无法","我不能","抱歉"],"mode":"any"}',
  status: 1
})

const dialogTitle = computed(() => {
  if (viewMode.value) return '查看用例'
  return form.id ? '编辑用例' : '新增用例'
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  content: [{ required: true, message: '请输入攻击内容', trigger: 'blur' }]
}

const statsVisible = ref(false)
const stats = ref(null)

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: { page: page.value, size: pageSize.value }
    }
    Object.keys(query).forEach(k => {
      if (query[k] !== null) params[k] = query[k]
    })
    const res = await testcaseApi.page(params)
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  stats.value = await testcaseApi.stats()
  statsVisible.value = true
}

const handleAdd = () => {
  viewMode.value = false
  Object.assign(form, {
    id: null,
    name: '',
    category: 1,
    risk_level: 2,
    attack_type: '',
    system_prompt: '',
    content: '',
    judge_config: '{"keywords":["我无法","我不能","抱歉"],"mode":"any"}',
    status: 1
  })
  dialogVisible.value = true
}

const handleView = async (row) => {
  const data = await testcaseApi.detail(row.id)
  Object.assign(form, data)
  viewMode.value = true
  dialogVisible.value = true
}

const handleEdit = async (row) => {
  const data = await testcaseApi.detail(row.id)
  Object.assign(form, data)
  viewMode.value = false
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  if (form.id) {
    await testcaseApi.update(form)
    ElMessage.success('更新成功')
  } else {
    await testcaseApi.add(form)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  loadData()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该用例吗？', '提示')
  await testcaseApi.delete(row.id)
  ElMessage.success('删除成功')
  loadData()
}

const handleStatusChange = async (row) => {
  await testcaseApi.batchStatus([row.id], row.status)
  ElMessage.success('状态更新成功')
}

onMounted(() => {
  loadData()
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
  white-space: pre-wrap;
  word-break: break-all;
  font-family: inherit;
}
</style>
