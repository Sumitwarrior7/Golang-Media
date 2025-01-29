import React from "react";

const CheckEmail = () => {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white shadow-md rounded-lg p-8 w-full max-w-md text-center">
        <h1 className="text-2xl font-bold text-gray-800 mb-4">
          Check Your Email
        </h1>
        <p className="text-gray-600 mb-6">
          A confirmation email has been sent to your inbox. Please check your email to verify your account and activate it.
        </p>
        <p className="text-sm text-gray-500">
          Didn't receive the email? Check your spam folder or try resending the email.
        </p>
      </div>
    </div>
  );
};

export default CheckEmail;
