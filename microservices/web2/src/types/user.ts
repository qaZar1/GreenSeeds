export type UserRole = 'admin' | 'operator';

export interface User {
  id: number;
  username: string;
  full_name: string;
  is_admin: boolean;
}

export type UserInput = {
  username: string;
  full_name: string;
  is_admin: boolean;
};

export const ROLE_LABELS = {
  true: 'Администратор',
  false: 'Оператор',
} as const;

export const ROLE_COLORS = {
  true: 'var(--status-warning-text)',
  false: 'var(--text-secondary)',
} as const;

export const ROLE_BG = {
  true: 'var(--status-warning-bg)',
  false: 'var(--bg-hover)',
} as const;