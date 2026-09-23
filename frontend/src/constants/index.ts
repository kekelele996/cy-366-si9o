// 与后端 internal/constants 对应的业务枚举（新增枚举值需前后端同步 ≥10 处文件）。
export const STATION_STATUS = {
  IDLE: 'idle',
  USING: 'using',
  FAULT: 'fault',
  RESERVED: 'reserved',
} as const

export const STATION_STATUS_TEXT: Record<string, string> = {
  idle: '空闲',
  using: '使用中',
  fault: '故障',
  reserved: '已预约',
}

export const STATION_STATUS_TYPE: Record<string, string> = {
  idle: 'success',
  using: 'primary',
  fault: 'danger',
  reserved: 'warning',
}

export const RESERVATION_STATUS = {
  PENDING: 'pending',
  WAITLISTED: 'waitlisted',
  CONFIRMED: 'confirmed',
  CHECKED_IN: 'checked_in',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
} as const

export const RESERVATION_STATUS_TEXT: Record<string, string> = {
  pending: '待确认',
  waitlisted: '候补中',
  confirmed: '已确认',
  checked_in: '已开机',
  completed: '已完成',
  cancelled: '已取消',
}

export const RESERVATION_STATUS_TYPE: Record<string, string> = {
  pending: 'warning',
  waitlisted: 'warning',
  confirmed: 'primary',
  checked_in: 'primary',
  completed: 'success',
  cancelled: 'default',
}

export const SESSION_STATUS = {
  ACTIVE: 'active',
  COMPLETED: 'completed',
} as const

export const SESSION_STATUS_TEXT: Record<string, string> = {
  active: '进行中',
  completed: '已结束',
}

export const TOURNAMENT_STATUS = {
  DRAFT: 'draft',
  OPEN: 'open',
  READY: 'ready',
  FINISHED: 'finished',
} as const

export const TOURNAMENT_STATUS_TEXT: Record<string, string> = {
  draft: '草稿',
  open: '报名中',
  ready: '已分组',
  finished: '已结束',
}

export const TOURNAMENT_STATUS_TYPE: Record<string, string> = {
  draft: 'default',
  open: 'success',
  ready: 'primary',
  finished: 'default',
}

export const USER_ROLE = {
  ADMIN: 'admin',
  STAFF: 'staff',
  MEMBER: 'member',
} as const

export const USER_ROLE_TEXT: Record<string, string> = {
  admin: '管理员',
  staff: '店员',
  member: '会员',
}

export const GAME_TYPE = {
  LOL: 'lol',
  CSGO: 'csgo',
  KOG: 'kog',
  OTHER: 'other',
} as const

export const GAME_TYPE_TEXT: Record<string, string> = {
  lol: '英雄联盟',
  csgo: 'CSGO',
  kog: '王者荣耀',
  other: '其他',
}

export const PAYMENT_METHOD = {
  BALANCE: 'balance',
  CASH: 'cash',
  WECHAT: 'wechat',
  ALIPAY: 'alipay',
} as const

export const PAYMENT_METHOD_TEXT: Record<string, string> = {
  balance: '余额',
  cash: '现金',
  wechat: '微信',
  alipay: '支付宝',
}

export const REGISTRATION_MODE = {
  SOLO: 'solo',
  TEAM: 'team',
} as const

export const REGISTRATION_MODE_TEXT: Record<string, string> = {
  solo: '个人报名',
  team: '战队报名',
}

export const STATION_TYPE = {
  SEAT: 'seat',
  BOX: 'box',
} as const

export const STATION_TYPE_TEXT: Record<string, string> = {
  seat: '机位',
  box: '包厢',
}

export const AREA_OPTIONS = ['A区', 'B区', '包厢区']
