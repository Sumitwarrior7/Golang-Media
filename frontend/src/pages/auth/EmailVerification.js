import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../../utils/Api";

const EmailVerification = () => {
  const { token = "" } = useParams(); // Get the token from the URL
  const navigate = useNavigate(); // For redirection

  const handleConfirm = async () => {
    try {
      const response = await api.put(`/users/activate/${token}`); // API call to activate the user
      if (response.status === 200) {
        // Redirect to dashboard on success
        navigate("/dashboard");
      } else {
        alert("Failed to verify email. Please try again.");
      }
    } catch (error) {
      console.error("Email verification error:", error);
      alert("Something went wrong. Please try again later.");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white shadow-md rounded-lg p-8 w-full max-w-md text-center">
        <h1 className="text-2xl font-bold text-gray-800 mb-4">
          Verify Your Email
        </h1>
        <p className="text-gray-600 mb-6">
          Please click the button below to verify your email address and
          activate your account.
        </p>
        <button
          onClick={handleConfirm}
          className="bg-blue-500 hover:bg-blue-600 text-white font-semibold px-6 py-3 rounded-lg shadow transition-all"
        >
          Verify Email
        </button>
      </div>
    </div>
  );
};

export default EmailVerification;
