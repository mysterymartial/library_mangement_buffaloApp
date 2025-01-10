const BASE_URL = 'http://localhost:3000';
const booksContainer = document.getElementById('books-container');
const addBookButton = document.getElementById('add-book-btn');
const searchInput = document.getElementById('search-query');

const MESSAGES = {
    ADD_SUCCESS: "Book added successfully! The book is now in the library.",
    ADD_FAIL: "Failed to add book. Please check your input and try again.",
    UPDATE_SUCCESS: "Book updated successfully! The changes have been saved.",
    UPDATE_FAIL: "Failed to update book. Please verify the information and try again.",
    DELETE_SUCCESS: "Book deleted successfully! The book has been removed from the library.",
    DELETE_FAIL: "Failed to delete book. Please try again later.",
    SEARCH_SUCCESS: "Search completed successfully!",
    SEARCH_FAIL: "Search operation failed. Please try again.",
    LOAD_SUCCESS: "Books loaded successfully!",
    LOAD_FAIL: "Failed to load books. Please refresh the page.",
    NETWORK_ERROR: "Network connection issue. Please check your internet connection."
};

const fetchConfig = {
    mode: 'cors',
    credentials: 'include',
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    }
};

async function addBook() {
    const title = document.getElementById('title').value;
    const author = document.getElementById('author').value;
    const isbn = document.getElementById('isbn').value;

    if (!title, !author, !isbn) {
        alert('All fields are required.');
        return;
    }

    if (!isValidISBN(isbn)) {
        alert('Invalid ISBN format.');
        return;
    }

    try {
        const response = await fetch(`${BASE_URL}/books/add`, {
            ...fetchConfig,
            method: 'POST',
            body: JSON.stringify({ title, author, isbn })
        });

        if (response.ok) {
            const newBook = await response.json();
            renderBook(newBook, false);
            clearForm();
            alert(MESSAGES.ADD_SUCCESS);
            saveBooksToLocalStorage();
        } else {
            const errorData = await response.json();
            alert(`${MESSAGES.ADD_FAIL}\nError: ${errorData.error}`);
        }
    } catch (error) {
        console.error('Error adding book:', error);
        alert(MESSAGES.ADD_FAIL);
    }
}

async function updateBook() {
    const title = document.getElementById('update-title').value;
    const author = document.getElementById('update-author').value;
    const isbn = document.getElementById('update-isbn').value;

    if (!title || !author || !isbn) {
        alert('Title, Author, and ISBN are required.');
        return;
    }

    if (!isValidISBN(isbn)) {
        alert('Invalid ISBN format.');
        return;
    }

    try {
        const response = await fetch(`${BASE_URL}/books/update`, {
            ...fetchConfig,
            method: 'PUT',
            body: JSON.stringify({ title, author, isbn }) // No status included
        });

        if (response.ok) {
            const updatedBook = await response.json();
            const bookElement = document.querySelector(`[data-isbn='${isbn}']`);
            if (bookElement) {
                bookElement.querySelector('strong').textContent = title;
                bookElement.querySelector('.author').textContent = `by ${author}`;
                bookElement.querySelector('.isbn').textContent = `ISBN: ${isbn}`;
                bookElement.querySelector('.status').textContent = `Status: ${updatedBook.status}`; // Update status from response
            }
            clearForm();
            alert(MESSAGES.UPDATE_SUCCESS);
            saveBooksToLocalStorage();
        } else {
            const errorData = await response.json();
            alert(`${MESSAGES.UPDATE_FAIL}\nError: ${errorData.error}`);
        }
    } catch (error) {
        console.error('Error updating book:', error);
        alert(MESSAGES.UPDATE_FAIL);
    }
}

function renderBook(book) {
    const existingBook = document.querySelector(`[data-isbn='${book.isbn}']`);
    if (existingBook) {
        existingBook.querySelector('strong').textContent = book.title;
        existingBook.querySelector('.author').textContent = `by ${book.author}`;
        existingBook.querySelector('.isbn').textContent = `ISBN: ${book.isbn}`;
        existingBook.querySelector('.status').textContent = `Status: ${book.status}`;
        return;
    }

    const bookItem = document.createElement('div');
    bookItem.classList.add('book-item');
    bookItem.setAttribute('data-isbn', book.isbn);
    bookItem.setAttribute('id', book.id);
    bookItem.innerHTML = `
        <div>
            <strong>${book.title}</strong> by <span class="author">${book.author}</span> (ISBN: <span class="isbn">${book.isbn}</span>) <span class="status">Status: ${book.status}</span>
        </div>
        <div>
            <button class="btn delete-btn" onclick="deleteBook('${book.id}')">Delete</button>
        </div>
    `;
    booksContainer.appendChild(bookItem);
}

function populateUpdateForm(title, author, isbn, status) {
    document.getElementById('update-title').value = title;
    document.getElementById('update-author').value = author;
    document.getElementById('update-isbn').value = isbn;
}

async function deleteBook(id) {
    try {
        const response = await fetch(`${BASE_URL}/books/remove/${id}`, {
            ...fetchConfig,
            method: 'DELETE'
        });

        if (response.ok) {
            const bookItem = document.getElementById(id);
            if (bookItem) {
                bookItem.remove();
                removeFromLocalStorage(id);
                alert(MESSAGES.DELETE_SUCCESS);
            }
        } else {
            throw new Error('Failed to delete from server');
        }
    } catch (error) {
        console.error('Error deleting book:', error);
        alert(MESSAGES.DELETE_FAIL);
    }
}

function removeFromLocalStorage(id) {
    const books = JSON.parse(localStorage.getItem('books') || '[]');
    const updatedBooks = books.filter(book => book.id !== id);
    localStorage.setItem('books', JSON.stringify(updatedBooks));
}

function clearForm() {
    ['title', 'author', 'isbn', 'update-title', 'update-author', 'update-isbn']
        .forEach(id => document.getElementById(id).value = '');
}

async function loadBooks() {
    try {
        const response = await fetch(`${BASE_URL}/books`, fetchConfig);
        if (!response.ok) throw new Error('Failed to fetch books');

        const books = await response.json();
        const localBooks = JSON.parse(localStorage.getItem('books') || '[]');

        // Sync local storage with fetched books
        const fetchedBookIds = new Set(books.map(book => book.id));
        const updatedLocalBooks = localBooks.filter(book => fetchedBookIds.has(book.id));

        // Remove books from UI that are not in the fetched books
        const localBookIds = new Set(localBooks.map(book => book.id));
        const booksToRemoveFromUI = Array.from(localBookIds).filter(id => !fetchedBookIds.has(id));
        booksToRemoveFromUI.forEach(id => {
            const bookItem = document.getElementById(id);
            if (bookItem) {
                bookItem.remove();
            }
        });

        booksContainer.innerHTML = ''; // Clear the current books in the container
        books.forEach(book => renderBook(book));

        localStorage.setItem('books', JSON.stringify(updatedLocalBooks));

        alert(MESSAGES.LOAD_SUCCESS);
    } catch (error) {
        console.error('Error loading books:', error);
        alert(MESSAGES.LOAD_FAIL);
    }
}

function saveBooksToLocalStorage() {
    const books = Array.from(document.querySelectorAll('.book-item')).map(item => ({
        id: item.getAttribute('id'),
        title: item.querySelector('strong').textContent,
        author: item.querySelector('.author').textContent.replace('by ', ''),
        isbn: item.querySelector('.isbn').textContent.replace('ISBN: ', ''),
        status: item.querySelector('.status').textContent.replace('Status: ', '')
    }));
    localStorage.setItem('books', JSON.stringify(books));
}

async function searchBooks(query) {
    try {
        const response = await fetch(`${BASE_URL}/books/search?query=${encodeURIComponent(query)}`, fetchConfig);
        if (!response.ok) throw new Error('Failed to search books');

        const books = await response.json();
        booksContainer.innerHTML = '';
        books.forEach(book => renderBook(book));
        alert(MESSAGES.SEARCH_SUCCESS);
        searchInput.value = '';
    } catch (error) {
        console.error('Error searching books:', error);
        alert(MESSAGES.SEARCH_FAIL);
    }
}

function filterBooks(query) {
    if (!query) {
        loadBooks();
        return;
    }

    searchBooks(query);
}

function isValidISBN(isbn) {
    const isbn10Regex = /^\d{9}[\dXx]$/;
    const isbn13Regex = /^\d{13}$/;
    isbn = isbn.replace(/-/g, '');
    return isbn10Regex.test(isbn) || isbn13Regex.test(isbn);
}

// Event Listeners
addBookButton.addEventListener('click', addBook);
document.getElementById('update-book-btn').addEventListener('click', updateBook);
searchInput.addEventListener('input', (e) => {
    const query = e.target.value.trim();
    filterBooks(query);
});
window.onload = loadBooks;