import { useState } from 'react';
import { orderService } from '../../services/orderService';

interface OrderFormProps {
  onOrderCreated: () => void;
}

export const OrderForm = ({ onOrderCreated }: OrderFormProps) => {
  const [orderType, setOrderType] = useState<'buy' | 'sell'>('buy');
  const [price, setPrice] = useState('');
  const [amount, setAmount] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await orderService.createOrder({
        order_type: orderType,
        price: parseFloat(price),
        amount: parseFloat(amount),
      });
      
      setPrice('');
      setAmount('');
      onOrderCreated();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create order');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
      <h3 className="text-xl font-bold text-white mb-4">Create Order</h3>
      
      {error && (
        <div className="bg-red-500 text-white p-3 rounded mb-4">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-gray-300 mb-2">Order Type</label>
          <div className="flex gap-2">
            <button
              type="button"
              onClick={() => setOrderType('buy')}
              className={`flex-1 py-2 px-4 rounded font-semibold ${
                orderType === 'buy'
                  ? 'bg-green-600 text-white'
                  : 'bg-gray-700 text-gray-300'
              }`}
            >
              Buy
            </button>
            <button
              type="button"
              onClick={() => setOrderType('sell')}
              className={`flex-1 py-2 px-4 rounded font-semibold ${
                orderType === 'sell'
                  ? 'bg-red-600 text-white'
                  : 'bg-gray-700 text-gray-300'
              }`}
            >
              Sell
            </button>
          </div>
        </div>

        <div className="mb-4">
          <label className="block text-gray-300 mb-2">Price</label>
          <input
            type="number"
            step="0.00000001"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            className="w-full px-4 py-2 bg-gray-700 text-white rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0.00"
            required
          />
        </div>

        <div className="mb-6">
          <label className="block text-gray-300 mb-2">Amount</label>
          <input
            type="number"
            step="0.00000001"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="w-full px-4 py-2 bg-gray-700 text-white rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0.00"
            required
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className={`w-full font-bold py-2 px-4 rounded disabled:opacity-50 ${
            orderType === 'buy'
              ? 'bg-green-600 hover:bg-green-700'
              : 'bg-red-600 hover:bg-red-700'
          } text-white`}
        >
          {loading ? 'Creating...' : `${orderType === 'buy' ? 'Buy' : 'Sell'}`}
        </button>
      </form>
    </div>
  );
};
