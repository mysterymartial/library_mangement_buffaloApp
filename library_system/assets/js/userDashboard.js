document.addEventListener('DOMContentLoaded', function () {
    const registerForm = document.getElementById('register-form');
    const checkoutForm = document.getElementById('checkout-form');
    const returnForm = document.getElementById('return-form');
    const reserveForm = document.getElementById('reserve-form');
    const messagesDiv = document.getElementById('messages');
    const booksContainer = document.getElementById('books-container');
    const checkoutSelect = document.getElementById('checkout-title');
    const returnSelect = document.getElementById('return-title');
    const reserveSelect = document.getElementById('reserve-title');

    registerForm.addEventListener('submit', async function (event) {
        event.preventDefault();
        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;

        if (!name || !email) {
            displayMessage('Name and email are required.', 'error');
            return;
        }

        try {
            const response = await fetch('http://localhost:3000/users/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name, email }),
            });

            const result = await response.json();
            if (response.ok) {
                displayMessage('User registered successfully', 'success');
                registerForm.reset();
            } else {
                if (result.error === 'Email already registered') {
                    displayMessage('Email already registered.', 'error');
                } else {
                    displayMessage(result.error, 'error');
                }
            }
        } catch (error) {
            displayMessage('Error registering user', 'error');
        }
    });

    checkoutForm.addEventListener('submit', async function (event) {
        event.preventDefault();
        const bookId = checkoutSelect.value;
        const userName = document.getElementById('checkout-username').value;

        if (!userName) {
            displayMessage('Username is required.', 'error');
            return;
        }

        try {
            const book = await getBookById(bookId);
            if (book && (book.status === 'available' || book.status === 'returned')) {
                const response = await fetch('http://localhost:3000/users/checkout', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ bookId, userName }),
                });

                const result = await response.json();
                if (response.ok) {
                    displayMessage('Book checked out successfully', 'success');
                    checkoutForm.reset();
                } else {
                    if (result.error === 'Invalid User Name') {
                        displayMessage('Invalid user name.', 'error');
                    } else {
                        displayMessage(result.error, 'error');
                    }
                }
            } else {
                displayMessage('Book is not available for checkout', 'error');
            }
        } catch (error) {
            displayMessage('Error checking out book', 'error');
        }
    });

    returnForm.addEventListener('submit', async function (event) {
        event.preventDefault();
        const bookId = returnSelect.value;
        const userName = document.getElementById('return-username').value;

        if (!userName) {
            displayMessage('Username is required.', 'error');
            return;
        }

        try {
            const response = await fetch('http://localhost:3000/users/return', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ bookId, userName }),
            });

            const result = await response.json();
            if (response.ok) {
                displayMessage('Book returned successfully', 'success');
                returnForm.reset();
            } else {
                if (result.error === 'Invalid User Name') {
                    displayMessage('Invalid user name.', 'error');
                } else {
                    displayMessage(result.error, 'error');
                }
            }
        } catch (error) {
            displayMessage('Error returning book', 'error');
        }
    });

    reserveForm.addEventListener('submit', async function (event) {
        event.preventDefault();
        const bookId = reserveSelect.value;
        const userName = document.getElementById('reserve-username').value;

        if (!userName) {
            displayMessage('Username is required.', 'error');
            return;
        }

        try {
            const book = await getBookById(bookId);
            if (book && (book.status === 'available' || book.status === 'returned')) {
                const response = await fetch('http://localhost:3000/users/reserve', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ bookId, userName }),
                });

                const result = await response.json();
                if (response.ok) {
                    displayMessage('Book reserved successfully', 'success');
                    reserveForm.reset();
                } else {
                    if (result.error === 'Invalid User Name') {
                        displayMessage('Invalid user name.', 'error');
                    } else {
                        displayMessage(result.error, 'error');
                    }
                }
            } else {
                displayMessage('Book is not available for reservation', 'error');
            }
        } catch (error) {
            displayMessage('Error reserving book', 'error');
        }
    });

    function displayMessage(message, type) {
        const div = document.createElement('div');
        div.classList.add(type);
        div.textContent = message;
        messagesDiv.appendChild(div);

        setTimeout(() => {
            div.remove();
        }, 3000);
    }

    // Load books and populate book list
    async function loadBooks() {
        try {
            const response = await fetch('http://localhost:3000/books', { credentials: 'include' });
            if (!response.ok) {
                console.error('Failed to fetch books from backend');
                return;
            }
            const books = await response.json();
            books.forEach(book => {
                renderBook(book);
                appendToSelect(book);
            });
        } catch (error) {
            console.error('Error loading books:', error);
        }
    }

    // Render book to the page
    function renderBook(book) {
        const existingBook = document.getElementById(book.id);
        if (existingBook) return;

        const bookItem = document.createElement('div');
        bookItem.classList.add('book-item');
        bookItem.setAttribute('id', book.id);
        bookItem.innerHTML = `
            <div>
                <strong>${book.title}</strong> by ${book.author} (ISBN: ${book.isbn})
            </div>
        `;
        booksContainer.appendChild(bookItem);
    }

    // Append book to select options
    function appendToSelect(book) {
        const option = document.createElement('option');
        option.value = book.id; // Use book ID as value
        option.textContent = `${book.title} by ${book.author}`;
        checkoutSelect.appendChild(option);
        returnSelect.appendChild(option.cloneNode(true));
        reserveSelect.appendChild(option.cloneNode(true));
    }

    // Fetch book details by ID
    async function getBookById(bookId) {
        try {
            const response = await fetch(`http://localhost:3000/books/getBookById/${bookId}`, { credentials: 'include' });
            if (!response.ok) {
                console.error('Failed to fetch book details from backend');
                return null;
            }
            return await response.json();
        } catch (error) {
            console.error('Error fetching book details:', error);
            return null;
        }
    }

    // Load books on page load
    loadBooks();
});