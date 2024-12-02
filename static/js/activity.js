document.addEventListener('DOMContentLoaded', () => {
  const totalBooksElement = document.getElementById('total-books');
  const booksAvailableElement = document.getElementById('books-available');
  const totalMembersElement = document.getElementById('total-members');
  const activeCheckoutsElement = document.getElementById('active-checkouts');
  const activityTableBody =
      document.querySelector('.activity-section table tbody');

  // Fetch dashboard stats
  fetch('/api/dashboard/stats', {
    method : 'GET',
    headers : {'Accept' : 'application/json'},
    credentials : 'same-origin'
  })
      .then(async response => {
        if (!response.ok) {
          const text = await response.text();
          console.error('Stats response not OK:', response.status, text);
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then(data => {
        if (data.success) {
          totalBooksElement.textContent = data.data.totalBooks;
          booksAvailableElement.textContent = data.data.availableBooks;
          totalMembersElement.textContent = data.data.totalMembers;
          activeCheckoutsElement.textContent = data.data.activeCheckouts;
        }
      })
      .catch(error => {
        console.error('Error fetching dashboard stats:', error);
      });

  // Fetch activity log
  fetch('/api/dashboard/activity', {
    method : 'GET',
    headers : {'Accept' : 'application/json'},
    credentials : 'same-origin'
  })
      .then(async response => {
        if (!response.ok) {
          const text = await response.text();
          console.error('Activity response not OK:', response.status, text);
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then(data => {
        console.log('Activity data:', data); // Debug log
        if (data.success && data.data.activities &&
            data.data.activities.length > 0) {
          activityTableBody.innerHTML = data.data.activities
                                            .map(activity => `
                <tr>
                    <td>${formatDateTime(activity.created_at)}</td>
                    <td>${escapeHtml(activity.action)}</td>
                    <td>${escapeHtml(activity.details)}</td>
                    <td>${escapeHtml(activity.username)}</td>
                </tr>
            `).join('');
        } else {
          activityTableBody.innerHTML =
              '<tr><td colspan="4">No recent activity</td></tr>';
        }
      })
      .catch(error => {
        console.error('Error fetching activity log:', error);
        activityTableBody.innerHTML =
            '<tr><td colspan="4">Error loading activity log</td></tr>';
      });

  // Helper functions
  function formatDateTime(dateTimeStr) {
    try {
      const date = new Date(dateTimeStr);
      return date.toLocaleString();
    } catch (e) {
      console.error('Error formatting date:', e);
      return dateTimeStr; // Return original string if parsing fails
    }
  }

  function escapeHtml(str) {
    if (!str)
      return '';
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
  }
});
