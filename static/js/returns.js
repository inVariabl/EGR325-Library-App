document.addEventListener('DOMContentLoaded', () => {
  // Elements
  const checkoutSearch = document.getElementById('checkoutSearch');
  const checkoutInfo = document.getElementById('checkoutInfo');
  const confirmReturn = document.getElementById('confirmReturn');
  const errorAlert = document.getElementById('errorAlert');
  const successAlert = document.getElementById('successAlert');
  const errorMessage = document.getElementById('errorMessage');

  // Initialize return date to today
  const returnDateInput = document.getElementById('returnDate');
  const today = new Date().toISOString().split('T')[0];
  returnDateInput.value = today;

  // Hide alerts and checkout info initially
  errorAlert.style.display = 'none';
  successAlert.style.display = 'none';
  checkoutInfo.style.display = 'none';

  // Checkout search with debounce
  let searchTimeout;
  checkoutSearch.addEventListener('input', (e) => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(async () => {
      const query = e.target.value.trim();
      if (query.length < 2) {
        checkoutInfo.style.display = 'none';
        return;
      }

      try {
        const response =
            await fetch(`/api/checkouts/search?q=${encodeURIComponent(query)}`);
        const data = await response.json();

        if (!data.success) {
          throw new Error(data.message || 'Error searching for checkout');
        }

        if (data.data.length > 0) {
          const checkout = data.data[0]; // Use the first result
          displayCheckoutInfo(checkout);
          checkoutInfo.style.display = 'block';
        } else {
          checkoutInfo.style.display = 'none';
          showError('No active checkouts found');
        }
      } catch (error) {
        console.error('Search error:', error);
        showError('Error searching for checkout: ' + error.message);
        checkoutInfo.style.display = 'none';
      }
    }, 300);
  });

  // Handle return
  confirmReturn.addEventListener('click', async () => {
    // Validate required fields
    const checkoutId = document.getElementById('memberId').dataset.checkoutId;
    const condition = document.getElementById('condition').value;
    const returnDate = document.getElementById('returnDate').value;

    if (!checkoutId || !condition || !returnDate) {
      showError('Please select a checkout and fill in all required fields');
      return;
    }

    try {
      const returnData = {
        checkout_id : parseInt(checkoutId),
        condition : condition,
        return_date : returnDate,
        notes : document.getElementById('notes').value
      };

      const response = await fetch('/api/returns', {
        method : 'POST',
        headers : {
          'Content-Type' : 'application/json',
        },
        body : JSON.stringify(returnData)
      });

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.message);
      }

      // Show success message and reset form
      showSuccess();
      resetForm();
      setTimeout(() => { window.location.href = '/dashboard'; }, 2000);
    } catch (error) {
      console.error('Return error:', error);
      showError('Error processing return: ' + error.message);
    }
  });

  // Helper functions
  function displayCheckoutInfo(checkout) {
    console.log('Displaying checkout info:', checkout); // Debug log
    document.getElementById('bookTitle').value = checkout.title;
    document.getElementById('isbn').value = checkout.isbn;
    document.getElementById('memberName').value = checkout.member_name;
    document.getElementById('memberId').value = checkout.member_id;
    document.getElementById('memberId').dataset.checkoutId =
        checkout.checkout_id;
    document.getElementById('checkoutDate').value =
        formatDate(checkout.checkout_date);
    document.getElementById('dueDate').value = formatDate(checkout.due_date);
  }

  function formatDate(dateString) {
    return new Date(dateString).toISOString().split('T')[0];
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
    checkoutSearch.value = '';
    checkoutInfo.style.display = 'none';
    document.getElementById('notes').value = '';
    document.getElementById('condition').value = '';
    returnDateInput.value = today;
  }
});
