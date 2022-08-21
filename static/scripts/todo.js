function getAllTodos() {
    const endpoint = "/api/todo/all"
    fetch(endpoint, {
        method: "GET"
    }).then((response) => response.json())
    .then((data) => {
        const tbody = document.getElementsByTagName('tbody')[0];
        tbody.querySelectorAll('*').forEach(kid => kid.remove());
        for (let i = 0; i < data.length; ++i) {
            const tr = document.createElement('tr');
            const task = document.createElement('td');
            const deadline = document.createElement('td');
            const deleteColumn = document.createElement('td');
            const markDoneColumn = document.createElement('td');
            const deleteButton = document.createElement('button');
            const markDoneButton = document.createElement('button');
            deleteButton.onclick = function() {
                deleteTodo(data[i].id)
            }
            markDoneButton.onclick = function() {
                markDone(data[i]);
            }
            deleteButton.innerText = "delete"
            markDoneButton.innerText = "done";
            deleteColumn.appendChild(deleteButton);
            markDoneColumn.appendChild(markDoneButton);
            task.innerText = data[i].task;
            deadline.innerText = trimTimestamp(data[i].deadline);
            tr.appendChild(task);
            tr.appendChild(deadline);
            tr.appendChild(deleteColumn);
            tr.appendChild(markDoneColumn);
            tbody.appendChild(tr);
        }
    }).catch((err) => console.log(err));
}

function trimTimestamp(input) {
    const suffix = "T00:00:00Z";
    return input.substring(0, input.length - suffix.length);
}

function getCurrentDate() {
    let today = new Date();
    const dd = String(today.getDate()).padStart(2, '0');
    const mm = String(today.getMonth() + 1).padStart(2, '0'); //January is 0!
    const yyyy = today.getFullYear();
    today = yyyy + '-' + mm + '-' + dd;
    return today;
}

function markDone(todo) {
    const endpoint = "/api/todo/"
    fetch(endpoint, {
        method: "PUT",
        body: JSON.stringify({
            id: todo.id,
            status: "done",
            deadline: todo.deadline,
            task: todo.task,
            completedAt: getCurrentDate()                    
        })
    }).then(() => {
        window.location.reload();
    })
    .catch((err) => console.log(err));
}

function deleteTodo(id) {
    const endpoint = "/api/todo/"
    fetch(endpoint, {
        method: "DELETE",
        body: JSON.stringify({
            id
        })
    }).then(() => {
        window.location.reload();
    })
    .catch((err) => console.log(err));
}

function createTodo(event) {

    event.preventDefault();
    const endpoint = "/api/todo/"
    const task = document.getElementById("todo-text").value;
    const deadline = document.getElementById("todo-date").value;
    fetch(endpoint, {
        method: "POST",
        body: JSON.stringify({
            task,
            deadline
        })
    }).then((response) => response.json())
    .then(() => {
        window.location.reload();
    }).catch((err) => console.log(err));
}

document.addEventListener('DOMContentLoaded', function() {
    document.getElementById("todo-submit").addEventListener('click', createTodo);
    getAllTodos();
});

