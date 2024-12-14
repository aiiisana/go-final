import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Register from './components/Register';
import Login from './components/Login';
import ProductsList from './components/ProductsList';
import Product from './components/Product';
import Category from './components/Category';
import Cart from './components/Cart';
// import Order from './components/Order';
// import Review from './components/Review';
// import Address from './components/Address';
import Admin from './components/Admin';
import CreateCategory from './components/CreateCategory';
import CreateProducts from './components/CreateProducts';
import CreateUser from './components/CreateUser';
import UserList from './components/UserList'; 
import Categories from './components/Categories';
import Header from './components/Header';
import Profile from './components/Profile';


function App() {
  return (
    <Router>
      <Header/>
      <Routes>
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/products" element={<ProductsList />} />
        <Route path="/products/:productId" element={<Product />} />
        <Route path="/categories" element={<Categories />} />
        <Route path="/category/:categoryId" element={<Category />} />
        <Route path="/cart" element={<Cart />} />
        {/* <Route path="/orders" element={<Order />} />
        <Route path="/reviews" element={<Review />} />
        <Route path="/address" element={<Address />} /> */}
        <Route path="/admin" element={<Admin />} />
        <Route path="/admin/categories/create" element={<CreateCategory />} />
        <Route path="/admin/products/create" element={<CreateProducts />} />
        <Route path="/admin/users/create" element={<CreateUser />} />

        <Route path="/admin/users" element={<UserList />} />
        <Route path="/" element={<Login />} />
      </Routes>
    </Router>
  );
}

export default App;