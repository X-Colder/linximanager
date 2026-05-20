<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'

onLaunch(() => {
  console.log('灵犀掌柜商家端启动')
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style>
page {
  background-color: #f9fafb;
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Hiragino Sans GB', sans-serif;
}

/* 通用变量（rpx单位适配750基准）*/
:root {
  --color-primary: #4a6cf7;
  --color-primary-light: #e8edff;
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-danger: #ef4444;
  --color-info: #3b82f6;
  --color-text-title: #1f2937;
  --color-text-body: #4b5563;
  --color-text-hint: #9ca3af;
  --color-border: #e5e7eb;
  --color-divider: #f3f4f6;
  --color-bg-card: #ffffff;
}
</style>
