import React, { useState } from "react";
import { Form, Button, Alert } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import "./login.css";

const Login = () => {
  const [inputUsername, setInputUsername] = useState("");
  const [inputPassword, setInputPassword] = useState("");
  const [showError, setShowError] = useState(false);
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    setLoading(true);

    try {
      // Шаг 1: Отправка данных для входа
      const loginResponse = await axios.post("http://localhost:8080/api/login", {
        username: inputUsername,
        password: inputPassword,
      });

      console.log("Login successful:", loginResponse.data);

      const token = loginResponse.data.token;

      localStorage.setItem("authToken", token);

      // Шаг 2: Получение информации о пользователе
      const userResponse = await axios.get("http://localhost:8080/api/profile", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      console.log("User info:", userResponse.data);

      // Извлечение роли из ответа
      const { role } = userResponse.data;

      // Шаг 3: Перенаправление на соответствующую страницу
      if (role === "admin") {
        navigate("/admin");
      } else {
        navigate("/products");
      }
    } catch (error) {
      console.error("Login error:", error.response ? error.response.data : error.message);
      setShowError(true);
    } finally {
      setLoading(false);
    }
  };

  const handleRegisterRedirect = () => {
    navigate('/register');
  };

  return (
    <div className="sign-in__wrapper">
      <Form className="shadow p-4 bg-white rounded" onSubmit={handleSubmit}>
        <div className="h4 mb-2 text-center">Sign In</div>
        {showError && (
          <Alert
            className="mb-2"
            variant="danger"
            onClose={() => setShowError(false)}
            dismissible
          >
            Incorrect Username or Password.
          </Alert>
        )}
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
        <Button className="w-100" variant="primary" type="submit" disabled={loading}>
          {loading ? "Logging In..." : "Log In"}
        </Button>
        <div className="d-grid justify-content-end">
          <Button className="text-muted px-0" variant="link" onClick={handleRegisterRedirect}>
            Create an account
          </Button>
        </div>
      </Form>
      <div className="text-center mt-3 text-muted">Golang Midterm 2024</div>
    </div>
  );
};

export default Login;