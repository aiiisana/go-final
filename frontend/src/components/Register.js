import React, { useState } from "react";
import { Form, Button, Alert } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import axios from 'axios';
import "./register.css";

const Register = () => {
    const [inputUsername, setInputUsername] = useState("");
    const [inputEmail, setInputEmail] = useState("");
    const [inputPassword, setInputPassword] = useState("");
    const [error, setError] = useState("");
    const [success, setSuccess] = useState(false);
    const navigate = useNavigate();
  
    const handleLoginRedirect = () => {
      navigate('/');
    };
  
    const handleSubmit = async (event) => {
      event.preventDefault();
  
      const newUser = {
        username: inputUsername, // здесь изменили на username
        email: inputEmail,
        password: inputPassword,
      };
  
      try {
        const response = await axios.post('http://localhost:8080/api/users', newUser);
        console.log("User registered successfully:", response.data);
        setSuccess(true);
        setError("");
      } catch (error) {
        console.error("Registration error:", error);
        setError("Error during registration. Please try again.");
        setSuccess(false);
      }
    };
  
    return (
      <div className="register__wrapper">
        {error && <Alert variant="danger" className="mb-2">{error}</Alert>}
        {success && <Alert variant="success" className="mb-2">Registration successful!</Alert>}
        <Form className="shadow p-4 bg-white rounded" onSubmit={handleSubmit}>
          <div className="h4 mb-2 text-center">Register</div>
          <Form.Group className="mb-2" controlId="Username">
            <Form.Label>Username</Form.Label>
            <Form.Control
              type="text"
              value={inputUsername}
              placeholder="Username"
              onChange={(e) => setInputUsername(e.target.value)}
              required
            />
          </Form.Group>
          <Form.Group className="mb-2" controlId="Email">
            <Form.Label>Email</Form.Label>
            <Form.Control
              type="email"
              value={inputEmail}
              placeholder="Email"
              onChange={(e) => setInputEmail(e.target.value)}
              required
            />
          </Form.Group>
          <Form.Group className="mb-2" controlId="Password">
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              value={inputPassword}
              placeholder="Password"
              onChange={(e) => setInputPassword(e.target.value)}
              required
            />
          </Form.Group>
          <Button className="w-100" variant="primary" type="submit">
            Register
          </Button>
        </Form>
        <div className="d-grid justify-content-end">
          <Button className="text-muted px-0" variant="link" onClick={handleLoginRedirect}>
            LOGIN
          </Button>
        </div>
      </div>
    );
  };
export default Register;