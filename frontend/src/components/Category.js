import React, { useState, useEffect } from "react";
import axios from "axios";
import { useParams, Link, useNavigate } from "react-router-dom";
import "./category.css";

const Category = () => {
  const { categoryId } = useParams();
  const navigate = useNavigate();
  const [category, setCategory] = useState(null);

  useEffect(() => {
    const fetchCategory = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/api/categories/${categoryId}`);
        setCategory(response.data);
      } catch (error) {
        console.error("Error fetching category details:", error);
      }
    };
    fetchCategory();
  }, [categoryId]);

  if (!category) {
    return <div>Loading...</div>;
  }

  return (
    <div className="category-details-container">
      <h3 className="category-details-title">Category Details</h3>
      <div className="category-details-card">
        <h4>{category.name}</h4>
        <p><strong>Description:</strong> {category.description}</p>
        {category.Products && category.Products.length > 0 && (
          <div className="category-products">
            <h5>Products in this category:</h5>
            <ul>
              {category.Products.map((product) => (
                <li key={product.ID}>
                  <Link to={`/products/${product.ID}`}>{product.name}</Link>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
      <button className="back-button" onClick={() => navigate(-1)}>
        ← Back
      </button>
    </div>
  );
};

export default Category;