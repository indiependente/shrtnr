db.createUser(
    {
        user: "frank",
        pwd: "password",
        roles: [
            {
                role: "readWrite",
                db: "shrtnr"
            },
        ],
    },
);
db.getCollection('urls').createIndex({ "url": 1 });
