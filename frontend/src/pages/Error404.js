import React from 'react';
import { Link } from 'react-router-dom';

const Error404 = () => {
  return (
    <div className="flex items-center justify-center h-screen flex-col">
      <h1 className="text-4xl font-bold">404</h1>
      <p className="text-lg mb-4">Page Not Found</p>
      <Link to="/" className="bg-blue-500 text-white px-4 py-2 rounded">
        Go to Home
      </Link>
    </div>
  );
};

export default Error404;
