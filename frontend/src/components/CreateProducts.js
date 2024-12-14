import React, { useState } from 'react';
import { Form, Button, Alert } from 'react-bootstrap';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const CreateProduct = () => {
  const [productName, setProductName] = useState('');
  const [productPrice, setProductPrice] = useState('');
  const [productDescription, setProductDescription] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();

    const newProduct = {
      name: productName,
      price: parseFloat(productPrice),
      description: productDescription,
    };

    try {
      const response = await axios.post(`http://localhost:8080/api/products`, newProduct);
      console.log('Product created successfully:', response.data);
      setSuccess(true);
      setError('');
      setProductName('');
      setProductPrice('');
      setProductDescription('');
      setTimeout(() => navigate('/admin'), 2000);
    } catch (error) {
      console.error('Error creating product:', error);
      setError('Error during product creation. Please try again.');
      setSuccess(false);
    }
  };

  return (
    <div className="create-product-wrapper">
      {error && <Alert variant="danger" className="mb-2">{error}</Alert>}
      {success && <Alert variant="success" className="mb-2">Product created successfully!</Alert>}

      <Form className="shadow p-4 bg-white rounded" onSubmit={handleSubmit}>
        <div className="h4 mb-2 text-center">Create New Product</div>
        
        <Form.Group className="mb-2" controlId="ProductName">
          <Form.Label>Product Name</Form.Label>
          <Form.Control
            type="text"
            value={productName}
            placeholder="Product Name"
            onChange={(e) => setProductName(e.target.value)}
            required
          />
        </Form.Group>

        <Form.Group className="mb-2" controlId="ProductPrice">
          <Form.Label>Price</Form.Label>
          <Form.Control
            type="number"
            value={productPrice}
            placeholder="Price"
            onChange={(e) => setProductPrice(e.target.value)}
            required
          />
        </Form.Group>

        <Form.Group className="mb-2" controlId="ProductDescription">
          <Form.Label>Description</Form.Label>
          <Form.Control
            as="textarea"
            value={productDescription}
            placeholder="Description"
            onChange={(e) => setProductDescription(e.target.value)}
            required
          />
        </Form.Group>

        <Button className="w-100" variant="primary" type="submit">
          Create Product
        </Button>
      </Form>
    </div>
  );
};

export default CreateProduct;