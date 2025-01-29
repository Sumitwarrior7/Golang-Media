import React, { createContext, useState, useEffect } from "react";
import { getUserFromToken, saveToken, removeToken, isAuthenticated } from "../utils/Auth";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);

  // Initialize user state when the app loads
  useEffect(() => {
    if (isAuthenticated()) {
      const userInfo = getUserFromToken();
      setUser(userInfo);
    }
  }, []);

  const login = (token) => {
    saveToken(token);
    const userInfo = getUserFromToken();
    setUser(userInfo);
  };

  const logout = () => {
    removeToken();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => React.useContext(AuthContext);

export default AuthContext;
