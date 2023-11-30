create table if not exists "Messages"(
    "Id" varchar(36) not null primary key,
    "From" varchar(255) not null,
    "Message" text not null,
    "Memory" BYTEA,
	"CreatedAt" timestamp not null default now()
);
