<template>
  <div class="page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>目标管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新增目标
          </el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" width="180" />
        <el-table-column prop="endpoint" label="接口地址" min-width="300" />
        <el-table-column prop="model" label="模型" width="150" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleTest(row)">测试</el-button>
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑目标' : '新增目标'" width="600">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="目标名称" />
        </el-form-item>
        <el-form-item label="接口地址" prop="endpoint">
          <el-input v-model="form.endpoint" placeholder="OpenAI兼容接口地址" />
        </el-form-item>
        <el-form-item label="API Key" prop="api_key">
          <el-input v-model="form.api_key" placeholder="API密钥" show-password />
        </el-form-item>
        <el-form-item label="模型" prop="model">
          <el-input v-model="form.model" placeholder="模型名称，如 gpt-4" />
        </el-form-item>
        <el-form-item label="超时时间" prop="timeout">
          <el-input-number v-model="form.timeout" :min="1000" :max="120000" :step="1000" />
          <span style="margin-left: 8px; color: #999;">毫秒</span>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { targetApi } from '../api'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const formRef = ref()
const form = reactive({
  id: null,
  name: '',
  endpoint: '',
  api_key: '',
  model: '',
  timeout: 30000,
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  endpoint: [{ required: true, message: '请输入接口地址', trigger: 'blur' }],
  api_key: [{ required: true, message: '请输入API Key', trigger: 'blur' }],
  model: [{ required: true, message: '请输入模型名称', trigger: 'blur' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await targetApi.page({
      page: { page: page.value, size: pageSize.value }
    })
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  Object.assign(form, {
    id: null,
    name: '',
    endpoint: '',
    api_key: '',
    model: '',
    timeout: 30000,
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = async (row) => {
  const data = await targetApi.detail(row.id)
  Object.assign(form, data)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  if (form.id) {
    await targetApi.update(form)
    ElMessage.success('更新成功')
  } else {
    await targetApi.add(form)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  loadData()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该目标吗？', '提示')
  await targetApi.delete(row.id)
  ElMessage.success('删除成功')
  loadData()
}

const handleTest = async (row) => {
  const loading = ElMessage({
    message: '正在测试连接...',
    duration: 0
  })
  try {
    await targetApi.test(row.id)
    loading.close()
    ElMessage.success('连接测试成功')
  } catch (e) {
    loading.close()
  }
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
</style>
