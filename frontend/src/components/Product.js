import React, { useState, useEffect } from "react";
import axios from "axios";
import { Link, useParams, useNavigate } from "react-router-dom";
import "./product.css";

const Product = () => {
    const { productId } = useParams();
    const navigate = useNavigate();
    const [product, setProduct] = useState(null);
  
    useEffect(() => {
      const fetchProduct = async () => {
        try {
          const response = await axios.get(`http://localhost:8080/api/products/${productId}`);
          setProduct(response.data);
        } catch (error) {
          console.error('Error fetching product details:', error);
        }
      };
      fetchProduct();
    }, [productId]);
  
    if (!product) {
      return <div>Loading...</div>;
    }
  
    return (
        <div className="product-details-container">
          <h3 className="product-details-title">Product Details</h3>
          <div className="product-details-card">
            <img
              src={product.image || "https://via.placeholder.com/300"}
              alt={product.name}
              className="product-details-img"
            />
            <div className="product-details-info">
              <h4>{product.name}</h4>
              <p><strong>Description:</strong> {product.description}</p>
              <p><strong>Price:</strong> ${product.price}</p>
              <p><strong>Stock:</strong> {product.stock}</p>
              <p>
                <strong>Category:</strong> 
                <Link to={`/category/${product.category_id}`}>View Category</Link>
              </p>
            </div>
          </div>
          <button className="back-button" onClick={() => navigate(-1)}>
            ← Back
          </button>
        </div>
      );
  };

export default Product;