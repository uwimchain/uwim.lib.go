# uwim.lib.go

<h2>Установка:</h2> 

```go
go get github.com/uwimchain/uwim.lib.go
```

<h2>Генерация мнемонической фразы</h2>

```go
mnemonic, err := Uwim.GenerateMnemonic()
```

Для генерации публичного, секретного ключей или адреса из мнемонической фразы можно использовать готовую мнемофразу.

<h2>Генерация Seed строки из мнемонической фразы</h2>
  
```go
seed, err := Uwim.SeedFromMnemonic(mnemonic)
```
<h2>Генерация секретного ключа из Seed строки или мнемонической фразы</h2>

```go
secret_key, err := Uwim.SecretKeyFromSeed(seed);<br><br>
secret_key, err := Uwim.SecretKeyFromMnemonic(mnemonic)
```
<h2>Генерация публичного из секретного ключа или мнемонической фразы</h2>

```go
public_key, err := Uwim.PublicKeyFromSecretKey(secret_key)
public_key, err := Uwim.PublicKeyFromMnemonic(mnemonic)
```
<h2>Генерация адреса пользователя из публичного ключа или мнемонической фразы</h2>

Для генерации адреса можно использовать публичный ключ или мнемоническую фразу, а также необходимо указать один из трёх доступных префиксов, если вы укажите какой-либо другой префикс, то функция вернёт ошибку<br><br>

<h3>Генерация адреса с префиксом "uw" - адрес кошелька пользователя</h3>

```go
uw_address, err := Uwim.AddressFromPublicKey(public_key, "uw")
uw_address, err := Uwim.AddressFromMnemonic(mnemonic, "uw")
```
<h3>Генерация адреса с префиксом "sc" - адрес смарт-контракта</h3>

```go
sc_address, err := Uwim.AddressFromPublicKey(public_key, "sc")
sc_address, err := Uwim.AddressFromMnemonic(mnemonic, "sc")
```
<h3>Генерация адреса с префиксом "nd" - адрес ноды</h3>

```go
nd_address, err := Uwim.AddressFromPublicKey(public_key, "nd")
nd_address, err := Uwim.AddressFromMnemonic(mnemonic, "nd")
```
<h2>Получение RAW строки транзакции для отправки в API блокчейна</h2>

Для того, чтобы сгенерировать RAW строку транзакции, вам необходимо указать такие данные как:

<ul>
  Мнемоническая фраза (отправителя транзакции);<br>
  Адрес отправителя (должен быть сгенерирован из мнемонической фразы или же подходить к ней)<br>
  Адрес получателя;<br>
  Количество монет, которое вы хотите перевести (для некоторых типов транзакции или подтипов транзакции, количество монет может быть рано нулю);<br>
  Адрес получателя;<br>
  Обозначение токена, монеты которого вы хотите перевести (например: "uwm")<br>
  Подтип пранзакции (например: "default_transaction")<br>
  Данные комментария к транзакции в формате JSON(для каждого типа или подтипа транзакции указываются свои данные комметрария или же не указываются совсем);<br>
  Тип пранзакции (Число 1 или 3);
</ul>
  
```go
transaction_raw, err := Uwim.GetRawTransaction(
    mnemonic,
    sender_address,
    recipient_addres,
    amount,
    token_label,
    transaction_comment_title,
    transaction_comment_data,
    transaction_type
)
```
