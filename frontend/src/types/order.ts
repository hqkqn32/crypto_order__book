export interface Order {
  id: number;
  user_id: number;
  username: string;
  order_type: 'buy' | 'sell';
  price: number;
  amount: number;
  status: string;
  created_at: string;
}

export interface CreateOrderRequest {
  order_type: 'buy' | 'sell';
  price: number;
  amount: number;
}

export interface OrderBook {
  buy_orders: Order[];
  sell_orders: Order[];
}