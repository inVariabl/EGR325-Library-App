// Login Form Handler
document.addEventListener('DOMContentLoaded', () => {
  const loginForm = document.getElementById('loginForm');
  const errorMessage = document.getElementById('error-message');

  if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();

      const formData = {
        username : loginForm.username.value,
        password : loginForm.password.value
      };

      try {
        const response = await fetch('/login', {
          method : 'POST',
          headers : {'Content-Type' : 'application/json'},
          body : JSON.stringify(formData)
        });

        const data = await response.json();

        if (data.success) {
          window.location.href = data.data.redirect;
        } else {
          if (errorMessage) {
            errorMessage.textContent = data.message;
            errorMessage.style.display = 'block';
          }
        }
      } catch (error) {
        console.error('Login error:', error);
        if (errorMessage) {
          errorMessage.textContent =
              'An error occurred during login. Please try again.';
          errorMessage.style.display = 'block';
        }
      }
    });
  }
});
