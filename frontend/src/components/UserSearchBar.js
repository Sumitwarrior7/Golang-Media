import React, { useState } from 'react';

export const UserSearchBar = ({ onSearch }) => {
  const [inputValue, setInputValue] = useState('');

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      onSearch(inputValue);
    }
  };

  const handleSearchClick = () => {
    onSearch(inputValue);
  };

  return (
    <div className="w-full max-w-2xl mb-4 flex items-center gap-2">
      <input
        type="text"
        placeholder="Search for users..."
        value={inputValue}
        className="w-full p-3 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        onChange={handleInputChange}
        onKeyDown={handleKeyDown}
      />
      <button
        onClick={handleSearchClick}
        className="px-4 py-2 bg-blue-500 text-white rounded-lg shadow-sm hover:bg-blue-600"
      >
        Search
      </button>
    </div>
  );
};
