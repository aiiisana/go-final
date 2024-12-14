import React, { useState, useEffect } from "react";
import axios from "axios";
import { Link } from "react-router-dom";
import "./categories.css";

const Categories = () => {
  const [categories, setCategories] = useState([]);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await axios.get("http://localhost:8080/api/categories");
        setCategories(response.data);
      } catch (error) {
        console.error("Error fetching categories:", error);
      }
    };
    fetchCategories();
  }, []);

  if (categories.length === 0) {
    return <div>Loading...</div>;
  }

  return (
    <div className="categories-container">
      <h3 className="categories-title">Categories</h3>
      <div className="categories-card-container">
        {categories.map((category) => (
          <div key={category.ID} className="category-card">
            <h4>{category.name}</h4>
            <p>{category.description}</p>
            <Link to={`/category/${category.ID}`} className="view-category-link">
              View Category →
            </Link>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Categories;