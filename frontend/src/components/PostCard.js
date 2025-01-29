import React, { useState } from 'react';
import api from '../utils/Api';
import { Edit, Trash, Save, XCircle } from 'lucide-react'; 

const PostCard = ({ post, onDelete, onEdit, allowActions = false }) => {
  const [newComment, setNewComment] = useState('');
  const [comments, setComments] = useState(post.Comments || []);
  const [isEditing, setIsEditing] = useState(false);
  const [editedTitle, setEditedTitle] = useState(post.Title);
  const [editedContent, setEditedContent] = useState(post.Content);

  const [editingCommentId, setEditingCommentId] = useState(null);
  const [editedCommentContent, setEditedCommentContent] = useState('');

  const handleAddComment = async () => {
    if (!newComment.trim()) return;
    try {
      const response = await api.post(`/posts/${post.Id}/comments`, { content: newComment });
      setComments([...comments, response.data.data]);
      setNewComment('');
    } catch (err) {
      console.error(err.response?.data?.message || 'Error adding comment');
    }
  };

  const handleEditComment = async (commentId) => {
    try {
      const response = await api.put(`/posts/${post.Id}/comments/${commentId}`, {
        content: editedCommentContent,
      });
      setComments(
        comments.map((comment) => (comment.Id === commentId ? response.data.data : comment))
      );
      setEditingCommentId(null);
      setEditedCommentContent('');
    } catch (err) {
      console.error(err.response?.data?.message || 'Error editing comment');
    }
  };

  const handleDeleteComment = async (commentId) => {
    try {
      await api.delete(`/posts/${post.Id}/comments/${commentId}`);
      setComments(comments.filter((comment) => comment.Id !== commentId));
    } catch (err) {
      console.error(err.response?.data?.message || 'Error deleting comment');
    }
  };

  const handleSaveEdit = () => {
    onEdit({ ...post, Title: editedTitle, Content: editedContent });
    setIsEditing(false);
  };

  return (
    <div className="p-6 max-w-4xl mx-auto border rounded-lg shadow-lg bg-gray-50">
      {isEditing ? (
        <div className="mb-4">
          <input
            type="text"
            className="w-full p-3 border rounded-md mb-3 bg-gray-100"
            value={editedTitle}
            onChange={(e) => setEditedTitle(e.target.value)}
            placeholder="Edit title"
          />
          <textarea
            className="w-full p-3 border rounded-md mb-3 bg-gray-100"
            value={editedContent}
            onChange={(e) => setEditedContent(e.target.value)}
            placeholder="Edit content"
          />
          <div className="flex gap-3">
            <button
              onClick={handleSaveEdit}
              className="px-4 py-2 bg-teal-500 text-white rounded-md hover:bg-teal-600"
            >
              <Save className="inline mr-1" />
              Save
            </button>
            <button
              onClick={() => setIsEditing(false)}
              className="px-4 py-2 bg-gray-400 text-white rounded-md hover:bg-gray-500"
            >
              <XCircle className="inline mr-1" />
              Cancel
            </button>
          </div>
        </div>
      ) : (
        <>
          <h3 className="text-3xl font-semibold text-gray-800 mb-4">{post.Title}</h3>
          <p className="text-lg text-gray-700 mb-6">{post.Content}</p>
          <small className="block text-gray-500 mb-6">
            Posted on: {new Date(post.CreatedAt).toLocaleDateString()}
          </small>
        </>
      )}

      <div className="mt-6">
        <h4 className="text-xl font-semibold text-gray-800 mb-4">Comments</h4>
        <ul className="space-y-3">
          {comments.map((comment) => (
            <li key={comment.Id} className="bg-white p-4 rounded-lg shadow-sm flex flex-col border">
              {editingCommentId === comment.Id ? (
                <div className="flex gap-3">
                  <input
                    type="text"
                    className="flex-grow p-2 border rounded-md"
                    value={editedCommentContent}
                    onChange={(e) => setEditedCommentContent(e.target.value)}
                    placeholder="Edit comment"
                  />
                  <button
                    onClick={() => handleEditComment(comment.Id)}
                    className="px-3 py-2 bg-teal-500 text-white rounded-md hover:bg-teal-600"
                  >
                    <Save className="inline mr-1" />
                    Save
                  </button>
                  <button
                    onClick={() => setEditingCommentId(null)}
                    className="px-3 py-2 bg-gray-400 text-white rounded-md hover:bg-gray-500"
                  >
                    <XCircle className="inline mr-1" />
                    Cancel
                  </button>
                </div>
              ) : (
                <div className="flex justify-between items-start">
                  <div>
                    <h5 className="font-semibold text-sm text-gray-800">{comment.UserId}</h5>
                    <p className="text-gray-700 mt-1 text-sm">{comment.Content}</p>
                  </div>
                  <div className="text-xs text-gray-500">
                    {new Date(comment.CreatedAt).toLocaleString()}
                  </div>
                </div>
              )}
              {editingCommentId !== comment.Id && (
                <div className="flex gap-2 mt-2">
                  <button
                    onClick={() => {
                      setEditingCommentId(comment.Id);
                      setEditedCommentContent(comment.Content);
                    }}
                    className="px-3 py-1 bg-yellow-400 text-white rounded-md hover:bg-yellow-500"
                  >
                    <Edit className="inline mr-1" />
                    Edit
                  </button>
                  <button
                    onClick={() => handleDeleteComment(comment.Id)}
                    className="px-3 py-1 bg-red-500 text-white rounded-md hover:bg-red-600"
                  >
                    <Trash className="inline mr-1" />
                    Delete
                  </button>
                </div>
              )}
            </li>
          ))}
        </ul>

        <div className="mt-6 flex gap-2">
          <input
            type="text"
            placeholder="Write a comment..."
            className="flex-grow p-3 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-teal-300"
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
          />
          <button
            onClick={handleAddComment}
            className="px-4 py-2 bg-teal-500 text-white font-semibold rounded-md shadow-sm hover:bg-teal-600"
          >
            Add
          </button>
        </div>
      </div>

      {allowActions && !isEditing && (
        <div className="mt-4 flex gap-4">
          <button
            onClick={() => setIsEditing(true)}
            className="p-2 bg-yellow-400 text-white rounded-md hover:bg-yellow-500"
          >
            <Edit className="inline mr-1" />
            Edit Post
          </button>
          <button
            onClick={() => onDelete(post.Id)}
            className="p-2 bg-red-500 text-white rounded-md hover:bg-red-600"
          >
            <Trash className="inline mr-1" />
            Delete Post
          </button>
        </div>
      )}
    </div>
  );

  // return (
  //   <div className="p-6 border rounded-lg shadow-md bg-white">
  //     {isEditing ? (
  //       <div className="mb-4">
  //         <input
  //           type="text"
  //           className="w-full p-3 border rounded mb-3"
  //           value={editedTitle}
  //           onChange={(e) => setEditedTitle(e.target.value)}
  //           placeholder="Edit title"
  //         />
  //         <textarea
  //           className="w-full p-3 border rounded mb-3"
  //           value={editedContent}
  //           onChange={(e) => setEditedContent(e.target.value)}
  //           placeholder="Edit content"
  //         />
  //         <div className="flex gap-3">
  //           <button
  //             onClick={handleSaveEdit}
  //             className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
  //           >
  //             Save
  //           </button>
  //           <button
  //             onClick={() => setIsEditing(false)}
  //             className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
  //           >
  //             Cancel
  //           </button>
  //         </div>
  //       </div>
  //     ) : (
  //       <>
  //         <h3 className="text-xl font-bold mb-2 text-gray-800">{post.Title}</h3>
  //         <p className="mb-4 text-gray-700">{post.Content}</p>
  //         <small className="block mb-4 text-gray-500">
  //           Posted on: {new Date(post.CreatedAt).toLocaleDateString()}
  //         </small>
  //       </>
  //     )}

  //     <div className="mt-6">
  //       <h4 className="text-lg font-bold mb-4 text-gray-800">Comments</h4>
  //       <ul className="space-y-4">
  //         {comments.map((comment) => (
  //           <li key={comment.Id} className="bg-gray-100 p-4 rounded-lg shadow-sm flex flex-col">
  //             {editingCommentId === comment.Id ? (
  //               <div className="flex gap-3">
  //                 <input
  //                   type="text"
  //                   className="flex-grow p-2 border rounded"
  //                   value={editedCommentContent}
  //                   onChange={(e) => setEditedCommentContent(e.target.value)}
  //                   placeholder="Edit comment"
  //                 />
  //                 <button
  //                   onClick={() => handleEditComment(comment.Id)}
  //                   className="px-3 py-2 bg-green-500 text-white rounded hover:bg-green-600"
  //                 >
  //                   Save
  //                 </button>
  //                 <button
  //                   onClick={() => setEditingCommentId(null)}
  //                   className="px-3 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
  //                 >
  //                   Cancel
  //                 </button>
  //               </div>
  //             ) : (
  //               <div className="flex justify-between items-start">
  //                 <div>
  //                   <h5 className="font-semibold text-sm text-gray-800">{comment.UserId}</h5>
  //                   <p className="text-gray-700 mt-1 text-sm">{comment.Content}</p>
  //                 </div>
  //                 <div className="text-xs text-gray-500">
  //                   {new Date(comment.CreatedAt).toLocaleString()}
  //                 </div>
  //               </div>
  //             )}
  //             {editingCommentId !== comment.Id && (
  //               <div className="flex gap-2 mt-2">
  //                 <button
  //                   onClick={() => {
  //                     setEditingCommentId(comment.Id);
  //                     setEditedCommentContent(comment.Content);
  //                   }}
  //                   className="px-3 py-1 bg-yellow-500 text-white rounded hover:bg-yellow-600"
  //                 >
  //                   Edit
  //                 </button>
  //                 <button
  //                   onClick={() => handleDeleteComment(comment.Id)}
  //                   className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
  //                 >
  //                   Delete
  //                 </button>
  //               </div>
  //             )}
  //           </li>
  //         ))}
  //       </ul>

  //       <div className="mt-6 flex gap-2">
  //         <input
  //           type="text"
  //           placeholder="Write a comment..."
  //           className="flex-grow p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring focus:ring-blue-200"
  //           value={newComment}
  //           onChange={(e) => setNewComment(e.target.value)}
  //         />
  //         <button
  //           onClick={handleAddComment}
  //           className="px-4 py-2 bg-blue-500 text-white font-semibold rounded-lg shadow-sm hover:bg-blue-600"
  //         >
  //           Add
  //         </button>
  //       </div>
  //     </div>

  //     {allowActions && !isEditing && (
  //       <div className="mt-4 flex gap-4">
  //         <button
  //           onClick={() => setIsEditing(true)}
  //           className="p-2 bg-yellow-500 text-white rounded hover:bg-yellow-600"
  //         >
  //           Edit Post
  //         </button>
  //         <button
  //           onClick={() => onDelete(post.Id)}
  //           className="p-2 bg-red-500 text-white rounded hover:bg-red-600"
  //         >
  //           Delete Post
  //         </button>
  //       </div>
  //     )}
  //   </div>
  // );
};

export default PostCard;
