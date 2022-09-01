const editor = new SimpleMDE({ element: document.getElementById("editor") });

function getNoteContent() {
    const id = (new URL(document.location)).searchParams.get('id');
    const title = document.getElementById("title");
    fetch(`/api/note?id=${id}`, {
        method: 'GET'
    }).then((res) => res.json())
    .then((data) => {
        editor.value(data[0].content);
        localStorage.setItem("title", data[0].title);
        title.innerText = data[0].title;
    }).catch(err => console.log(err));
}

function updateNoteContent() {
    const id = (new URL(document.location)).searchParams.get('id');
    const message = document.getElementById("message");
    fetch('/api/note', {
        method: 'PUT',
        body: JSON.stringify({
            id: id,
            title: localStorage.getItem("title") || "",
            content: editor.value()
        })
    }).then((res) => res.json())
    .then(() => {
        message.innerText = "saved!"
        message.style.color = "green";
    }).catch(err => {
        console.log(err);
        message.innerText = "could not save!";
        message.style.color = "red";
    });

    setTimeout(() => {
        message.innerText = "";
    }, 1000);
}

document.addEventListener('DOMContentLoaded', function() {
    document.getElementById("note-update").addEventListener('click', updateNoteContent);
    document.getElementById("go-back").addEventListener('click', () => { window.history.back() });
    getNoteContent();
});