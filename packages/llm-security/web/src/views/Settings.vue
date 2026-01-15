<template>
  <div class="page">
    <el-card>
      <template #header>系统设置</template>

      <el-form :model="form" label-width="150px" style="max-width: 600px;">
        <el-form-item label="并发执行数">
          <el-input-number v-model="form.concurrency" :min="1" :max="20" />
          <div class="tips">同时执行的测试用例数量</div>
        </el-form-item>
        <el-form-item label="单用例超时(ms)">
          <el-input-number v-model="form.timeout" :min="5000" :max="120000" :step="1000" />
          <div class="tips">单个用例的执行超时时间</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saving">保存设置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 16px;">
      <template #header>系统信息</template>

      <el-descriptions :column="2" border v-if="sysInfo">
        <el-descriptions-item label="系统版本">{{ sysInfo.version }}</el-descriptions-item>
        <el-descriptions-item label="测试用例数">{{ sysInfo.testcase_count }}</el-descriptions-item>
        <el-descriptions-item label="目标数">{{ sysInfo.target_count }}</el-descriptions-item>
        <el-descriptions-item label="任务数">{{ sysInfo.task_count }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card style="margin-top: 16px;">
      <template #header>关于</template>
      <div class="about">
        <h3>LLM Security BAS</h3>
        <p>大模型安全有效性验证平台</p>
        <p>用于验证大模型应用的安全防护能力，支持多种攻击类型的测试。</p>
        <el-divider />
        <p><strong>主要功能：</strong></p>
        <ul>
          <li>目标管理 - 配置待测试的大模型接口</li>
          <li>测试用例 - 内置 58+ 攻击用例，支持自定义</li>
          <li>任务管理 - 创建和执行安全测试任务</li>
          <li>测试报告 - 查看防护效果分析报告</li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { systemApi } from '../api'

const form = reactive({
  concurrency: 5,
  timeout: 30000
})

const saving = ref(false)
const sysInfo = ref(null)

const loadConfig = async () => {
  try {
    const concurrency = await systemApi.getConfig('executor.concurrency')
    const timeout = await systemApi.getConfig('executor.timeout')
    form.concurrency = parseInt(concurrency?.value) || 5
    form.timeout = parseInt(timeout?.value) || 30000
  } catch (e) {
    console.error('Load config error:', e)
  }
}

const loadSysInfo = async () => {
  try {
    sysInfo.value = await systemApi.info()
  } catch (e) {
    console.error('Load sys info error:', e)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await systemApi.updateConfig('executor.concurrency', String(form.concurrency))
    await systemApi.updateConfig('executor.timeout', String(form.timeout))
    ElMessage.success('设置已保存')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadConfig()
  loadSysInfo()
})
</script>

<style scoped>
.page {
  height: 100%;
}
.tips {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
.about {
  line-height: 1.8;
}
.about h3 {
  margin: 0 0 8px 0;
}
.about ul {
  padding-left: 20px;
}
</style>
