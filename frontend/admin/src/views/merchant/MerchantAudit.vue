<template>
  <div class="merchant-audit">
    <h1 class="page-title">入驻审核</h1>

    <div class="audit-layout">
      <!-- 左侧：待审核列表 -->
      <el-card class="audit-list-card" v-loading="listLoading">
        <template #header>
          <div class="list-header">
            <el-tabs v-model="activeTab" @tab-click="handleTabChange">
              <el-tab-pane label="待审核" name="pending">
                <template #label>
                  待审核
                  <el-badge v-if="pendingCount > 0" :value="pendingCount" class="tab-badge" />
                </template>
              </el-tab-pane>
              <el-tab-pane label="已处理" name="handled" />
            </el-tabs>
          </div>
        </template>

        <div class="merchant-list">
          <div
            v-for="item in list"
            :key="item.id"
            class="merchant-item"
            :class="{ active: selectedId === item.id }"
            @click="selectMerchant(item.id)"
          >
            <el-avatar :size="44" class="merchant-avatar">
              {{ item.shop_name[0] }}
            </el-avatar>
            <div class="merchant-info">
              <div class="merchant-name">{{ item.shop_name }}</div>
              <div class="merchant-meta">
                {{ INDUSTRY_MAP[item.industry] }} · {{ formatDate(item.created_at) }}
              </div>
            </div>
            <el-tag :type="AUDIT_TAG[item.audit_status]" size="small">
              {{ AUDIT_MAP[item.audit_status] }}
            </el-tag>
          </div>

          <el-empty v-if="!list.length && !listLoading" description="暂无数据" />
        </div>
      </el-card>

      <!-- 右侧：详情面板 -->
      <el-card class="audit-detail-card" v-loading="detailLoading">
        <template v-if="detail">
          <div class="detail-header">
            <el-avatar :size="64">{{ detail.shop_name[0] }}</el-avatar>
            <div class="detail-title">
              <h2>{{ detail.shop_name }}</h2>
              <p>{{ detail.address }}</p>
            </div>
          </div>

          <el-descriptions :column="2" border class="detail-descriptions">
            <el-descriptions-item label="联系电话">{{ detail.contact_phone }}</el-descriptions-item>
            <el-descriptions-item label="行业分类">{{ INDUSTRY_MAP[detail.industry] }}</el-descriptions-item>
            <el-descriptions-item label="套餐版本">{{ PLAN_MAP[detail.version_plan] }}</el-descriptions-item>
            <el-descriptions-item label="申请时间">{{ formatDate(detail.created_at) }}</el-descriptions-item>
          </el-descriptions>

          <div class="license-section">
            <h3 class="section-title">证照资料</h3>
            <div class="license-images">
              <div class="license-item">
                <p class="license-label">营业执照</p>
                <el-image
                  :src="detail.license_url || ''"
                  :preview-src-list="[detail.license_url]"
                  fit="cover"
                  class="license-img"
                >
                  <template #error>
                    <div class="img-placeholder">暂无图片</div>
                  </template>
                </el-image>
              </div>
              <div class="license-item">
                <p class="license-label">身份证</p>
                <el-image
                  :src="detail.id_card_url || ''"
                  :preview-src-list="[detail.id_card_url]"
                  fit="cover"
                  class="license-img"
                >
                  <template #error>
                    <div class="img-placeholder">暂无图片</div>
                  </template>
                </el-image>
              </div>
            </div>
          </div>

          <!-- 审核操作 -->
          <div v-if="detail.audit_status === 'pending'" class="audit-actions">
            <el-input
              v-model="auditRemark"
              placeholder="驳回原因（通过时可不填）"
              type="textarea"
              :rows="3"
              class="remark-input"
            />
            <div class="action-btns">
              <el-button
                type="success"
                size="large"
                :loading="approving"
                @click="handleAudit('approved')"
              >审核通过</el-button>
              <el-button
                type="danger"
                size="large"
                :loading="rejecting"
                @click="handleAudit('rejected')"
              >驳回申请</el-button>
            </div>
          </div>

          <div v-else class="audit-result">
            <el-alert
              :type="detail.audit_status === 'approved' ? 'success' : 'error'"
              :title="detail.audit_status === 'approved' ? '审核已通过' : '审核已驳回'"
              :description="detail.audit_remark"
              show-icon
              :closable="false"
            />
          </div>
        </template>

        <el-empty v-else description="请从左侧选择商家查看详情" />
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { merchantApi, type Merchant, type MerchantDetail } from '@/api/merchant'

const INDUSTRY_MAP: Record<string, string> = {
  catering: '餐饮', retail: '零售', fresh: '生鲜', bakery: '烘焙'
}
const PLAN_MAP: Record<string, string> = {
  basic: '基础版', pro: '专业版', chain: '连锁版'
}
const AUDIT_MAP: Record<string, string> = {
  pending: '待审核', approved: '已通过', rejected: '已驳回'
}
const AUDIT_TAG: Record<string, '' | 'success' | 'warning' | 'danger'> = {
  pending: 'warning', approved: 'success', rejected: 'danger'
}

const activeTab = ref<'pending' | 'handled'>('pending')
const listLoading = ref(false)
const detailLoading = ref(false)
const approving = ref(false)
const rejecting = ref(false)

const list = ref<Merchant[]>([])
const total = ref(0)
const pendingCount = ref(0)
const selectedId = ref<number | null>(null)
const detail = ref<MerchantDetail | null>(null)
const auditRemark = ref('')

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('zh-CN')
}

async function loadList() {
  listLoading.value = true
  try {
    const audit_status = activeTab.value === 'pending' ? 'pending' : undefined
    const res = await merchantApi.list({ audit_status, page: 1, page_size: 50 })
    list.value = res.list
    total.value = res.total
    if (activeTab.value === 'pending') pendingCount.value = res.total
  } finally {
    listLoading.value = false
  }
}

async function selectMerchant(id: number) {
  selectedId.value = id
  detailLoading.value = true
  try {
    detail.value = await merchantApi.detail(id)
    auditRemark.value = ''
  } finally {
    detailLoading.value = false
  }
}

function handleTabChange() {
  list.value = []
  detail.value = null
  selectedId.value = null
  loadList()
}

async function handleAudit(action: 'approved' | 'rejected') {
  if (action === 'rejected' && !auditRemark.value.trim()) {
    ElMessage.warning('驳回时请填写原因')
    return
  }
  if (!detail.value) return

  const flag = action === 'approved' ? approving : rejecting
  flag.value = true
  try {
    await merchantApi.audit(detail.value.id, {
      audit_status: action,
      audit_remark: auditRemark.value
    })
    ElMessage.success(action === 'approved' ? '已通过审核' : '已驳回申请')
    detail.value.audit_status = action
    loadList()
  } finally {
    flag.value = false
  }
}

onMounted(() => loadList())
</script>

<style scoped>
.merchant-audit { max-width: 1440px; }

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-title);
  margin-bottom: var(--space-lg);
}

.audit-layout {
  display: grid;
  grid-template-columns: 360px 1fr;
  gap: var(--space-lg);
  align-items: start;
}

.audit-list-card { height: calc(100vh - 200px); overflow: hidden; }

.list-header :deep(.el-tabs__nav-wrap::after) { display: none; }

.merchant-list { overflow-y: auto; max-height: calc(100vh - 280px); }

.merchant-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.15s;
}

.merchant-item:hover { background: var(--color-primary-light); }
.merchant-item.active { background: var(--color-primary-light); }

.merchant-avatar { background: var(--color-primary); flex-shrink: 0; }

.merchant-info { flex: 1; min-width: 0; }

.merchant-name {
  font-weight: 600;
  color: var(--color-text-title);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.merchant-meta { font-size: 12px; color: var(--color-text-hint); margin-top: 2px; }

.detail-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--color-divider);
}

.detail-title h2 { font-size: 20px; font-weight: 600; color: var(--color-text-title); }
.detail-title p  { font-size: 14px; color: var(--color-text-hint); margin-top: 4px; }

.detail-descriptions { margin-bottom: 24px; }

.section-title { font-size: 16px; font-weight: 600; color: var(--color-text-title); margin-bottom: 12px; }

.license-images { display: flex; gap: 24px; }

.license-label { font-size: 12px; color: var(--color-text-hint); margin-bottom: 8px; }

.license-img {
  width: 200px;
  height: 130px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.img-placeholder {
  width: 200px;
  height: 130px;
  background: var(--color-divider);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-hint);
  border-radius: 8px;
}

.audit-actions {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--color-divider);
}

.remark-input { margin-bottom: 16px; }

.action-btns { display: flex; gap: 16px; }
.action-btns .el-button { flex: 1; }

.audit-result { margin-top: 24px; }

.tab-badge { margin-left: 6px; }
</style>
