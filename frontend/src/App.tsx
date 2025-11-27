import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { Navbar } from './components/Layout/Navbar';
import { PrivateRoute } from './components/Layout/PrivateRoute';
import { Login } from './components/Auth/Login';
import { Register } from './components/Auth/Register';
import { OrderBook } from './components/OrderBook/OrderBook';

function App() {
  return (
    <Router>
      <AuthProvider>
        <div className="min-h-screen bg-gray-900">
          <Navbar />
          <Routes>
            <Route path="/" element={<Navigate to="/login" replace />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route
              path="/orderbook"
              element={
                <PrivateRoute>
                  <OrderBook />
                </PrivateRoute>
              }
            />
          </Routes>
        </div>
      </AuthProvider>
    </Router>
  );
}

export default App;
