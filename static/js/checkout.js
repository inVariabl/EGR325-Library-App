// static/js/checkout.js
document.addEventListener('DOMContentLoaded', () => {
  // Elements
  const bookSearch = document.getElementById('bookSearch');
  const memberSearch = document.getElementById('memberSearch');
  const bookInfo = document.getElementById('bookInfo');
  const confirmCheckout = document.getElementById('confirmCheckout');
  const errorAlert = document.getElementById('errorAlert');
  const successAlert = document.getElementById('successAlert');
  const errorMessage = document.getElementById('errorMessage');

  // Initialize checkout date to today
  const checkoutDateInput = document.getElementById('checkoutDate');
  const today = new Date().toISOString().split('T')[0];
  checkoutDateInput.value = today;

  // Set default due date to 14 days from today
  const dueDateInput = document.getElementById('dueDate');
  const defaultDueDate = new Date();
  defaultDueDate.setDate(defaultDueDate.getDate() + 14);
  dueDateInput.value = defaultDueDate.toISOString().split('T')[0];

  // Hide alerts and book info initially
  errorAlert.style.display = 'none';
  successAlert.style.display = 'none';
  bookInfo.style.display = 'none';

  // Book search with debounce
  let bookSearchTimeout;
  bookSearch.addEventListener('input', (e) => {
    clearTimeout(bookSearchTimeout);
    bookSearchTimeout = setTimeout(async () => {
      const query = e.target.value.trim();
      if (query.length < 2) {
        bookInfo.style.display = 'none';
        return;
      }

      try {
        const response =
            await fetch(`/api/books/search?q=${encodeURIComponent(query)}`);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();

        if (!data.success) {
          throw new Error(data.message || 'Error searching for book');
        }

        if (data.data && data.data.length > 0) {
          const book = data.data[0]; // Use the first result
          displayBookInfo(book);
          bookInfo.style.display = 'block';
        } else {
          bookInfo.style.display = 'none';
          showError('No books found matching your search');
        }
      } catch (error) {
        console.error('Book search error:', error);
        showError('Error searching for book: ' + error.message);
        bookInfo.style.display = 'none';
      }
    }, 300);
  });

  // Member search with debounce
  let memberSearchTimeout;
  memberSearch.addEventListener('input', (e) => {
    clearTimeout(memberSearchTimeout);
    memberSearchTimeout = setTimeout(async () => {
      const query = e.target.value.trim();
      if (query.length < 2) {
        clearMemberInfo();
        return;
      }

      try {
        const response =
            await fetch(`/api/members/search?q=${encodeURIComponent(query)}`);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();

        if (!data.success) {
          throw new Error(data.message || 'Error searching for member');
        }

        if (data.data && data.data.length > 0) {
          const member = data.data[0]; // Use the first result
          displayMemberInfo(member);
        } else {
          clearMemberInfo();
          showError('No members found matching your search');
        }
      } catch (error) {
        console.error('Member search error:', error);
        showError('Error searching for member: ' + error.message);
        clearMemberInfo();
      }
    }, 300);
  });

  // Handle checkout
  confirmCheckout.addEventListener('click', async () => {
    // Validate required fields
    const bookId = document.getElementById('isbn').dataset.bookId;
    const memberId = document.getElementById('memberName').dataset.memberId;
    const dueDate = dueDateInput.value;

    if (!bookId || !memberId || !dueDate) {
      showError('Please select a book and member, and set a due date');
      return;
    }

    try {
      const checkoutData = {
        book_id : parseInt(bookId),
        member_id : parseInt(memberId),
        due_date : dueDate,
        notes : document.getElementById('notes').value
      };

      const response = await fetch('/api/checkout', {
        method : 'POST',
        headers : {
          'Content-Type' : 'application/json',
        },
        body : JSON.stringify(checkoutData)
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.message || 'Error processing checkout');
      }

      // Show success message and reset form
      showSuccess();
      resetForm();
      setTimeout(() => { window.location.href = '/dashboard'; }, 2000);
    } catch (error) {
      console.error('Checkout error:', error);
      showError('Error processing checkout: ' + error.message);
    }
  });

  // Helper functions
  function displayBookInfo(book) {
    document.getElementById('isbn').value = book.isbn;
    document.getElementById('isbn').dataset.bookId = book.id;
    document.getElementById('title').value = book.title;
    document.getElementById('author').value = book.author || '';
    document.getElementById('status').value = book.status || 'unknown';
    document.getElementById('shelfLocation').value = book.shelf_location || '';
  }

  function displayMemberInfo(memberResponse) {
    const member = memberResponse.member;
    document.getElementById('memberName').value = member.name;
    document.getElementById('memberName').dataset.memberId = member.member_id;
    document.getElementById('email').value = member.email || '';
    document.getElementById('phone').value = member.phone_number || '';
  }

  function clearMemberInfo() {
    document.getElementById('memberName').value = '';
    document.getElementById('memberName').dataset.memberId = '';
    document.getElementById('email').value = '';
    document.getElementById('phone').value = '';
  }

  function showError(message) {
    errorMessage.textContent = message;
    errorAlert.style.display = 'block';
    successAlert.style.display = 'none';
  }

  function showSuccess() {
    errorAlert.style.display = 'none';
    successAlert.style.display = 'block';
  }

  function resetForm() {
    bookSearch.value = '';
    memberSearch.value = '';
    bookInfo.style.display = 'none';
    clearMemberInfo();
    document.getElementById('notes').value = '';
    checkoutDateInput.value = today;
    dueDateInput.value = defaultDueDate.toISOString().split('T')[0];
  }
});
