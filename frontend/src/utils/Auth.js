import {jwtDecode} from 'jwt-decode';

// Save the token in localStorage
export const saveToken = (token) => {
  localStorage.setItem('token', token);
};

// Remove the token from localStorage
export const removeToken = () => {
  localStorage.removeItem('token');
};

// Get the token from localStorage
export const getToken = () => {
  return localStorage.getItem('token');
};

// Decode the token to get the user information
export const getUserFromToken = () => {
  const token = getToken();
  if (!token) return null;
  console.log("token is :", token)

  try {
    const decoded = jwtDecode(token);
    return decoded;
  } catch (error) {
    console.error('Invalid token', error);
    return null;
  }
};

// Check if the user is logged in
export const isAuthenticated = () => {
  const token = getToken();
  if (!token) return false;

  try {
    jwtDecode(token); // Check if token is valid
    return true;
  } catch (error) {
    console.error('Invalid token', error);
    return false;
  }
};
