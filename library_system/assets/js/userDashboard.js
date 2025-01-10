document.addEventListener('DOMContentLoaded', function () {
    const registerForm = document.getElementById('register-form');
    const checkoutForm = document.getElementById('checkout-form');
    const returnForm = document.getElementById('return-form');
    const reserveForm = document.getElementById('reserve-form');
    const booksContainer = document.getElementById('books-container');
    const checkoutSelect = document.getElementById('checkout-title');
    const returnSelect = document.getElementById('return-title');
    const reserveSelect = document.getElementById('reserve-title');

    const MESSAGES = {
        REGISTER_SUCCESS: "Registration successful! Welcome to the library system.",
        REGISTER_FAIL: "Registration failed.",
        CHECKOUT_SUCCESS: "Book checked out successfully!",
        CHECKOUT_FAIL: "Checkout failed.",
        RETURN_SUCCESS: "Book returned successfully!",
        RETURN_FAIL: "Return failed.",
        RESERVE_SUCCESS: "Book reserved successfully!",
        RESERVE_FAIL: "Reservation failed.",
        LOAD_SUCCESS: "Books loaded successfully!",
        LOAD_FAIL: "Failed to load books. Please refresh the page.",
        NETWORK_ERROR: "Network connection issue. Please check your internet connection."
    };

    async function handleFormSubmission(endpoint, formData, form, successMessage) {
        try {
            const response = await fetch(`http://localhost:3000/users/${endpoint}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
                credentials: 'include'
            });

            if (response.ok) {
                const result = await response.json();
                alert(successMessage);
                form.reset();
                if (endpoint !== 'register') {
                    await loadBooks();
                }
                return true;
            }

            const errorData = await response.json();
            alert(errorData.error);
            return false;
        } catch (error) {
            console.error(`${endpoint} error:`, error);
            alert(MESSAGES.NETWORK_ERROR);
            return false;
        }
    }

    function isValidEmail(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }

    registerForm.addEventListener('submit', async function(event) {
        event.preventDefault();
        const formData = {
            name: document.getElementById('name').value.trim(),
            email: document.getElementById('email').value.trim()
        };

        if (!formData.name || !formData.email) {
            alert('Name and email are required.');
            return;
        }

        if (!isValidEmail(formData.email)) {
            alert('Please enter a valid email address.');
            return;
        }

        await handleFormSubmission('register', formData, registerForm, MESSAGES.REGISTER_SUCCESS);
    });

    checkoutForm.addEventListener('submit', async function(event) {
        event.preventDefault();
        const formData = {
            book_id: checkoutSelect.value,
            email: document.getElementById('checkout-email').value.trim()
        };

        if (!formData.book_id) {
            alert('Please select a book.');
            return;
        }

        if (!isValidEmail(formData.email)) {
            alert('Please enter a valid email address.');
            return;
        }

        await handleFormSubmission('checkout', formData, checkoutForm, MESSAGES.CHECKOUT_SUCCESS);
    });

    returnForm.addEventListener('submit', async function(event) {
        event.preventDefault();
        const formData = {
            book_id: returnSelect.value,
            email: document.getElementById('return-email').value.trim()
        };

        if (!formData.book_id) {
            alert('Please select a book.');
            return;
        }

        if (!isValidEmail(formData.email)) {
            alert('Please enter a valid email address.');
            return;
        }

        await handleFormSubmission('return', formData, returnForm, MESSAGES.RETURN_SUCCESS);
    });

    reserveForm.addEventListener('submit', async function(event) {
        event.preventDefault();
        const formData = {
            book_id: reserveSelect.value,
            email: document.getElementById('reserve-email').value.trim()
        };

        if (!formData.book_id) {
            alert('Please select a book.');
            return;
        }

        if (!isValidEmail(formData.email)) {
            alert('Please enter a valid email address.');
            return;
        }

        await handleFormSubmission('reserve', formData, reserveForm, MESSAGES.RESERVE_SUCCESS);
    });

    async function loadBooks() {
        try {
            const response = await fetch('http://localhost:3000/books', {
                credentials: 'include'
            });

            if (!response.ok) {
                throw new Error('Failed to fetch books');
            }

            const books = await response.json();
            updateBookDisplay(books);
        } catch (error) {
            console.error('Error loading books:', error);
            alert(MESSAGES.LOAD_FAIL);
        }
    }

    function updateBookDisplay(books) {
        booksContainer.innerHTML = '';
        checkoutSelect.innerHTML = '<option value="">Select a book</option>';
        returnSelect.innerHTML = '<option value="">Select a book</option>';
        reserveSelect.innerHTML = '<option value="">Select a book</option>';

        books.forEach(book => {
            renderBook(book);
            appendToSelect(book);
        });
    }

    function renderBook(book) {
        const bookItem = document.createElement('div');
        bookItem.classList.add('book-item');
        bookItem.setAttribute('data-id', book.id); // Set the data attribute for internal usage
        const status = book.status.toLowerCase();
        bookItem.innerHTML = `
            <div class="book-details">
                <h3>${book.title}</h3>
                <p>Author: ${book.author}</p>
                <p>ISBN: ${book.isbn}</p>
                <p>Status: ${status}</p>
            </div>
        `;
        booksContainer.appendChild(bookItem);
    }

    function appendToSelect(book) {
        const option = document.createElement('option');
        option.value = book.id;
        option.textContent = `${book.title} by ${book.author}`;

        const status = book.status.toLowerCase();

        if (status === 'available') {
            checkoutSelect.appendChild(option.cloneNode(true));
            reserveSelect.appendChild(option.cloneNode(true));
        }
        if (status === 'borrowed') {
            returnSelect.appendChild(option.cloneNode(true));
        }
    }

    loadBooks();
});