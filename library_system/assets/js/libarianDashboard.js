const booksContainer = document.getElementById('books-container');
const addBookButton = document.getElementById('add-book-btn');
const searchInput = document.getElementById('search-query');

async function addBook() {
    const title = document.getElementById('title').value;
    const author = document.getElementById('author').value;
    const isbn = document.getElementById('isbn').value;

    if (!title || !author || !isbn) {
        alert('All fields are required.');
        return;
    }

    const bookData = { title, author, isbn };

    try {
        const response = await fetch('http://localhost:3000/books/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookData),
            mode: 'cors', // Add CORS mode
            credentials: 'include', // Ensure cookies/session are included in requests
        });

        if (response.ok) {
            const newBook = await response.json();
            renderBook(newBook);
            clearForm();
            alert('Book added successfully!');
            saveBooksToLocalStorage();
        } else {
            const errorData = await response.json();
            alert('Error adding book: ' + errorData.error);
        }
    } catch (error) {
        console.error('Error adding book:', error);
    }
}

async function updateBook() {
    const title = document.getElementById('update-title').value;
    const author = document.getElementById('update-author').value;
    const isbn = document.getElementById('update-isbn').value;
    const bookId = document.getElementById('update-book-id').value;

    if (!title || !author || !isbn || !bookId) {
        alert('All fields are required.');
        return;
    }

    const bookData = { title, author, isbn };

    try {
        const response = await fetch(`http://localhost:3000/books/update/${bookId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookData),
            mode: 'cors', // Add CORS mode
            credentials: 'include', // Ensure cookies/session are included in requests
        });

        if (response.ok) {
            const updatedBook = await response.json();
            renderBook(updatedBook);
            clearForm();
            alert('Book updated successfully!');
            saveBooksToLocalStorage();
        } else {
            const errorData = await response.json();
            alert('Error updating book: ' + errorData.error);
        }
    } catch (error) {
        console.error('Error updating book:', error);
    }
}

function renderBook(book) {
    const existingBook = document.getElementById(book.id);
    if (existingBook) {
        existingBook.querySelector('strong').textContent = book.title;
        existingBook.querySelector('.author').textContent = `by ${book.author}`;
        existingBook.querySelector('.isbn').textContent = `ISBN: ${book.isbn}`;
        return;
    }

    const bookItem = document.createElement('div');
    bookItem.classList.add('book-item');
    bookItem.setAttribute('id', book.id);
    bookItem.innerHTML = `
        <div>
            <strong>${book.title}</strong> by <span class="author">${book.author}</span> (ISBN: <span class="isbn">${book.isbn}</span>)
        </div>
        <div>
            <button class="btn delete-btn" onclick="deleteBook('${book.id}')">Delete</button>
        </div>
    `;
    booksContainer.appendChild(bookItem);
}

async function deleteBook(bookId) {
    try {
        const response = await fetch(`http://localhost:3000/books/remove/${bookId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            mode: 'cors', // Add CORS mode
            credentials: 'include', // Ensure cookies/session are included in requests
        });

        if (response.ok) {
            alert('Book deleted successfully!');
            const bookItem = document.getElementById(bookId);
            if (bookItem) bookItem.remove();
            saveBooksToLocalStorage();
        } else {
            const errorData = await response.json();
            alert('Error deleting book: ' + errorData.error);
        }
    } catch (error) {
        console.error('Error deleting book:', error);
    }
}

function clearForm() {
    document.getElementById('title').value = '';
    document.getElementById('author').value = '';
    document.getElementById('isbn').value = '';
    document.getElementById('update-title').value = '';
    document.getElementById('update-author').value = '';
    document.getElementById('update-isbn').value = '';
    document.getElementById('update-book-id').value = '';
}

async function loadBooks() {
    try {
        const storedBooks = localStorage.getItem('books');
        if (storedBooks) {
            const books = JSON.parse(storedBooks);
            books.forEach(book => renderBook(book));
        } else {
            const response = await fetch('http://localhost:3000/books', {
                mode: 'cors', // Add CORS mode
                credentials: 'include', // Ensure cookies/session are included in requests
            });
            if (!response.ok) {
                console.error('Failed to fetch books from backend');
                return;
            }
            const books = await response.json();
            books.forEach(book => renderBook(book));
            saveBooksToLocalStorage();
        }
    } catch (error) {
        console.error('Error loading books:', error);
    }
}

function saveBooksToLocalStorage() {
    const books = [];
    const bookItems = document.querySelectorAll('.book-item');
    bookItems.forEach(item => {
        const title = item.querySelector('strong').textContent;
        const author = item.querySelector('.author').textContent.replace('by ', '');
        const isbn = item.querySelector('.isbn').textContent.replace('ISBN: ', '');
        const id = item.id;
        books.push({ id, title, author, isbn });
    });
    localStorage.setItem('books', JSON.stringify(books));
}

async function searchBooks() {
    const query = searchInput.value.trim();

    if (!query) {
        alert('Please enter a search query.');
        return;
    }

    try {
        const response = await fetch(`http://localhost:3000/books/search?query=${query}`, {
            mode: 'cors', // Add CORS mode
            credentials: 'include', // Ensure cookies/session are included in requests
        });
        if (!response.ok) {
            console.error('Failed to search books');
            return;
        }
        const books = await response.json();
        clearBooksContainer();
        books.forEach(book => renderBook(book));
    } catch (error) {
        console.error('Error searching books:', error);
    }
}

function clearBooksContainer() {
    booksContainer.innerHTML = '';
}

addBookButton.addEventListener('click', addBook);
document.getElementById('search-books-btn').addEventListener('click', searchBooks);
document.getElementById('update-book-btn').addEventListener('click', updateBook);
window.onload = loadBooks;