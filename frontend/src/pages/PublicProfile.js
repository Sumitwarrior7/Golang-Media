import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import api from '../utils/Api';
import { User, Mail, Calendar, Edit3, FileText, Users, MessageCircle, Tag } from 'lucide-react';


const PublicProfile = () => {
  const { userId } = useParams(); // User userId from URL params
  const [profile, setProfile] = useState(null); // User profile data
  const [posts, setPosts] = useState([]); // Posts by the user
  const [isFollowing, setIsFollowing] = useState(false); // Follow/unfollow state
  const [loading, setLoading] = useState(true); // Loading state

  // Fetch user profile and posts
  useEffect(() => {
    const fetchData = async () => {
      try {
        const { data: profileData } = await api.get(`/users/${userId}`);
        setProfile(profileData.data);

        const { data: postsData } = await api.get(`/posts/user/${userId}`);
        setPosts(postsData.data);

        const { data: followedUsers } = await api.get('/users/followed-users');
        console.log("followedUsers.data :", followedUsers.data)
        const isUserFollowed = followedUsers.data.some(
          (followedUser) => followedUser.UserId === parseInt(userId)
        );
        setIsFollowing(isUserFollowed);
      } catch (error) {
        console.error('Error fetching profile or posts:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [userId]);

  const handleFollowToggle = async () => {
    try {
      if (isFollowing) {
        await api.put(`/users/${userId}/unfollow`);
      } else {
        await api.put(`/users/${userId}/follow`);
      }
      setIsFollowing(!isFollowing); // Toggle follow state
    } catch (error) {
      console.error('Error updating follow status:', error);
    }
  };

  if (loading) return <p>Loading...</p>;

  return (
    <div className="max-w-4xl mx-auto p-4">
      {/* Profile Section */}
      <div className="p-6 bg-white shadow-lg rounded-lg border border-gray-200">
        <div className="flex items-center space-x-4">
          <div className="bg-gray-100 p-4 rounded-full">
            <User className="text-gray-600 w-8 h-8" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-800">{profile.username}</h1>
            <p className="text-gray-600 flex items-center space-x-2">
              <Mail className="w-4 h-4" />
              <span>{profile.email}</span>
            </p>
            <p className="text-gray-500 flex items-center space-x-2 mt-1">
              <Calendar className="w-4 h-4" />
              <span>Joined: {new Date(profile.created_at).toLocaleDateString()}</span>
            </p>
          </div>
        </div>
        <p className="mt-4 text-gray-700">{profile.bio}</p>

        {/* Follow/Unfollow Button */}
        <button
          onClick={handleFollowToggle}
          className={`mt-6 px-6 py-2 rounded-lg text-white font-medium ${
            isFollowing ? 'bg-red-500 hover:bg-red-600' : 'bg-blue-500 hover:bg-blue-600'
          }`}
        >
          {isFollowing ? 'Unfollow' : 'Follow'}
        </button>
      </div>

      {/* Posts Section */}
      <div className="mt-8">
        <h2 className="text-xl font-bold flex items-center space-x-2">
          <FileText className="w-6 h-6 text-gray-700" />
          <span>Posts</span>
        </h2>
        {posts.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mt-4">
            {posts.map((post) => (
              <div key={post.id} className="p-6 bg-white shadow-lg rounded-lg border border-gray-200 transition-transform transform hover:scale-105">
                <h3 className="font-semibold text-lg text-gray-800 flex items-center space-x-2">
                  <Edit3 className="w-5 h-5 text-gray-700" />
                  <span>{post.Title}</span>
                </h3>
                
                <p className="text-gray-600 mt-2">{post.Content}</p>
                
                <div className="flex items-center text-sm text-gray-500 mt-4 space-x-4">
                  <div className="flex items-center space-x-1">
                    <MessageCircle className="w-4 h-4 text-gray-600" />
                    <span>{post.CommentCount} Comments</span>
                  </div>
                  <div className="flex items-center space-x-1">
                    <Calendar className="w-4 h-4 text-gray-600" />
                    <span>{new Date(post.CreatedAt).toLocaleDateString()}</span>
                  </div>
                </div>
                
                <div className="mt-4 flex flex-wrap gap-2">
                  {post.Tags.map((tag, index) => (
                    <span key={index} className="inline-flex items-center px-3 py-1 text-xs font-medium bg-blue-100 text-blue-800 rounded-full">
                      <Tag className="w-4 h-4 mr-1" />
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500 mt-4">No posts available.</p>
        )}
      </div>
    </div>
  );
};

export default PublicProfile;
