import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';

const Admin = () => {
  const [users, setUsers] = useState([]);
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api'; // URL из переменной окружения

  const token = localStorage.getItem("authToken"); // или получите из sessionStorage или cookies
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError('');
      try {
        const [usersResponse, productsResponse, categoriesResponse] = await Promise.all([
          axios.get(`${API_URL}/users`),
          axios.get(`${API_URL}/products`),
          axios.get(`${API_URL}/categories`),
        ]);
  
        setUsers(usersResponse.data);
        setProducts(productsResponse.data);
        setCategories(categoriesResponse.data);
      } catch (error) {
        console.error('Error fetching data:', error);
        setError('Error fetching data. Please try again later.');
      } finally {
        setLoading(false);
      }
    };
  
    fetchData();
  }, [API_URL]); // Добавляем API_URL в массив зависимостей

  const handleDeleteUser = async (userId) => {
    try {
      await axios.delete(`${API_URL}/users/${userId}`);
      setUsers(users.filter((user) => user.id !== userId));
    } catch (error) {
      console.error('Error deleting user:', error);
      setError('Error deleting user. Please try again later.');
    }
  };

  const handleDeleteProduct = async (productId) => {
    try {
      await axios.delete(`${API_URL}/products/${productId}`);
      setProducts(products.filter((product) => product.id !== productId));
    } catch (error) {
      console.error('Error deleting product:', error);
      setError('Error deleting product. Please try again later.');
    }
  };

  const handleDeleteCategory = async (categoryId) => {
    try {
      await axios.delete(`${API_URL}/categories/${categoryId}`);
      setCategories(categories.filter((category) => category.id !== categoryId));
    } catch (error) {
      console.error('Error deleting category:', error);
      setError('Error deleting category. Please try again later.');
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="admin-dashboard">
      <h1>Admin Dashboard</h1>
      {error && <div className="error">{error}</div>}

      <div className="section">
        <h2>Users</h2>
        <ul>
        {users.map((user, index) => (
          <li key={`${user.id}-${index}`}>
            <span>{user.username}</span>
            <Link to={`/admin/users/${user.id}`}>Edit</Link>
            <button onClick={() => handleDeleteUser(user.id)}>Delete</button>
          </li>
        ))}
        </ul>
        <Link to="/admin/users/create">Create New User</Link>
      </div>

      <div className="section">
        <h2>Products</h2>
        <ul>
        {products.map((product, index) => (
          <li key={`${product.id}-${index}`}>
            <span>{product.name}</span>
            <Link to={`/admin/products/${product.id}`}>Edit</Link>
            <button onClick={() => handleDeleteProduct(product.id)}>Delete</button>
          </li>
        ))}
        </ul>
        <Link to="/admin/products/create">Add New Product</Link>
      </div>

      <div className="section">
        <h2>Categories</h2>
        <ul>
          {categories.map((category) => (
            <li key={category.id}>
              <span>{category.name}</span>
              <Link to={`/admin/categories/${category.id}`}>Edit</Link>
              <button onClick={() => handleDeleteCategory(category.id)}>Delete</button>
            </li>
          ))}
        </ul>
        <Link to="/admin/categories/create">Add New Category</Link>
      </div>
    </div>
  );
};

export default Admin;