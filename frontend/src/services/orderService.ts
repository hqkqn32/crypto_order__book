import api from './api';

import type { CreateOrderRequest, OrderBook } from '../types/order';
export const orderService = {
  getOrderBook: async (): Promise<OrderBook> => {
    const response = await api.get<OrderBook>('/orders');
    return response.data;
  },

  createOrder: async (data: CreateOrderRequest) => {
    const response = await api.post('/orders', data);
    return response.data;
  },

  getMyOrders: async () => {
    const response = await api.get('/orders/my');
    return response.data;
  },
};
