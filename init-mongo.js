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
