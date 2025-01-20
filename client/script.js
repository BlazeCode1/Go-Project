const form = document.getElementById('book-form');
const responseMessage = document.getElementById('response-message');
const booksList = document.getElementById('books-list');
const input = document.getElementById('book-name');
const deleteResponse = document.getElementById('delete-response-message');

// Modal elements
const editModal = document.getElementById('edit-modal');
const editBookNameInput = document.getElementById('edit-book-name');
const saveEditButton = document.getElementById('save-edit-btn');
const cancelEditButton = document.getElementById('cancel-edit-btn');

let editingBookId = null; // To store the id of the book being edited

// Function to fetch and display the list of books
const fetchBooks = async () => {
    try {
        const res = await fetch('/book', {
            method: 'GET',
        });
        const books = await res.json();

        // Clear the current list
        booksList.innerHTML = '';

        // Fill the list with fetched books
        books.forEach(book => {
            booksList.innerHTML += `
                <div class="flex justify-between px-8 py-4">
                    <li class="flex-grow" id="${book.id}">${book.book_name}</li>
                    <div class="space-x-2">
                        <button
                            onclick='openEditModal("${book.id}", "${book.book_name}")'>
                            ✏️ 
                        </button>
                        <button
                            data-id="${book.id}"
                            onclick='deleteBook("${book.id}")'
                            >
                            ❌
                        </button>
                    </div>
                </div>
            `;
        });
        input.value = '';
    } catch (err) {
        console.error('Failed to fetch books:', err);
    }
};

// Open the edit modal
function openEditModal(bookId, bookName) {
    editingBookId = bookId;
    editBookNameInput.value = bookName;
    editModal.classList.remove('hidden');
}

// Close the edit modal
function closeEditModal() {
    editingBookId = null;
    editModal.classList.add('hidden');
}

// Save edited book
async function saveEdit() {
    if (!editingBookId) return;

    const newBookName = editBookNameInput.value;

    try {
        const res = await fetch(`/book`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id: editingBookId, book_name: newBookName }),
        });

        deleteResponse.textContent = "Book has been updated successfully.";
        closeEditModal();
        fetchBooks();


    } catch (err) {
        console.error('Failed to update book:', err);
    }
}

// Delete a book
async function deleteBook(bookId) {
    try {
        const res = await fetch(`/book`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id: bookId }),
        });
        await fetchBooks();
        deleteResponse.textContent = "Book has been deleted successfully.";

    } catch (err) {
        console.error('Failed to delete book:', err);
    }
}

// Event listener for modal buttons
saveEditButton.addEventListener('click', saveEdit);
cancelEditButton.addEventListener('click', closeEditModal);

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
        await fetchBooks();
    } catch (err) {
        responseMessage.textContent = 'Failed to submit book name.';
    }
});

// Fetch the initial list of books when the page loads
fetchBooks();