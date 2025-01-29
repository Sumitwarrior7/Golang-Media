import React from 'react';

const Input = ({ label, value, onChange, type = 'text', className, ...props }) => {
  return (
    <div className="mb-4">
      {label && <label className="block text-gray-700 mb-1">{label}</label>}
      <input
        type={type}
        value={value}
        onChange={onChange}
        className={`w-full p-2 border rounded ${className}`}
        {...props}
      />
    </div>
  );
};

export default Input;
