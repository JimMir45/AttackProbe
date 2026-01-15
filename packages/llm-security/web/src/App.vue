<template>
  <el-container class="app-container">
    <el-aside width="220px" class="aside">
      <div class="logo">
        <el-icon :size="24"><Aim /></el-icon>
        <span>LLM Security BAS</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#001529"
        text-color="#ffffffa6"
        active-text-color="#fff"
      >
        <el-menu-item index="/target">
          <el-icon><Aim /></el-icon>
          <span>目标管理</span>
        </el-menu-item>
        <el-menu-item index="/testcase">
          <el-icon><Document /></el-icon>
          <span>测试用例</span>
        </el-menu-item>
        <el-menu-item index="/task">
          <el-icon><VideoPlay /></el-icon>
          <span>任务管理</span>
        </el-menu-item>
        <el-menu-item index="/report">
          <el-icon><DataAnalysis /></el-icon>
          <span>测试报告</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>系统设置</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <span class="title">大模型安全有效性验证平台</span>
        <span class="version">v1.0.0</span>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => '/' + route.path.split('/')[1])
</script>

<style scoped>
.app-container {
  height: 100vh;
}

.aside {
  background-color: #001529;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  border-bottom: 1px solid #ffffff1a;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.title {
  font-size: 18px;
  font-weight: 500;
}

.version {
  color: #999;
  font-size: 12px;
}

.main {
  background: #f0f2f5;
  padding: 24px;
}
</style>
