import React, { useState } from 'react';
import { Form, Button, Alert } from 'react-bootstrap';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const CreateCategory = () => {
  const [categoryName, setCategoryName] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();

    const newCategory = {
      name: categoryName,
    };

    try {
      const response = await axios.post('/api/categories', newCategory);
      console.log('Category created successfully:', response.data);
      setSuccess(true);
      setError('');
      setCategoryName('');
      setTimeout(() => navigate('/admin'), 2000);
    } catch (error) {
      console.error('Error creating category:', error);
      setError('Error during category creation. Please try again.');
      setSuccess(false);
    }
  };

  return (
    <div className="create-category-wrapper">
      {error && <Alert variant="danger" className="mb-2">{error}</Alert>}
      {success && <Alert variant="success" className="mb-2">Category created successfully!</Alert>}
      <Form className="shadow p-4 bg-white rounded" onSubmit={handleSubmit}>
        <div className="h4 mb-2 text-center">Create New Category</div>
        <Form.Group className="mb-2" controlId="CategoryName">
          <Form.Label>Category Name</Form.Label>
          <Form.Control
            type="text"
            value={categoryName}
            placeholder="Category Name"
            onChange={(e) => setCategoryName(e.target.value)}
            required
          />
        </Form.Group>
        <Button className="w-100" variant="primary" type="submit">
          Create Category
        </Button>
      </Form>
    </div>
  );
};

export default CreateCategory;