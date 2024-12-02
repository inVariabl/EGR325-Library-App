document.addEventListener('DOMContentLoaded', () => {
  const memberSearch = document.getElementById('memberSearch');
  const membersTableBody = document.getElementById('membersTableBody');
  const addMemberBtn = document.getElementById('addMemberBtn');
  const errorAlert = document.getElementById('errorAlert');
  const successAlert = document.getElementById('successAlert');
  const errorMessage = document.getElementById('errorMessage');
  const successMessage = document.getElementById('successMessage');

  // Load initial members list
  loadMembers();

  // Search members with debounce
  let searchTimeout;
  memberSearch.addEventListener('input', (e) => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      const query = e.target.value.trim();

      if (query.length < 2) {
        loadMembers(); // Load all members if search is cleared
        return;
      }

      fetch(`/api/members/search?q=${encodeURIComponent(query)}`, {
        method : 'GET',
        headers : {'Accept' : 'application/json'},
        credentials : 'same-origin'
      })
          .then(response => response.json())
          .then(data => {
            console.log('Search results:', data); // Debug log
            if (data.success) {
              renderMembersTable(data.data);
            } else {
              showError(data.message || 'Error searching members');
            }
          })
          .catch(error => {
            console.error('Error searching members:', error);
            showError('Error loading members: ' + error.message);
          });
    }, 300);
  });

  function loadMembers() {
    fetch('/api/members', {
      method : 'GET',
      headers : {'Accept' : 'application/json'},
      credentials : 'same-origin'
    })
        .then(response => response.json())
        .then(data => {
          if (data.success) {
            renderMembersTable(data.data);
          } else {
            showError('Error loading members');
          }
        })
        .catch(error => {
          console.error('Error:', error);
          showError('Error loading members');
        });
  }

  function renderMembersTable(members) {
    console.log('Rendering members:', members); // Debug log

    if (!members || members.length === 0) {
      membersTableBody.innerHTML = `
                <tr>
                    <td colspan="6" style="text-align: center;">No members found</td>
                </tr>
            `;
      return;
    }

    membersTableBody.innerHTML =
        members
            .map(memberResponse => {
              // Extract the actual member data
              const member = memberResponse.member || memberResponse;
              console.log('Processing member:', member); // Debug log

              return `
                <tr>
                    <td>${escapeHtml(member.member_id)}</td>
                    <td>${escapeHtml(member.name)}</td>
                    <td>${escapeHtml(member.email)}</td>
                    <td>${escapeHtml(member.phone_number || '')}</td>
                    <td>${formatDate(member.membership_date)}</td>
                    <td>
                        <div class="button-group">
                            <button class="outline secondary" onclick="editMember(${
                  member.member_id})">Edit</button> </div>
                    </td>
                </tr>
            `;
            })
            .join('');
  }

  // Helper functions
  function escapeHtml(unsafe) {
    if (unsafe === undefined || unsafe === null)
      return '';
    return unsafe.toString()
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
  }

  function formatDate(dateString) {
    if (!dateString)
      return '';
    try {
      return new Date(dateString).toLocaleDateString();
    } catch (e) {
      console.error('Error formatting date:', dateString, e);
      return dateString;
    }
  }

  function showError(message) {
    console.error('Error:', message); // Debug log
    errorMessage.textContent = message;
    errorAlert.style.display = 'block';
    successAlert.style.display = 'none';
  }

  function showSuccess(message) {
    successMessage.textContent = message;
    successAlert.style.display = 'block';
    errorAlert.style.display = 'none';
  }

  // Add new elements
  const addMemberModal = document.getElementById('addMemberModal');
  const addMemberForm = document.getElementById('addMemberForm');

  // Show add member modal
  addMemberBtn.addEventListener('click', () => {
    addMemberForm.reset();
    addMemberModal.showModal();
  });

  // Handle add member form submission
  addMemberForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    // Reset error messages
    document.querySelectorAll('small.error')
        .forEach(el => el.style.display = 'none');
    let hasError = false;

    // Validate name
    const name = addMemberForm.name.value.trim();
    if (!name) {
      document.getElementById('nameError').textContent = 'Name is required';
      document.getElementById('nameError').style.display = 'block';
      hasError = true;
    }

    // Validate email
    const email = addMemberForm.email.value.trim();
    const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!email) {
      document.getElementById('emailError').textContent = 'Email is required';
      document.getElementById('emailError').style.display = 'block';
      hasError = true;
    } else if (!emailPattern.test(email)) {
      document.getElementById('emailError').textContent =
          'Please enter a valid email address';
      document.getElementById('emailError').style.display = 'block';
      hasError = true;
    }

    // Validate phone number (optional but must be valid if provided)
    const phone = addMemberForm.phone.value.trim();
    const phonePattern = /^\(\d{3}\) \d{3}-\d{4}$/;
    if (phone && !phonePattern.test(phone)) {
      document.getElementById('phoneError').textContent =
          'Please use format: (123) 456-7890';
      document.getElementById('phoneError').style.display = 'block';
      hasError = true;
    }

    if (hasError)
      return;

    try {
      const formData = {
        name : name,
        email : email,
        phone_number : phone,
        address : addMemberForm.address.value.trim()
      };

      const response = await fetch('/api/members', {
        method : 'POST',
        headers : {
          'Content-Type' : 'application/json',
          'Accept' : 'application/json'
        },
        credentials : 'same-origin',
        body : JSON.stringify(formData)
      });

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.message || 'Error adding member');
      }

      showSuccess('Member added successfully');
      addMemberModal.close();
      loadMembers(); // Reload the members list
    } catch (error) {
      console.error('Error adding member:', error);
      showError('Error adding member: ' + error.message);
    }
  });

  // Add auto-formatting for phone number
  document.getElementById('phone').addEventListener('input', function(e) {
    let x =
        e.target.value.replace(/\D/g, '').match(/(\d{0,3})(\d{0,3})(\d{0,4})/);
    e.target.value =
        !x[2] ? x[1] : '(' + x[1] + ') ' + x[2] + (x[3] ? '-' + x[3] : '');
  });

  // Edit member elements
  const editMemberModal = document.getElementById('editMemberModal');
  const editMemberForm = document.getElementById('editMemberForm');

  // Handle edit member form submission
  editMemberForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    // Reset error messages
    document.querySelectorAll('small.error')
        .forEach(el => el.style.display = 'none');
    let hasError = false;

    // Validate name
    const name = editMemberForm.name.value.trim();
    if (!name) {
      document.getElementById('editNameError').textContent = 'Name is required';
      document.getElementById('editNameError').style.display = 'block';
      hasError = true;
    }

    // Validate email
    const email = editMemberForm.email.value.trim();
    const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!email) {
      document.getElementById('editEmailError').textContent =
          'Email is required';
      document.getElementById('editEmailError').style.display = 'block';
      hasError = true;
    } else if (!emailPattern.test(email)) {
      document.getElementById('editEmailError').textContent =
          'Please enter a valid email address';
      document.getElementById('editEmailError').style.display = 'block';
      hasError = true;
    }

    // Validate phone number (optional but must be valid if provided)
    const phone = editMemberForm.phone.value.trim();
    const phonePattern = /^\(\d{3}\) \d{3}-\d{4}$/;
    if (phone && !phonePattern.test(phone)) {
      document.getElementById('editPhoneError').textContent =
          'Please use format: (123) 456-7890';
      document.getElementById('editPhoneError').style.display = 'block';
      hasError = true;
    }

    if (hasError)
      return;

    const memberId = document.getElementById('editMemberId').value;

    try {
      const formData = {
        name : name,
        email : email,
        phone_number : phone,
        address : editMemberForm.address.value.trim()
      };

      const response = await fetch(`/api/members/${memberId}`, {
        method : 'PUT',
        headers : {
          'Content-Type' : 'application/json',
          'Accept' : 'application/json'
        },
        credentials : 'same-origin',
        body : JSON.stringify(formData)
      });

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.message || 'Error updating member');
      }

      showSuccess('Member updated successfully');
      editMemberModal.close();
      loadMembers(); // Reload the members list
    } catch (error) {
      console.error('Error updating member:', error);
      showError('Error updating member: ' + error.message);
    }
  });

  // Add phone formatting to edit form
  document.getElementById('editPhone').addEventListener('input', function(e) {
    let x =
        e.target.value.replace(/\D/g, '').match(/(\d{0,3})(\d{0,3})(\d{0,4})/);
    e.target.value =
        !x[2] ? x[1] : '(' + x[1] + ') ' + x[2] + (x[3] ? '-' + x[3] : '');
  });

  // Replace the global editMember function with this one
  window.editMember = async function(memberId) {
    try {
      const response = await fetch(`/api/members/${memberId}`, {
        method : 'GET',
        headers : {'Accept' : 'application/json'},
        credentials : 'same-origin'
      });

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.message || 'Error fetching member details');
      }

      const member = data.data;

      // Populate the edit form
      document.getElementById('editMemberId').value = member.member_id;
      document.getElementById('editName').value = member.name;
      document.getElementById('editEmail').value = member.email;
      document.getElementById('editPhone').value = member.phone_number || '';
      document.getElementById('editAddress').value = member.address || '';

      // Show the modal
      editMemberModal.showModal();
    } catch (error) {
      console.error('Error loading member details:', error);
      showError('Error loading member details: ' + error.message);
    }
  };
});
