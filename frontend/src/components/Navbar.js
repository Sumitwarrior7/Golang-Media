import React, { useContext, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import AuthContext from '../context/AuthContext';
import { Users as UsersIcon } from 'lucide-react'; // Importing Lucide icon for users

const Navbar = () => {
  const navigate = useNavigate();
  const authContext = useContext(AuthContext);
  const [isOpen, setIsOpen] = useState(false);

  if (!authContext) {
    console.error("AuthContext is not provided properly.");
    return (
      <nav className="bg-blue-600 p-3 text-white"> {/* Reduced padding here */}
        <div className="max-w-7xl mx-auto px-2 sm:px-6 lg:px-8">
          <div className="relative flex items-center justify-between h-14"> {/* Reduced height */}
            <div className="flex-1 flex items-center justify-start">
              <Link to="/" className="text-2xl font-bold">
                GolangMedia
              </Link>
            </div>
          </div>
        </div>
      </nav>
    );
  }

  const { user, logout } = authContext;

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="bg-blue-600 p-3 text-white"> {/* Reduced padding here */}
      <div className="max-w-7xl mx-auto px-2 sm:px-6 lg:px-8">
        <div className="relative flex items-center justify-between h-14"> {/* Reduced height */}
          <div className="flex-1 flex items-center justify-start">
            <Link to="/" className="text-2xl font-bold">
              GolangMedia
            </Link>
          </div>
          <div className="absolute inset-y-0 right-0 flex items-center sm:hidden">
            <button
              onClick={() => setIsOpen(!isOpen)}
              className="inline-flex items-center justify-center p-2 rounded-md text-white hover:text-gray-700 hover:bg-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-white"
            >
              <span className="sr-only">Open main menu</span>
              <svg
                className="block h-6 w-6"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M4 6h16M4 12h16M4 18h16"
                />
              </svg>
            </button>
          </div>
          <div className="hidden sm:flex sm:items-center sm:space-x-6">
            {user ? (
              <>
                <Link
                  to="/users"
                  className="flex items-center text-white hover:text-gray-200 space-x-2" 
                >
                  <UsersIcon className="w-5 h-5" /> 
                  <span>Users</span>
                </Link>
                <Link to="/dashboard" className="text-white hover:text-gray-200">
                  Dashboard
                </Link>
                <button
                  onClick={handleLogout}
                  className="bg-red-500 px-4 py-2 rounded text-white hover:bg-red-600"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-white hover:text-gray-200">
                  Login
                </Link>
                <Link
                  to="/register"
                  className="bg-green-500 px-4 py-2 rounded text-white hover:bg-green-600"
                >
                  Register
                </Link>
              </>
            )}
          </div>
        </div>
      </div>

      <div className={`${isOpen ? 'block' : 'hidden'} sm:hidden`}>
        <div className="px-2 pt-2 pb-3 space-y-1">
          {user ? (
            <>
              <Link
                to="/dashboard"
                className="block text-white hover:bg-blue-600 px-3 py-2 rounded-md text-base font-medium"
              >
                Dashboard
              </Link>
              <Link
                to="/users"
                className="block text-white hover:bg-blue-600 px-3 py-2 rounded-md text-base font-medium"
              >
                Users
              </Link>
              <button
                onClick={handleLogout}
                className="w-full text-left bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-md text-base font-medium"
              >
                Logout
              </button>
            </>
          ) : (
            <>
              <Link
                to="/login"
                className="block text-white hover:bg-blue-600 px-3 py-2 rounded-md text-base font-medium"
              >
                Login
              </Link>
              <Link
                to="/register"
                className="block text-white bg-green-500 hover:bg-green-600 px-3 py-2 rounded-md text-base font-medium"
              >
                Register
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
