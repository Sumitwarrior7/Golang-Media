import React from "react";
import { Routes, Route, useLocation } from "react-router-dom";
import Login from "./pages/auth/Login";
import Register from "./pages/auth/Register";
import EmailVerification from "./pages/auth/EmailVerification";
import CheckEmail from "./pages/auth/CheckEmail";
import Dashboard from "./pages/Dashboard";
import Feed from "./pages/Feed";
import PublicProfile from "./pages/PublicProfile";
import Error404 from "./pages/Error404";
import Navbar from "./components/Navbar";
import { AuthProvider } from "./context/AuthContext";
import ProtectedRoute from "./components/ProtectedRoute";
import PublicRoute from "./components/PublicRoute";
import "./input.css";
import PostDetailsPage from "./pages/Post";
import Users from "./pages/Users";


function App() {
  const location = useLocation();

  // Routes where Navbar should not be displayed
  const noNavbarRoutes = ["/login", "/register", "/confirm/:token", "/check-email"];
  const hideNavbar = noNavbarRoutes.includes(location.pathname);

  return (
    <AuthProvider>
      {/* Conditionally render Navbar */}
      {!hideNavbar && <Navbar />}
      <Routes>
        {/* Public Routes */}
        <Route
          path="/login"
          element={
            <PublicRoute>
              <Login />
            </PublicRoute>
          }
        />
        <Route
          path="/register"
          element={
            <PublicRoute>
              <Register />
            </PublicRoute>
          }
        />
        <Route
          path="/confirm/:token"
          element={
            <PublicRoute>
              <EmailVerification />
            </PublicRoute>
          }
        />
        <Route
          path="/check-email"
          element={
            <PublicRoute>
              <CheckEmail />
            </PublicRoute>
          }
        />

        {/* Protected Routes */}
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <Feed />
            </ProtectedRoute>
          }
        />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/profile/:id"
          element={
            <ProtectedRoute>
              <PublicProfile />
            </ProtectedRoute>
          }
        />
        <Route
          path="/post/:postId"
          element={
            <ProtectedRoute>
              <PostDetailsPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/users"
          element={
            <ProtectedRoute>
              <Users />
            </ProtectedRoute>
          }
        />
        <Route
          path="/user/:userId"
          element={
            <ProtectedRoute>
              <PublicProfile />
            </ProtectedRoute>
          }
        />

        {/* Catch-All Route */}
        <Route path="*" element={<Error404 />} />
      </Routes>
    </AuthProvider>
  );
}

export default App;
