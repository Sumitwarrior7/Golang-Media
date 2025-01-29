import React, { useState } from 'react';
import api from '../utils/Api';
import { User, Mail, Calendar } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import Button from './Button';

export const UserCard = ({ user }) => {
  const navigate = useNavigate();

  const handleSubmit = () => {
    console.log("user :", user)
    navigate(`/user/${user.id}`);
  };
  return (
    <div className="p-5 bg-white shadow-lg rounded-2xl flex flex-col items-start border border-gray-200 hover:shadow-xl transition-shadow">
      <div className="flex items-center space-x-4">
        <div className="bg-gray-100 p-2 rounded-full">
          <User className="text-gray-600 w-6 h-6" />
        </div>
        <h2 className="text-lg font-semibold text-gray-800">{user.username}</h2>
      </div>
      <div className="mt-2 flex items-center space-x-2 text-sm text-gray-600">
        <Mail className="w-4 h-4" />
        <p>{user.email}</p>
      </div>
      <div className="mt-2 flex items-center space-x-2 text-sm text-gray-500">
        <Calendar className="w-4 h-4" />
        <p>Joined: {new Date(user.created_at).toLocaleDateString()}</p>
      </div>
      <Button 
        onClick={handleSubmit} 
        className="mt-4 w-full px-4 py-2 rounded-lg text-white font-medium"
      >
        More
      </Button>
    </div>
  );
};


