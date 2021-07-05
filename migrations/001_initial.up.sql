CREATE TABLE "users" (
     "uuid" UUID NOT NULL,
     "email" character varying(256) NOT NULL,
     "password" character varying(256) NOT NULL,
     "name" character varying(256) NOT NULL,
     "cpf" character varying(11) NOT NULL,
     "wallet" float default 0 NOT NULL,
     PRIMARY KEY ("uuid")
);

CREATE TABLE "shopkeepers"(
    "uuid" UUID NOT NULL,
    "email" character varying(256) NOT NULL,
    "password" character varying(256) NOT NULL,
    "name" character varying(256) NOT NULL,
    "cnpj" character varying(14) NOT NULL,
    "wallet" float NOT NULL,
    PRIMARY KEY ("uuid")
);

CREATE TABLE "transfers"(
   "uuid" UUID NOT NULL,
   "value" float NOT NULL,
   "payer" UUID NOT NULL,
   "payee" UUID NOT NULL,
   "status" character varying(10) NOT NULL,
   PRIMARY KEY ("uuid")
)