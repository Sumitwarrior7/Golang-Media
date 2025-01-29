import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import api from '../utils/Api';

const PublicProfile = () => {
  const { id } = useParams();
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const { data } = await api.get(`/users/user/${id}`);
        setProfile(data);
      } catch (err) {
        console.error(err.response.data.message);
      }
    };
    fetchProfile();
  }, [id]);

  if (!profile) return <p>Loading...</p>;

  return (
    <div className="p-4">
      <h1 className="text-3xl font-bold">{profile.name}</h1>
      <p className="text-gray-700">{profile.bio}</p>
      <div className="mt-4">
        <h2 className="text-xl font-bold">Posts</h2>
        {profile.posts.map((post) => (
          <div key={post.id} className="border rounded p-2 mt-2">
            <h3 className="font-bold">{post.title}</h3>
            <p>{post.content}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default PublicProfile;
