import React, { useState, useEffect } from "react";
import axios from "axios";
import { Card, Button } from "react-bootstrap";
import { Link } from "react-router-dom";
import "./products_list.css";  // Импортируем наш файл CSS

const Products = () => {
  const [products, setProducts] = useState([]);

  useEffect(() => {
    const fetchProducts = async () => {
      const response = await axios.get("http://localhost:8080/api/products");
      setProducts(response.data);
    };
    fetchProducts();
  }, []);

  return (
    <div className="products-container">
      <h3 className="products-title">Products</h3>
      <div className="card-container">
        {products.map((product) => (
          <Card key={product.ID} className="product-card">
            <Card.Img variant="top" src={product.image} className="product-img" />
            <Card.Body>
              <Card.Title className="product-title">{product.name}</Card.Title>
              <Card.Text>{product.description}</Card.Text>
              <Link to={`/products/${product.ID}`}>
                <Button className="product-button">View Details</Button>
              </Link>
            </Card.Body>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default Products;