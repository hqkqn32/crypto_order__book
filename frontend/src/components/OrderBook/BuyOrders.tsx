import type { Order } from '../../types/order';
interface BuyOrdersProps {
  orders: Order[];
}

export const BuyOrders = ({ orders }: BuyOrdersProps) => {
  return (
    <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
      <h3 className="text-xl font-bold text-green-500 mb-4">Buy Orders</h3>
      
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="text-gray-400 border-b border-gray-700">
              <th className="text-left py-2">Price</th>
              <th className="text-left py-2">Amount</th>
              <th className="text-left py-2">Total</th>
              <th className="text-left py-2">User</th>
            </tr>
          </thead>
          <tbody>
            {orders.length === 0 ? (
              <tr>
                <td colSpan={4} className="text-center text-gray-500 py-4">
                  No buy orders
                </td>
              </tr>
            ) : (
              orders.map((order) => (
                <tr key={order.id} className="border-b border-gray-700 hover:bg-gray-700">
                  <td className="py-2 text-green-400">${order.price.toFixed(2)}</td>
                  <td className="py-2 text-white">{order.amount.toFixed(8)}</td>
                  <td className="py-2 text-white">${(order.price * order.amount).toFixed(2)}</td>
                  <td className="py-2 text-gray-400">{order.username}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
