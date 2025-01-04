// const BookManager = (function() {
//     const API_BASE_URL = '/books';  // Simplified to match Buffalo routes
//     const NOTIFICATION_DURATION = 3000;
//
//     // Add missing functions
//     const loadBooks = async () => {
//         try {
//             const books = await makeApiCall('/');
//             displayBooks(books);
//             logger.success('Books loaded successfully');
//         } catch (error) {
//             showNotification('Failed to load books', 'error');
//         }
//     };
//
//     const addBook = async (bookData) => {
//         try {
//             await makeApiCall('/', {
//                 method: 'POST',
//                 body: JSON.stringify(bookData)
//             });
//             await loadBooks(); // Refresh book list
//             showNotification('Book added successfully');
//             logger.success('Book added');
//         } catch (error) {
//             showNotification('Failed to add book', 'error');
//         }
//     };
//
//     const updateBook = async (bookId, bookData) => {
//         try {
//             await makeApiCall(`/${bookId}`, {
//                 method: 'PUT',
//                 body: JSON.stringify(bookData)
//             });
//             await loadBooks();
//             showNotification('Book updated successfully');
//         } catch (error) {
//             showNotification('Failed to update book', 'error');
//         }
//     };
//
//     const searchBooks = async (query) => {
//         try {
//             const books = await makeApiCall(`/search?q=${encodeURIComponent(query)}`);
//             displayBooks(books);
//         } catch (error) {
//             showNotification('Search failed', 'error');
//         }
//     };
//
//     const showNotification = (message, type = 'success') => {
//         const notification = document.getElementById('notification');
//         if (!notification) return;
//
//         notification.textContent = message;
//         notification.className = `notification ${type}`;
//         notification.style.display = 'block';
//
//         setTimeout(() => {
//             notification.style.display = 'none';
//         }, NOTIFICATION_DURATION);
//     };
//
//     const debounce = (func, wait) => {
//         let timeout;
//         return (...args) => {
//             clearTimeout(timeout);
//             timeout = setTimeout(() => func.apply(this, args), wait);
//         };
//     };
//
//     // Return expanded public API
//     return {
//         init,
//         editBook,
//         removeBook: async (bookId) => {
//             if (!confirm('Are you sure you want to delete this book?')) return;
//             try {
//                 await makeApiCall(`/${bookId}`, { method: 'DELETE' });
//                 showNotification('Book deleted successfully!');
//                 await loadBooks();
//             } catch (error) {
//                 showNotification('Failed to delete book', 'error');
//             }
//         },
//         // Add these to make functions accessible
//         addBook,
//         updateBook,
//         loadBooks,
//         searchBooks
//     };
// })();
//
// // Initialize when DOM is ready
// document.addEventListener('DOMContentLoaded', () => {
//     BookManager.init();
//     // Make BookManager globally available
//     window.BookManager = BookManager;
// });
