<template>
  <view class="page">
    <!-- Tab 切换 -->
    <view class="tab-bar">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: currentTab === tab.value }"
        @tap="currentTab = tab.value"
      >
        {{ tab.label }}
      </view>
    </view>

    <!-- 补货建议 -->
    <scroll-view v-if="currentTab === 'replenishment'" scroll-y class="tab-content">
      <view class="info-bar">
        <text class="info-text">预测周期: 未来7天 · 周末晴天预计客流↑</text>
      </view>

      <view
        v-for="item in replenishItems"
        :key="item.product_id"
        class="replenish-card"
      >
        <view class="card-header">
          <view class="checkbox" :class="{ checked: item.selected }" @tap="item.selected = !item.selected">
            <text v-if="item.selected" class="check-icon">✓</text>
          </view>
          <text class="product-name">{{ item.product_name }}</text>
          <text v-if="item.urgency === 'high'" class="urgency-tag">紧急</text>
        </view>
        <view class="card-detail">
          <view class="detail-row">
            <text class="detail-key">当前库存</text>
            <text class="detail-val">{{ item.current_stock }}{{ item.unit }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-key">日均需求</text>
            <text class="detail-val">{{ item.daily_demand_min }}-{{ item.daily_demand_max }}{{ item.unit }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-key">安全库存</text>
            <text class="detail-val">{{ item.safety_stock }}{{ item.unit }}</text>
          </view>
          <view class="detail-row highlight">
            <text class="detail-key">建议补货</text>
            <text class="detail-val primary">→ {{ item.suggested_qty }}{{ item.unit }}</text>
          </view>
        </view>
        <view v-if="item.remark" class="card-remark">💡 {{ item.remark }}</view>
      </view>

      <!-- 底部汇总 -->
      <view class="summary-bar">
        <text class="summary-text">已选 {{ selectedCount }} 项 · 预计采购额: ¥{{ estimatedCost }}</text>
        <button class="summary-btn" @tap="generateOrder">生成采购单</button>
      </view>
    </scroll-view>

    <!-- 促销方案 -->
    <scroll-view v-if="currentTab === 'promo'" scroll-y class="tab-content">
      <view class="promo-target">
        <text class="promo-target__name">🧃 果汁（500ml）</text>
        <text class="promo-target__info">当前库存: 45瓶 | 正常周转: 20瓶 | 积压2.25倍</text>
      </view>
      <view class="promo-hint">💡 AI为你生成了 {{ promoPlans.length }} 个方案：</view>
      <PromoCompare
        :plans="promoPlans"
        :recommend-index="1"
        @select="handlePromoSelect"
      />
    </scroll-view>

    <!-- AI顾问 -->
    <view v-if="currentTab === 'chat'" class="tab-content chat-page">
      <scroll-view scroll-y class="chat-messages" :scroll-into-view="lastMsgId">
        <view
          v-for="msg in chatMessages"
          :key="msg.id"
          :id="`msg-${msg.id}`"
          class="chat-bubble"
          :class="msg.role === 'user' ? 'bubble--user' : 'bubble--ai'"
        >
          <text class="bubble-text">{{ msg.content }}</text>
        </view>
        <view v-if="aiThinking" class="chat-bubble bubble--ai">
          <text class="ai-thinking">AI思考中...</text>
        </view>
      </scroll-view>
      <view class="chat-input-bar">
        <input
          v-model="chatInput"
          class="chat-input"
          placeholder="问问AI顾问..."
          confirm-type="send"
          @confirm="sendMessage"
        />
        <button class="send-btn" @tap="sendMessage">发送</button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import PromoCompare, { type PromoPlan } from '@/components/PromoCompare.vue'
import http from '@/api/request'

const tabs = [
  { label: '补货建议', value: 'replenishment' },
  { label: '促销方案', value: 'promo' },
  { label: 'AI顾问', value: 'chat' }
]

const currentTab = ref('replenishment')

// 补货建议数据
interface ReplenishItem {
  product_id: number
  product_name: string
  current_stock: number
  daily_demand_min: number
  daily_demand_max: number
  safety_stock: number
  suggested_qty: number
  unit: string
  urgency: 'high' | 'normal' | 'low'
  remark?: string
  selected: boolean
}

const replenishItems = ref<ReplenishItem[]>([])
const selectedCount = computed(() => replenishItems.value.filter(i => i.selected).length)
const estimatedCost = computed(() => {
  // 实际应从后端获取单价计算
  return replenishItems.value.filter(i => i.selected).length * 120
})

// 促销方案
const promoPlans = ref<PromoPlan[]>([
  { name: '方案A - 保利润', discount: 0.9, originalPrice: 4.5, promoPrice: 4.0, estimatedDays: 14, estimatedProfit: 500 },
  { name: '方案B', label: '速度与利润最佳平衡', discount: 0.7, originalPrice: 4.5, promoPrice: 3.15, estimatedDays: 3, estimatedProfit: 200 },
  { name: '方案C - 快速清仓', discount: 0.5, originalPrice: 4.5, promoPrice: 2.25, estimatedDays: 1, estimatedProfit: -50, reason: '⚠️ 仅建议临期最后清仓' }
])

// AI顾问
interface ChatMessage {
  id: number
  role: 'user' | 'ai'
  content: string
}

const chatMessages = ref<ChatMessage[]>([
  { id: 1, role: 'ai', content: '你好，我是灵犀AI顾问。根据你的库存和销售数据，我可以帮你分析经营情况、给出补货和促销建议，有什么想问的吗？' }
])
const chatInput = ref('')
const aiThinking = ref(false)
const lastMsgId = ref('')

async function sendMessage() {
  const text = chatInput.value.trim()
  if (!text) return
  chatInput.value = ''

  const userMsg: ChatMessage = { id: Date.now(), role: 'user', content: text }
  chatMessages.value.push(userMsg)
  lastMsgId.value = `msg-${userMsg.id}`

  aiThinking.value = true
  try {
    const res = await http.post<{ reply: string }>('/merchant/ai/chat', { message: text })
    chatMessages.value.push({ id: Date.now() + 1, role: 'ai', content: res.reply })
    lastMsgId.value = `msg-${Date.now() + 1}`
  } catch {
    chatMessages.value.push({ id: Date.now() + 1, role: 'ai', content: '抱歉，AI服务暂时不可用，请稍后再试。' })
  } finally {
    aiThinking.value = false
  }
}

function handlePromoSelect(plan: PromoPlan) {
  uni.showModal({
    title: '确认执行促销',
    content: `确定执行「${plan.name}」方案？执行后将自动更新店铺价格并加入联盟曝光池。`,
    success: ({ confirm }) => {
      if (confirm) {
        http.post('/merchant/ai/promotions/execute', { plan }).then(() => {
          uni.showToast({ title: '促销已启动', icon: 'success' })
        })
      }
    }
  })
}

function generateOrder() {
  const selected = replenishItems.value.filter(i => i.selected)
  if (!selected.length) {
    uni.showToast({ title: '请选择补货商品', icon: 'none' })
    return
  }
  uni.showToast({ title: '采购单已生成', icon: 'success' })
}

async function loadReplenishment() {
  try {
    const res = await http.get<ReplenishItem[]>('/merchant/ai/replenishment')
    replenishItems.value = res.map(i => ({ ...i, selected: true }))
  } catch {
    replenishItems.value = [
      { product_id: 1, product_name: '鸡翅', current_stock: 2, daily_demand_min: 8, daily_demand_max: 12, safety_stock: 2, suggested_qty: 13, unit: 'kg', urgency: 'high', remark: '周末烧烤需求高，建议多备', selected: true },
      { product_id: 2, product_name: '青岛啤酒', current_stock: 24, daily_demand_min: 30, daily_demand_max: 40, safety_stock: 24, suggested_qty: 192, unit: '瓶', urgency: 'high', remark: '鸡翅缺货可能影响啤酒销量', selected: true },
      { product_id: 3, product_name: '竹签', current_stock: 500, daily_demand_min: 50, daily_demand_max: 80, safety_stock: 100, suggested_qty: 0, unit: '根', urgency: 'low', selected: false }
    ]
  }
}

onMounted(() => loadReplenishment())
</script>

<style scoped>
.page { min-height: 100vh; background: #f9fafb; display: flex; flex-direction: column; }

.tab-bar {
  display: flex;
  background: #fff;
  padding: 0 24rpx;
  border-bottom: 1rpx solid #f3f4f6;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  font-size: 28rpx;
  color: #9ca3af;
  border-bottom: 4rpx solid transparent;
  transition: all 0.2s;
}

.tab-item.active {
  color: #4a6cf7;
  font-weight: 600;
  border-bottom-color: #4a6cf7;
}

.tab-content {
  flex: 1;
  padding: 24rpx;
  height: calc(100vh - 120rpx);
}

/* 补货 */
.info-bar {
  background: #eff6ff;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  margin-bottom: 20rpx;
}
.info-text { font-size: 24rpx; color: #3b82f6; }

.replenish-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.04);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 16rpx;
}

.checkbox {
  width: 40rpx;
  height: 40rpx;
  border-radius: 8rpx;
  border: 2rpx solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.checkbox.checked { background: #4a6cf7; border-color: #4a6cf7; }
.check-icon { font-size: 24rpx; color: #fff; }

.product-name { font-size: 30rpx; font-weight: 600; color: #1f2937; flex: 1; }

.urgency-tag {
  font-size: 22rpx;
  color: #ef4444;
  background: #fef2f2;
  padding: 4rpx 12rpx;
  border-radius: 100rpx;
}

.card-detail {
  background: #f9fafb;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 8rpx 0;
  border-bottom: 1rpx solid #f3f4f6;
}

.detail-row:last-child { border-bottom: none; }
.detail-row.highlight { background: #e8edff; margin: 0 -20rpx -16rpx; padding: 12rpx 20rpx; border-radius: 0 0 12rpx 12rpx; }

.detail-key { font-size: 26rpx; color: #9ca3af; }
.detail-val { font-size: 26rpx; color: #1f2937; font-weight: 500; }
.detail-val.primary { color: #4a6cf7; font-weight: 700; }

.card-remark {
  font-size: 24rpx;
  color: #4b5563;
  margin-top: 12rpx;
  padding: 10rpx 16rpx;
  background: #fffbeb;
  border-radius: 10rpx;
}

.summary-bar {
  position: sticky;
  bottom: 0;
  background: #fff;
  padding: 20rpx 24rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1rpx solid #f3f4f6;
  box-shadow: 0 -4rpx 16rpx rgba(0,0,0,0.06);
  margin: 0 -24rpx;
}

.summary-text { font-size: 26rpx; color: #4b5563; }

.summary-btn {
  background: #4a6cf7;
  color: #fff;
  border-radius: 100rpx;
  font-size: 26rpx;
  padding: 14rpx 32rpx;
  border: none;
}

/* 促销 */
.promo-target {
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
}
.promo-target__name { display: block; font-size: 30rpx; font-weight: 700; color: #1f2937; }
.promo-target__info { display: block; font-size: 24rpx; color: #9ca3af; margin-top: 8rpx; }
.promo-hint { font-size: 26rpx; color: #4b5563; margin-bottom: 16rpx; font-weight: 500; }

/* 聊天 */
.chat-page { display: flex; flex-direction: column; padding: 0 !important; }

.chat-messages {
  flex: 1;
  padding: 24rpx;
  height: calc(100vh - 230rpx);
}

.chat-bubble {
  max-width: 80%;
  padding: 20rpx 24rpx;
  border-radius: 20rpx;
  margin-bottom: 20rpx;
  font-size: 28rpx;
  line-height: 1.6;
}

.bubble--ai {
  background: #fff;
  color: #1f2937;
  align-self: flex-start;
  border-bottom-left-radius: 4rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06);
}

.bubble--user {
  background: #4a6cf7;
  color: #fff;
  margin-left: auto;
  border-bottom-right-radius: 4rpx;
}

.ai-thinking { color: #9ca3af; }

.chat-input-bar {
  display: flex;
  gap: 16rpx;
  padding: 16rpx 24rpx;
  background: #fff;
  border-top: 1rpx solid #f3f4f6;
}

.chat-input {
  flex: 1;
  height: 72rpx;
  background: #f3f4f6;
  border-radius: 36rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
  color: #1f2937;
}

.send-btn {
  background: #4a6cf7;
  color: #fff;
  border-radius: 36rpx;
  font-size: 26rpx;
  padding: 0 28rpx;
  height: 72rpx;
  border: none;
}
</style>
