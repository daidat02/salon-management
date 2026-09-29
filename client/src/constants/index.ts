export const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE ||
  process.env.NEXT_PUBLIC_API_URL ||
  'http://localhost:8080/api/v1';

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: `${API_BASE}/login`,
    REGISTER: `${API_BASE}/register`,
    FORGOT_PASSWORD: `${API_BASE}/forgot-password`,
  },

  CUSTOMERS: {
    GET_CUSTOMERS: `${API_BASE}/customers`,
  },

  PRODUCTS: {
    GET_PRODUCTS: `${API_BASE}/products`,
  },

  SERVICES: {
    GET_SERVICES: `${API_BASE}/services`,
  },

  CATEGORIES: {
    GET_CATEGORIES: `${API_BASE}/categories`,
  },

  INVENTORY: {
    GET_DOCUMENTS: `${API_BASE}/inventory/`,
    GET_TRANSACTIONS: `${API_BASE}/inventory/transactions`,
  },

  ORDERS: {
    GET_ORDERS: `${API_BASE}/orders`,
    GET_POS_ITEMS: `${API_BASE}/orders/pos-items`,
  },
};
