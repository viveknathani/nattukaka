function getById(id) { 
    return document.getElementById(id); 
}

function login(event) {

    event.preventDefault();
    let email = getById("login-email").value;
    let password = getById("login-password").value;

    const endpoint = "/api/user/login/"
    fetch(endpoint, {
        method: "POST",
        body: JSON.stringify({
            email,
            password
        })
    }).then((response) => response.json())
    .then((data) => {
        if (data.message === "ok") {
            localStorage.setItem("isAuthenticated", "true")
            redirectIfNeeded();
        } else {
            msg = getById("login-message");
            msg.innerText = "oops! " + data.message;
            msg.style.color = 'red';
            
        }
    }).catch((err) => console.log(err));
}

document.addEventListener('DOMContentLoaded', function() {
    document.querySelector('#login-submit').addEventListener('click', login);
});
Footer
