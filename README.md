# PL_lab3
GOLang (Client/server)

Для запуска сервера необходимо написать только порт.
Для запуска клиента нужно написать "ip:port -n (количество запросов)"

в случае совпадения ключей клиента и сервера у первого в консоли выведится <<<Key match>>>
  
Для параллельного подключения клиентов можно просто дописать go Client(...) и вывод, оставлять в программе не стал, т.к. логи невозможно читать
UPD: неправильно понял задание, поэтому дописал немного кода, а именно обработка трёх запросов одного клиента(до этого коннект у клиента закрывался после отправки одного сообщения).
UPD2: неправильно написал функцию рандома(было [0..9], сейчас как в питоне [0..1]*9)
