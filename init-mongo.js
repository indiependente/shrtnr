// Create user
dbAdmin = db.getSiblingDB("admin");
dbAdmin.createUser({
  user: "frank",
  pwd: "password",
  roles: [{ role: "userAdminAnyDatabase", db: "admin" }],
  mechanisms: ["SCRAM-SHA-1"],
});

// Authenticate user
dbAdmin.auth({
  user: "frank",
  pwd: "password",
  mechanisms: ["SCRAM-SHA-1"],
  digestPassword: true,
});

// Create DB and collection
db = new Mongo().getDB("shrtnr");
db.createCollection("urls", { capped: false });
db.getCollection("urls").createIndex({ "url": 1 });
