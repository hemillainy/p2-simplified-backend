CREATE TABLE "users" (
     "uuid" UUID NOT NULL,
     "email" character varying(256) NOT NULL,
     "password" character varying(256) NOT NULL,
     "name" character varying(256) NOT NULL,
     "document" character varying(15) NOT NULL,
     "wallet" float default 0 NOT NULL,
     "common_user" bool default true,
     PRIMARY KEY ("uuid"),
     UNIQUE("email", "document")
);

CREATE TABLE "transfers"(
   "uuid" UUID NOT NULL,
   "value" float NOT NULL,
   "payer" UUID NOT NULL,
   "payee" UUID NOT NULL,
   "status" character varying(10) NOT NULL,
   PRIMARY KEY ("uuid")
)