<template>
  <view class="page">
    <!-- 店铺信息头部 -->
    <view class="shop-header">
      <image class="shop-avatar" :src="shopInfo.logo_url || ''" mode="aspectFill" />
      <view class="shop-info">
        <text class="shop-name">{{ shopInfo.shop_name }}</text>
        <text class="shop-plan">{{ PLAN_MAP[shopInfo.version_plan] || '基础版' }}</text>
      </view>
    </view>

    <!-- 菜单列表 -->
    <view class="menu-section">
      <view
        v-for="item in menuItems"
        :key="item.label"
        class="menu-item"
        @tap="item.action()"
      >
        <text class="menu-icon">{{ item.icon }}</text>
        <text class="menu-label">{{ item.label }}</text>
        <view class="menu-right">
          <text v-if="item.badge" class="menu-badge">{{ item.badge }}</text>
          <text class="menu-arrow">›</text>
        </view>
      </view>
    </view>

    <!-- 退出登录 -->
    <button class="logout-btn" @tap="handleLogout">退出登录</button>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import http from '@/api/request'

const PLAN_MAP: Record<string, string> = {
  basic: '基础版', pro: '专业版', chain: '连锁版'
}

interface ShopInfo {
  shop_name: string
  logo_url: string
  version_plan: string
  plan_expire_at: string
}

const shopInfo = ref<ShopInfo>({
  shop_name: '我的店铺',
  logo_url: '',
  version_plan: 'basic',
  plan_expire_at: ''
})

const menuItems = [
  { icon: '🏪', label: '店铺设置', action: () => uni.navigateTo({ url: '/pages/shop-settings/index' }) },
  { icon: '🤝', label: '流量联盟', badge: '', action: () => uni.navigateTo({ url: '/pages/alliance/index' }) },
  { icon: '💳', label: '我的订阅', badge: '', action: () => uni.navigateTo({ url: '/pages/subscription/index' }) },
  { icon: '💰', label: '收入统计', action: () => uni.navigateTo({ url: '/pages/income/index' }) },
  { icon: '📱', label: '我的小程序', action: () => uni.navigateTo({ url: '/pages/miniapp/index' }) },
  { icon: '🔔', label: '通知设置', action: () => uni.navigateTo({ url: '/pages/notification-settings/index' }) },
  { icon: '❓', label: '帮助中心', action: () => uni.navigateTo({ url: '/pages/help/index' }) }
]

function handleLogout() {
  uni.showModal({
    title: '提示',
    content: '确定要退出登录吗？',
    success: ({ confirm }) => {
      if (confirm) {
        uni.removeStorageSync('access_token')
        uni.removeStorageSync('refresh_token')
        uni.reLaunch({ url: '/pages/login/index' })
      }
    }
  })
}

onMounted(async () => {
  try {
    shopInfo.value = await http.get<ShopInfo>('/merchant/shop')
  } catch {
    // 使用默认值
  }
})
</script>

<style scoped>
.page { min-height: 100vh; background: #f9fafb; }

.shop-header {
  background: linear-gradient(135deg, #4a6cf7, #3451d1);
  padding: 48rpx 32rpx;
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.shop-avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background: rgba(255,255,255,0.2);
  border: 4rpx solid rgba(255,255,255,0.4);
}

.shop-name  { display: block; font-size: 34rpx; font-weight: 700; color: #fff; }
.shop-plan  { display: block; font-size: 24rpx; color: rgba(255,255,255,0.7); margin-top: 6rpx; }

.menu-section {
  background: #fff;
  margin: 24rpx;
  border-radius: 20rpx;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f3f4f6;
  gap: 20rpx;
}

.menu-item:last-child { border-bottom: none; }

.menu-icon  { font-size: 36rpx; flex-shrink: 0; }
.menu-label { flex: 1; font-size: 28rpx; color: #1f2937; }

.menu-right { display: flex; align-items: center; gap: 8rpx; }

.menu-badge {
  background: #ef4444;
  color: #fff;
  font-size: 20rpx;
  padding: 2rpx 10rpx;
  border-radius: 100rpx;
}

.menu-arrow { font-size: 32rpx; color: #9ca3af; }

.logout-btn {
  margin: 0 24rpx 40rpx;
  height: 88rpx;
  background: #fff;
  color: #ef4444;
  font-size: 30rpx;
  border-radius: 20rpx;
  border: none;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.04);
}
</style>
