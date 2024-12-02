document.addEventListener('DOMContentLoaded', () => {
  const totalBooksElement = document.getElementById('total-books');
  const booksAvailableElement = document.getElementById('books-available');
  const totalMembersElement = document.getElementById('total-members');
  const activeCheckoutsElement = document.getElementById('active-checkouts');

  // Fetch dashboard stats
  fetch('/api/dashboard/stats', {
    method : 'GET',
    headers : {'Accept' : 'application/json'},
    credentials : 'same-origin'
  })
      .then(async response => {
        if (!response.ok) {
          const text = await response.text();
          console.error('Response not OK:', response.status, text);
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const contentType = response.headers.get('content-type');
        if (!contentType || !contentType.includes('application/json')) {
          const text = await response.text();
          console.error('Invalid content type:', contentType,
                        'Response:', text);
          throw new TypeError("Response was not JSON");
        }
        return response.json();
      })
      .then(data => {
        console.log('Received data:', data); // Debug log
        if (data.success) {
          totalBooksElement.textContent = data.data.totalBooks;
          booksAvailableElement.textContent = data.data.availableBooks;
          totalMembersElement.textContent = data.data.totalMembers;
          activeCheckoutsElement.textContent = data.data.activeCheckouts;
        } else {
          console.error('Success false in response:', data);
          throw new Error(data.message || 'Unknown error');
        }
      })
      .catch(error => {
        console.error('Error fetching dashboard stats:', error);
        // Set error state in UI
        totalBooksElement.textContent = 'Error';
        booksAvailableElement.textContent = 'Error';
        totalMembersElement.textContent = 'Error';
        activeCheckoutsElement.textContent = 'Error';
      });
});
