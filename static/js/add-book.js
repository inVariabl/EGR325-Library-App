// Add Book form handler
document.getElementById('addBookForm').addEventListener('submit', async (e) => {
  e.preventDefault();

  const form = e.target;
  const submitButton = form.querySelector('button[type="submit"]');
  const errorAlert = document.getElementById('errorAlert');
  const successAlert = document.getElementById('successAlert');
  const errorMessage = document.getElementById('errorMessage');

  // Hide alerts
  errorAlert.style.display = 'none';
  successAlert.style.display = 'none';

  // Set loading state
  submitButton.setAttribute('aria-busy', 'true');
  submitButton.disabled = true;

  try {
    const formData = {
      isbn : form.isbn.value,
      title : form.title.value,
      author : form.author.value,
      publisher : form.publisher.value,
      publication_year : parseInt(form.publication_year.value) || null,
      category : form.category.value,
      language : form.language.value,
      pages : parseInt(form.pages.value) || null,
      available_copies : parseInt(form.available_copies.value) || 1,
      total_copies : parseInt(form.total_copies.value) || 1,
      shelf_location : form.shelf_location.value,
    };

    const response = await fetch('/api/books', {
      method : 'POST',
      headers : {
        'Content-Type' : 'application/json',
      },
      credentials : 'same-origin',
      body : JSON.stringify(formData),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Failed to add book');
    }

    // Show success message
    successAlert.style.display = 'block';
    form.reset();

    // Hide success message after 3 seconds
    setTimeout(() => { successAlert.style.display = 'none'; }, 3000);

  } catch (error) {
    errorMessage.textContent = error.message;
    errorAlert.style.display = 'block';
  } finally {
    // Reset loading state
    submitButton.removeAttribute('aria-busy');
    submitButton.disabled = false;
  }
});

// Logout handler
async function handleLogout(e) {
  e.preventDefault();

  try {
    const response = await fetch(
        '/api/logout', {method : 'POST', credentials : 'same-origin'});

    if (response.ok) {
      window.location.href = '/login';
    } else {
      console.error('Logout failed');
    }
  } catch (error) {
    console.error('Error during logout:', error);
  }
}
