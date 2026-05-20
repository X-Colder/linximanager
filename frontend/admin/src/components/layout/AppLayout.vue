<template>
  <div class="app-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <template v-if="!sidebarCollapsed">
          <span class="logo-icon">&#9775;</span>
          <span class="logo-text">灵犀掌柜</span>
        </template>
        <template v-else>
          <span class="logo-icon">&#9775;</span>
        </template>
      </div>

      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        :collapse="sidebarCollapsed"
        router
        background-color="transparent"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <template #title>数据大盘</template>
        </el-menu-item>

        <el-sub-menu index="merchant">
          <template #title>
            <el-icon><UserFilled /></el-icon>
            <span>商家管理</span>
          </template>
          <el-menu-item index="/merchant/audit">入驻审核</el-menu-item>
          <el-menu-item index="/merchant/list">商家列表</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="miniapp">
          <template #title>
            <el-icon><Monitor /></el-icon>
            <span>小程序</span>
          </template>
          <el-menu-item index="/template">模板管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="finance">
          <template #title>
            <el-icon><Wallet /></el-icon>
            <span>财务管理</span>
          </template>
          <el-menu-item index="/billing/plans">计费方案</el-menu-item>
          <el-menu-item index="/billing/records">账单记录</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="operation">
          <template #title>
            <el-icon><Promotion /></el-icon>
            <span>运营活动</span>
          </template>
          <el-menu-item index="/activity/list">活动管理</el-menu-item>
          <el-menu-item index="/notification">消息推送</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </aside>

    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 顶部导航 -->
      <header class="topbar">
        <div class="topbar-left">
          <el-button
            text
            class="collapse-btn"
            @click="sidebarCollapsed = !sidebarCollapsed"
          >
            <el-icon size="20"><Fold v-if="!sidebarCollapsed" /><Expand v-else /></el-icon>
          </el-button>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="topbar-right">
          <el-badge :value="3" class="notice-badge">
            <el-button text>
              <el-icon size="20"><Bell /></el-icon>
            </el-button>
          </el-badge>

          <el-dropdown @command="handleUserCommand">
            <div class="user-info">
              <el-avatar :size="32" class="user-avatar">
                {{ auth.user?.nickname?.[0] || '管' }}
              </el-avatar>
              <span class="user-name">{{ auth.user?.nickname || '管理员' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人设置</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- 页面内容 -->
      <main class="content-area">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  DataAnalysis, UserFilled, Monitor, Wallet, Promotion,
  Fold, Expand, Bell, ArrowDown
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const sidebarCollapsed = ref(false)

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => (route.meta.title as string) || '')

async function handleUserCommand(command: string) {
  if (command === 'logout') {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    auth.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }
}
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: var(--sidebar-width);
  background: var(--color-bg-card);
  border-right: 1px solid var(--color-border);
  transition: width 0.25s ease;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow: hidden;
}

.sidebar.collapsed {
  width: var(--sidebar-collapsed-width);
}

.sidebar-header {
  height: var(--topbar-height);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 0 20px;
  border-bottom: 1px solid var(--color-divider);
  font-weight: 700;
}

.logo-icon {
  font-size: 24px;
  color: var(--color-primary);
  flex-shrink: 0;
}

.logo-text {
  font-size: 18px;
  color: var(--color-text-title);
  white-space: nowrap;
}

.sidebar-menu {
  flex: 1;
  border-right: none !important;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background-color: var(--color-primary-light);
  color: var(--color-primary);
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.topbar {
  height: var(--topbar-height);
  background: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  flex-shrink: 0;
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  padding: 8px;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.notice-badge {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--radius-btn);
  transition: background 0.15s;
}

.user-info:hover {
  background: var(--color-primary-light);
}

.user-name {
  font-size: 14px;
  color: var(--color-text-title);
}

.content-area {
  flex: 1;
  padding: var(--space-xl);
  overflow-y: auto;
  background: var(--color-bg-page);
}
</style>
