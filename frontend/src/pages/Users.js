import React, { useState, useEffect } from 'react';
import api from '../utils/Api';
import { UserSearchBar } from '../components/UserSearchBar';
import { UserCard } from '../components/UserCard';

const Users = () => {
  const [users, setUsers] = useState([]);
  const [query, setQuery] = useState('');
  const [offset, setOffset] = useState(0);
  const [limit] = useState(12); // Number of users per page
  const [userCount, setUserCount] = useState(0);
  const testLimit = limit+1;

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await api.get('/all-users', {
          params: { 
            search: query, 
            offset: offset, 
            limit: testLimit
          },
        });

        console.log("Total Users :", response.data.data)
        // If there are no users or search term dont contain any user
        if(!response.data.data) {
          setUsers([])
          setUserCount(0); 
          return
        }

        // Create a new array up to the limit
        var limitedUsers = [];
        for (let i = 0; i < Math.min(limit, response.data.data.length); i++) {
          limitedUsers.push(response.data.data[i]);
        }
        setUsers(limitedUsers);
        setUserCount(response.data.data.length); 
      } catch (error) {
        console.error('Error fetching users:', error);
      }
    };
    console.log("offset:", offset, "userCount:", userCount);
    fetchUsers();
  }, [query, offset, limit]);

  const handleNext = () => {
    if (userCount === testLimit) {
      setOffset(offset + limit);
    }
  };

  const handlePrevious = () => {
    if (offset - limit >= 0) {
      setOffset(offset - limit);
    }
  };

  return (
    <div className="p-4 min-h-screen bg-gray-50 flex flex-col items-center">
      <h1 className="text-2xl font-semibold mb-4 text-gray-800">User Management</h1>
      <UserSearchBar onSearch={setQuery} />
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6 w-full max-w-4xl">
        {users.map((user) => (
          <UserCard key={user.id} user={user} />
        ))}
      </div>
      <div className="mt-6 flex justify-between w-full max-w-4xl">
        <button
          onClick={handlePrevious}
          disabled={offset === 0}
          className="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-400"
        >
          Previous
        </button>
        <button
          onClick={handleNext}
          disabled={userCount < testLimit}
          className="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-400"
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default Users;
