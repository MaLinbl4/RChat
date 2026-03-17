import { GetContacts, GetMessagesByContact, SaveMessage } from '../wailsjs/go/main/App';

let currentContact = ""; // Переменная хранит, с кем мы сейчас говорим

// Функция переключения чата
window.selectContact = function(name) {
    currentContact = name;
    document.getElementById('current-chat-name').innerText = name;
    
    // Загружаем историю из Бэкенда (Go)
    GetMessagesByContact(name).then(messages => {
        const chatWindow = document.getElementById('chat-window');
        chatWindow.innerHTML = ''; // Чистим экран
        
        if (messages) {
            messages.forEach(m => {
                // Расшифровываем и выводим...
                renderMessage(m.text, m.type); 
            });
        }
    });
}

// При отрисовке контактов добавим событие клика
function renderContacts(contacts) {
    const list = document.getElementById('contacts-list');
    list.innerHTML = '';
    contacts.forEach(c => {
        const item = document.createElement('div');
        item.className = 'contact-item';
        item.innerHTML = `<div class="avatar">${c.name[0]}</div><div><b>${c.name}</b></div>`;
        item.onclick = () => window.selectContact(c.name);
        list.appendChild(item);
    });
}