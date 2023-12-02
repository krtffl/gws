# GWS 

gws is a simple application built with GO and HTMX that allows users to send positive messages to a beloved one that is 
undergoing a difficult situation, or for a special ocasion. the messages can optionally include an image, and can be accesssed
only after a secret passcode is correctly guessed.

## prerequisites

- docker compose

## run 

- clone repository

```bash
git clone https://github.com/krtffl/gws
```

- create yout own config

```bash
cp ./config/config.default.yaml ./config/config.yaml
```

and update your port, database connection parameters and challenges (secret passcodes)

- run the application and the database

```bash
make run
```

## usage

- open the app in the browser
- click on the box to leave your message
- leave your message

- click on the box to see the messages if you are the recipient
- enter the correct password
- read the messages

## database

the application uses a postgresql database for no specific reason. feel free to update the repository
and switch to a database of your choice. 

it stores the images directly in the database which might not be best way to do it


