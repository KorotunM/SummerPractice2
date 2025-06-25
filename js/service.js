import { generateCode, checkCode } from './sms.js';

let currentCode = null;
let currentPhone = null;

const screenPhone = document.getElementById('screen-phone');
const screenCode  = document.getElementById('screen-code');
const inputPhone  = document.getElementById('input-phone');
const inputCode   = document.getElementById('input-code');
const btnSend     = document.getElementById('btn-send');
const btnVerify   = document.getElementById('btn-verify');

btnSend.addEventListener('click', async () => {
  const phone = inputPhone.value.trim().replace('+', '');

  if (!phone) {
    alert('Номер телефона не может быть пустым.');
    return;
  }

  const code = generateCode();
  currentCode = code;
  currentPhone = phone;

  btnSend.disabled = true;
  btnSend.textContent = 'Отправка...';

  let result;
  try {
    const resp = await fetch('http://localhost:3000/api/send-sms', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ phone, code })
    });
    result = await resp.json();
  } catch (err) {
    result = { success: false, error: err.message };
  }

  btnSend.disabled = false;
  btnSend.textContent = 'Отправить код';

  if (result.success) {
    screenPhone.style.display = 'none';
    screenCode.style.display = 'block';
  } else {
    alert('Ошибка при отправке SMS: ' + result.error);
  }
});

btnVerify.addEventListener('click', () => {
  const input = inputCode.value.trim();
  if (!input) {
    alert('Введите код из SMS.');
    return;
  }
  if (checkCode(input, currentCode)) {
    alert('Успешно, код верный.');

    currentCode = null;
    currentPhone = null;
    inputPhone.value = '';
    inputCode.value = '';
    screenCode.style.display = 'none';
    screenPhone.style.display = 'block';
  } else {
    alert('Неверный код. Попробуйте ещё раз или запросите новый.');
  }
});
