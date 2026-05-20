import request from './request'
import type { PaginatedData } from './request'

export interface MerchantListParams {
  page?: number
  page_size?: number
  keyword?: string
  industry?: string
  audit_status?: string
  status?: string
}

export interface Merchant {
  id: number
  shop_name: string
  industry: string
  logo_url: string
  address: string
  contact_phone: string
  audit_status: 'pending' | 'approved' | 'rejected'
  version_plan: 'basic' | 'pro' | 'chain'
  plan_expire_at: string
  status: 'active' | 'frozen'
  created_at: string
  user_id: number
}

export interface MerchantDetail extends Merchant {
  license_url: string
  id_card_url: string
  shop_photos: string[]
  audit_remark: string
  miniapp_status: string
}

export interface AuditParams {
  audit_status: 'approved' | 'rejected'
  audit_remark?: string
}

export interface DashboardStats {
  total_merchants: number
  active_merchants: number
  pending_audit: number
  total_orders_today: number
  total_gmv_today: number
  merchant_growth: number
  order_growth: number
}

export const merchantApi = {
  list: (params: MerchantListParams) =>
    request.get<unknown, PaginatedData<Merchant>>('/admin/merchants', { params }),

  detail: (id: number) =>
    request.get<unknown, MerchantDetail>(`/admin/merchants/${id}`),

  audit: (id: number, params: AuditParams) =>
    request.put(`/admin/merchants/${id}/audit`, params),

  freeze: (id: number, frozen: boolean) =>
    request.put(`/admin/merchants/${id}/freeze`, { frozen }),

  stats: (id: number) =>
    request.get(`/admin/merchants/${id}/stats`),

  dashboard: () =>
    request.get<unknown, DashboardStats>('/admin/dashboard')
}
