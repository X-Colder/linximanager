<template>
  <view class="page">
    <view v-if="cartItems.length" class="cart-content">
      <!-- 商品列表 -->
      <view class="shop-group">
        <view class="shop-name-row">
          <view class="select-all" @tap="toggleSelectAll">
            <view class="checkbox" :class="{ checked: allSelected }">
              <text v-if="allSelected" class="check-icon">✓</text>
            </view>
            <text class="select-label">全选</text>
          </view>
          <text class="shop-label">{{ shopName }}</text>
        </view>

        <view v-for="item in cartItems" :key="item.id" class="cart-item">
          <view class="checkbox" :class="{ checked: item.selected }" @tap="item.selected = !item.selected">
            <text v-if="item.selected" class="check-icon">✓</text>
          </view>

          <image class="item-image" :src="item.image_url || ''" mode="aspectFill" />

          <view class="item-info">
            <text class="item-name">{{ item.name }}</text>
            <text v-if="item.spec" class="item-spec">{{ item.spec }}</text>
            <view class="item-price-row">
              <text class="item-price">¥{{ item.price }}</text>
              <view class="quantity-ctrl">
                <view class="qty-btn" @tap="decreaseQty(item)">-</view>
                <text class="qty-num">{{ item.quantity }}</text>
                <view class="qty-btn qty-btn--plus" @tap="increaseQty(item)">+</view>
              </view>
            </view>
          </view>
        </view>
      </view>

      <!-- 优惠券 -->
      <view class="coupon-section" @tap="selectCoupon">
        <text class="coupon-label">优惠券</text>
        <view class="coupon-right">
          <text v-if="selectedCoupon" class="coupon-discount">-¥{{ couponDiscount }}</text>
          <text v-else class="coupon-hint">{{ availableCoupons }}张可用</text>
          <text class="coupon-arrow">›</text>
        </view>
      </view>

      <!-- 费用汇总 -->
      <view class="summary-section">
        <view class="summary-row">
          <text class="summary-key">商品合计</text>
          <text class="summary-val">¥{{ subtotal.toFixed(1) }}</text>
        </view>
        <view class="summary-row" v-if="couponDiscount > 0">
          <text class="summary-key">优惠券</text>
          <text class="summary-val summary-val--discount">-¥{{ couponDiscount }}</text>
        </view>
        <view class="summary-row summary-row--total">
          <text class="summary-key">应付金额</text>
          <text class="summary-val summary-val--total">¥{{ payAmount.toFixed(1) }}</text>
        </view>
      </view>

      <!-- 取货时间 -->
      <view class="pickup-section" @tap="selectPickupTime">
        <text class="pickup-label">⏰ 自提时间</text>
        <view class="pickup-right">
          <text class="pickup-time">{{ pickupTimeText }}</text>
          <text class="pickup-edit">修改 ›</text>
        </view>
      </view>
    </view>

    <!-- 空购物车 -->
    <view v-else class="empty-cart">
      <text class="empty-icon">🛒</text>
      <text class="empty-text">购物车还是空的</text>
      <button class="go-shop-btn" @tap="goShopping">去逛逛</button>
    </view>

    <!-- 底部结算栏 -->
    <view v-if="cartItems.length" class="checkout-bar">
      <view class="select-count">
        <text>全选({{ selectedItems.length }})</text>
      </view>
      <button
        class="checkout-btn"
        :disabled="!selectedItems.length"
        @tap="checkout"
      >
        去结算 ¥{{ payAmount.toFixed(1) }}
      </button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import http from '@/api/request'

interface CartItem {
  id: number
  name: string
  image_url: string
  price: number
  quantity: number
  spec?: string
  selected: boolean
}

const shopName = ref('老王烧烤')
const cartItems = ref<CartItem[]>([])
const selectedCoupon = ref<{ value: number; min_amount: number } | null>(null)
const availableCoupons = ref(0)
const pickupTimeText = ref('今天 18:00-18:30')

const selectedItems = computed(() => cartItems.value.filter(i => i.selected))
const allSelected = computed(() => cartItems.value.length > 0 && cartItems.value.every(i => i.selected))
const subtotal = computed(() => selectedItems.value.reduce((s, i) => s + i.price * i.quantity, 0))
const couponDiscount = computed(() => {
  if (!selectedCoupon.value || subtotal.value < selectedCoupon.value.min_amount) return 0
  return selectedCoupon.value.value
})
const payAmount = computed(() => Math.max(0, subtotal.value - couponDiscount.value))

function toggleSelectAll() {
  const target = !allSelected.value
  cartItems.value.forEach(i => { i.selected = target })
}

function increaseQty(item: CartItem) {
  item.quantity++
  saveCart()
}

function decreaseQty(item: CartItem) {
  if (item.quantity <= 1) {
    uni.showModal({
      title: '提示', content: '确定删除该商品？',
      success: ({ confirm }) => {
        if (confirm) {
          cartItems.value = cartItems.value.filter(i => i.id !== item.id)
          saveCart()
        }
      }
    })
  } else {
    item.quantity--
    saveCart()
  }
}

function saveCart() {
  uni.setStorageSync('cart', JSON.stringify(cartItems.value))
}

function selectCoupon() {
  uni.navigateTo({ url: '/pages/coupons/index?mode=select' })
}

function selectPickupTime() {
  uni.showActionSheet({
    itemList: ['今天 18:00-18:30', '今天 19:00-19:30', '明天 11:00-11:30'],
    success: ({ tapIndex }) => {
      const options = ['今天 18:00-18:30', '今天 19:00-19:30', '明天 11:00-11:30']
      pickupTimeText.value = options[tapIndex]
    }
  })
}

function goShopping() {
  uni.switchTab({ url: '/pages/index/index' })
}

async function checkout() {
  if (!selectedItems.value.length) return

  const orderData = {
    items: selectedItems.value.map(i => ({
      product_id: i.id,
      quantity: i.quantity,
      unit_price: i.price
    })),
    coupon_id: selectedCoupon.value ? 1 : undefined,
    pickup_time: pickupTimeText.value
  }

  try {
    uni.showLoading({ title: '创建订单...' })
    const res = await http.post<{ order_id: number; pay_url: string }>('/orders', orderData, true)
    // 清空已结算商品
    cartItems.value = cartItems.value.filter(i => !i.selected)
    saveCart()
    uni.navigateTo({ url: `/pages/orders/index?orderId=${res.order_id}` })
  } catch {
    // 错误已处理
  }
}

function loadCart() {
  const saved = uni.getStorageSync('cart')
  if (saved) {
    cartItems.value = JSON.parse(saved)
  }
}

onMounted(() => loadCart())
onShow(() => loadCart())
</script>

<style scoped>
.page { min-height: 100vh; background: #f9fafb; padding-bottom: 120rpx; }

.cart-content { padding: 24rpx; }

.shop-group {
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
}

.shop-name-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 20rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid #f3f4f6;
}

.select-all { display: flex; align-items: center; gap: 10rpx; }

.checkbox {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 2rpx solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkbox.checked { background: #4a6cf7; border-color: #4a6cf7; }
.check-icon { font-size: 22rpx; color: #fff; }

.select-label { font-size: 26rpx; color: #4b5563; }
.shop-label   { flex: 1; font-size: 28rpx; font-weight: 600; color: #1f2937; }

.cart-item {
  display: flex;
  align-items: flex-start;
  gap: 16rpx;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f3f4f6;
}

.cart-item:last-child { border-bottom: none; }

.item-image {
  width: 120rpx;
  height: 120rpx;
  border-radius: 12rpx;
  background: #f3f4f6;
  flex-shrink: 0;
}

.item-info { flex: 1; min-width: 0; }
.item-name { display: block; font-size: 28rpx; color: #1f2937; font-weight: 500; }
.item-spec { display: block; font-size: 24rpx; color: #9ca3af; margin-top: 4rpx; }

.item-price-row { display: flex; align-items: center; justify-content: space-between; margin-top: 12rpx; }
.item-price { font-size: 30rpx; font-weight: 700; color: #1f2937; }

.quantity-ctrl { display: flex; align-items: center; gap: 16rpx; }

.qty-btn {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  border: 2rpx solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  color: #4b5563;
}

.qty-btn--plus { background: #4a6cf7; border-color: #4a6cf7; color: #fff; }
.qty-num { font-size: 28rpx; font-weight: 600; color: #1f2937; min-width: 40rpx; text-align: center; }

/* 优惠券 */
.coupon-section {
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx 28rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.coupon-label { font-size: 28rpx; color: #1f2937; }
.coupon-right { display: flex; align-items: center; gap: 8rpx; }
.coupon-discount { font-size: 26rpx; color: #ef4444; font-weight: 600; }
.coupon-hint     { font-size: 26rpx; color: #9ca3af; }
.coupon-arrow    { font-size: 32rpx; color: #9ca3af; }

/* 费用汇总 */
.summary-section {
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx 28rpx;
  margin-bottom: 20rpx;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  padding: 10rpx 0;
  border-bottom: 1rpx solid #f3f4f6;
}

.summary-row:last-child { border-bottom: none; }

.summary-row--total { padding-top: 16rpx; }

.summary-key { font-size: 26rpx; color: #4b5563; }
.summary-val { font-size: 26rpx; color: #1f2937; }
.summary-val--discount { color: #ef4444; font-weight: 600; }
.summary-val--total    { font-size: 32rpx; font-weight: 700; color: #1f2937; }

/* 取货时间 */
.pickup-section {
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx 28rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.pickup-label { font-size: 28rpx; color: #1f2937; }
.pickup-right { display: flex; align-items: center; gap: 12rpx; }
.pickup-time  { font-size: 26rpx; color: #4b5563; }
.pickup-edit  { font-size: 24rpx; color: #4a6cf7; }

/* 空购物车 */
.empty-cart {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 120rpx;
  gap: 20rpx;
}

.empty-icon { font-size: 96rpx; }
.empty-text { font-size: 28rpx; color: #9ca3af; }

.go-shop-btn {
  margin-top: 16rpx;
  background: #4a6cf7;
  color: #fff;
  border-radius: 100rpx;
  font-size: 28rpx;
  padding: 16rpx 48rpx;
  border: none;
}

/* 结算栏 */
.checkout-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  padding: 16rpx 32rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1rpx solid #f3f4f6;
  box-shadow: 0 -4rpx 16rpx rgba(0,0,0,0.06);
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
}

.select-count { font-size: 26rpx; color: #4b5563; }

.checkout-btn {
  background: #4a6cf7;
  color: #fff;
  border-radius: 100rpx;
  font-size: 28rpx;
  font-weight: 600;
  padding: 16rpx 40rpx;
  border: none;
}

.checkout-btn[disabled] {
  background: #d1d5db;
}
</style>
