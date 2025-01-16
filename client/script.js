const form = document.getElementById('book-form');
const responseMessage = document.getElementById('response-message');
const booksList = document.getElementById('books-list');
const input = document.getElementById('book-name');
const deleteResponse = document.getElementById('delete-response-message')
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

        // Fill the list with fetched books
        books.forEach(book => {
            booksList.innerHTML += `
        <div class="flex justify-between px-8 py-4">
        <li id="${book.id}">${book.book_name}</li> 
    
        <button data-id="${book.id}" onclick='deleteBook("${book.id}")' >❌</button>
        </div>
`;
            input.value = '';
        }
        )

        ;
    } catch (err) {
        console.error('Failed to fetch books:', err);
    }
};

async function deleteBook(bookId){
    try {
        const res = await fetch(`/book`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id: bookId }),
        });
        const data = await res.json();
        deleteResponse.textContent = "Book has been deleted successfully";


    }
    catch (err) {
        console.error('Failed to delete book:', err);
    }



    await fetchBooks();
}

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