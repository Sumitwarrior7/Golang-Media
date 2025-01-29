import React, { useEffect, useState, useContext } from 'react';
import AuthContext from '../context/AuthContext';
import api from '../utils/Api';
import FeedCard from '../components/FeedCard';
import { UserCard } from '../components/UserCard';
import { Edit, Mail, User, Users as UsersIcon, XCircle, Save, Trash, Calendar, UserPen, Newspaper} from 'lucide-react';

const Dashboard = () => {
  const { user, logout } = useContext(AuthContext);
  const [profile, setProfile] = useState({});
  const [followers, setFollowers] = useState([]);
  const [posts, setPosts] = useState([]);
  const [editMode, setEditMode] = useState(false);
  const [formData, setFormData] = useState({ name: '', email: '' });

  // Fetch profile and followers data independently
  const fetchProfile = async () => {
    try {
      // Fetch user details
      const currentUserResponse = await api.get('/users/current-user');
      setProfile(currentUserResponse.data.data);
      setFormData({ 
        username: currentUserResponse.data.data.username, 
        email: currentUserResponse.data.data.email 
      });

      // Fetch followed users
      const followedUsersResponse = await api.get('/users/followed-users');
      if (followedUsersResponse.data.data) {
        // console.log("fu :", followedUsersResponse.data.data)
        const formattedFollowers = followedUsersResponse.data.data.map(user => ({
          username: user.Username,
          email: user.Email,
          id: user.UserId,
          created_at: user.CreatedAt
        }));
        
        setFollowers(formattedFollowers);
      } else {
        setFollowers([]);
      }
    } catch (err) {
      console.error(err.response?.data?.message || 'Error fetching profile data');
    }
  };

  const fetchPosts = async () => {
    try {
      const response = await api.get('/posts');
      // console.log('fetched posts :', response.data.data[0].Tags)
      console.log(response.data.data)
      setPosts(response.data.data);
    } catch (err) {
      console.error(err.response?.data?.message || 'Error fetching posts');
    }
  };

  // Handle form changes
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  // Handle profile update
  const handleUpdateProfile = async (e) => {
    e.preventDefault();
    try {
      await api.patch(`/users/${user.id}`, formData);
      setProfile({ ...profile, ...formData });
      setEditMode(false);
    } catch (err) {
      console.error(err.response?.data?.message || 'Error updating profile');
    }
  };

  // Handle unfollow user
  const handleUnfollow = async (userId) => {
    try {
      await api.put(`/users/${userId}/unfollow`);
      setFollowers(followers.filter((f) => f.id !== userId));
    } catch (err) {
      console.error(err.response?.data?.message || 'Error unfollowing user');
    }
  };

  // Fetch user profile and posts on mount or when `user` changes
  useEffect(() => {
    if (user) {
      fetchProfile();
      fetchPosts();
    }
  }, [user]);

  return (
    <div className="p-4 min-h-screen bg-gray-50 flex flex-col items-center">
      <h1 className="text-3xl font-semibold mb-4 text-gray-800">Dashboard</h1>

      {/* Profile Section */}
      <div className="my-6 w-full max-w-xl text-center bg-white shadow-lg rounded-2xl p-6 border border-gray-200">
        <h2 className="text-3xl font-bold mb-4 flex items-center justify-center gap-2 text-gray-800">
          <UserPen size={28} /> Profile
        </h2>
        {editMode ? (
          <form onSubmit={handleUpdateProfile} className="mt-4">
            <input
              type="text"
              name="username"
              placeholder="Username"
              value={formData.username}
              onChange={handleInputChange}
              className="w-full p-3 mb-4 border rounded-lg text-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
              required
            />
            <input
              type="email"
              name="email"
              placeholder="Email"
              value={formData.email}
              onChange={handleInputChange}
              className="w-full p-3 mb-4 border rounded-lg text-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
              required
            />
            <div className="flex justify-center gap-4">
              <button className="p-3 bg-green-500 text-white rounded-lg text-lg flex items-center gap-2 hover:bg-green-600 transition">
                <Save size={20} /> Save
              </button>
              <button
                type="button"
                onClick={() => setEditMode(false)}
                className="p-3 bg-red-500 text-white rounded-lg text-lg flex items-center gap-2 hover:bg-red-600 transition"
              >
                <XCircle size={20} /> Cancel
              </button>
            </div>
          </form>
        ) : (
          <div className="text-xl text-gray-700 space-y-3">
            <p className="flex items-center gap-2">
            <User size={20} className="text-gray-500"/>{profile.username}
            </p>
            <p className="flex items-center gap-2">
              <Mail size={20} className="text-gray-500" /> {profile.email}
            </p>
            <p className="flex items-center  gap-2">
              <Calendar size={20} className="text-gray-500" /> Joined: {new Date(profile.created_at).toLocaleDateString()}
            </p>
            <button
              onClick={() => setEditMode(true)}
              className="p-3 bg-blue-500 text-white rounded-lg text-lg flex items-center gap-2 hover:bg-blue-600 transition mt-4"
            >
              <Edit size={20} /> Edit Profile
            </button>
          </div>
        )}
      </div>

      {/* Followed Users Section */}
      <div className="my-6 w-full max-w-4xl flex flex-col items-center">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4 flex items-center">
          <User className="mr-2" size={24} /> Followed Users
        </h2>

        {followers.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 w-full">
            {followers.map((follower) => (
              <UserCard key={follower.id} user={follower} />
            ))}
          </div>
        ) : (
          <p className="text-lg text-gray-500">No followed users.</p>
        )}
      </div>

      {/* Posts Section */}
      <div className="my-6 w-full max-w-3xl">
        <h2 className="text-3xl font-bold mb-4 flex items-center justify-center gap-2 text-gray-800">
          <Newspaper size={28} />  Your Posts
        </h2>
        {posts.length > 0 ? (
          <ul className="space-y-6">
            {posts.map((post) => (
              <div key={post.Id} className="w-full rounded-lg bg-white shadow-sm">
                <FeedCard post={post} />
              </div>
            ))}
          </ul>
        ) : (
          <p className="text-lg text-gray-500">No posts yet.</p>
        )}
      </div>
    </div>
  );

  // return (
  //   <div className="p-4 min-h-screen bg-gray-50 flex flex-col items-center">
  //     <h1 className="text-2xl font-semibold mb-4 text-gray-800">Dashboard</h1>

  //     <div className="my-6">
  //       <h2 className="text-xl font-semibold">Profile</h2>
  //       {editMode ? (
  //         <form onSubmit={handleUpdateProfile} className="mt-4">
  //           <input
  //             type="text"
  //             name="username"
  //             placeholder="Username"
  //             value={formData.username}
  //             onChange={handleInputChange}
  //             className="w-full p-2 mb-4 border rounded-md"
  //             required
  //           />
  //           <input
  //             type="email"
  //             name="email"
  //             placeholder="Email"
  //             value={formData.email}
  //             onChange={handleInputChange}
  //             className="w-full p-2 mb-4 border rounded-md"
  //             required
  //           />
  //           <button className="p-2 bg-green-500 text-white rounded-md mr-4">
  //             <Edit className="inline-block mr-2" size={16} /> Save
  //           </button>
  //           <button
  //             type="button"
  //             onClick={() => setEditMode(false)}
  //             className="p-2 bg-red-500 text-white rounded-md"
  //           >
  //             <Trash className="inline-block mr-2" size={16} /> Cancel
  //           </button>
  //         </form>
  //       ) : (
  //         <div>
  //           <p>
  //             <strong>Name:</strong> {profile.username}
  //           </p>
  //           <p>
  //             <strong>Email:</strong> {profile.email}
  //           </p>
  //           <button
  //             onClick={() => setEditMode(true)}
  //             className="p-2 bg-blue-500 text-white rounded-md mt-4"
  //           >
  //             <Edit className="inline-block mr-2" size={16} /> Edit Profile
  //           </button>
  //         </div>
  //       )}
  //     </div>

  //     <div className="my-6">
  //       <h2 className="text-2xl font-semibold text-gray-800 mb-4 flex items-center">
  //         <User className="mr-2" size={24} /> Followed Users
  //       </h2>
        
  //       {followers.length > 0 ? (
  //         <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6 w-full max-w-4xl">
  //           {followers.map((follower) => (
  //             <UserCard key={follower.id} user={follower} />
  //           ))}
  //         </div>
  //       ) : (
  //         <p className="text-base text-gray-500">No followed users.</p>
  //       )}
  //     </div>

  //     <div className="my-6">
  //       <h2 className="text-xl font-semibold mb-4">Your Posts</h2>

  //       {posts.length > 0 ? (
  //         <ul className="space-y-6">
  //           {posts.map((post) => (
  //             <FeedCard key={post.Id} post={post} />
  //           ))}
  //         </ul>
  //       ) : (
  //         <p className="text-base text-gray-500">No posts yet.</p>
  //       )}
  //     </div>

  //   </div>
  // )
};

export default Dashboard;
