db = db.getSiblingDB("desafio");
db.auth('root', 'root');

db = db.getSiblingDB('desafio');


db.users.insertMany([
  {
    username: "admin",
    password: "123",
    is_admin: true,
  },
  {
    username: "user",
    password: "123",
    is_admin: false,
  }
]);
