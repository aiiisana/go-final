import React from 'react';
import { Link } from 'react-router-dom';
import { Navbar, Nav, Dropdown } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import "./header.css";  // Импортируем файл стилей

const Header = () => {
  const navigate = useNavigate();

  // Функция для выхода
  const handleLogout = () => {
    // Логика выхода (например, удаление токена из локального хранилища)
    localStorage.removeItem('authToken');
    navigate('/login');
  };

  return (
    <Navbar expand="lg">
      <Navbar.Collapse id="basic-navbar-nav">
        <Nav className="ml-auto navbar-nav">
          <Nav.Item>
            <Link to="/products" className="nav-link">
              Products
            </Link>
          </Nav.Item>
          <Nav.Item>
            <Link to="/categories" className="nav-link">
              Categories
            </Link>
          </Nav.Item>
          <Nav.Item>
            <Link to="/cart" className="nav-link">
              Cart
            </Link>
          </Nav.Item>

          <Dropdown align="end">
            <Dropdown.Toggle variant="link" id="dropdown-account">
              Account
            </Dropdown.Toggle>
            <Dropdown.Menu>
              <Dropdown.Item as={Link} to="/profile">
                View Profile
              </Dropdown.Item>
              <Dropdown.Item as="button" onClick={handleLogout}>
                Logout
              </Dropdown.Item>
            </Dropdown.Menu>
          </Dropdown>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
};

export default Header;