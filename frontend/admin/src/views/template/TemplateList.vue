<template>
  <div class="template-list">
    <div class="page-header">
      <h1 class="page-title">模板管理</h1>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">新建模板</el-button>
    </div>

    <el-card v-loading="loading">
      <el-table :data="list" stripe>
        <el-table-column prop="name" label="模板名称" min-width="160" />
        <el-table-column prop="industry" label="行业" width="100">
          <template #default="{ row }">{{ INDUSTRY_MAP[row.industry] || row.industry }}</template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '已启用' : '已停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ new Date(row.created_at).toLocaleDateString('zh-CN') }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default>
            <el-button text type="primary" size="small">编辑</el-button>
            <el-button text type="danger" size="small">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showCreateDialog" title="新建模板" width="500px">
      <p class="placeholder-text">模板配置表单（待实现）</p>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'

const INDUSTRY_MAP: Record<string, string> = {
  catering: '餐饮', retail: '零售', fresh: '生鲜', bakery: '烘焙'
}

const loading = ref(false)
const showCreateDialog = ref(false)
const list = ref([
  { id: 1, name: '餐饮基础模板', industry: 'catering', description: '适用于中小型餐饮门店', status: 'active', created_at: '2026-01-01' },
  { id: 2, name: '烘焙专业模板', industry: 'bakery', description: '含BOM物料清单配置', status: 'active', created_at: '2026-01-15' },
  { id: 3, name: '生鲜零售模板', industry: 'fresh', description: '含批次管理和保质期预警', status: 'active', created_at: '2026-02-01' }
])

onMounted(() => {
  // 实际应从接口拉取
})
</script>

<style scoped>
.template-list { max-width: 1440px; }

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-lg);
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-title);
}

.placeholder-text { color: var(--color-text-hint); text-align: center; padding: 24px 0; }
</style>
