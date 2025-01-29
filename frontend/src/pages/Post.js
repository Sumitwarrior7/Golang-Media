import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../utils/Api';
import PostCard from '../components/PostCard';

const PostDetailsPage = () => {
  const { postId } = useParams();
  const [post, setPost] = useState(null);

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const response = await api.get(`/posts/${postId}`);
        console.log("resp :", response.data)
        setPost(response.data.data || null);
      } catch (err) {
        console.error('Error fetching post details:', err.response?.data?.message || err.message);
      }
    };
    fetchPost();
  }, [postId]);

  if (!post) {
    return <p>Loading post details...</p>;
  }

  return (
    <div className="max-w-3xl mx-auto p-6">
      <PostCard post={post} allowActions={true} />
    </div>
  );
};

export default PostDetailsPage;
