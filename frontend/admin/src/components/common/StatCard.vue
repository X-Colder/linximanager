<template>
  <div class="stat-card" :class="`stat-card--${color}`">
    <div class="stat-card__header">
      <span class="stat-card__label">{{ label }}</span>
      <el-icon class="stat-card__icon" :size="20">
        <component :is="icon" />
      </el-icon>
    </div>
    <div class="stat-card__value">{{ formattedValue }}</div>
    <div v-if="growth !== undefined" class="stat-card__growth">
      <span :class="growth >= 0 ? 'growth--up' : 'growth--down'">
        {{ growth >= 0 ? '↑' : '↓' }}{{ Math.abs(growth) }}%
      </span>
      <span class="growth-label">较昨日</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'

interface Props {
  label: string
  value: number | string
  icon?: Component
  color?: 'primary' | 'success' | 'warning' | 'danger'
  growth?: number
  prefix?: string
  suffix?: string
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary'
})

const formattedValue = computed(() => {
  const v = props.value
  if (typeof v === 'number') {
    return `${props.prefix || ''}${v.toLocaleString()}${props.suffix || ''}`
  }
  return v
})
</script>

<style scoped>
.stat-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  padding: var(--space-lg);
  border-top: 3px solid transparent;
}

.stat-card--primary { border-top-color: var(--color-primary); }
.stat-card--success { border-top-color: var(--color-success); }
.stat-card--warning { border-top-color: var(--color-warning); }
.stat-card--danger  { border-top-color: var(--color-danger); }

.stat-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.stat-card__label {
  font-size: 14px;
  color: var(--color-text-hint);
}

.stat-card__icon {
  color: var(--color-text-hint);
}

.stat-card__value {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-title);
  line-height: 36px;
}

.stat-card__growth {
  margin-top: 8px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.growth--up   { color: var(--color-success); }
.growth--down { color: var(--color-danger); }
.growth-label { color: var(--color-text-hint); }
</style>
