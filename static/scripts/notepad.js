const getByID = function (id) {
    return document.getElementById(id);
}

const saveContent = () => {
    let content = getByID('main').value;
    localStorage.setItem('content', content);
}

document.addEventListener('DOMContentLoaded', function() {
    getByID('save').addEventListener('click', saveContent);
    if (localStorage.getItem('content') !== undefined) {
        getByID('main').value = localStorage.getItem('content');
    }
});