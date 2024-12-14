import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Button } from 'react-bootstrap';
import axios from 'axios';
import './profile.css';

const Profile = () => {
  const navigate = useNavigate();

  const [userData, setUserData] = useState({
    username: '',
    email: '',
    address: '',
    id: null, 
  });

  const [isEditing, setIsEditing] = useState(false); 
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null); 

  useEffect(() => {
    const fetchProfileData = async () => {
      try {
        const token = localStorage.getItem('authToken');

        const response = await axios.get('http://localhost:8080/api/profile', {
          headers: {
            'Authorization': `Bearer ${token}`,
          }
        });
        const profileData = response.data;
        const userId = profileData.user_id || profileData.ID;
        setUserData({
          ...profileData,
          id: userId,
        });
        setLoading(false);
      } catch (err) {
        console.error('Error fetching profile data:', err);
        setError('Error fetching profile data');
        setLoading(false);
      }
    };

    fetchProfileData();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setUserData({
      ...userData,
      [name]: value,
    });
  };

  const handleSave = async (e) => {
    e.preventDefault();
    try {
      const updatedData = {
        username: userData.username,
        email: userData.email,
        address: userData.address,
      };

      console.log('Saving profile data:', updatedData); 

      const response = await axios.put(`http://localhost:8080/api/users/${userData.id}`, updatedData, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
          'Content-Type': 'application/json'
        },
      });

      console.log('Profile updated:', response.data);
      setIsEditing(false);
    } catch (err) {
      console.error('Error saving profile data:', err.response ? err.response.data : err.message);
      setError('Error saving profile data');
    }
};

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    navigate('/login');
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <div className="profile-container">
      <h3 className="profile-title">Profile</h3>
      <div className="profile-info">
        {isEditing ? (
          <Form onSubmit={handleSave}>
            <Form.Group controlId="username">
              <Form.Label>Username</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter username"
                name="username"
                value={userData.username}
                onChange={handleChange}
              />
            </Form.Group>

            <Form.Group controlId="email">
              <Form.Label>Email</Form.Label>
              <Form.Control
                type="email"
                placeholder="Enter email"
                name="email"
                value={userData.email}
                onChange={handleChange}
              />
            </Form.Group>

            <Form.Group controlId="address">
              <Form.Label>Address</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter address"
                name="address"
                value={userData.address}
                onChange={handleChange}
              />
            </Form.Group>

            <Button variant="primary" type="submit">
              Save Changes
            </Button>
          </Form>
        ) : (
          <div className="profile-details">
            <p><strong>Username:</strong> {userData.username}</p>
            <p><strong>Email:</strong> {userData.email}</p>
            <p><strong>Address:</strong> {userData.address}</p>
            <Button variant="secondary" onClick={() => setIsEditing(true)}>
              Edit Profile
            </Button>
          </div>
        )}
      </div>
      <div className="profile-actions">
        <Button variant="danger" onClick={handleLogout}>
          Logout
        </Button>
      </div>
    </div>
  );
};

export default Profile;