import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Button } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';

const Cart = () => {
  const navigate = useNavigate();
  const [cart, setCart] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const token = localStorage.getItem('authToken');

  // Функция для получения user_id из токена
  const getUserIdFromToken = () => {
    if (!token) return null;
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.user_id;
  };

  const userId = getUserIdFromToken();

  // Функция для получения корзины
  const fetchCart = async () => {
    if (!userId) return;

    try {
      const response = await axios.get(`http://localhost:8080/api/cart/${userId}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      setCart(response.data);
      setLoading(false);
    } catch (err) {
      setError('Error fetching cart data');
      setLoading(false);
    }
  };

  // Функция для создания корзины
  const createCart = async () => {
    if (!userId) return;

    try {
      const response = await axios.post(`http://localhost:8080/api/cart`, {
        user_id: userId,
      }, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      setCart(response.data);
    } catch (err) {
      setError('Error creating cart');
    }
  };

  // Загружаем корзину при монтировании компонента
  useEffect(() => {
    fetchCart();
  }, [userId]);

  // Рендерим компонент
  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <div>
      {cart ? (
        <div>
          <h3>Your Cart</h3>
          <ul>
            {cart.cart_items.map(item => (
              <li key={item.ID}>
                Product ID: {item.product_id}, Quantity: {item.quantity}
              </li>
            ))}
          </ul>
        </div>
      ) : (
        <div>
          <h3>No Cart Found</h3>
          <Button onClick={createCart}>Create Cart</Button>
        </div>
      )}
    </div>
  );
};

export default Cart;