const form = document.getElementById('book-form');
const responseMessage = document.getElementById('response-message');
const booksList = document.getElementById('books-list');

// Function to fetch and display the list of books
const fetchBooks = async () => {
    try {
        // GET the list of books from /books endpoint
        const res = await fetch('/book', {
            method: 'GET'
        });
        const books = await res.json();

        // Clear the current list
        booksList.innerHTML = '';

        // Populate the list with fetched books
        books.forEach(book => {
            const listItem = document.createElement('li');
            listItem.textContent = book.book_name;
            booksList.appendChild(listItem);
        });
    } catch (err) {
        console.error('Failed to fetch books:', err);
    }
};

// Event listener for form submission
form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const bookName = document.getElementById('book-name').value;

    try {
        const res = await fetch('/book', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ book_name: bookName }),
        });

        const data = await res.json();
        responseMessage.textContent = data.message;

        // After adding a new book, fetch the updated list
        await fetchBooks();
    } catch (err) {
        responseMessage.textContent = 'Failed to submit book name.';
    }
});

// Fetch the initial list of books when the page loads
fetchBooks();