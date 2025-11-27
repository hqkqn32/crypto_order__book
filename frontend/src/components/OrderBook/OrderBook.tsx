import { useState, useEffect, useCallback } from 'react';
import { orderService } from '../../services/orderService';
import { BuyOrders } from './BuyOrders';
import { SellOrders } from './SellOrders';
import { OrderForm } from './OrderForm';
import { useWebSocket } from '../../hooks/useWebSocket';
import type { Order } from '../../types/order';

export const OrderBook = () => {
  const [buyOrders, setBuyOrders] = useState<Order[]>([]);
  const [sellOrders, setSellOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchOrders = useCallback(async () => {
    try {
      const data = await orderService.getOrderBook();
      setBuyOrders(data.buy_orders || []);
      setSellOrders(data.sell_orders || []);
    } catch (error) {
      console.error('Failed to fetch orders:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchOrders();
  }, [fetchOrders]);

  const handleNewOrder = useCallback((order: Order) => {
    if (order.order_type === 'buy') {
      setBuyOrders((prev) => {
        // Eğer order zaten varsa ekleme
        if (prev.some(o => o.id === order.id)) return prev;
        return [order, ...prev];
      });
    } else {
      setSellOrders((prev) => {
        // Eğer order zaten varsa ekleme
        if (prev.some(o => o.id === order.id)) return prev;
        return [order, ...prev];
      });
    }
  }, []);

  useWebSocket(handleNewOrder);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-white text-2xl">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 py-8">
      <div className="container mx-auto px-4">
        <h1 className="text-4xl font-bold text-white mb-8 text-center">
          Crypto Orderbook
        </h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
          <div className="lg:col-span-2">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <BuyOrders orders={buyOrders} />
              <SellOrders orders={sellOrders} />
            </div>
          </div>
          
          <div>
            <OrderForm onOrderCreated={() => {}} />
          </div>
        </div>
      </div>
    </div>
  );
};