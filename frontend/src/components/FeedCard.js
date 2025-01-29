import React from "react";
import { useNavigate } from "react-router-dom";
import { MessageCircle, Tag, Eye } from 'lucide-react';

const FeedCard = ({ post }) => {
  const navigate = useNavigate();
  const commentCount = post.CommentCount || 0;

  return (
    <div className="bg-white shadow-lg rounded-lg p-6 mb-6 hover:shadow-xl transition-shadow duration-300">
      {/* Post Title */}
      <h2 className="text-2xl font-bold text-gray-800 mb-4 truncate">{post.Title}</h2>

      {/* Post Content */}
      <p className="text-gray-600 mb-4 line-clamp-3">{post.Content}</p>

      {/* Tags Section */}
      {post.Tags && post.Tags.length > 0 && (
        <div className="mb-4 flex flex-wrap gap-2">
          {post.Tags.map((tag, index) => (
            <span
              key={index}
              className="bg-blue-100 text-blue-600 text-xs font-medium px-2 py-1 rounded-full"
            >
              <Tag className="inline-block w-3 h-3 mr-1" />
              {tag}
            </span>
          ))}
        </div>
      )}

      <div className="flex justify-between items-center">
        {/* Read More Button with Eye Icon */}
        <button
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center gap-2"
          onClick={() => navigate(`/post/${post.Id}`)}
        >
          <Eye className="w-4 h-4" />
          Read More
        </button>

        {/* Comments Section with MessageCircle Icon */}
        <span className="text-sm text-gray-500 flex items-center gap-1">
          <MessageCircle className="w-4 h-4 text-gray-500" />
          {commentCount} {commentCount === 1 ? 'Comment' : 'Comments'}
        </span>
      </div>
    </div>
  );
};

export default FeedCard;
