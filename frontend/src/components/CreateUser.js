import React, { useState } from 'react';
import { Form, Button, Alert } from 'react-bootstrap';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const CreateUser = () => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('user');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();

    const newUser = {
      username,
      email,
      password,
      role, 
    };

    try {
      const response = await axios.post('/api/admin/users', newUser);
      console.log('User created successfully:', response.data);
      setSuccess(true);
      setError('');
      setUsername('');
      setEmail('');
      setPassword('');
      setRole('user');
      setTimeout(() => navigate('/admin'), 2000); 
    } catch (error) {
      console.error('Error creating user:', error);
      setError('Error during user creation. Please try again.');
      setSuccess(false);
    }
  };

  return (
    <div className="create-user-wrapper">
      {error && <Alert variant="danger" className="mb-2">{error}</Alert>}
      {success && <Alert variant="success" className="mb-2">User created successfully!</Alert>}
      <Form className="shadow p-4 bg-white rounded" onSubmit={handleSubmit}>
        <div className="h4 mb-2 text-center">Create New User</div>
        <Form.Group className="mb-2" controlId="Username">
          <Form.Label>Username</Form.Label>
          <Form.Control
            type="text"
            value={username}
            placeholder="Username"
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </Form.Group>
        <Form.Group className="mb-2" controlId="Email">
          <Form.Label>Email</Form.Label>
          <Form.Control
            type="email"
            value={email}
            placeholder="Email"
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </Form.Group>
        <Form.Group className="mb-2" controlId="Password">
          <Form.Label>Password</Form.Label>
          <Form.Control
            type="password"
            value={password}
            placeholder="Password"
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </Form.Group>
        <Form.Group className="mb-2" controlId="Role">
          <Form.Label>Role</Form.Label>
          <Form.Control
            as="select"
            value={role}
            onChange={(e) => setRole(e.target.value)}
            required
          >
            <option value="user">User</option>
            <option value="admin">Admin</option>
          </Form.Control>
        </Form.Group>
        <Button className="w-100" variant="primary" type="submit">
          Create User
        </Button>
      </Form>
    </div>
  );
};

export default CreateUser;