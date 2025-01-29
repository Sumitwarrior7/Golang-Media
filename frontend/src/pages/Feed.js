import React, { useState, useEffect } from 'react';
import api from '../utils/Api';
import FeedCard from '../components/FeedCard';
import SearchBar from '../components/SearchBar';
import { PlusCircle } from 'lucide-react';

const Feed = () => {
  const [posts, setPosts] = useState([]);
  const [newPost, setNewPost] = useState({
    title: '',
    content: '',
    tags: '',
  });
  const [isCreatingPost, setIsCreatingPost] = useState(false);

  const fetchPosts = async (query = '') => {
    try {
      const response = await api.get(`/users/feed${query ? `?search=${query}` : ''}`);
      setPosts(response.data.data);
    } catch (err) {
      console.error(err.response?.data?.message || 'Error fetching feed');
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  const handleCreatePost = async () => {
    if (!newPost.title || !newPost.content) {
      alert('Title and Content are required!');
      return;
    }

    try {
      const postData = {
        title: newPost.title,
        content: newPost.content,
        tags: newPost.tags.split(',').map((tag) => tag.trim()),
      };
      const response = await api.post("/posts", postData);
      const newPostData = response.data.data;
      if (newPostData) {
        setPosts([newPostData, ...posts]);
      }
      setNewPost({ title: '', content: '', tags: '' });
      setIsCreatingPost(false);
    } catch (err) {
      console.error(err.response?.data?.message || 'Error creating post');
    }
  };

  return (
    <div className="flex flex-col items-center p-6  bg-gray-50">
      <h1 className="text-3xl font-bold mb-6 text-center">Your Feed</h1>

      {/* Create Post Section */}
      <div className="w-[900px] bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-xl font-semibold mb-4 flex items-center space-x-2">
          
          <span>Create a New Post</span>
        </h2>
        {isCreatingPost ? (
          <div className="space-y-4">
            <input
              type="text"
              className="w-full p-2 border rounded-md"
              placeholder="Post Title"
              value={newPost.title}
              onChange={(e) => setNewPost({ ...newPost, title: e.target.value })}
            />
            <textarea
              className="w-full p-2 border rounded-md"
              placeholder="Post Content"
              rows="4"
              value={newPost.content}
              onChange={(e) => setNewPost({ ...newPost, content: e.target.value })}
            ></textarea>
            <input
              type="text"
              className="w-full p-2 border rounded-md"
              placeholder="Tags (comma-separated)"
              value={newPost.tags}
              onChange={(e) => setNewPost({ ...newPost, tags: e.target.value })}
            />
            <div className="flex space-x-4">
              <button
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg"
                onClick={handleCreatePost}
              >
                Submit Post
              </button>
              <button
                className="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg"
                onClick={() => setIsCreatingPost(false)}
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <button
            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center space-x-2"
            onClick={() => setIsCreatingPost(true)}
          >
            <PlusCircle className="w-5 h-5" />
            <span>Create</span>
          </button>
        )}
      </div>

      {/* Search Bar */}
      <div className="w-[900px] mb-4">
        <SearchBar onSearch={fetchPosts} />
      </div>

      {/* Feed Posts */}
      <div className="w-[800px]">
        {posts.length > 0 ? (
          posts.map((post) => <FeedCard key={post.Id} post={post} className="mb-4 pb-4"/>)
        ) : (
          <p className="text-center text-gray-500">No posts available. Create one now!</p>
        )}
      </div>
    </div>
  );
};

export default Feed;