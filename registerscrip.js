function register() {
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const phoneNumber = document.getElementById('phoneNumber').value;
    const city = document.getElementById('city').value;
    const state = document.getElementById('state').value;
    const password = document.getElementById('password').value;
    const errormsg = document.getElementById('result');
  
    const newUser = {
      username: username,
      email: email,
      phoneNumber: phoneNumber,
      city: city,
      state: state,
      password: password
    };
  
    fetch("http://localhost:8082/people", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(newUser)
    })
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    })
    .then(data => {
     errormsg.innerHTML="person created  successfully"
    })
    .catch(error => {
      errormsg.innerHTML('error:', error);
    });
  
    return false;
  }
  