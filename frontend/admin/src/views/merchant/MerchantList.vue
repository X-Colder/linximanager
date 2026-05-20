<template>
  <div class="merchant-list">
    <div class="page-header">
      <h1 class="page-title">商家列表</h1>
    </div>

    <!-- 筛选栏 -->
    <el-card class="filter-card">
      <el-form :model="filterForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="filterForm.keyword"
            placeholder="店铺名/手机号"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="行业">
          <el-select v-model="filterForm.industry" placeholder="全部" clearable style="width: 120px">
            <el-option label="餐饮" value="catering" />
            <el-option label="零售" value="retail" />
            <el-option label="生鲜" value="fresh" />
            <el-option label="烘焙" value="bakery" />
          </el-select>
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="filterForm.audit_status" placeholder="全部" clearable style="width: 120px">
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="正常" value="active" />
            <el-option label="已冻结" value="frozen" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card" v-loading="loading">
      <el-table :data="list" stripe row-key="id">
        <el-table-column prop="shop_name" label="店铺名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="industry" label="行业" width="90">
          <template #default="{ row }">
            {{ INDUSTRY_MAP[row.industry] || row.industry }}
          </template>
        </el-table-column>
        <el-table-column prop="contact_phone" label="联系电话" width="130" />
        <el-table-column prop="version_plan" label="套餐" width="90">
          <template #default="{ row }">
            <el-tag :type="PLAN_TAG[row.version_plan]" size="small">
              {{ PLAN_MAP[row.version_plan] || row.version_plan }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="audit_status" label="审核" width="100">
          <template #default="{ row }">
            <el-tag :type="AUDIT_TAG[row.audit_status]" size="small">
              {{ AUDIT_MAP[row.audit_status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '正常' : '已冻结' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="入驻时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" size="small" @click="viewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === 'active'"
              text type="warning" size="small"
              @click="toggleFreeze(row, true)"
            >冻结</el-button>
            <el-button
              v-else
              text type="success" size="small"
              @click="toggleFreeze(row, false)"
            >解冻</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @change="loadList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { merchantApi, type Merchant } from '@/api/merchant'

const INDUSTRY_MAP: Record<string, string> = {
  catering: '餐饮', retail: '零售', fresh: '生鲜', bakery: '烘焙'
}
const PLAN_MAP: Record<string, string> = {
  basic: '基础版', pro: '专业版', chain: '连锁版'
}
const PLAN_TAG: Record<string, '' | 'success' | 'warning' | 'danger' | 'info'> = {
  basic: 'info', pro: 'success', chain: 'warning'
}
const AUDIT_MAP: Record<string, string> = {
  pending: '待审核', approved: '已通过', rejected: '已驳回'
}
const AUDIT_TAG: Record<string, '' | 'success' | 'warning' | 'danger' | 'info'> = {
  pending: 'warning', approved: 'success', rejected: 'danger'
}

const loading = ref(false)
const list = ref<Merchant[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const filterForm = reactive({
  keyword: '',
  industry: '',
  audit_status: '',
  status: ''
})

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

async function loadList() {
  loading.value = true
  try {
    const res = await merchantApi.list({
      page: page.value,
      page_size: pageSize.value,
      ...Object.fromEntries(Object.entries(filterForm).filter(([, v]) => v))
    })
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadList()
}

function handleReset() {
  Object.assign(filterForm, { keyword: '', industry: '', audit_status: '', status: '' })
  page.value = 1
  loadList()
}

function viewDetail(_row: Merchant) {
  // TODO: 跳转到详情页或打开 drawer
}

async function toggleFreeze(row: Merchant, frozen: boolean) {
  await ElMessageBox.confirm(
    `确定要${frozen ? '冻结' : '解冻'}商家「${row.shop_name}」吗？`,
    '操作确认',
    { type: 'warning' }
  )
  try {
    await merchantApi.freeze(row.id, frozen)
    ElMessage.success(frozen ? '已冻结' : '已解冻')
    row.status = frozen ? 'frozen' : 'active'
  } catch {
    // 错误已由拦截器处理
  }
}

onMounted(() => loadList())
</script>

<style scoped>
.merchant-list {
  max-width: 1440px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-title);
  margin-bottom: var(--space-lg);
}

.filter-card {
  margin-bottom: var(--space-md);
}

.table-card {
  /* 无额外样式需要 */
}

.pagination-wrap {
  margin-top: var(--space-md);
  display: flex;
  justify-content: flex-end;
}
</style>
